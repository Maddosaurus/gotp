package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Maddosaurus/gotp/gotp"
	"google.golang.org/grpc"
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

func (s *gOTPServer) loadFeatures() {
	if err := json.Unmarshal(exampleData, &s.savedEntries); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
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

// Example data to have something to serve right now
var exampleData = []byte(`[{
		"uuid": "123abc",
		"name": "github",
		"secret_token": "11111"
	}, {
		"uuid": "22223abc",
		"name": "twitch",
		"secret_token": "22222"
	},{
		"uuid": "333333abc",
		"name": "google",
		"secret_token": "33333 "
	}]`)
