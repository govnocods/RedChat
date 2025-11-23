package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/govnocods/RedChat/internal/logger"
)

type Client struct {
	ID       int
	Username string
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg map[string]any
		if err := json.Unmarshal(message, &msg); err != nil {
			logger.WithError(err).
				With("client_id", c.ID).
				With("username", c.Username).
				Warn("Invalid message JSON")
			continue
		}

		content, ok := msg["content"].(string)
		if !ok || content == "" {
			logger.Warn("Invalid or missing content field",
				"client_id", c.ID,
				"username", c.Username,
			)
			continue
		}

		err = c.Hub.messageService.SaveMessage(c.ID, []byte(content))
		if err != nil {
			logger.WithError(err).
				With("client_id", c.ID).
				With("username", c.Username).
				Error("Failed to save message")
		}

		msg["username"] = c.Username
		msg["sender_id"] = c.ID
		msg["created_at"] = time.Now().Format(time.RFC3339)

		finalMsg, err := json.Marshal(msg)
		if err != nil {
			logger.WithError(err).
				With("client_id", c.ID).
				With("username", c.Username).
				Error("Failed to marshal message")
			continue
		}

		c.Hub.Broadcast <- finalMsg
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logger.WithError(err).
					With("client_id", c.ID).
					With("username", c.Username).
					Error("WebSocket write error")
				return
			}
		}
	}
}
