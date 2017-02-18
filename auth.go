package admin

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"google.golang.org/api/identitytoolkit/v3"
	"google.golang.org/api/iterator"
)

// Auth type
type Auth struct {
	app       *App
	client    *identitytoolkit.Service
	keysMutex *sync.RWMutex
	keys      map[string]*rsa.PublicKey
	keysExp   time.Time
}

const (
	keysEndpoint        = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"
	customTokenAudience = "https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit"
)

func newAuth(app *App) (*Auth, error) {
	client, err := identitytoolkit.New(oauth2.NewClient(context.Background(), app.tokenSource))
	if err != nil {
		return nil, err
	}
	return &Auth{
		app:       app,
		client:    client,
		keysMutex: &sync.RWMutex{},
	}, nil
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
func (auth *Auth) GetUsers(userIDs []string) ([]*UserRecord, error) {
	r := auth.client.Relyingparty.GetAccountInfo(&identitytoolkit.IdentitytoolkitRelyingpartyGetAccountInfoRequest{
		LocalId: userIDs,
	})
	resp, err := r.Do()
	if err != nil {
		return nil, err
	}
	if len(resp.Users) == 0 {
		return nil, nil
	}
	return toUserRecords(resp.Users), nil
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
	r := auth.client.Relyingparty.GetAccountInfo(&identitytoolkit.IdentitytoolkitRelyingpartyGetAccountInfoRequest{
		Email: emails,
	})
	resp, err := r.Do()
	if err != nil {
		return nil, err
	}
	if len(resp.Users) == 0 {
		return nil, nil
	}
	return toUserRecords(resp.Users), nil
}

// DeleteUser deletes an user by user id
func (auth *Auth) DeleteUser(userID string) error {
	if len(userID) == 0 {
		return ErrRequireUID
	}

	r := auth.client.Relyingparty.DeleteAccount(&identitytoolkit.IdentitytoolkitRelyingpartyDeleteAccountRequest{
		LocalId: userID,
	})
	_, err := r.Do()
	if err != nil {
		return err
	}
	return nil
}

func (auth *Auth) createUserAutoID(user *User) (string, error) {
	r := auth.client.Relyingparty.SignupNewUser(&identitytoolkit.IdentitytoolkitRelyingpartySignupNewUserRequest{
		Disabled:      user.Disabled,
		DisplayName:   user.DisplayName,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Password:      user.Password,
		PhotoUrl:      user.PhotoURL,
	})
	resp, err := r.Do()
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	if len(resp.LocalId) == 0 {
		return "", errors.New("firebaseauth: create account error")
	}
	return resp.LocalId, nil
}

func (auth *Auth) createUserCustomID(user *User) error {
	r := auth.client.Relyingparty.UploadAccount(&identitytoolkit.IdentitytoolkitRelyingpartyUploadAccountRequest{
		AllowOverwrite: false,
		SanityCheck:    true,
		Users: []*identitytoolkit.UserInfo{
			&identitytoolkit.UserInfo{
				LocalId:       user.UserID,
				Email:         user.Email,
				EmailVerified: user.EmailVerified,
				RawPassword:   user.Password,
				DisplayName:   user.DisplayName,
				Disabled:      user.Disabled,
				PhotoUrl:      user.PhotoURL,
			},
		},
	})
	resp, err := r.Do()
	if err != nil {
		return err
	}
	if len(resp.Error) > 0 {
		return errors.New("firebaseauth: create user error")
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
	MaxResults    int64
}

// ListUsers creates list account cursor for retrieves accounts
// MaxResults can change later after create cursor
func (auth *Auth) ListUsers(maxResults int64) *ListAccountCursor {
	return &ListAccountCursor{MaxResults: maxResults, auth: auth}
}

// Next retrieves next users from cursor which limit to MaxResults
// then move cursor to the next users
func (cursor *ListAccountCursor) Next() ([]*UserRecord, error) {
	r := cursor.auth.client.Relyingparty.DownloadAccount(&identitytoolkit.IdentitytoolkitRelyingpartyDownloadAccountRequest{
		MaxResults:    cursor.MaxResults,
		NextPageToken: cursor.nextPageToken,
	})
	resp, err := r.Do()
	if err != nil {
		return nil, err
	}
	if len(resp.Users) == 0 {
		return nil, iterator.Done
	}
	cursor.nextPageToken = resp.NextPageToken
	return toUserRecords(resp.Users), nil
}

// UpdateUser updates an existing user
func (auth *Auth) UpdateUser(user *User) (*UserRecord, error) {
	r := auth.client.Relyingparty.SetAccountInfo(&identitytoolkit.IdentitytoolkitRelyingpartySetAccountInfoRequest{
		LocalId:       user.UserID,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Password:      user.Password,
		DisplayName:   user.DisplayName,
		DisableUser:   user.Disabled,
		PhotoUrl:      user.PhotoURL,
	})
	resp, err := r.Do()
	if err != nil {
		return nil, err
	}

	res, err := auth.GetUser(resp.LocalId)
	if err != nil {
		return nil, err
	}
	return res, nil
}
