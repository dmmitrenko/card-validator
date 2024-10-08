package main

import (
	"github.com/dmmitrenko/card-validator/internal/config"
	"github.com/dmmitrenko/card-validator/internal/grpc"
)

func main() {
	cfg := config.NewConfig()
	grpcServer := grpc.NewGRPCServer(":" + cfg.Port)
	grpcServer.Run()
}
