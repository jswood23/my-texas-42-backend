package sockets

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"my-texas-42-backend/games"
	"my-texas-42-backend/logger"
	"my-texas-42-backend/models"
	"my-texas-42-backend/util"
	"strconv"
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
		err := conn.Close()

		if err != nil {
			logger.Error("Failed to close ws connection: " + err.Error())
		}

		delete(cm.connections, userID)
	}
}

// SendMessage sends a message to a specific user by user ID.
func (cm *ConnectionManager) SendMessage(userID models.UserID, message models.WSOutgoingMessageAPIModel) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if conn, ok := cm.connections[userID]; ok {
		messageByte, _ := json.Marshal(message)
		return conn.WriteMessage(websocket.TextMessage, messageByte)
	}
	return nil
}

func (cm *ConnectionManager) SendMessageToGame(message models.WSOutgoingMessageAPIModel, game models.GlobalGameState) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for _, playerID := range append(game.Team1PlayerIDs, game.Team2PlayerIDs...) {
		if conn, ok := cm.connections[playerID]; ok {
			messageByte, _ := json.Marshal(message)
			err := conn.WriteMessage(websocket.TextMessage, messageByte)

			if err != nil {
				logger.Error("Failed to send ws message to user with ID " + strconv.Itoa(int(playerID)) + ": " + err.Error())
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

			games.HandleChatMessage(userID, data)
		} else if result.Action == "play_turn" {
			println("play turn")
		} else if result.Action == "refresh_player_game" {
			println("refresh player game")
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
