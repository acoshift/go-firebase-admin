package firebase

import (
	"errors"
)

// Error constants
var (
	// ErrRequireServiceAccount
	ErrRequireServiceAccount = errors.New("firebase: requires service account")

	// ErrRequireUID
	ErrRequireUID = errors.New("firebaseauth: require user id")

	// ErrNotImplemented
	ErrNotImplemented = errors.New("firebase: feature not yet implemented")

	// ErrUserNotFound
	ErrUserNotFound = errors.New("firebaseauth: user not found")

	// ErrAuthentication
	ErrAuthentication = errors.New("firebase: Authentication Error")
)

// Official FCM Error constants
// see https://firebase.google.com/docs/cloud-messaging/admin/errors
// see https://firebase.google.com/docs/cloud-messaging/http-server-ref#error-codes
var (
	// ErrMissingRegistration
	//
	// Check that the request contains a registration token (in the registration_id in or
	// registration_ids field in JSON).
	ErrMissingRegistration = errors.New("firebaseFCM: missing registration token")

	// ErrInvalidRegistration
	//
	// Check the format of the registration token you pass to the server. Make sure it matches
	// the registration token the client app receives from registering with Firebase Notifications.
	// Do not truncate or add additional characters.
	ErrInvalidRegistration = errors.New("firebaseFCM: invalid registration token")

	// ErrNotRegistered
	//
	// An existing registration token may cease to be valid in a number of scenarios, including:
	//
	// If the client app unregisters with FCM.
	//
	// If the client app is automatically unregistered, which can happen if the user uninstalls
	// the application. For example, on iOS, if the APNS Feedback Service reported the APNS
	// token as invalid.
	//
	// If the registration token expires (for example, Google might decide to refresh
	// registration tokens, or the APNS token has expired for iOS devices).
	//
	// If the client app is updated but the new version is not configured to receive messages.
	// For all these cases, remove this registration token from the app server and stop using it
	// to send messages.
	ErrNotRegistered = errors.New("firebaseFCM: unregistered device")

	// ErrInvalidPackageName
	//
	// Make sure the message was addressed to a registration token whose package name
	// matches the value passed in the request.
	ErrInvalidPackageName = errors.New("firebaseFCM: invalid package name")

	// ErrMismatchSenderID
	//
	// A registration token is tied to a certain group of senders. When a client app registers for
	// FCM, it must specify which senders are allowed to send messages. You should use one
	// of those sender IDs when sending messages to the client app. If you switch to a different
	// sender, the existing registration tokens won't work.
	ErrMismatchSenderID = errors.New("firebaseFCM: mismatched sender id")

	// ErrInvalidParameters
	//
	// Check that the provided parameters have the right name and type.
	ErrInvalidParameters = errors.New("firebaseFCM: invalid parameters")

	// ErrMessageTooBig
	//
	// Check that the total size of the payload data included in a message does not exceed
	// FCM limits: 4096 bytes for most messages, or 2048 bytes in the case of messages to
	// topics. This includes both the keys and the values.
	ErrMessageTooBig = errors.New("firebaseFCM: message is too big")

	// ErrInvalidDataKey
	//
	// Check that the payload data does not contain a key (such as from, or gcm, or any value
	// prefixed by google) that is used internally by FCM. Note that some words (such as
	// collapse_key) are also used by FCM but are allowed in the payload, in which case the
	// payload value will be overridden by the FCM value.
	ErrInvalidDataKey = errors.New("firebaseFCM: invalid data key")

	// ErrInvalidTTL
	//
	// Check that the value used in time_to_live is an integer representing a duration in
	// seconds between 0 and 2,419,200 (4 weeks).
	ErrInvalidTTL = errors.New("firebaseFCM: invalid time to live")

	// ErrUnavailable
	//
	// The server couldn't process the request in time. Retry the same request, but you must:
	//
	// Honor the Retry-After header if it is included in the response from the
	// FCM Connection Server.
	//
	// Implement exponential back-off in your retry mechanism. (e.g. if you waited one
	// second before the first retry, wait at least two second before the next one, then 4
	// seconds and so on). If you're sending multiple messages, delay each one
	// independently by an additional random amount to avoid issuing a new request for all
	// messages at the same time.
	//
	// Senders that cause problems risk being blacklisted.
	ErrUnavailable = errors.New("firebaseFCM: timeout")

	// ErrInternalServerError
	//
	// The server encountered an error while trying to process the request. You could retry the
	// same request following the requirements listed in "Timeout" (see row above). If the
	// error persists, please report the problem in the android-gcm group.
	ErrInternalServerError = errors.New("firebaseFCM: internal server error")

	// ErrDeviceMessageRateExceeded
	//
	// The rate of messages to a particular device is too high. If an iOS app sends messages at
	// a rate exceeding APNs limits, it may receive this error message
	//
	// Reduce the number of messages sent to this device and use exponential backoff to retry
	// sending.
	ErrDeviceMessageRateExceeded = errors.New("firebaseFCM: device message rate exceeded")

	// ErrTopicsMessageRateExceeded
	//
	// The rate of messages to subscribers to a particular topic is too high. Reduce the number
	// of messages sent for this topic and use exponential backoff to retry sending.
	ErrTopicsMessageRateExceeded = errors.New("firebaseFCM: topics message rate exceeded")

	// ErrInvalidApnsCredential
	//
	// A message targeted to an iOS device could not be sent because the required APNs SSL
	// certificate was not uploaded or has expired. Check the validity of your development and
	// production certificates.
	ErrInvalidApnsCredential = errors.New("firebaseFCM: invalid APNs credential")
)

// ErrTokenInvalid is the invalid token error
type ErrTokenInvalid struct {
	s string
}

// Error implements error interface
func (err *ErrTokenInvalid) Error() string {
	return err.s
}

// ErrInvalidMessage is the invalid Message error
type ErrInvalidMessage struct {
	s string
}

// Error implements error interface
func (err *ErrInvalidMessage) Error() string {
	return err.s
}
