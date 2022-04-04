package main

import (
	"context"
	"errors"
	"log"
	"net"

	pgsql "github.com/Maddosaurus/gotp/db"
	cm "github.com/Maddosaurus/gotp/lib"
	pb "github.com/Maddosaurus/gotp/proto/gotp"
	"google.golang.org/grpc"
)

type gOTPServer struct {
	pb.UnimplementedGOTPServer
	db *pgsql.PgSQL
}

func (s *gOTPServer) ListEntries(uuid *pb.UUID, stream pb.GOTP_ListEntriesServer) error {
	entries, err := s.db.GetAllEntries()
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := stream.Send(entry); err != nil {
			return err
		}
	}
	return nil
}

func (s *gOTPServer) AddEntry(ctx context.Context, newEntry *pb.OTPEntry) (*pb.OTPEntry, error) {
	// FIXME: Add verification of UUID / Secret
	// FIXME: Salt & Hash Secret!
	if err := s.db.AddEntry(newEntry); err != nil {
		return nil, err
	}
	return newEntry, nil
}

func (s *gOTPServer) UpdateEntry(ctx context.Context, candidate *pb.OTPEntry) (*pb.OTPEntry, error) {
	// FIXME: Add verification of UUID / Secret
	// FIXME: Salt & Hash Secret!
	if tes, _ := s.db.GetEntry(&candidate.Uuid); tes == nil {
		return nil, errors.New("Update candidate not found in DB!")
	}
	if err := s.db.UpdateEntry(candidate); err != nil {
		return nil, err
	}
	return candidate, nil
}

func (s *gOTPServer) DeleteEntry(ctx context.Context, candidate *pb.OTPEntry) (*pb.OTPEntry, error) {
	if err := s.db.DeleteEntry(candidate); err != nil {
		return nil, err
	}
	return candidate, nil
}

func newServer() *gOTPServer {
	s := &gOTPServer{}
	s.db = &pgsql.PgSQL{}
	s.db.InitDB()
	return s
}

func main() {
	grpc_endpoint := cm.Getenv("GOTP_GRPC_ENDPOINT", ":50051")
	lis, err := net.Listen("tcp", grpc_endpoint)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGOTPServer(grpcServer, newServer())
	log.Printf("gRPC server listening at: %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
