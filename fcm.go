package admin

// FCM type
type FCM struct {
	app *App
}

func newFCM(app *App) *FCM {
	return &FCM{
		app: app,
	}
}
