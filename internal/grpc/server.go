package grpc

import (
	"log"
	"net"

	"github.com/dmmitrenko/card-validator/cards"
	"github.com/dmmitrenko/card-validator/internal/grpc/handler"
	"github.com/dmmitrenko/card-validator/internal/grpc/middleware"
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
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryInterceptor()),
	)

	apiClient := cards.NewApiClient("https://api.chargeblast.io/bin/")
	cardValidator := cards.NewCardValidator(apiClient)
	handler.NewCardValidatorHandler(grpcServer, cardValidator)

	log.Println("Starting gRPC server on", s.address)
	return grpcServer.Serve(listener)
}
