package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Maddosaurus/gotp/gotp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	port = flag.Int("port", 50051, "Listen port of the server")
)

type gOTPServer struct {
	pb.UnimplementedGOTPServer
	savedEntries []*pb.OTPEntry
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
}

func newServer() *gOTPServer {
	s := &gOTPServer{}
	s.loadFeatures()
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGOTPServer(grpcServer, newServer())
	log.Printf("Server listening at: %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
