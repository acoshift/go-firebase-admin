package admin

import (
	"errors"
)

// Error constants
var (
	ErrRequireServiceAccount = errors.New("firebase: requires service account")
	ErrRequireUID            = errors.New("firebaseauth: require user id")
	ErrNotImplement          = errors.New("firebase: feature not implement yet")
)

// ErrTokenInvalid is the invalid token error
type ErrTokenInvalid struct {
	s string
}

// Error implements error interface
func (err *ErrTokenInvalid) Error() string {
	return err.s
}
