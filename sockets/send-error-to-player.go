package sockets

import (
	"my-texas-42-backend/logger"
	"my-texas-42-backend/models"
)

func sendErrorToPlayer(username string, errorMessage string) {
	wsMessage := models.WSOutgoingMessageAPIModel{
		MessageType: models.WSMessageTypeGameError,
		Message:     errorMessage,
		Username:    "(System)",
		GameData:    nil,
	}
	err := GetConnectionManager().SendMessage(username, wsMessage)
	if err != nil {
		logger.Error("Failed to send error message to player " + username + ": " + err.Error())
	}
}
