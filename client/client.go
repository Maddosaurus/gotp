package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	pb "github.com/Maddosaurus/gotp/gotp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "Server address in the format of host:port")
)

func printEntries(client pb.GOTPClient, uuid *pb.UUID) {
	log.Printf("Getting all entries for UUID %v", uuid)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListEntries(ctx, uuid)
	if err != nil {
		log.Fatalf("%v.ListEntries(_) = _, %v", client, err)
	}
	for {
		entry, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListEntries(_) = _, %v", client, err)
		}
		log.Printf("Feature: uuid: %v, name: %v, secret_token: %v", entry.Uuid, entry.Name, entry.SecretToken)
	}
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials())) //FIXME: Use TLS!
	if err != nil {
		log.Fatalf("fail to dial %v", err)
	}
	defer conn.Close()
	client := pb.NewGOTPClient(conn)

	printEntries(client, &pb.UUID{Uuid: "123123"})
}
