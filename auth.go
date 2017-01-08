package admin

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Auth type
type Auth struct {
	app       *App
	keysMutex *sync.RWMutex
	keys      map[string]*rsa.PublicKey
	keysExp   time.Time
}

const keysEndpoint = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"

func newAuth(app *App) *Auth {
	return &Auth{
		app:       app,
		keysMutex: &sync.RWMutex{},
	}
}

// CreateCustomToken creates a custom token used for client to authenticate
// with firebase server using signInWithCustomToken
// https://firebase.google.com/docs/auth/admin/create-custom-tokens
func (auth *Auth) CreateCustomToken(userID string, claims interface{}) (string, error) {
	if auth.app.jwtConfig == nil || auth.app.privateKey == nil {
		return "", ErrRequireServiceAccount
	}
	now := time.Now()
	payload := &Claims{
		Issuer:    auth.app.jwtConfig.Email,
		Subject:   auth.app.jwtConfig.Email,
		Audience:  "https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit",
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Hour).Unix(),
		UserID:    userID,
		Claims:    claims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	return token.SignedString(auth.app.privateKey)
}

// VerifyIDToken validates given idToken
// return Claims for that token only valid token
func (auth *Auth) VerifyIDToken(idToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(idToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Auth: Firebase ID token has incorrect algorithm. Expected \"RSA\" but got \"%#v\"", token.Header["alg"])
		}
		kid := token.Header["kid"].(string)
		if kid == "" {
			return nil, fmt.Errorf("Auth: Firebase ID token has no \"kid\" claim")
		}
		key := auth.selectKey(kid)
		if key == nil {
			return nil, fmt.Errorf("Auth: Firebase ID token has \"kid\" claim which does not correspond to a known public key. Most likely the ID token is expired, so get a fresh token from your client app and try again")
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if !claims.verifyAudience(auth.app.projectID) {
			return nil, fmt.Errorf("Auth: Firebase ID token has incorrect \"aud\" (audience) claim. Expected \"%s\" but got \"%s\"", auth.app.projectID, claims.Audience)
		}
		if !claims.verifyIssuer("https://securetoken.google.com/" + auth.app.projectID) {
			return nil, fmt.Errorf("Auth: Firebase ID token has incorrect \"iss\" (issuer) claim. Expected \"https://securetoken.google.com/%s\" but got \"%s\"", auth.app.projectID, claims.Issuer)
		}
		if claims.Subject == "" {
			return nil, fmt.Errorf("Auth: Firebase ID token has an empty string \"sub\" (subject) claim")
		}
		if len(claims.Subject) > 128 {
			return nil, fmt.Errorf("Auth: Firebase ID token has \"sub\" (subject) claim longer than 128 characters")
		}

		claims.UserID = claims.Subject

		return claims, nil
	}
	return nil, fmt.Errorf("Auth: invalid token")
}

func (auth *Auth) fetchKeys() error {
	auth.keysMutex.Lock()
	defer auth.keysMutex.Unlock()
	resp, err := http.Get(keysEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	auth.keysExp, _ = time.Parse(time.RFC1123, resp.Header.Get("Expires"))

	m := map[string]string{}
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return err
	}
	ks := map[string]*rsa.PublicKey{}
	for k, v := range m {
		p, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(v))
		if p != nil {
			ks[k] = p
		}
	}
	auth.keys = ks
	return nil
}

func (auth *Auth) selectKey(kid string) *rsa.PublicKey {
	auth.keysMutex.RLock()
	if auth.keysExp.IsZero() || auth.keysExp.Before(time.Now()) || len(auth.keys) == 0 {
		auth.keysMutex.RUnlock()
		if err := auth.fetchKeys(); err != nil {
			return nil
		}
		auth.keysMutex.RLock()
	}
	defer auth.keysMutex.RUnlock()
	return auth.keys[kid]
}

// GetAccountInfoByUID retrieves an account info by user id
func (auth *Auth) GetAccountInfoByUID(uid string) (*User, error) {
	users, err := auth.GetAccountInfoByUIDs([]string{uid})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

// GetAccountInfoByUIDs retrieves account infos by user ids
func (auth *Auth) GetAccountInfoByUIDs(uids []string) ([]*User, error) {
	var resp getAccountInfoResponse
	err := auth.app.invokeRequest(http.MethodPost, getAccountInfo, &getAccountInfoRequest{LocalIDs: uids}, &resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Users) == 0 {
		return nil, nil
	}
	return resp.Users, nil
}

// GetAccountInfoByEmail retrieves account info by email
func (auth *Auth) GetAccountInfoByEmail(email string) (*User, error) {
	users, err := auth.GetAccountInfoByEmails([]string{email})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

// GetAccountInfoByEmails retrieves account infos by emails
func (auth *Auth) GetAccountInfoByEmails(emails []string) ([]*User, error) {
	var resp getAccountInfoResponse
	err := auth.app.invokeRequest(http.MethodPost, getAccountInfo, &getAccountInfoRequest{Emails: emails}, &resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Users) == 0 {
		return nil, nil
	}
	return resp.Users, nil
}

// DeleteAccount deletes an account by user id
func (auth *Auth) DeleteAccount(uid string) error {
	if uid == "" {
		return ErrRequireUID
	}

	var resp deleteAccountResponse
	return auth.app.invokeRequest(http.MethodPost, deleteAccount, &deleteAccountRequest{LocalID: uid}, &resp)
}

// CreateAccount creates an account
// if not provides LocalID, firebase server will auto generate
// created user id can get from first param of result
func (auth *Auth) CreateAccount(user *Account) (string, error) {
	var err error
	if user.LocalID == "" {
		var resp signupNewUserResponse
		if user.RawPassword != "" && user.Password == "" {
			user.Password = user.RawPassword
			user.RawPassword = ""
		}
		err = auth.app.invokeRequest(http.MethodPost, signupNewUser, &signupNewUserRequest{user}, &resp)
		if err != nil {
			return "", err
		}
		if resp.LocalID == "" {
			return "", errors.New("Auth: create account error")
		}
		return resp.LocalID, nil
	}
	var resp uploadAccountResponse
	if user.RawPassword == "" && user.Password != "" {
		user.RawPassword = user.Password
		user.Password = ""
	}
	err = auth.app.invokeRequest(http.MethodPost, uploadAccount, &uploadAccountRequest{
		Users:          []*Account{user},
		AllowOverwrite: false,
		SanityCheck:    true,
	}, &resp)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", errors.New("Auth: upload account error")
	}
	return user.LocalID, nil
}

// ListAccountCursor type
type ListAccountCursor struct {
	nextPageToken string
	auth          *Auth
	MaxResults    int
}

// ListAccount creates list account cursor for retrieves accounts
// MaxResults can change later after create cursor
func (auth *Auth) ListAccount(maxResults int) *ListAccountCursor {
	return &ListAccountCursor{MaxResults: maxResults, auth: auth}
}

// Next retrieves next users from cursor which limit to MaxResults
// then move cursor to the next users
func (cursor *ListAccountCursor) Next() ([]*User, error) {
	var resp downloadAccountResponse
	err := cursor.auth.app.invokeRequest(http.MethodPost, downloadAccount, &downloadAccountRequest{MaxResults: cursor.MaxResults, NextPageToken: cursor.nextPageToken}, &resp)
	if err != nil {
		return nil, err
	}
	cursor.nextPageToken = resp.NextPageToken
	return resp.Users, nil
}

// UpdateAccount updates an existing account
func (auth *Auth) UpdateAccount(account *UpdateAccount) (string, error) {
	var resp setAccountInfoResponse
	err := auth.app.invokeRequest(http.MethodPost, setAccountInfo, &setAccountInfoRequest{account}, &resp)
	if err != nil {
		return "", err
	}
	return resp.LocalID, nil
}
