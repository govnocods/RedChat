package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool{return true},
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		Hub: hub,
		Conn: conn, 
		Send: make(chan []byte),
		ID: r.RemoteAddr,
	}

	hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
