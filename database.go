package admin

// Database type
type Database struct {
	app *FirebaseApp
}

// ServerValue
var (
	ServerValueTimestamp interface{} = map[string]string{".sv": "timestamp"}
)

func newDatabase(app *FirebaseApp) *Database {
	return &Database{
		app: app,
	}
}

// Ref returns a Reference for a path
func (database *Database) Ref(path string) *Reference {
	return &Reference{database: database, path: path}
}
