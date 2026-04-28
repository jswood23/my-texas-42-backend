package sockets

import (
	"my-texas-42-backend/logger"
	"my-texas-42-backend/models"
)

func messagePlayersInGame(game *models.GlobalGameState, messageType string, message string) {
	usernames := game.GetAllConnectedUsernames()

	for _, username := range usernames {
		wsMessage := models.WSOutgoingMessageAPIModel{
			MessageType: messageType,
			Message:     message,
			Username:    "(System)",
			GameData:    game.GetPlayerGameState(username),
		}
		err := GetConnectionManager().SendMessage(username, wsMessage)
		if err != nil {
			logger.Error("Failed to send message to player " + username + ": " + err.Error())
			continue
		}
	}
}
