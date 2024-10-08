package handler_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/dmmitrenko/card-validator/cards"
	"github.com/dmmitrenko/card-validator/cards/mocks"
	"github.com/dmmitrenko/card-validator/domain"
	"github.com/dmmitrenko/card-validator/internal/grpc/handler"
	"github.com/dmmitrenko/card-validator/internal/grpc/middleware"
	"github.com/dmmitrenko/card-validator/internal/grpc/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestCardValidator_InvalidBin_Integration(t *testing.T) {
	mockApiClient := mocks.NewMockApiClient(func(iin string) error {
		return domain.ErrUnknownINN
	})

	server, lis := startTestServer(t, mockApiClient)
	defer server.Stop()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewCardValidatorClient(conn)

	req := &proto.CardValidationRequest{
		CardNumber:      "1234567890123452",
		ExpirationMonth: 12,
		ExpirationYear:  2025,
	}

	resp, err := client.ValidateCard(context.Background(), req)
	if err != nil {
		t.Fatalf("error during gRPC call: %v", err)
	}

	assert.False(t, resp.Valid, "Card should be valid according to Luhn but with invalid IIN")
	assert.True(t, resp.Error.Code == domain.ErrUnknownINN.ErrorCode())
}

func TestCardValidator_ThirdPartyServiceUnavailable_Integration(t *testing.T) {
	mockApiClient := mocks.NewMockApiClient(func(iin string) error {
		return fmt.Errorf("error response from IIN API.")
	})

	server, lis := startTestServer(t, mockApiClient)
	defer server.Stop()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewCardValidatorClient(conn)

	req := &proto.CardValidationRequest{
		CardNumber:      "1234567890123452",
		ExpirationMonth: 12,
		ExpirationYear:  2025,
	}

	res, err := client.ValidateCard(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, status.Errorf(codes.Internal, "internal server error"), err, "interceptor should return internal server error for unhandled errors")
}

func startTestServer(t *testing.T, mockApiClient *mocks.MockApiClient) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryInterceptor()),
	)

	cardValidator := cards.NewCardValidator(mockApiClient)
	handler.NewCardValidatorHandler(server, cardValidator)
	go func() {
		if err := server.Serve(lis); err != nil {
			t.Fatalf("failed to serve: %v", err)
		}
	}()

	return server, lis
}
