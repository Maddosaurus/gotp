package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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

// TODO: Refactor?
func getAllEntries(client pb.GOTPClient) []pb.OTPEntry {
	var entries []pb.OTPEntry
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListEntries(ctx, &pb.UUID{Uuid: ""})
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
		entries = append(entries, *entry)
	}
	return entries
}

func printEntry(entry *pb.OTPEntry) {
	if entry.Type == pb.OTPEntry_TOTP {
		totp := otp.NewDefaultTOTP(entry.SecretToken)
		log.Printf("%v \tTOTP: %v", entry.Name, totp.Now())
		entry.Counter++ // FIXME: Update the counter on the server side ASAP!
	} else {
		hotp := otp.NewDefaultHOTP(entry.SecretToken)
		log.Printf("\tHOTP Mode - counter: %v", entry.Counter)
		log.Printf("%v \tHOTP: %v", entry.Name, hotp.At(int(entry.Counter)))
	}
}

func printOTP(client pb.GOTPClient) {
	var entries []pb.OTPEntry
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListEntries(ctx, &pb.UUID{Uuid: ""})
	if err != nil {
		log.Fatalf("%v.ListEntries(_) = _, %v", client, err)
	}
	i := 0
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Select OTP Entry:")
	for {
		entry, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListEntries(_) = _, %v", client, err)
		}
		entries = append(entries, *entry)
		fmt.Printf("%v - %v\n", i, entry.Name)
		i++
	}
	fmt.Printf("%v - Return to Main Menu\n--->", i)
	for {
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)

		if strings.Compare(strconv.Itoa(i), input) == 0 {
			break
		}

		ini, _ := strconv.Atoi(input)
		printEntry(&entries[ini])
		break
	}
}

func main() {
	flag.Parse()
	reader := bufio.NewReader(os.Stdin)

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials())) //FIXME: Use TLS!
	if err != nil {
		log.Fatalf("fail to dial %v", err)
	}
	defer conn.Close()
	client := pb.NewGOTPClient(conn)

	log.Printf("%v", getAllEntries(client))

	for {
		fmt.Print("1 - Get OTP\n2 - Add Entry\n3 - Quit\n---> ")
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)

		if strings.Compare("1", input) == 0 {
			printOTP(client)
		}

		if strings.Compare("3", input) == 0 {
			fmt.Println("Goodbye")
			os.Exit(0)
		}
	}

	//printEntries(client, &pb.UUID{Uuid: ""})
}
