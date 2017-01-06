package admin

// FirebaseDatabase type
type FirebaseDatabase struct {
	app *FirebaseApp
}

// ServerValue
var (
	ServerValueTimestamp interface{} = map[string]string{".sv": "timestamp"}
)

func newFirebaseDatabase(app *FirebaseApp) *FirebaseDatabase {
	return &FirebaseDatabase{
		app: app,
	}
}

// Ref returns a Reference for a path
func (database *FirebaseDatabase) Ref(path string) *Reference {
	return &Reference{database: database, path: path}
}
