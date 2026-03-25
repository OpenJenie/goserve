package auth

import (
	"errors"
	"fmt"
)

// authError is used to pass an error during the request through the
// application with auth specific context.
type authError struct {
	msg    string
	status int
}

// NewAuthError creates an AuthError for the provided message.
func NewAuthError(format string, args ...any) error {
	return &authError{
		msg:    fmt.Sprintf(format, args...),
		status: 0,
	}
}

// NewAuthErrorWithStatus creates an AuthError with an explicit HTTP status.
func NewAuthErrorWithStatus(status int, format string, args ...any) error {
	return &authError{
		msg:    fmt.Sprintf(format, args...),
		status: status,
	}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (ae *authError) Error() string {
	return ae.msg
}

// IsAuthError checks if an error of type AuthError exists.
func IsAuthError(err error) bool {
	var ae *authError
	return errors.As(err, &ae)
}

// GetAuthError returns the auth error if present.
func GetAuthError(err error) *authError {
	var ae *authError
	if !errors.As(err, &ae) {
		return nil
	}
	return ae
}

// Status returns the configured HTTP status for the auth error.
func (ae *authError) Status() int {
	return ae.status
}
