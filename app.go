package admin

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"

	jwtgo "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/googleapi"
)

// App holds information about application configuration
type App struct {
	projectID      string
	serviceAccount string
	jwtConfig      *jwt.Config
	privateKey     *rsa.PrivateKey
	databaseURL    string
	tokenSource    oauth2.TokenSource
}

// AppOptions is the firebase app options for initialize app
type AppOptions struct {
	ProjectID      string
	ServiceAccount []byte
	DatabaseURL    string
}

// InitializeApp initializes firebase application with options
func InitializeApp(options AppOptions) (*App, error) {
	var err error

	app := App{
		projectID:   options.ProjectID,
		databaseURL: options.DatabaseURL,
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
		app.tokenSource = app.jwtConfig.TokenSource(context.Background())
	} else {
		app.tokenSource, err = google.DefaultTokenSource(context.Background(), scopes...)
		if err != nil {
			return nil, err
		}
	}

	return &app, nil
}

// Auth creates new Auth instance
// each instance has the save firebase app instance
// but difference public keys instance
// better create only one instance
func (app *App) Auth() *Auth {
	return newAuth(app)
}

// Database creates new Database instance
func (app *App) Database() *Database {
	return newDatabase(app)
}

func (app *App) invokeRequest(method string, api apiMethod, requestData interface{}, response interface{}) error {
	if app.tokenSource == nil {
		return ErrRequireServiceAccount
	}
	ctx, cancel := getContext()
	defer cancel()
	client := oauth2.NewClient(ctx, app.tokenSource)

	var resp *http.Response
	var err error
	path := baseURL + string(api)
	if method == http.MethodPost {
		var requestBytes []byte
		requestBytes, err = json.Marshal(requestData)
		if err != nil {
			return err
		}
		resp, err = client.Post(path, "application/json", bytes.NewReader(requestBytes))
	} else {
		resp, err = client.Get(path)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = googleapi.CheckResponse(resp); err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return err
	}
	return nil
}
