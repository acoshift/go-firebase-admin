package firebase

import (
	"fmt"
	"time"
)

// Token is the firebase access token
type Token struct {
	Issuer        string `json:"iss"`
	Name          string `json:"name"`
	ID            string `json:"id"`
	Audience      string `json:"aud"`
	AuthTime      int64  `json:"auth_time"`
	UserID        string `json:"user_id"`
	Subject       string `json:"sub"`
	IssuedAt      int64  `json:"iat"`
	ExpiresAt     int64  `json:"exp"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	PhoneNumber   string `json:"phone_number"`
	Firebase      struct {
		Identities struct {
			Phone []string `json:"phone"`
			Email []string `json:"email"`
		} `json:"identities"`
		SignInProvider string `json:"sign_in_provider"`
	} `json:"firebase"`
}

// Valid implements jwt-go Claims interface
// for validates time based claims, such as IssuedAt, and ExpiresAt
// But not verify token signature and header
func (t *Token) Valid() error {
	now := time.Now().Unix()

	if !t.verifyExpiresAt(now) {
		delta := time.Unix(now, 0).Sub(time.Unix(t.ExpiresAt, 0))
		return fmt.Errorf("token is expired by %v", delta)
	}

	if !t.verifyIssuedAt(now) {
		return fmt.Errorf("token used before issued")
	}

	return nil
}

func (t *Token) verifyExpiresAt(now int64) bool {
	return now <= t.ExpiresAt
}

func (t *Token) verifyIssuedAt(now int64) bool {
	return now >= t.IssuedAt
}

func (t *Token) verifyAudience(aud string) bool {
	return t.Audience == aud
}

func (t *Token) verifyIssuer(iss string) bool {
	return t.Issuer == iss
}
