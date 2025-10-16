// Package apperrors defines a shared set of base error types used throughout the application.
package apperrors

import (
	"errors"

	"github.com/hay-kot/httpkit/errtrace"
)

var ErrBlockedTokenRequest = errtrace.New("token requests are blocked")

// As is a generic implementation of the errors.As function.
func As[T error](err error) bool {
	var re T
	return errors.As(err, &re)
}

type RateLimitError struct {
	msg string
}

func (e RateLimitError) Error() string {
	return e.msg
}

func NewRateLimitError(msg string) error {
	return RateLimitError{msg: msg}
}

var ErrNoBody = errors.New("no body in request")

func IsNoBody(err error) bool {
	return errors.Is(err, ErrNoBody)
}
