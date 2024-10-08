package cards_test

import (
	"context"
	"testing"
	"time"

	"github.com/dmmitrenko/card-validator/cards"
	"github.com/dmmitrenko/card-validator/cards/mocks"
	"github.com/dmmitrenko/card-validator/domain"
	"github.com/stretchr/testify/assert"
)

func TestCardValidator(t *testing.T) {
	mockApiClient := mocks.NewMockApiClient(func(iin string) error {
		return nil
	})
	v := cards.NewCardValidator(mockApiClient)

	t.Run("InvalidCardNumberStartsWithZero", func(t *testing.T) {
		card := &domain.Card{
			Number:          "0123456789012",
			ExpirationMonth: 12,
			ExpirationYear:  time.Now().Year() + 1,
		}

		err := v.Validate(context.Background(), card)
		assert.Equal(t, domain.ErrCardNumber, err)
	})

	t.Run("InvalidCardNumberTooShort", func(t *testing.T) {
		card := &domain.Card{
			Number:          "123456",
			ExpirationMonth: 12,
			ExpirationYear:  time.Now().Year() + 1,
		}

		err := v.Validate(context.Background(), card)
		assert.Equal(t, domain.ErrCardNumber, err)
	})

	t.Run("InvalidCardNumberTooLong", func(t *testing.T) {
		card := &domain.Card{
			Number:          "12345678901234567890",
			ExpirationMonth: 12,
			ExpirationYear:  time.Now().Year() + 1,
		}

		err := v.Validate(context.Background(), card)
		assert.Equal(t, domain.ErrCardNumber, err)
	})

	t.Run("InvalidCardNumberWithNonDigit", func(t *testing.T) {
		card := &domain.Card{
			Number:          "1234abcd5678",
			ExpirationMonth: 12,
			ExpirationYear:  time.Now().Year() + 1,
		}

		err := v.Validate(context.Background(), card)
		assert.Equal(t, domain.ErrCardNumber, err)
	})

	t.Run("InvalidMonth", func(t *testing.T) {
		card := &domain.Card{
			Number:          "4111111111111111",
			ExpirationMonth: 13,
			ExpirationYear:  time.Now().Year() + 1,
		}

		err := v.Validate(context.Background(), card)
		assert.Equal(t, domain.ErrMonthNumber, err)
	})

	t.Run("InvalidYear", func(t *testing.T) {
		card := &domain.Card{
			Number:          "4111111111111111",
			ExpirationMonth: 12,
			ExpirationYear:  -1,
		}

		err := v.Validate(context.Background(), card)
		assert.Equal(t, domain.ErrYearNumber, err)
	})

	t.Run("ExpiredCard", func(t *testing.T) {
		card := &domain.Card{
			Number:          "4111111111111111",
			ExpirationMonth: 12,
			ExpirationYear:  time.Now().Year() - 1,
		}

		err := v.Validate(context.Background(), card)
		assert.Equal(t, domain.ErrExpiredCard, err)
	})

	t.Run("InvalidLuhnCheck", func(t *testing.T) {
		card := &domain.Card{
			Number:          "4111111111111112",
			ExpirationMonth: 12,
			ExpirationYear:  time.Now().Year() + 1,
		}

		err := v.Validate(context.Background(), card)
		assert.Equal(t, domain.ErrLuhnAlgorithm, err)
	})

	t.Run("ValidCard", func(t *testing.T) {
		card := &domain.Card{
			Number:          "4111111111111111",
			ExpirationMonth: 12,
			ExpirationYear:  time.Now().Year() + 1,
		}

		err := v.Validate(context.Background(), card)
		assert.Nil(t, err)
	})
}
