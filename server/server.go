package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	cm "github.com/Maddosaurus/gotp/lib"
	pb "github.com/Maddosaurus/gotp/proto/gotp"
	_ "github.com/jackc/pgx/v4/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gOTPServer struct {
	pb.UnimplementedGOTPServer
	savedEntries []*pb.OTPEntry
	db           *sql.DB
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

	// FIXME: REMOVE ALL OF THIS!
	e := s.savedEntries[2]
	log.Printf("Serializing entry: %v", e)
	result, err := s.db.Exec("INSERT INTO gotp (uuid, otptype, name, secret_token, counter, update_time) VALUES ($1, $2, $3, $4, $5, $6)",
		e.Uuid, e.Type, e.Name, e.SecretToken, e.Counter, e.UpdateTime.AsTime())
	if err != nil {
		log.Fatalf("Could not insert row: %v", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Cout not get affected rows: %v", err)
	}
	log.Printf("%v rows affected", rows)
}

func initDBConnection() *sql.DB {
	db, err := sql.Open("pgx", "postgresql://postgres:passpass@localhost:5432/gotp")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to reach database %v", err)
	}
	log.Print("Database connection established!")
	return db
}

func newServer() *gOTPServer {
	s := &gOTPServer{}
	s.db = initDBConnection()
	s.loadFeatures()
	s.getDBEntry()
	return s
}

func (s *gOTPServer) getDBEntry() pb.OTPEntry {
	row := s.db.QueryRow("SELECT * FROM gotp LIMIT 1")
	r := pb.OTPEntry{}
	t := time.Time{}
	if err := row.Scan(&r.Uuid, &r.Type, &r.Name, &r.SecretToken, &r.Counter, &t); err != nil {
		log.Fatalf("Could not scan row: %v", err)
	}
	r.UpdateTime = timestamppb.New(t)
	log.Printf("Deserialized entry: %v", r)
	return r
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
