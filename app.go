package admin

import (
	"context"
	"crypto/rsa"
	"net/http"

	jwtgo "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// App holds information about application configuration
type App struct {
	projectID            string
	jwtConfig            *jwt.Config
	privateKey           *rsa.PrivateKey
	databaseURL          string
	databaseAuthVariable interface{}
	client               *http.Client
	apiKey               string
}

// AppOptions is the firebase app options for initialize app
type AppOptions struct {
	ProjectID                    string
	ServiceAccount               []byte
	DatabaseURL                  string
	DatabaseAuthVariableOverride interface{}
	APIKey                       string
}

// InitializeApp initializes firebase application with options
func InitializeApp(ctx context.Context, options AppOptions) (*App, error) {
	var err error

	app := App{
		projectID:            options.ProjectID,
		databaseURL:          options.DatabaseURL,
		databaseAuthVariable: options.DatabaseAuthVariableOverride,
		apiKey:               options.APIKey,
	}

	if options.ServiceAccount != nil {
		app.jwtConfig, err = google.JWTConfigFromJSON(options.ServiceAccount, scopes...)
		if err != nil {
			return nil, err
		}
		app.privateKey, err = jwtgo.ParseRSAPrivateKeyFromPEM(app.jwtConfig.PrivateKey)
		if err != nil {
			return nil, err
		}
		app.client = app.jwtConfig.Client(ctx)
	} else {
		app.client, err = google.DefaultClient(ctx, scopes...)
		if err != nil {
			return nil, err
		}
	}

	return &app, nil
}

// Auth creates new Auth instance
// each instance has the same firebase app instance
// but difference public keys instance
// better create only one instance
func (app *App) Auth() *Auth {
	return newAuth(app)
}

// Database creates new Database instance
func (app *App) Database() *Database {
	return newDatabase(app)
}

// FCM creates new FCM instance
func (app *App) FCM() *FCM {
	return newFCM(app)
}
