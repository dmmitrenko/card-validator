package cards

import (
	"context"
	"time"
	"unicode"

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
	number := card.Number

	if number[0] == '0' || len(number) < 12 || len(number) > 19 {
		return domain.ErrCardNumber
	}

	if !s.isValidNumber(number) {
		return domain.ErrCardNumber
	}

	if s.isExpired(card.ExpirationMonth, card.ExpirationYear) {
		return domain.ErrExpiredCard
	}

	if !s.luhnCheck(card.Number) {
		return domain.ErrLuhnAlgorithm
	}

	return nil
}

func (s *CardValidator) isValidNumber(cardNumber string) bool {
	for _, r := range cardNumber {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
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
