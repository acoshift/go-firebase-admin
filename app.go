package firebase

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"

	jwtgo "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
)

// App holds information about application configuration
type App struct {
	projectID            string
	privateKey           *rsa.PrivateKey
	clientEmail          string
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
func InitializeApp(ctx context.Context, options AppOptions, opts ...option.ClientOption) (*App, error) {
	opts = append([]option.ClientOption{option.WithScopes(scopes...)}, opts...)

	var err error

	app := App{
		projectID:            options.ProjectID,
		databaseURL:          options.DatabaseURL,
		databaseAuthVariable: options.DatabaseAuthVariableOverride,
		apiKey:               options.APIKey,
	}

	// check service account first
	if len(options.ServiceAccount) > 0 {
		jwtConfig, err := google.JWTConfigFromJSON(options.ServiceAccount, scopes...)
		if err != nil {
			return nil, err
		}
		app.privateKey, err = jwtgo.ParseRSAPrivateKeyFromPEM(jwtConfig.PrivateKey)
		if err != nil {
			return nil, err
		}
		app.client = jwtConfig.Client(ctx)
		app.clientEmail = jwtConfig.Email

		// load project id
		var serviceAccount struct {
			ProjectID string `json:"project_id"`
		}
		json.Unmarshal(options.ServiceAccount, &serviceAccount)
		app.projectID = serviceAccount.ProjectID
	} else {
		// load service account using google's option
		app.client, _, err = transport.NewHTTPClient(ctx, opts...)
		if err != nil {
			app.client = http.DefaultClient
		}

		cred, err := transport.Creds(ctx, opts...)
		if err != nil {
			return nil, err
		}

		if len(app.projectID) == 0 {
			app.projectID = cred.ProjectID
		}

		// load private key from google credential
		var serviceAccount struct {
			PrivateKey  string `json:"private_key"`
			ClientEmail string `json:"client_email"`
		}
		json.Unmarshal(cred.JSON, &serviceAccount)
		if len(serviceAccount.PrivateKey) > 0 {
			app.privateKey, err = jwtgo.ParseRSAPrivateKeyFromPEM([]byte(serviceAccount.PrivateKey))
			if err != nil {
				return nil, err
			}
			app.clientEmail = serviceAccount.ClientEmail
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
