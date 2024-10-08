package cards

import (
	"context"
	"time"

	"github.com/dmmitrenko/card-validator/domain"
)

type Validator interface {
	Validate(context.Context, *domain.Card) error
}

type CardValidator struct {
}

func NewCardValidator() *CardValidator {
	return &CardValidator{}
}

func (s *CardValidator) Validate(ctx context.Context, card *domain.Card) error {
	if s.isExpired(card.ExpirationMonth, card.ExpirationYear) {
		return nil
	}

	if s.luhnCheck(card.Number) {
		return nil
	}

	return nil
}

func (s *CardValidator) isExpired(expMonth int, expYear int) bool {
	cardExpirationDate := time.Date(expYear, time.Month(expMonth)+1, 0, 23, 59, 59, 0, time.UTC)
	currentDate := time.Now()

	return currentDate.After(cardExpirationDate)
}

func (s *CardValidator) luhnCheck(cardNumber string) bool {
	sum := 0
	isSecond := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')

		if isSecond {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isSecond = !isSecond
	}

	return sum%10 == 0
}
