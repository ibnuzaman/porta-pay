package errors

import (
	"errors"
	"fmt"
)

var (
	// Domain errors
	ErrBookingNotFound  = errors.New("booking not found")
	ErrInvalidBookingID = errors.New("invalid booking id")
	ErrInvalidQuantity  = errors.New("quantity must be greater than 0")
	ErrBookingExpired   = errors.New("booking has expired")
	ErrBookingConfirmed = errors.New("cannot modify confirmed booking")

	// Infrastructure errors
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrQueryFailed        = errors.New("database query failed")
)

// BookingError represents a booking-specific error
type BookingError struct {
	Code    string
	Message string
	Err     error
}

func (e *BookingError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *BookingError) Unwrap() error {
	return e.Err
}

// NewBookingError creates a new booking error
func NewBookingError(code, message string, err error) *BookingError {
	return &BookingError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
