package app

import (
	"net/http"

	"github.com/govnocods/RedChat/internal/websocket"
)

func (a *App) routes() {
	// Публичные маршруты (без аутентификации)
	a.Router.HandleFunc("/api/register", a.Handlers.RegisterHandler)
	a.Router.HandleFunc("/api/login", a.Handlers.AuthHandler)

	a.Router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/login.html")
	})

	a.Router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/register.html")
	})

	// Защищенные маршруты (требуют аутентификации)
	chatHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/chat.html")
	})
	a.Router.Handle("/", a.Middlewares.AuthMiddleware(chatHandler))

	wsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(a.Hub, w, r)
	})
	a.Router.Handle("/ws", a.Middlewares.AuthMiddleware(wsHandler))
}
