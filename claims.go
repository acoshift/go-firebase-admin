package admin

import (
	"fmt"
	"time"
)

// Claims is the firebase authentication token claims
type Claims struct {
	Issuer    string      `json:"iss,omitempty"`
	Subject   string      `json:"sub,omitempty"`
	Audience  string      `json:"aud,omitempty"`
	IssuedAt  int64       `json:"iat,omitempty"`
	ExpiresAt int64       `json:"exp,omitempty"`
	UserID    string      `json:"uid,omitempty"`
	Claims    interface{} `json:"claims,omitempty"`
}

// Valid implements jwt-go Claims interface
// for validates time based claims, such as IssuedAt, and ExpiresAt
// But not verify token signature and header
func (c *Claims) Valid() error {
	now := time.Now().Unix()

	if !c.verifyExpiresAt(now) {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		return fmt.Errorf("token is expired by %v", delta)
	}

	if !c.verifyIssuedAt(now) {
		return fmt.Errorf("token used before issued")
	}

	return nil
}

func (c *Claims) verifyExpiresAt(now int64) bool {
	return now <= c.ExpiresAt
}

func (c *Claims) verifyIssuedAt(now int64) bool {
	return now >= c.IssuedAt
}

func (c *Claims) verifyAudience(aud string) bool {
	return c.Audience == aud
}

func (c *Claims) verifyIssuer(iss string) bool {
	return c.Issuer == iss
}
