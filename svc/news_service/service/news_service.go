package service

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/adairxie/delinkcious/pb/news_service/pb"
	nm "github.com/adairxie/delinkcious/pkg/news_manager"
	"google.golang.org/grpc"
)

func Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "6060"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	natsHostname := os.Getenv("NATS_CLUSTER_SERVICE_HOST")
	natsPort := os.Getenv("NATS_CLUSTER_SERVICE_PORT")
	svc, err := nm.NewNewsManager(natsHostname, natsPort)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNewsServer(grpcServer, newNewsServer(svc))

	fmt.Printf("News service is listening on port %s...\n", port)
	err = grpcServer.Serve(listener)
	fmt.Println("Serve() failed", err)
}
