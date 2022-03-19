package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	pb "github.com/Maddosaurus/gotp/gotp"
	otp "github.com/xlzd/gotp"
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
		if entry.Type == pb.OTPEntry_HOTP {
			hotp := otp.NewDefaultHOTP(entry.SecretToken)
			log.Printf("\tHOTP Mode - counter: %v", entry.Counter)
			// FYI: This is the lookahead / skew window. This is configured at the server side!
			// We don't have to do anything here, the server synchronizes the counter on its side.
			// We basically increment it by 1 after a code generation event (e.g. user touch)
			// This needs to be synced with the gOTP server, though!
			log.Printf("\tHOTP: %v", hotp.At(1))
			// log.Printf("\tHOTP: %v", hotp.At(2))
			// log.Printf("\tHOTP: %v", hotp.At(3))
			// log.Printf("\tHOTP: %v", hotp.At(4))
			log.Printf("\tProvisioning URI: %v", hotp.ProvisioningUri("me", "gOTP", 0))
			log.Printf("\tUpdate Time: %v", entry.UpdateTime.AsTime())
		}
		if entry.Uuid == "1234" {
			totp := otp.NewDefaultTOTP(entry.SecretToken)
			log.Printf("\tTOTP: %v", totp.Now())
		}
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
