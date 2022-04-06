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

	cm "github.com/Maddosaurus/pallas/lib"
	gw "github.com/Maddosaurus/pallas/proto/pallas"
)

func allowedOrigin(origin string) bool {
	// src: https://fale.io/blog/2021/07/28/cors-headers-with-grpc-gateway
	// FIXME: Do *real* CORS origin check!

	// if viper.GetString("cors") == "*" {
	//     return true
	// }
	// if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
	//     return true
	// }
	// return false
	return true
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigin(r.Header.Get("Origin")) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func run() error {
	grpc_endpoint := cm.Getenv("PALLAS_GRPC_ENDPOINT", "server:50051")
	serv_endpoint := cm.Getenv("PALLAS_REST_SERVE_ENDPOINT", ":8081")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterOtpHandlerFromEndpoint(ctx, mux, grpc_endpoint, opts)
	if err != nil {
		return err
	}

	log.Printf("REST API server listening at %v", serv_endpoint)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(serv_endpoint, cors(mux))
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
