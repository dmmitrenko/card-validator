package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmmitrenko/card-validator/cards"
	"github.com/dmmitrenko/card-validator/internal/config"
	"github.com/dmmitrenko/card-validator/internal/grpc/handler"
	"github.com/dmmitrenko/card-validator/internal/grpc/middleware"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	address string
	server  *grpc.Server
}

func NewGRPCServer(address string) *gRPCServer {
	return &gRPCServer{
		address: address,
		server:  grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryInterceptor())),
	}
}

func (s *gRPCServer) Run() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cfg := config.NewConfig()
	apiClient := cards.NewApiClient(cfg.APIURL)
	cardValidator := cards.NewCardValidator(apiClient)
	handler.NewCardValidatorHandler(s.server, cardValidator)

	log.Println("Starting gRPC server on", s.address)

	go func() {
		if err := s.server.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	s.waitForShutdown()
	return nil
}

func (s *gRPCServer) waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down gRPC server...")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.server.GracefulStop()
	log.Println("gRPC server stopped gracefully")
}
