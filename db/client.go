package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	cm "github.com/Maddosaurus/gotp/lib"
	pb "github.com/Maddosaurus/gotp/proto/gotp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PgSQL struct {
	db *sql.DB
}

func (p *PgSQL) InitDB() (*PgSQL, error) {
	server := cm.Getenv("GOTP_DB_SERVER", "localhost")
	port := cm.Getenv("GOTP_DB_PORT", "5432")
	db_name := cm.Getenv("GOTP_DB_NAME", "gotp")
	db_user := cm.Getenv("GOTP_DB_USER", "postgres")
	db_pass := cm.Getenv("GOTP_DB_PASS", "passpass")

	connection_string := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		db_user, db_pass, server, port, db_name,
	)

	log.Printf("Connecting to db: postgresql://%s:[REDACTED]@%s:%s/%s", db_user, server, port, db_name)

	db, err := sql.Open("pgx", connection_string)
	if err != nil {
		log.Printf("Could not connect to database: %v", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Printf("Unable to reach database %v", err)
		return nil, err
	}
	log.Print("Database connection established!")
	p.db = db
	return p, nil
}

func (p *PgSQL) AddEntry(entry *pb.OTPEntry) error {
	log.Printf("Serializing entry: %v", entry)
	result, err := p.db.Exec("INSERT INTO gotp (uuid, otptype, name, secret_token, counter, update_time) VALUES ($1, $2, $3, $4, $5, $6)",
		entry.Uuid, entry.Type, entry.Name, entry.SecretToken, entry.Counter, entry.UpdateTime.AsTime())
	if err != nil {
		log.Printf("Could not insert row: %v", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("Cout not get affected rows: %v", err)
		return err
	}
	log.Printf("%v rows affected", rows)
	return nil
}

func (p *PgSQL) GetEntry(uuid *string) (*pb.OTPEntry, error) {
	row := p.db.QueryRow("SELECT * FROM gotp LIMIT 1")
	r := &pb.OTPEntry{}
	t := time.Time{}
	if err := row.Scan(&r.Uuid, &r.Type, &r.Name, &r.SecretToken, &r.Counter, &t); err != nil {
		log.Printf("Could not scan row: %v", err)
		return nil, err
	}
	r.UpdateTime = timestamppb.New(t)
	log.Printf("Deserialized entry: %v", r)
	return r, nil
}
