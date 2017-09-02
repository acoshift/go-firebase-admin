package firebase

// customClaims is the firebase authentication custom token claims
type customClaims struct {
	Issuer    string      `json:"iss,omitempty"`
	Subject   string      `json:"sub,omitempty"`
	Audience  string      `json:"aud,omitempty"`
	IssuedAt  int64       `json:"iat,omitempty"`
	ExpiresAt int64       `json:"exp,omitempty"`
	UserID    string      `json:"uid,omitempty"`
	Claims    interface{} `json:"claims,omitempty"`
}

func (c *customClaims) Valid() error {
	return nil
}
