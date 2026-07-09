package websocket

import (
	"errors"
	"sync"

	"github.com/gofiber/websocket/v2"
)

var Clients = map[string]*websocket.Conn{}
var clientsMu sync.RWMutex

func WebSocketHandler(c *websocket.Conn) {

	userID := c.Query("user_id")

	// Store active websocket connections in memory for this process.
	// This is why clients must connect to the notification worker process.
	clientsMu.Lock()
	Clients[userID] = c
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(Clients, userID)
		clientsMu.Unlock()
	}()

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}

func SendToUser(userID string, payload interface{}) error {
	// This only finds users connected to the current process.
	// If websocket moves to another VM later, replace this with Redis/RabbitMQ pub-sub.
	clientsMu.RLock()
	conn, ok := Clients[userID]
	clientsMu.RUnlock()

	if !ok {
		return errors.New("user is not connected")
	}

	return conn.WriteJSON(payload)
}
