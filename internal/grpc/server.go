package grpc

import (
	"log"
	"net"

	"github.com/dmmitrenko/card-validator/cards"
	"github.com/dmmitrenko/card-validator/internal/grpc/handler"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	address string
}

func NewGRPCServer(address string) *gRPCServer {
	return &gRPCServer{address: address}
}

func (s *gRPCServer) Run() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	cardValidator := cards.NewCardValidator()
	handler.NewCardValidatorHandler(grpcServer, cardValidator)

	log.Println("Starting gRPC server on", s.address)
	return grpcServer.Serve(listener)
}
