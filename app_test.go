package firebase_test

import (
	"context"
	"os"
	"testing"

	"github.com/acoshift/go-firebase-admin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

type config struct {
	ProjectID                    string      `yaml:"projectId"`
	DatabaseURL                  string      `yaml:"databaseURL"`
	DatabaseAuthVariableOverride interface{} `yaml:"DatabaseAuthVariableOverride"`
	APIKey                       string      `yaml:"apiKey"`
}

func initApp(t *testing.T) *firebase.App {
	// t.Helper()

	// load config from env
	c := config{
		ProjectID:   os.Getenv("PROJECT_ID"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		APIKey:      os.Getenv("API_KEY"),
	}

	app, _ := firebase.InitializeApp(context.Background(), firebase.AppOptions{
		ProjectID:                    c.ProjectID,
		DatabaseURL:                  c.DatabaseURL,
		DatabaseAuthVariableOverride: c.DatabaseAuthVariableOverride,
		APIKey: c.APIKey,
	}, option.WithCredentialsFile("private/service_account.json"))
	return app
}

func TestAuth(t *testing.T) {

	app := initApp(t)
	firAuth := app.Auth()

	assert.NotNil(t, app)
	assert.NotNil(t, firAuth)
}

func TestDatabase(t *testing.T) {

	app := initApp(t)
	firDatabase := app.Database()

	assert.NotNil(t, app)
	assert.NotNil(t, firDatabase)
}

func TestFCM(t *testing.T) {

	app := initApp(t)
	firFCM := app.FCM()

	assert.NotNil(t, app)
	assert.NotNil(t, firFCM)
}
