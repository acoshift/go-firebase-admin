package firebase_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/acoshift/go-firebase-admin"
	"github.com/stretchr/testify/assert"
)

type config struct {
	ProjectID                    string      `yaml:"projectId"`
	ServiceAccount               []byte      `yaml:"serviceAccount"`
	DatabaseURL                  string      `yaml:"databaseURL"`
	DatabaseAuthVariableOverride interface{} `yaml:"DatabaseAuthVariableOverride"`
	APIKey                       string      `yaml:"apiKey"`
}

func initApp(t *testing.T) *firebase.App {
	t.Helper()

	// load config from env
	c := config{
		ProjectID:      os.Getenv("PROJECT_ID"),
		ServiceAccount: []byte(os.Getenv("SERVICE_ACCOUNT")),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		APIKey:         os.Getenv("API_KEY"),
	}

	// if service account is in separate file service_account.json
	if len(c.ServiceAccount) <= 0 {
		serviceAccount, _ := ioutil.ReadFile("private/service_account.json")
		c.ServiceAccount = serviceAccount
	}

	app, _ := firebase.InitializeApp(context.Background(), firebase.AppOptions(c))
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
