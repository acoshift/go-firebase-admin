package firebase

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"google.golang.org/api/identitytoolkit/v3"
	"google.golang.org/api/iterator"
)

// Auth type
type Auth struct {
	app       *App
	client    *identitytoolkit.RelyingpartyService
	keysMutex *sync.RWMutex
	keys      map[string]*rsa.PublicKey
	keysExp   time.Time

	Leeway time.Duration
}

const (
	keysEndpoint        = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"
	customTokenAudience = "https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit"
)

func newAuth(app *App) *Auth {
	gitClient, _ := identitytoolkit.New(app.client)
	return &Auth{
		app:       app,
		client:    gitClient.Relyingparty,
		keysMutex: &sync.RWMutex{},
	}
}

// CreateCustomToken creates a custom token used for client to authenticate
// with firebase server using signInWithCustomToken
// See https://firebase.google.com/docs/auth/admin/create-custom-tokens
func (auth *Auth) CreateCustomToken(userID string, claims interface{}) (string, error) {
	if auth.app.privateKey == nil {
		return "", ErrRequireServiceAccount
	}
	now := time.Now()
	payload := &customClaims{
		Issuer:    auth.app.clientEmail,
		Subject:   auth.app.clientEmail,
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
// See https://firebase.google.com/docs/auth/admin/verify-id-tokens
func (auth *Auth) VerifyIDToken(idToken string) (*Token, error) {
	token, err := jwt.ParseWithClaims(idToken, &Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, &ErrTokenInvalid{fmt.Sprintf("firebaseauth: Firebase ID token has incorrect algorithm. Expected \"RSA\" but got \"%#v\"", token.Header["alg"])}
		}
		kid, _ := token.Header["kid"].(string)
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

	claims, ok := token.Claims.(*Token)
	if !ok || !token.Valid {
		return nil, &ErrTokenInvalid{"firebaseauth: invalid token"}
	}

	now := time.Now().Unix()
	if !claims.verifyExpiresAt(now) {
		delta := time.Unix(now, 0).Sub(time.Unix(claims.ExpiresAt, 0))
		return nil, &ErrTokenInvalid{fmt.Sprintf("token is expired by %v", delta)}
	}
	if !claims.verifyIssuedAt(now + int64(auth.Leeway/time.Second)) {
		return nil, &ErrTokenInvalid{fmt.Sprintf("token used before issued")}
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

	return claims, nil
}

func (auth *Auth) fetchKeys() error {
	auth.keysMutex.Lock()
	defer auth.keysMutex.Unlock()
	resp, err := auth.app.client.Get(keysEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	auth.keysExp, _ = time.Parse(time.RFC1123, resp.Header.Get("Expires"))

	m := make(map[string]string)
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return err
	}
	ks := make(map[string]*rsa.PublicKey)
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
func (auth *Auth) GetUser(ctx context.Context, uid string) (*UserRecord, error) {
	users, err := auth.GetUsers(ctx, []string{uid})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users[0], nil
}

// GetUsers retrieves users by user ids
func (auth *Auth) GetUsers(ctx context.Context, userIDs []string) ([]*UserRecord, error) {
	resp, err := auth.client.GetAccountInfo(&identitytoolkit.IdentitytoolkitRelyingpartyGetAccountInfoRequest{
		LocalId: userIDs,
	}).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	return toUserRecords(resp.Users), nil
}

// GetUserByEmail retrieves user by email
func (auth *Auth) GetUserByEmail(ctx context.Context, email string) (*UserRecord, error) {
	users, err := auth.GetUsersByEmail(ctx, []string{email})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users[0], nil
}

// GetUsersByEmail retrieves users by emails
func (auth *Auth) GetUsersByEmail(ctx context.Context, emails []string) ([]*UserRecord, error) {
	resp, err := auth.client.GetAccountInfo(&identitytoolkit.IdentitytoolkitRelyingpartyGetAccountInfoRequest{
		Email: emails,
	}).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	return toUserRecords(resp.Users), nil
}

// GetUserByPhoneNumber retrieves user by phoneNumber
func (auth *Auth) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*UserRecord, error) {
	users, err := auth.GetUsersByPhoneNumber(ctx, []string{phoneNumber})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users[0], nil
}

// GetUsersByPhoneNumber retrieves users by phoneNumber
func (auth *Auth) GetUsersByPhoneNumber(ctx context.Context, phoneNumbers []string) ([]*UserRecord, error) {
	resp, err := auth.client.GetAccountInfo(&identitytoolkit.IdentitytoolkitRelyingpartyGetAccountInfoRequest{
		PhoneNumber: phoneNumbers,
	}).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	return toUserRecords(resp.Users), nil
}

// DeleteUser deletes an user by user id
func (auth *Auth) DeleteUser(ctx context.Context, userID string) error {
	if len(userID) == 0 {
		return ErrRequireUID
	}

	_, err := auth.client.DeleteAccount(&identitytoolkit.IdentitytoolkitRelyingpartyDeleteAccountRequest{
		LocalId: userID,
	}).Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}

func (auth *Auth) createUserAutoID(ctx context.Context, user *User) (string, error) {
	resp, err := auth.client.SignupNewUser(&identitytoolkit.IdentitytoolkitRelyingpartySignupNewUserRequest{
		Disabled:      user.Disabled,
		DisplayName:   user.DisplayName,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Password:      user.Password,
		PhotoUrl:      user.PhotoURL,
		PhoneNumber:   user.PhoneNumber,
	}).Context(ctx).Do()
	if err != nil {
		return "", err
	}
	return resp.LocalId, nil
}

func (auth *Auth) createUserCustomID(ctx context.Context, user *User) error {
	resp, err := auth.client.UploadAccount(&identitytoolkit.IdentitytoolkitRelyingpartyUploadAccountRequest{
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
				PhoneNumber:   user.PhoneNumber,
			},
		},
	}).Context(ctx).Do()
	if err != nil {
		return err
	}
	if len(resp.Error) > 0 {
		// TODO: merge errors into one error
		return errors.New("firebaseauth: create user error")
	}
	return nil
}

// CreateUser creates an user and return created user id
// if not provides UserID, firebase server will auto generate it
func (auth *Auth) CreateUser(ctx context.Context, user *User) (string, error) {
	var err error
	var userID string

	if len(user.UserID) == 0 {
		userID, err = auth.createUserAutoID(ctx, user)
	} else {
		userID = user.UserID
		err = auth.createUserCustomID(ctx, user)
	}
	if err != nil {
		return "", err
	}
	return userID, nil
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
func (cursor *ListAccountCursor) Next(ctx context.Context) ([]*UserRecord, error) {
	resp, err := cursor.auth.client.DownloadAccount(&identitytoolkit.IdentitytoolkitRelyingpartyDownloadAccountRequest{
		MaxResults:    cursor.MaxResults,
		NextPageToken: cursor.nextPageToken,
	}).Context(ctx).Do()
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
func (auth *Auth) UpdateUser(ctx context.Context, user *User) error {
	_, err := auth.client.SetAccountInfo(&identitytoolkit.IdentitytoolkitRelyingpartySetAccountInfoRequest{
		LocalId:       user.UserID,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Password:      user.Password,
		DisplayName:   user.DisplayName,
		DisableUser:   user.Disabled,
		PhotoUrl:      user.PhotoURL,
		PhoneNumber:   user.PhoneNumber,
	}).Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}

// SendPasswordResetEmail sends password reset for the given user
// Only useful for the Email/Password provider
func (auth *Auth) SendPasswordResetEmail(ctx context.Context, email string) error {
	_, err := auth.client.GetOobConfirmationCode(&identitytoolkit.Relyingparty{
		Email:       email,
		RequestType: "PASSWORD_RESET",
	}).Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}

// VerifyPassword verifies given email and password,
// return user id if success
func (auth *Auth) VerifyPassword(ctx context.Context, email, password string) (string, error) {
	resp, err := auth.client.VerifyPassword(&identitytoolkit.IdentitytoolkitRelyingpartyVerifyPasswordRequest{
		Email:    email,
		Password: password,
	}).Context(ctx).Do()
	if err != nil {
		return "", err
	}
	return resp.LocalId, nil
}

// CreateAuthURI creates auth uri for provider sign in
// returns auth uri for redirect
func (auth *Auth) CreateAuthURI(ctx context.Context, providerID string, continueURI string, sessionID string) (string, error) {
	resp, err := auth.client.CreateAuthUri(&identitytoolkit.IdentitytoolkitRelyingpartyCreateAuthUriRequest{
		ProviderId:   providerID,
		ContinueUri:  continueURI,
		AuthFlowType: "CODE_FLOW",
		SessionId:    sessionID,
	}).Context(ctx).Do()
	if err != nil {
		return "", err
	}
	return resp.AuthUri, nil
}

// VerifyAuthCallbackURI verifies callback uri after user redirect back from CreateAuthURI
// returns UserInfo if success
func (auth *Auth) VerifyAuthCallbackURI(ctx context.Context, callbackURI string, sessionID string) (*UserInfo, error) {
	resp, err := auth.client.VerifyAssertion(&identitytoolkit.IdentitytoolkitRelyingpartyVerifyAssertionRequest{
		RequestUri: callbackURI,
		SessionId:  sessionID,
	}).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	return &UserInfo{
		UserID:      resp.LocalId,
		DisplayName: resp.DisplayName,
		Email:       resp.Email,
		PhotoURL:    resp.PhotoUrl,
		ProviderID:  resp.ProviderId,
	}, nil
}
