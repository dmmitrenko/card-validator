package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrBadParamInput       = errors.New("given data is invalid")
	ErrLuhnAlgorithm       = errors.New("your card doesn't match Luhn's algorithm")
	ErrExpiredCard         = errors.New("your card is expired")
)
