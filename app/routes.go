package app

import (
	"net/http"

	"github.com/govnocods/RedChat/internal/websocket"
)

func (a *App) routes() {
	a.Router.HandleFunc("/api/register", a.Handlers.RegisterHandler)
	a.Router.HandleFunc("/api/login", a.Handlers.AuthHandler)

	a.Router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/login.html")
	})

	a.Router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/register.html")
	})

	a.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/chat.html")
	})

	a.Router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(a.Hub, w, r)
	})
}
