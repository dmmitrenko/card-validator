package handler

import (
	"context"

	"github.com/dmmitrenko/card-validator/cards"
	"github.com/dmmitrenko/card-validator/domain"
	"github.com/dmmitrenko/card-validator/internal/grpc/proto"
	"google.golang.org/grpc"
)

type CardValidatorHandler struct {
	cardValidator cards.Validator
	proto.UnimplementedCardValidatorServer
}

func NewCardValidatorHandler(grpc *grpc.Server, cardValidator cards.Validator) {
	gRPCHandler := &CardValidatorHandler{
		cardValidator: cardValidator,
	}

	proto.RegisterCardValidatorServer(grpc, gRPCHandler)
}

func (h *CardValidatorHandler) ValidateCard(ctx context.Context, req *proto.CardValidationRequest) (*proto.CardValidationResponse, error) {
	card := &domain.Card{
		Number:          req.CardNumber,
		ExpirationMonth: int(req.ExpirationMonth),
		ExpirationYear:  int(req.ExpirationYear),
	}

	err := h.cardValidator.Validate(ctx, card)

	if err != nil {
		if codedErr, ok := err.(domain.CodedError); ok {
			return &proto.CardValidationResponse{
				Valid: false,
				Error: &proto.ErrorResponse{
					Message: err.Error(),
					Code:    codedErr.ErrorCode(),
				},
			}, nil
		}

		return nil, err
	}

	res := &proto.CardValidationResponse{
		Valid: true,
		Error: nil,
	}

	return res, nil
}
