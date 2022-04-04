package main

import (
	"context"
	"log"
	"net"

	pgsql "github.com/Maddosaurus/gotp/db"
	cm "github.com/Maddosaurus/gotp/lib"
	pb "github.com/Maddosaurus/gotp/proto/gotp"
	_ "github.com/jackc/pgx/v4/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gOTPServer struct {
	pb.UnimplementedGOTPServer
	savedEntries []*pb.OTPEntry
	db           *pgsql.PgSQL
}

func (s *gOTPServer) ListEntries(uuid *pb.UUID, stream pb.GOTP_ListEntriesServer) error {
	for _, entry := range s.savedEntries {
		if err := stream.Send(entry); err != nil {
			return err
		}
	}
	return nil
}

func (s *gOTPServer) AddEntry(ctx context.Context, newEntry *pb.OTPEntry) (*pb.OTPEntry, error) {
	s.savedEntries = append(s.savedEntries, newEntry)
	return newEntry, nil
}

func (s *gOTPServer) UpdateEntry(ctx context.Context, candidate *pb.OTPEntry) (*pb.OTPEntry, error) {
	for i, entry := range s.savedEntries {
		if entry.Uuid == candidate.Uuid {
			s.savedEntries[i] = candidate
		}
	}
	return candidate, nil
}

func (s *gOTPServer) DeleteEntry(ctx context.Context, candidate *pb.OTPEntry) (*pb.OTPEntry, error) {
	for i, entry := range s.savedEntries {
		if entry.Uuid == candidate.Uuid {
			if i+1 < len(s.savedEntries) {
				s.savedEntries = append(s.savedEntries[:i], s.savedEntries[i+1])
			} else {
				s.savedEntries = append(s.savedEntries[:i])
			}

		}
	}
	return candidate, nil
}

func (s *gOTPServer) loadFeatures() {
	s.savedEntries = []*pb.OTPEntry{
		{
			Type:        pb.OTPEntry_TOTP,
			Uuid:        "1234",
			Name:        "Site1",
			SecretToken: "JBSWY3DPEHPK3PX3",
			UpdateTime:  timestamppb.Now(),
		}, {
			Type:        pb.OTPEntry_TOTP,
			Uuid:        "45678",
			Name:        "Twitch",
			SecretToken: "JBSWY3DPEHPK3PX4",
			UpdateTime:  timestamppb.Now(),
		}, {
			Type:        pb.OTPEntry_HOTP,
			Uuid:        "1234567dfcg",
			Name:        "CustomSite",
			SecretToken: "4S62BZNFXXSZLCRO",
			Counter:     1,
			UpdateTime:  timestamppb.Now(),
		},
	}

	e := s.savedEntries[2]
	log.Printf("Calling for entry: %v", e)
	err := s.db.AddEntry(e)
	if err != nil {
		log.Printf("Found error, continuing")
	}

}

func newServer() *gOTPServer {
	s := &gOTPServer{}
	s.db = &pgsql.PgSQL{}
	s.db.InitDB()
	s.loadFeatures()
	uu := "1234"
	s.db.GetEntry(&uu)
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
