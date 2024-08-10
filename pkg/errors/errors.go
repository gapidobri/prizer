package errors

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var (
	internalCode atomic.Uint32
)

type ApiError struct {
	// HTTP status code
	statusCode int

	// Error code
	code string

	// Message displayed in the http response
	message string

	// Code used for comparison
	internalCode uint32

	// Optional wrapped error
	err error
}

func (e ApiError) StatusCode() int {
	return e.statusCode
}

func (e ApiError) Code() string {
	return e.code
}

func (e ApiError) Message() string {
	return e.message
}

// Error returns the error message
//
// Do NOT use this for returning error messages to the user, because it returns the wrapped errors as well!
// Use the Message field instead.
func (e ApiError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.message, e.err)
	}
	return e.message
}

func (e ApiError) Unwrap() error {
	return e.err
}

func (e ApiError) Is(target error) bool {
	// Check if error is of type apiError
	if t, ok := target.(ApiError); ok {
		return e.internalCode == t.internalCode
	}
	// Else check underlying error
	return errors.Is(e.err, target)
}

func New(statusCode int, code string, message string) ApiError {
	return ApiError{
		internalCode: internalCode.Add(1),
		statusCode:   statusCode,
		message:      message,
	}
}

func (e ApiError) New(code string, message string) ApiError {
	err := New(e.statusCode, code, message)
	err.err = e
	return err
}

func (e ApiError) With(message string) ApiError {
	e.message = message
	return e
}
