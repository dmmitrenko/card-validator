package domain

import "fmt"

type CodedError interface {
	error
	ErrorCode() string
}

type ErrorWithCode struct {
	Code    int
	Message string
}

func (e *ErrorWithCode) Error() string {
	return fmt.Sprintf(e.Message)
}

func (e *ErrorWithCode) ErrorCode() string {
	return fmt.Sprintf("%03d", e.Code)
}

var (
	ErrCardNumber    = &ErrorWithCode{Code: 1, Message: "invalid card number"}
	ErrMonthNumber   = &ErrorWithCode{Code: 2, Message: "invalid month"}
	ErrYearNumber    = &ErrorWithCode{Code: 3, Message: "invalid year"}
	ErrLuhnAlgorithm = &ErrorWithCode{Code: 4, Message: "card doesn't match Luhn's algorithm"}
	ErrExpiredCard   = &ErrorWithCode{Code: 5, Message: "card is expired"}
)
