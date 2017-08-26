package firebase

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"google.golang.org/api/identitytoolkit/v3"
)

// Providers
const (
	Google   string = "google.com"
	Facebook string = "facebook.com"
	Github   string = "github.com"
	Twitter  string = "twitter.com"
)

type (
	// UserRecord is the firebase authentication user
	// See https://firebase.google.com/docs/reference/functions/functions.auth.UserRecord
	UserRecord struct {
		// The user's uid, unique to the Firebase project
		UserID string
		// The user's primary email, if set.
		Email string
		// Whether or not the user's primary email is verified.
		EmailVerified bool
		// The user's display name.
		DisplayName string
		// The user's primary phone number.
		PhoneNumber string
		// The user's photo URL.
		PhotoURL string
		// Whether or not the user is disabled: true for disabled; false for enabled.
		Disabled bool
		//Additional metadata about the user.
		Metadata UserMetadata
		// An array of providers (for example, Google, Facebook) linked to the user.
		ProviderData []*UserInfo
	}

	// UserMetadata is the metadata for user
	// See https://firebase.google.com/docs/reference/functions/functions.auth.UserMetadata
	UserMetadata struct {
		// The date the user was created.
		CreatedAt time.Time
		// The date the user last signed in.
		LastSignedInAt time.Time
	}

	// UserInfo is the user provider information
	// See https://firebase.google.com/docs/reference/functions/functions.auth.UserInfo
	UserInfo struct {
		// The user identifier for the linked provider.
		UserID string
		// The email for the linked provider.
		Email string
		// The display name for the linked provider.
		DisplayName string
		// The phone number for the linked provider.
		PhoneNumber string
		// The photo URL for the linked provider.
		PhotoURL string
		// The linked provider ID (for example, "google.com" for the Google provider).
		ProviderID string
	}

	// Message is the FCM message
	// See https://firebase.google.com/docs/cloud-messaging/http-server-ref#notification-payload-support
	Message struct {

		// To this parameter specifies the recipient of a message.
		//
		// The value can be a device's registration token, a device group's notification key, or a single topic
		// (prefixed with /topics/). To send to multiple topics, use the condition parameter.
		To string `json:"to,omitempty"`

		// This parameter specifies the recipient of a multicast message, a message sent to more than
		// one registration token.
		//
		// The value should be an array of registration tokens to which to send the multicast message.
		// The array must contain at least 1 and at most 1000 registration tokens. To send a message to a
		// single device, use the to parameter.
		//
		// Multicast messages are only allowed using the HTTP JSON format.
		RegistrationIDs []string `json:"registration_ids,omitempty"`

		// This parameter specifies a logical expression of conditions that determine the message target.
		//
		// Supported condition: Topic, formatted as "'yourTopic' in topics". This value is case-insensitive.
		//
		// Supported operators: &&, ||. Maximum two operators per topic message supported.
		Condition string `json:"condition,omitempty"`

		// This parameter identifies a group of messages (e.g., with collapse_key: "Updates
		// Available") that can be collapsed, so that only the last message gets sent when delivery can
		// be resumed. This is intended to avoid sending too many of the same messages when the
		// device comes back online or becomes active.
		//
		// Note that there is no guarantee of the order in which messages get sent.
		//
		// Note: A maximum of 4 different collapse keys is allowed at any given time. This means a FCM
		// connection server can simultaneously store 4 different send-to-sync messages per client app. If
		// you exceed this number, there is no guarantee which 4 collapse keys the FCM connection server
		// will keep.
		CollapseKey string `json:"collapse_key,omitempty"`

		// Sets the priority of the message. Valid values are "normal" and "high." On iOS, these correspond
		// to APNs priorities 5 and 10.
		//
		// By default, notification messages are sent with high priority, and data messages are sent with
		// normal priority. Normal priority optimizes the client app's battery consumption and should be
		// used unless immediate delivery is required. For messages with normal priority, the app may
		// receive the message with unspecified delay.
		//
		// When a message is sent with high priority, it is sent immediately, and the app can wake a
		// sleeping device and open a network connection to your server.
		//
		// For more information, see Setting the priority of a message.
		// See https://firebase.google.com/docs/cloud-messaging/concept-options#setting-the-priority-of-a-message
		Priority string `json:"priority,omitempty"`

		// On iOS, use this field to represent content-available in the APNs payload. When a
		// notification or message is sent and this is set to true, an inactive client app is awoken. On
		// Android, data messages wake the app by default. On Chrome, currently not supported.
		ContentAvailable bool `json:"content_available,omitempty"`

		// Currently for iOS 10+ devices only. On iOS, use this field to represent mutable-content in the
		// APNS payload. When a notification is sent and this is set to true, the content of the notification
		// can be modified before it is displayed, using a Notification Service app extension. This
		// parameter will be ignored for Android and web.
		// See https://developer.apple.com/documentation/usernotifications/unnotificationserviceextension
		MutableContent bool `json:"mutable_content,omitempty"`

		// This parameter specifies how long (in seconds) the message should be kept in FCM storage if
		// the device is offline. The maximum time to live supported is 4 weeks, and the default value is 4
		// weeks. For more information, see Setting the lifespan of a message.
		// See https://firebase.google.com/docs/cloud-messaging/concept-options#ttl
		TimeToLive int `json:"time_to_live,omitempty"`

		// This parameter specifies the package name of the application where the registration tokens
		// must match in order to receive the message.
		RestrictedPackageName string `json:"restricted_package_name,omitempty"`

		// This parameter, when set to true, allows developers to test a request without actually sending
		// a message.
		//
		// The default value is false.
		DryRun bool `json:"dry_run,omitempty"`

		// Data parameter specifies the custom key-value pairs of the message's payload.
		//
		// For example, with data:{"score":"3x1"}:
		//
		// On iOS, if the message is sent via APNS, it represents the custom data fields.
		// If it is sent via FCM connection server, it would be represented as key value dictionary
		// in AppDelegate application:didReceiveRemoteNotification:.
		//
		// On Android, this would result in an intent extra named score with the string value 3x1.
		//
		// The key should not be a reserved word ("from" or any word starting with "google" or "gcm").
		// Do not use any of the words defined in this table (such as collapse_key).
		//
		// Values in string types are recommended. You have to convert values in objects
		// or other non-string data types (e.g., integers or booleans) to string.
		//
		Data interface{} `json:"data,omitempty"`

		// This parameter specifies the predefined, user-visible key-value pairs of the notification payload.
		// See Notification payload support for detail. For more information about notification message
		// and data message options, see Message types.
		Notification Notification `json:"notification,omitempty"`
	}

	// Notification notification message payload
	Notification struct {

		// The notification's title.
		//
		// This field is not visible on iOS phones and tablets.
		Title string `json:"title,omitempty"`

		// The notification's body text.
		Body string `json:"body,omitempty"`

		// The notification's channel id (new in Android O).
		// See https://developer.android.com/preview/features/notification-channels.html
		//
		// The app must create a channel with this ID before any notification with this key is received.
		//
		// If you don't send this key in the request, or if the channel id provided has not yet been
		// created by your app, FCM uses the channel id specified in your app manifest.
		AndroidChannelID string `json:"android_channel_id,omitempty"`

		// The notification's icon.
		//
		// Sets the notification icon to myicon for drawable resource myicon. If you don't send this
		// key in the request, FCM displays the launcher icon specified in your app manifest.
		Icon string `json:"icon,omitempty"`

		// The sound to play when the device receives the notification.
		//
		// Sound files can be in the main bundle of the client app or in the Library/Sounds folder of the
		// app's data container. See the iOS Developer Library for more information.
		// See https://developer.apple.com/library/content/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/SupportingNotificationsinYourApp.html#//apple_ref/doc/uid/TP40008194-CH4-SW10
		Sound string `json:"sound,omitempty"`

		// The value of the badge on the home screen app icon.
		//
		// If not specified, the badge is not changed.
		//
		// If set to 0, the badge is removed.
		Badge string `json:"badge,omitempty"`

		// Identifier used to replace existing notifications in the notification drawer.
		//
		// If not specified, each request creates a new notification.
		//
		// If specified and a notification with the same tag is already being shown, the new
		// notification replaces the existing one in the notification drawer.
		Tag string `json:"tag,omitempty"`

		// The notification's icon color, expressed in #rrggbb format.
		Color string `json:"color,omitempty"`

		// The action associated with a user click on the notification.
		//
		// If specified, an activity with a matching intent filter is launched when a user clicks on the
		// notification.
		ClickAction string `json:"click_action,omitempty"`

		// The key to the body string in the app's string resources to use to localize the body text to
		// the user's current localization.
		//
		// See String Resources for more information.
		// https://developer.android.com/guide/topics/resources/string-resource.html
		BodyLocKey string `json:"body_loc_key,omitempty"`

		// Variable string values to be used in place of the format specifiers in body_loc_key to use
		// to localize the body text to the user's current localization.
		//
		// See Formatting and Styling for more information.
		// https://developer.android.com/guide/topics/resources/string-resource.html#FormattingAndStyling
		BodyLocArgs string `json:"body_loc_args,omitempty"`

		// The key to the title string in the app's string resources to use to localize the title text to
		// the user's current localization.
		//
		// See String Resources for more information.
		// https://developer.android.com/guide/topics/resources/string-resource.html
		TitleLocKey string `json:"title_loc_key,omitempty"`

		// Variable string values to be used in place of the format specifiers in title_loc_key to use
		// to localize the title text to the user's current localization.
		//
		// See Formatting and Styling for more information.
		// https://developer.android.com/guide/topics/resources/string-resource.html#FormattingAndStyling
		TitleLocArgs string `json:"title_loc_args,omitempty"`
	}

	// Response is the FCM server's response
	// See https://firebase.google.com/docs/cloud-messaging/http-server-ref#interpret-downstream
	Response struct {
		// Unique ID (number) identifying the multicast message.
		MulticastID int64 `json:"multicast_id"`

		// Number of messages that were processed without an error.
		Success int `json:"success"`

		// Number of messages that could not be processed.
		Failure int `json:"failure"`

		// Number of results that contain a canonical registration token. A canonical registration ID is the
		// registration token of the last registration requested by the client app. This is the ID that the server
		// should use when sending messages to the device.
		CanonicalIDs int `json:"canonical_ids"`

		// Array of objects representing the status of the messages processed. The objects are listed in the same
		// order as the request (i.e., for each registration ID in the request, its result is listed in the same index in
		// the response).
		//
		// message_id: String specifying a unique ID for each successfully processed message.
		//
		// registration_id: Optional string specifying the canonical registration token for the client app
		// that the message was processed and sent to. Sender should use this value as the registration token
		// for future requests. Otherwise, the messages might be rejected.
		//
		// error: String specifying the error that occurred when processing the message for the recipient. The
		// possible values can be found in table 9.
		// See https://firebase.google.com/docs/cloud-messaging/http-server-ref#table9
		Results []Result `json:"results"`
	}

	// Result representing the status of the messages processed.
	// See https://firebase.google.com/docs/cloud-messaging/http-server-ref#interpret-downstream
	Result struct {
		// The topic message ID when FCM has successfully received the request and will attempt to deliver to
		// all subscribed devices.
		MessageID string `json:"message_id"`

		// This parameter specifies the canonical registration token for the client app that the message was
		// processed and sent to. Sender should replace the registration token with this value on future
		// requests; otherwise, the messages might be rejected.
		RegistrationID string `json:"registration_id"`

		// Error that occurred when processing the message. The possible values can be found in table 9.
		Error error `json:"error"`
	}

	// Topic is the FCM topic
	// See https://developers.google.com/instance-id/reference/server#manage_relationship_maps_for_multiple_app_instances
	Topic struct {

		// To this parameter specifies the The topic name.
		//
		To string `json:"to,omitempty"`

		// This parameter specifies The array of IID tokens for the app instances you want to add or remove
		RegistrationTokens []string `json:"registration_tokens,omitempty"`
	}
)

// toFirebaseErr map an error string returned by firebase to error
// TODO find a best way to use directly error
func toFirebaseErr(firebaseErr string) error {

	var mapping = map[string]error{
		"MissingRegistration":       ErrMissingRegistration,
		"InvalidRegistration":       ErrInvalidRegistration,
		"NotRegistered":             ErrNotRegistered,
		"InvalidPackageName":        ErrInvalidPackageName,
		"MismatchSenderId":          ErrMismatchSenderID,
		"InvalidParameters":         ErrInvalidParameters,
		"MessageTooBig":             ErrMessageTooBig,
		"InvalidDataKey":            ErrInvalidDataKey,
		"InvalidTtl":                ErrInvalidTTL,
		"Unavailable":               ErrUnavailable,
		"InternalServerError":       ErrInternalServerError,
		"DeviceMessageRateExceeded": ErrDeviceMessageRateExceeded,
		"TopicsMessageRateExceeded": ErrTopicsMessageRateExceeded,
		"InvalidApnsCredential":     ErrInvalidApnsCredential,
	}

	return mapping[firebaseErr]
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (r *Result) UnmarshalJSON(data []byte) error {
	var result struct {
		MessageID      string `json:"message_id"`
		RegistrationID string `json:"registration_id"`
		Error          string `json:"error"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	r.MessageID = result.MessageID
	r.RegistrationID = result.RegistrationID
	r.Error = toFirebaseErr(result.Error)

	return nil
}

// Validate returns an error if the payload message is not valid.
func (payload *Message) Validate() error {

	// validate number of Condition operators
	opCnt := strings.Count(payload.Condition, "&&") + strings.Count(payload.Condition, "||")
	if opCnt > 2 {
		return &ErrInvalidMessage{"firebaseFCM: Too many operators for conditions only support up to two operators per expression"}
	}

	// validate recipient is not empty
	if payload.To == "" && payload.Condition == "" && len(payload.RegistrationIDs) == 0 {
		return &ErrInvalidMessage{"firebaseFCM: A recipient is missing"}
	}

	// validate max RegistrationIDs
	if len(payload.RegistrationIDs) > 1000 {
		return &ErrInvalidMessage{"firebaseFCM: Too many registrations ids max size is 1000"}
	}

	// validate TTL
	if payload.TimeToLive > 2419200 {
		return &ErrInvalidMessage{"firebaseFCM: TTL is greater than 2419200"}
	}

	// TODO validate size
	// 4096 bytes for messages or 2048 bytes for topics

	return nil
}

// Validate returns an error if the payload topic is not valid.
func (payload *Topic) Validate() error {

	// validate recipient is not empty
	if payload.To == "" && strings.HasPrefix(payload.To, "/topics/") {
		return &ErrInvalidMessage{"firebaseFCM: A recipient is missing"}
	}

	// validate name
	var validTopic = regexp.MustCompile(`[a-zA-Z0-9-_.~%]+`)
	if !validTopic.MatchString(payload.To) {
		return &ErrInvalidMessage{"firebaseFCM: Topic name is invalid must be : [a-zA-Z0-9-_.~%]+"}
	}

	// validate max RegistrationIDs
	if len(payload.RegistrationTokens) < 1 {
		return &ErrInvalidMessage{"firebaseFCM: registration_token is missing"}
	}

	return nil
}

func parseDate(t int64) time.Time {
	if t == 0 {
		return time.Time{}
	}
	return time.Unix(t/1000, 0)
}

func toUserRecord(user *identitytoolkit.UserInfo) *UserRecord {
	return &UserRecord{
		UserID:        user.LocalId,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		DisplayName:   user.DisplayName,
		PhotoURL:      user.PhotoUrl,
		PhoneNumber:   user.PhoneNumber,
		Disabled:      user.Disabled,
		Metadata: UserMetadata{
			CreatedAt:      parseDate(user.CreatedAt),
			LastSignedInAt: parseDate(user.LastLoginAt),
		},
		ProviderData: toUserInfos(user.ProviderUserInfo),
	}
}

func toUserRecords(users []*identitytoolkit.UserInfo) []*UserRecord {
	result := make([]*UserRecord, len(users))
	for i, user := range users {
		result[i] = toUserRecord(user)
	}
	return result
}

func toUserInfo(info *identitytoolkit.UserInfoProviderUserInfo) *UserInfo {
	return &UserInfo{
		UserID:      info.RawId,
		Email:       info.Email,
		DisplayName: info.DisplayName,
		PhotoURL:    info.PhotoUrl,
		PhoneNumber: info.PhoneNumber,
		ProviderID:  info.ProviderId,
	}
}

func toUserInfos(infos []*identitytoolkit.UserInfoProviderUserInfo) []*UserInfo {
	result := make([]*UserInfo, len(infos))
	for i, info := range infos {
		result[i] = toUserInfo(info)
	}
	return result
}

// User use for create new user
// use Password when create user (plain text)
// use RawPassword when update user (hashed password)
type User struct {
	UserID        string
	Email         string
	EmailVerified bool
	Password      string
	DisplayName   string
	PhotoURL      string
	PhoneNumber   string
	Disabled      bool
}

// UpdateAccount use for update existing account
type UpdateAccount struct {
	// UserID is the existing user id to update
	UserID        string `json:"localId,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"emailVerified,omitempty"`
	Password      string `json:"password,omitempty"`
	DisplayName   string `json:"displayName,omitempty"`
	PhotoURL      string `json:"photoUrl,omitempty"`
	Disabled      bool   `json:"disableUser,omitempty"`
}

var scopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/firebase.database",
	"https://www.googleapis.com/auth/identitytoolkit",
}
