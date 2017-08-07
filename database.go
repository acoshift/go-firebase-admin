package firebase

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

// Ref returns a Reference representing the location in the Database corresponding to the provided path.
// If no path is provided, the Reference will point to the root of the Database.
// See https://firebase.google.com/docs/reference/admin/node/admin.database.Database#ref
func (database *Database) Ref(path string) *Reference {
	path = _path.Clean(path)
	return &Reference{database: database, path: path}
}

// RefFromURL returns a Reference representing the location in the Database corresponding to the provided
// Firebase URL.
//
// An error is returned if the URL is not a valid Firebase Database URL or it has a different domain than
// the current Database instance.
//
// Note that all query parameters (orderBy, limitToLast, etc.) are ignored and are not applied to the
// returned Reference.
// See https://firebase.google.com/docs/reference/admin/node/admin.database.Database#refFromURL
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

// GoOnline Reconnects to the server and synchronizes the offline Database state with the server state.
//
// See https://firebase.google.com/docs/reference/admin/node/admin.database.Database#goOnline
func (database *Database) GoOnline() {
	panic(ErrNotImplemented)
}

// GoOffline Disconnects from the server (all Database operations will be completed offline).
//
// See https://firebase.google.com/docs/reference/admin/node/admin.database.Database#goOffline
func (database *Database) GoOffline() {
	panic(ErrNotImplemented)
}
