package admin

import (
	"context"
	"time"
)

// User is the firebase auth user
type User struct {
	LocalID           string              `json:"localId"`
	Email             string              `json:"email"`
	EmailVerified     bool                `json:"emailVerified"`
	ProviderUserInfo  []*ProviderUserInfo `json:"providerUserInfo"`
	PasswordHash      string              `json:"passwordHash"`
	PasswordUpdatedAt float64             `json:"passwordUpdatedAt"`
	ValidSince        string              `json:"validSince"`
	LastLoginAt       string              `json:"lastLoginAt"`
	CreatedAt         string              `json:"createdAt"`
}

// ProviderUserInfo type
type ProviderUserInfo struct {
	ProviderID  string `json:"providerId"`
	DisplayName string `json:"displayName"`
	PhotoURL    string `json:"photoUrl"`
	FederatedID string `json:"federatedId"`
	Email       string `json:"email"`
	RawID       string `json:"rawId"`
	ScreenName  string `json:"screenName"`
}

type apiResponse struct {
	Kind  string    `json:"kind"`
	Error *apiError `json:"error,omitempty"`
	Users []*User   `json:"users,omitempty"`
}

type apiError struct {
	Message string `json:"message"`
}

type getAccountInfoRequest struct {
	LocalID []string `json:"localId,omitempty"`
	Email   []string `json:"email,omitempty"`
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

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
