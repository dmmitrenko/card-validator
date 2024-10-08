package middleware

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			log.Printf("gRPC method %s failed with an error: %v", info.FullMethod, err)

			if _, ok := status.FromError(err); !ok {
				return nil, status.Errorf(codes.Internal, "internal server error")
			}
			return nil, err
		}

		return resp, nil
	}
}
