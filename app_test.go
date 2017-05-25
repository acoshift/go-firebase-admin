package admin_test

import (
	"context"
	"io/ioutil"

	admin "github.com/acoshift/go-firebase-admin"
	"gopkg.in/yaml.v2"
)

type config struct {
	ProjectID      string `yaml:"projectId"`
	ServiceAccount []byte `yaml:"serviceAccount"`
	DatabaseURL    string `yaml:"databaseURL"`
}

func initApp() *admin.App {
	// load config from ./private/config.yaml
	bs, _ := ioutil.ReadFile("private/config.yaml")
	var c config
	yaml.Unmarshal(bs, &c)

	app, _ := admin.InitializeApp(context.Background(), admin.AppOptions(c))
	return app
}
