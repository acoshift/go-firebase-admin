package admin

import (
	"net"
	"net/http"
	"time"
)

// Database type
type Database struct {
	app       *FirebaseApp
	transport *http.Transport
	client    *http.Client
}

// ServerValue
var (
	ServerValueTimestamp interface{} = map[string]string{".sv": "timestamp"}
)

func newDatabase(app *FirebaseApp) *Database {
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
	return &Reference{database: database, path: path}
}
