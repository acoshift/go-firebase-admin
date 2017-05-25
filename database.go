package admin

import (
	"fmt"
	_url "net/url"
	_path "path"
)

// Database type
type Database struct {
	app *App
}

// ServerValue
var (
	ServerValueTimestamp interface{} = struct {
		SV string `json:".sv"`
	}{"timestamp"}
)

func newDatabase(app *App) *Database {
	return &Database{
		app: app,
	}
}

// Ref returns a Reference for a path
func (database *Database) Ref(path string) *Reference {
	path = _path.Clean(path)
	return &Reference{database: database, path: path}
}

// RefFromURL returns a Reference from an url
func (database *Database) RefFromURL(url string) (*Reference, error) {
	u, err := _url.Parse(url)
	if err != nil {
		return nil, err
	}
	if u.Host != database.app.databaseURL {
		return nil, fmt.Errorf("firebasedatabase: invalid host %v", u.Host)
	}
	return database.Ref(u.Path), nil
}
