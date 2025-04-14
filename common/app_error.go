package common

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError is a structured error type that holds additional context.
type AppError struct {
	Code    int    `json:"code"`              // e.g., HTTP status code or custom code.
	Message string `json:"message"`           // User-friendly message.
	Err     error  `json:"-"`                 // The underlying error, if any.
}

// Error returns the string representation of the error.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) RootError() error {
	return e.Err
}

// WrapError wraps an existing, real internal error with additional context.
func WrapError(err error, message string, code int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     fmt.Errorf("%s: %w", message, err),
	}
}

// NewError handles business logic/"soft" errors
func NewError(message string, code int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func ErrDB(err error) *AppError {
	return WrapError(err, "database error", http.StatusInternalServerError)
}

func ErrInvalidRequest(err error) *AppError {
	return WrapError(err, "invalid request", http.StatusBadRequest)
}

func ErrUnauthorized(err error) *AppError {
	return WrapError(err, "unauthorized", http.StatusUnauthorized)
}

func ErrNotFound(resource string, id interface{}) *AppError {
	return NewError(fmt.Sprintf("%s with id %v not found", resource, id), http.StatusNotFound)
}

func ErrCannotCreate(resource string, err error) *AppError {
	return WrapError(err, fmt.Sprintf("cannot create %s", resource), http.StatusInternalServerError)
}

func ErrCannotUpdate(resource string, err error) *AppError {
	return WrapError(err, fmt.Sprintf("cannot update %s", resource), http.StatusInternalServerError)
}

func ErrCannotDelete(resource string, err error) *AppError {
	return WrapError(err, fmt.Sprintf("cannot delete %s", resource), http.StatusInternalServerError)
}

// A sentinel error that can be used to indicate no record was found.
var RecordNotFound = errors.New("record not found")
