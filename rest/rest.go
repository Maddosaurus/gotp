package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	cm "github.com/Maddosaurus/gotp/lib"
	gw "github.com/Maddosaurus/gotp/proto/gotp"
)

func run() error {
	gotp_endpoint := cm.Getenv("GOTP_GRPC_ENDPOINT", "server:50051")
	serv_endpoint := cm.Getenv("GOTP_REST_ENDPOINT", ":8081")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterGOTPHandlerFromEndpoint(ctx, mux, gotp_endpoint, opts)
	if err != nil {
		return err
	}

	log.Printf("REST API server listening at %v", serv_endpoint)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(serv_endpoint, mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
