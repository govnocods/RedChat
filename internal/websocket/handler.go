package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/govnocods/RedChat/internal/auth"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool{return true},
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, "Missing auth token", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		ID: claims.UserId,
		Username: claims.Username,
		Hub: hub,
		Conn: conn, 
		Send: make(chan []byte),
	}

	hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
