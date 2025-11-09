package websocket

import (
	"encoding/json"

	"github.com/govnocods/RedChat/internal/logger"
	"github.com/govnocods/RedChat/internal/service"
	"github.com/govnocods/RedChat/models"
)

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client

	messageService *service.MessageService
}

func NewHub(messageService *service.MessageService) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),

		messageService: messageService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			logger.Info("Client connected",
				"client_id", client.ID,
				"username", client.Username,
				"total_clients", len(h.Clients),
			)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				logger.Info("Client disconnected",
					"client_id", client.ID,
					"username", client.Username,
					"total_clients", len(h.Clients),
				)

			}

		case message := <-h.Broadcast:
			h.handleBroadcast(message)
		}
	}
}

func (h *Hub) handleBroadcast(message []byte) {
	var msg models.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		logger.WithError(err).Warn("Invalid message format in broadcast")
		return
	}

	for client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.Clients, client)
		}
	}
}
