package admin

import (
	"context"
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

// UserRecord is the firebase authentication user
type UserRecord struct {
	UserID        string
	Email         string
	EmailVerified bool
	DisplayName   string
	PhotoURL      string
	Disabled      bool
	Metadata      UserMetadata
	ProviderData  []*UserInfo
}

// UserMetadata is the metadata for user
type UserMetadata struct {
	CreatedAt      time.Time
	LastSignedInAt time.Time
}

// UserInfo is the user provider information
type UserInfo struct {
	UserID      string
	Email       string
	DisplayName string
	PhotoURL    string
	ProviderID  string
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

const (
	baseURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/"
	timeout = time.Second * 10000
)

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
