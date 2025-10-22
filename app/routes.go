package app

func (a *App) routes() {
	a.Router.HandleFunc("/register", a.Handlers.RegisterHandler)
	a.Router.HandleFunc("/login", a.Handlers.AuthHAndler)
}
