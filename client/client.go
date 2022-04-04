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

	pb "github.com/Maddosaurus/gotp/proto/gotp"
	"github.com/gofrs/uuid"
	otp "github.com/xlzd/gotp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "Server address in the format of host:port")
)

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

func printEntry(client pb.GOTPClient, entry *pb.OTPEntry) {
	log.Printf("Entry: %v // Updated at: %v", entry.Name, entry.UpdateTime.AsTime())
	if entry.Type == pb.OTPEntry_TOTP {
		totp := otp.NewDefaultTOTP(entry.SecretToken)
		log.Printf("TOTP: %v", totp.Now())
	} else {
		hotp := otp.NewDefaultHOTP(entry.SecretToken)
		log.Printf("\tHOTP Mode - counter: %v", entry.Counter)
		log.Printf("HOTP: %v", hotp.At(int(entry.Counter)))
		entry.Counter++
		entry.UpdateTime = timestamppb.Now()
		log.Printf("Updating Server side")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := client.UpdateEntry(ctx, entry)
		if err != nil {
			log.Fatalf("%v.UpdateEntry(_) = _, %v", client, err)
		}
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
		printEntry(client, &entries[ini])
		break
	}
}

func addEntry(client pb.GOTPClient) {
	uid, err := uuid.NewV4()
	if err != nil {
		log.Println("Error while creating entry UUID!")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Entry Name:")
	entry_name, _ := reader.ReadString('\n')
	entry_name = strings.Replace(entry_name, "\n", "", -1)

	fmt.Println("Seed:")
	seed, _ := reader.ReadString('\n')
	seed = strings.Replace(seed, "\n", "", -1)

	new_entry := pb.OTPEntry{
		Uuid:        uid.String(),
		Name:        entry_name,
		SecretToken: seed,
	}

	fmt.Println("OTP Type (HOTP or TOTP): ")
	entry_type, _ := reader.ReadString('\n')
	entry_type = strings.Replace(entry_type, "\n", "", -1)

	if strings.Compare("HOTP", strings.ToUpper(entry_type)) == 0 {
		fmt.Println("Starting Counter (default: 1): ")
		counter, _ := reader.ReadString('\n')
		counter = strings.Replace(counter, "\n", "", -1)
		counter_i, _ := strconv.Atoi(counter)
		new_entry.Type = pb.OTPEntry_HOTP
		new_entry.Counter = uint64(counter_i)
	} else {
		new_entry.Type = pb.OTPEntry_TOTP
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = client.AddEntry(ctx, &new_entry)
	if err != nil {
		log.Fatalf("%v.AddEntry(_) = _, %v", client, err)
	}
}

func deleteEntry(client pb.GOTPClient) {
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
		candidate := entries[ini]
		client.DeleteEntry(ctx, &candidate)
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
		fmt.Print("1 - Get OTP\n2 - Add Entry\n3 - Delete Entry\n4 - Quit\n---> ")
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)

		if strings.Compare("1", input) == 0 {
			printOTP(client)
		}

		if strings.Compare("2", input) == 0 {
			addEntry(client)
		}

		if strings.Compare("3", input) == 0 {
			deleteEntry(client)
		}

		if strings.Compare("4", input) == 0 {
			fmt.Println("Goodbye")
			os.Exit(0)
		}
	}
}
