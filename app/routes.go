package app

import (
	"net/http"

	"github.com/govnocods/RedChat/internal/websocket"
)

func (a *App) routes() {
	a.Router.HandleFunc("/register", a.Handlers.RegisterHandler)
	a.Router.HandleFunc("/login", a.Handlers.AuthHAndler)

	a.Router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(a.Hub, w, r)
	})

	a.Router.Handle("/", http.FileServer(http.Dir("./web")))
}
