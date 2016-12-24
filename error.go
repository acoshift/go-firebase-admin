package admin

import (
	"errors"
)

// Error constants
var (
	ErrRequireServiceAccount = errors.New("firebase: requires service account")
	ErrRequireUID            = errors.New("firebaseauth: require user id")
)
