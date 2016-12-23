package admin

import (
	"crypto/rsa"

	jwtgo "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// FirebaseApp type
type FirebaseApp struct {
	projectID      string
	serviceAccount string
	jwtConfig      *jwt.Config
	privateKey     *rsa.PrivateKey
	databaseURL    string
}

type options struct {
	ProjectID      string
	ServiceAccount []byte
	DatabaseURL    string
}

// OptionFunc type
type OptionFunc func(*options)

// InitializeApp initializes firebase app
func InitializeApp(opts ...OptionFunc) (*FirebaseApp, error) {
	var err error
	opt := &options{}
	for _, setter := range opts {
		setter(opt)
	}

	app := FirebaseApp{
		projectID:   opt.ProjectID,
		databaseURL: opt.DatabaseURL,
	}

	if opt.ServiceAccount != nil {
		app.jwtConfig, err = google.JWTConfigFromJSON(opt.ServiceAccount)
		if err != nil {
			return nil, err
		}
		app.privateKey, err = jwtgo.ParseRSAPrivateKeyFromPEM(app.jwtConfig.PrivateKey)
		if err != nil {
			return nil, err
		}
	}

	return &app, nil
}

// ProjectID sets project id to options
func ProjectID(projectID string) OptionFunc {
	return func(arg *options) {
		arg.ProjectID = projectID
	}
}

// ServiceAccount sets service account to options
func ServiceAccount(serviceAccount []byte) OptionFunc {
	return func(arg *options) {
		arg.ServiceAccount = serviceAccount
	}
}

// DatabaseURL sets database url to options
func DatabaseURL(url string) OptionFunc {
	return func(arg *options) {
		arg.DatabaseURL = url
	}
}

// Auth creates new FirebaseAuth instance
func (app *FirebaseApp) Auth() *FirebaseAuth {
	return newFirebaseAuth(app)
}
