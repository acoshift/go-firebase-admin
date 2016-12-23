package admin

import (
	"errors"
)

// Errors
var (
	ErrRequireServiceAccount = errors.New("firebase: requires service account")
	ErrRequireUID            = errors.New("firebaseauth: require user id")
)
