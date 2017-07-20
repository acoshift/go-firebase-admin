package admin_test

import (
	"context"
	"io/ioutil"

	admin "github.com/acoshift/go-firebase-admin"
	"gopkg.in/yaml.v2"
)

type config struct {
	ProjectID                    string      `yaml:"projectId"`
	ServiceAccount               []byte      `yaml:"serviceAccount"`
	DatabaseURL                  string      `yaml:"databaseURL"`
	DatabaseAuthVariableOverride interface{} `yaml:"DatabaseAuthVariableOverride"`
	APIKey                       string      `yaml:"apiKey"`
}

func initApp() *admin.App {
	// load config from ./private/config.yaml
	bs, _ := ioutil.ReadFile("private/config.yaml")
	var c config
	yaml.Unmarshal(bs, &c)

	// if service account is in separate file service_account.json
	if len(c.ServiceAccount) <= 0 {
		serviceAccount, _ := ioutil.ReadFile("private/service_account.json")
		c.ServiceAccount = serviceAccount
	}

	app, _ := admin.InitializeApp(context.Background(), admin.AppOptions(c))
	return app
}
