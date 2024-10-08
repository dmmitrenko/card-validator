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

	if s.luhnSum(card.Number, false)%10 != 0 {
		return nil
	}

	return nil
}

func (s *CardValidator) isExpired(expMonth int, expYear int) bool {
	cardExpirationDate := time.Date(expYear, time.Month(expMonth)+1, 0, 23, 59, 59, 0, time.UTC)
	currentDate := time.Now()

	return currentDate.After(cardExpirationDate)
}

func (s *CardValidator) luhnSum(cardNumber string, isSecond bool) int {
	if len(cardNumber) == 0 {
		return 0
	}

	lastDigit := int(cardNumber[len(cardNumber)-1] - '0')

	if isSecond {
		lastDigit *= 2
		if lastDigit > 9 {
			lastDigit -= 9
		}
	}

	return lastDigit + s.luhnSum(cardNumber[:len(cardNumber)-1], !isSecond)
}
