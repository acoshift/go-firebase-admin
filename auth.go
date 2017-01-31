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
	"google.golang.org/api/iterator"
)

// Auth type
type Auth struct {
	app       *App
	keysMutex *sync.RWMutex
	keys      map[string]*rsa.PublicKey
	keysExp   time.Time
}

const (
	keysEndpoint        = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"
	customTokenAudience = "https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit"
)

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
		Audience:  customTokenAudience,
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
			return nil, &ErrTokenInvalid{fmt.Sprintf("firebaseauth: Firebase ID token has incorrect algorithm. Expected \"RSA\" but got \"%#v\"", token.Header["alg"])}
		}
		kid := token.Header["kid"].(string)
		if kid == "" {
			return nil, &ErrTokenInvalid{"firebaseauth: Firebase ID token has no \"kid\" claim"}
		}
		key := auth.selectKey(kid)
		if key == nil {
			return nil, &ErrTokenInvalid{"firebaseauth: Firebase ID token has \"kid\" claim which does not correspond to a known public key. Most likely the ID token is expired, so get a fresh token from your client app and try again"}
		}
		return key, nil
	})
	if err != nil {
		return nil, &ErrTokenInvalid{err.Error()}
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, &ErrTokenInvalid{"firebaseauth: invalid token"}
	}
	if !claims.verifyAudience(auth.app.projectID) {
		return nil, &ErrTokenInvalid{fmt.Sprintf("firebaseauth: Firebase ID token has incorrect \"aud\" (audience) claim. Expected \"%s\" but got \"%s\"", auth.app.projectID, claims.Audience)}
	}
	if !claims.verifyIssuer("https://securetoken.google.com/" + auth.app.projectID) {
		return nil, &ErrTokenInvalid{fmt.Sprintf("firebaseauth: Firebase ID token has incorrect \"iss\" (issuer) claim. Expected \"https://securetoken.google.com/%s\" but got \"%s\"", auth.app.projectID, claims.Issuer)}
	}
	if claims.Subject == "" {
		return nil, &ErrTokenInvalid{"firebaseauth: Firebase ID token has an empty string \"sub\" (subject) claim"}
	}
	if len(claims.Subject) > 128 {
		return nil, &ErrTokenInvalid{"firebaseauth: Firebase ID token has \"sub\" (subject) claim longer than 128 characters"}
	}

	claims.UserID = claims.Subject
	return claims, nil
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

// GetUser retrieves an user by user id
func (auth *Auth) GetUser(uid string) (*UserRecord, error) {
	users, err := auth.GetUsers([]string{uid})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

// GetUsers retrieves users by user ids
func (auth *Auth) GetUsers(uids []string) ([]*UserRecord, error) {
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

// GetUserByEmail retrieves user by email
func (auth *Auth) GetUserByEmail(email string) (*UserRecord, error) {
	users, err := auth.GetUsersByEmail([]string{email})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

// GetUsersByEmail retrieves users by emails
func (auth *Auth) GetUsersByEmail(emails []string) ([]*UserRecord, error) {
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

// DeleteUser deletes an user by user id
func (auth *Auth) DeleteUser(uid string) error {
	if len(uid) == 0 {
		return ErrRequireUID
	}

	return auth.app.invokeRequest(http.MethodPost, deleteAccount, &deleteAccountRequest{LocalID: uid}, &deleteAccountResponse{})
}

func (auth *Auth) createUserAutoID(user *User) (string, error) {
	var resp signupNewUserResponse
	if len(user.RawPassword) > 0 && len(user.Password) == 0 {
		// signup new user need password
		user.Password = user.RawPassword
		user.RawPassword = ""
	}
	err := auth.app.invokeRequest(http.MethodPost, signupNewUser, &signupNewUserRequest{user}, &resp)
	if err != nil {
		return "", err
	}
	if len(resp.LocalID) == 0 {
		return "", errors.New("firebaseauth: create account error")
	}
	return resp.LocalID, nil
}

func (auth *Auth) createUserCustomID(user *User) error {
	var resp uploadAccountResponse
	if len(user.RawPassword) == 0 && len(user.Password) > 0 {
		// upload account use raw password
		user.RawPassword = user.Password
		user.Password = ""
	}
	err := auth.app.invokeRequest(http.MethodPost, uploadAccount, &uploadAccountRequest{
		Users:          []*User{user},
		AllowOverwrite: false,
		SanityCheck:    true,
	}, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return errors.New("firebaseauth: upload account error")
	}
	return nil
}

// CreateUser creates an user
// if not provides UserID, firebase server will auto generate
func (auth *Auth) CreateUser(user *User) (*UserRecord, error) {
	var err error
	var userID string

	if len(user.UserID) == 0 {
		userID, err = auth.createUserAutoID(user)
	} else {
		userID = user.UserID
		err = auth.createUserCustomID(user)
	}
	if err != nil {
		return nil, err
	}

	res, err := auth.GetUser(userID)
	if err != nil {
		return nil, err
	}
	return res, nil
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
func (cursor *ListAccountCursor) Next() ([]*UserRecord, error) {
	var resp downloadAccountResponse
	err := cursor.auth.app.invokeRequest(http.MethodPost, downloadAccount, &downloadAccountRequest{
		MaxResults:    cursor.MaxResults,
		NextPageToken: cursor.nextPageToken,
	}, &resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Users) == 0 {
		return nil, iterator.Done
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
