package main

import "github.com/dmmitrenko/card-validator/internal/grpc"

func main() {
	grpcServer := grpc.NewGRPCServer(":9000")
	grpcServer.Run()
}
