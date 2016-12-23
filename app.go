package admin

// FirebaseApp type
type FirebaseApp struct {
	projectID string
}

type options struct {
	ProjectID string
}

// Option func
type Option func(*options)

// InitializeApp initializes firebase app
func InitializeApp(opts ...Option) (*FirebaseApp, error) {
	opt := &options{}
	for _, setter := range opts {
		setter(opt)
	}

	app := FirebaseApp{
		projectID: opt.ProjectID,
	}

	return &app, nil
}

// ProjectID sets project id to options
func ProjectID(projectID string) Option {
	return func(arg *options) {
		arg.ProjectID = projectID
	}
}

// Auth creates new FirebaseAuth instance
func (app *FirebaseApp) Auth() *FirebaseAuth {
	return &FirebaseAuth{app: app}
}
