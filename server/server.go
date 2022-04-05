package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	pgsql "github.com/Maddosaurus/gotp/db"
	cm "github.com/Maddosaurus/gotp/lib"
	pb "github.com/Maddosaurus/gotp/proto/gotp"
	"github.com/gofrs/uuid"
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
	// FIXME: Encrypt Secret!
	if err := verifyEntry(newEntry); err != nil {
		return nil, fmt.Errorf("AddEntry: error verifying entry! %w", err)
	}
	if err := s.db.AddEntry(newEntry); err != nil {
		return nil, err
	}
	return newEntry, nil
}

func (s *gOTPServer) UpdateEntry(ctx context.Context, candidate *pb.OTPEntry) (*pb.OTPEntry, error) {
	// FIXME: Encrypt Secret!
	if tes, _ := s.db.GetEntry(&candidate.Uuid); tes == nil {
		return nil, errors.New("Update candidate not found in DB!")
	}
	if err := verifyEntry(candidate); err != nil {
		return nil, fmt.Errorf("AddEntry: error verifying entry! %w", err)
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

func verifyEntry(candiate *pb.OTPEntry) error {
	if _, err := uuid.FromString(candiate.Uuid); err != nil {
		return fmt.Errorf("verifyEntry: failed to verify UUID: %w", err)
	}
	if len(candiate.SecretToken) != 16 || strings.Compare(candiate.SecretToken, strings.ToUpper(candiate.SecretToken)) != 0 {
		return errors.New("Error while verifying token! Ensure it is 16 chars of upper case ASCII!")
	}
	// FIXME: Add OTP verification, but no error handling :/
	// https://github.com/xlzd/gotp/issues/18
	return nil
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
