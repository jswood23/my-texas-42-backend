package sockets

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"my-texas-42-backend/models"
	"sync"
)

type ConnectionManager struct {
	connections models.ConnectionMap
	mu          sync.Mutex
}

var manager = &ConnectionManager{
	connections: make(models.ConnectionMap),
}

func GetConnectionManager() *ConnectionManager {
	return manager
}

// AddConnection adds a new WebSocket connection with the associated user ID.
func (cm *ConnectionManager) AddConnection(userID models.UserID, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[userID] = conn
}

// RemoveConnection removes a WebSocket connection by user ID.
func (cm *ConnectionManager) RemoveConnection(userID models.UserID) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if conn, ok := cm.connections[userID]; ok {
		conn.Close()
		delete(cm.connections, userID)
	}
}

// SendMessage sends a message to a specific user by user ID.
func (cm *ConnectionManager) SendMessage(userID models.UserID, message models.WebsocketMessageAPIModel) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if conn, ok := cm.connections[userID]; ok {
		messageByte, _ := json.Marshal(message)
		return conn.WriteMessage(websocket.TextMessage, messageByte)
	}
	return nil
}

// BroadcastMessage sends a message to all connected users.
func (cm *ConnectionManager) BroadcastMessage(message models.WebsocketMessageAPIModel) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for _, conn := range cm.connections {
		messageByte, _ := json.Marshal(message)
		conn.WriteMessage(websocket.TextMessage, messageByte)
	}
}

// HandleIncomingMessages handles incoming messages from a specific user.
func (cm *ConnectionManager) HandleIncomingMessages(userID models.UserID) {
	cm.mu.Lock()
	conn, ok := cm.connections[userID]
	cm.mu.Unlock()
	if !ok {
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			cm.RemoveConnection(userID)
			break
		}

		// TODO: Process the message (e.g., broadcast it, handle commands, etc.)
		println(message)
	}
}

func (cm *ConnectionManager) GetConnectionCount() int {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return len(cm.connections)
}
