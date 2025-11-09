package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/govnocods/RedChat/internal/auth"
	"github.com/govnocods/RedChat/internal/logger"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool{return true},
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		logger.Warn("WebSocket connection failed: missing auth token",
			"ip", r.RemoteAddr,
		)
		http.Error(w, "Missing auth token", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		logger.Warn("WebSocket connection failed: invalid token",
			"ip", r.RemoteAddr,
		)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.WithError(err).
			With("user_id", claims.UserId).
			With("username", claims.Username).
			With("ip", r.RemoteAddr).
			Error("Failed to upgrade WebSocket connection")
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
