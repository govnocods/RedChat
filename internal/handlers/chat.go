package handlers

import (
    "net/http"
    "github.com/govnocods/RedChat/internal/websocket"
)

type ChatHandler struct {
    Hub *websocket.Hub
}

func NewChatHandler(hub *websocket.Hub) *ChatHandler {
    return &ChatHandler{Hub: hub}
}

func (h *ChatHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
    websocket.ServeWS(h.Hub, w, r)
}