package admin

import (
	"net"
	"net/http"
	_path "path"
	"time"
)

// Database type
type Database struct {
	app       *App
	transport *http.Transport
	client    *http.Client
}

// ServerValue
var (
	ServerValueTimestamp interface{} = struct {
		SV string `json:".sv"`
	}{"timestamp"}
)

func newDatabase(app *App) *Database {
	tr := &http.Transport{
		IdleConnTimeout: time.Minute * 5,
		MaxIdleConns:    20,
		Dial: func(network, address string) (net.Conn, error) {
			c, err := net.Dial(network, address)
			return c, err
		},
	}
	return &Database{
		app:       app,
		transport: tr,
		client:    &http.Client{Transport: tr},
	}
}

// Ref returns a Reference for a path
func (database *Database) Ref(path string) *Reference {
	path = _path.Clean(path)
	return &Reference{database: database, path: path}
}
