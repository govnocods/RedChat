package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   int
	Username string
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
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
			log.Println("Invalid message JSON:", err)
			continue
		}

		content, _ := msg["content"].(string)

		err = c.Hub.messageService.SaveMessage(c.ID, []byte(content))
		if err != nil {
			log.Println("Failed to save message:", err)
		}

		msg["username"] = c.Username
		msg["sender_id"] = c.ID

		finalMsg, _ := json.Marshal(msg)

		// Рассылаем всем клиентам
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
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
