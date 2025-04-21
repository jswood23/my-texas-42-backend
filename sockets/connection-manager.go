package sockets

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"my-texas-42-backend/games"
	"my-texas-42-backend/logger"
	"my-texas-42-backend/models"
	"my-texas-42-backend/util"
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
func (cm *ConnectionManager) AddConnection(username string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[username] = conn
}

// RemoveConnection removes a WebSocket connection by user ID.
func (cm *ConnectionManager) RemoveConnection(username string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if conn, ok := cm.connections[username]; ok {
		err := conn.Close()

		if err != nil {
			logger.Error("Failed to close ws connection: " + err.Error())
		}

		delete(cm.connections, username)
	}
}

// SendMessage sends a message to a specific user by user ID.
func (cm *ConnectionManager) SendMessage(username string, message models.WSOutgoingMessageAPIModel) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if conn, ok := cm.connections[username]; ok {
		messageByte, _ := json.Marshal(message)
		return conn.WriteMessage(websocket.TextMessage, messageByte)
	}
	return nil
}

func (cm *ConnectionManager) SendMessageToGame(message models.WSOutgoingMessageAPIModel, game models.GlobalGameState) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for _, playerUsername := range append(game.GameState.Team1UserNames, game.GameState.Team2UserNames...) {
		if conn, ok := cm.connections[playerUsername]; ok {
			messageByte, _ := json.Marshal(message)
			err := conn.WriteMessage(websocket.TextMessage, messageByte)

			if err != nil {
				logger.Error("Failed to send ws message to user " + playerUsername + ": " + err.Error())
			}
		}
	}
}

// BroadcastMessage sends a message to all connected users.
func (cm *ConnectionManager) BroadcastMessage(message models.WSOutgoingMessageAPIModel) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for _, conn := range cm.connections {
		messageByte, _ := json.Marshal(message)
		err := conn.WriteMessage(websocket.TextMessage, messageByte)

		if err != nil {
			logger.Error("Failed to broadcast ws message: " + err.Error())
		}
	}
}

// handleIncomingMessages handles incoming messages from a specific user.
func (cm *ConnectionManager) handleIncomingMessages(username string) {
	cm.mu.Lock()
	conn, ok := cm.connections[username]
	cm.mu.Unlock()
	if !ok {
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			cm.RemoveConnection(username)
			break
		}

		var result models.WSIncomingMessageAPIModel
		err = json.Unmarshal(message, &result)

		if err != nil {
			logger.Error("Failed to unmarshal incoming ws message: " + err.Error())
		}

		if result.Action == "send_chat_message" {
			data, err := util.ConvertStringMapToType[models.WSSendChatMessageAPIModel](result.Data)
			if err != nil {
				logger.Error("Failed to cast data to WSSendChatMessageAPIModel: " + err.Error())
				continue
			}

			games.HandleChatMessage(username, data)
		} else if result.Action == "play_turn" {
			println("play turn")
		} else if result.Action == "refresh_player_game_state" {
			refreshGameStateForPlayer(cm, username)
		} else if result.Action == "switch_teams" {
			println("switch teams")
		}
	}
}

func (cm *ConnectionManager) GetConnectionCount() int {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return len(cm.connections)
}
