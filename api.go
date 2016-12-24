package admin

import (
	"context"
	"time"
)

// User is the firebase authentication user
type User struct {
	LocalID           string      `json:"localId,omitempty"`
	Email             string      `json:"email,omitempty"`
	EmailVerified     bool        `json:"emailVerified,omitempty"`
	ProviderUserInfo  []*UserInfo `json:"providerUserInfo,omitempty"`
	PasswordHash      string      `json:"passwordHash,omitempty"`
	PasswordUpdatedAt float64     `json:"passwordUpdatedAt,omitempty"`
	ValidSince        string      `json:"validSince,omitempty"`
	Disabled          bool        `json:"disabled,omitempty"`
	LastLoginAt       string      `json:"lastLoginAt,omitempty"`
	CreatedAt         string      `json:"createdAt,omitempty"`
}

// UserInfo is the user provider information
type UserInfo struct {
	ProviderID  string `json:"providerId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	PhotoURL    string `json:"photoUrl,omitempty"`
	FederatedID string `json:"federatedId,omitempty"`
	Email       string `json:"email,omitempty"`
	RawID       string `json:"rawId,omitempty"`
	ScreenName  string `json:"screenName,omitempty"`
}

// Account use for create account
type Account struct {
	LocalID       string `json:"localId,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"emailVerified,omitempty"`
	Password      string `json:"password,omitempty"`
	RawPassword   string `json:"rawPassword,omitempty"`
	DisplayName   string `json:"displayName,omitempty"`
	PhotoURL      string `json:"photoUrl,omitempty"`
	Disabled      bool   `json:"disabled,omitempty"`
}

// UpdateAccount use for update existing account
type UpdateAccount struct {
	// LocalID is the existing user id to update
	LocalID       string `json:"localId,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"emailVerified,omitempty"`
	Password      string `json:"password,omitempty"`
	DisplayName   string `json:"displayName,omitempty"`
	PhotoURL      string `json:"photoUrl,omitempty"`
	Disabled      bool   `json:"disableUser,omitempty"`
}

type getAccountInfoRequest struct {
	LocalIDs []string `json:"localId,omitempty"`
	Emails   []string `json:"email,omitempty"`
}

type getAccountInfoResponse struct {
	Users []*User `json:"users,omitempty"`
}

type deleteAccountRequest struct {
	LocalID string `json:"localId,omitempty"`
}

type deleteAccountResponse struct {
}

type uploadAccountRequest struct {
	Users          []*Account `json:"users"`
	AllowOverwrite bool       `json:"allowOverwrite"`
	SanityCheck    bool       `json:"sanityCheck"`
}

type uploadAccountResponse struct {
	Error []*struct {
		Index   int    `json:"index,omitempty"`
		Message string `json:"message,omitempty"`
	} `json:"error,omitempty"`
}

type signupNewUserRequest struct {
	*Account
}

type signupNewUserResponse struct {
	LocalID string `json:"localId,omitempty"`
}

type downloadAccountRequest struct {
	MaxResults    int    `json:"maxResults,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}

type downloadAccountResponse struct {
	Users         []*User `json:"users,omitempty"`
	NextPageToken string  `json:"nextPageToken,omitempty"`
}

type setAccountInfoRequest struct {
	*UpdateAccount
}

type setAccountInfoResponse struct {
	LocalID string `json:"localId,omitempty"`
}

var scopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/firebase.database",
	"https://www.googleapis.com/auth/identitytoolkit",
}

const (
	baseURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/"
	timeout = time.Second * 10000
)

type apiMethod string

const (
	getAccountInfo   apiMethod = "getAccountInfo"
	setAccountInfo   apiMethod = "setAccountInfo"
	deleteAccount    apiMethod = "deleteAccount"
	uploadAccount    apiMethod = "uploadAccount"
	signupNewUser    apiMethod = "signupNewUser"
	downloadAccount  apiMethod = "downloadAccount"
	getOOBCode       apiMethod = "getOobConfirmationCode"
	getProjectConfig apiMethod = "getProjectConfig"
)

type httpMethod string

const (
	httpGet  httpMethod = "GET"
	httpPost httpMethod = "POST"
)

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
