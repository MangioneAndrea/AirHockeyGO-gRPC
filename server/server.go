package main

import (
	"fmt"
	"log"
	"net"

	"andrea.mangione.me/pong/server/positionpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Server running")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	positionpb.RegisterPositionServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}
