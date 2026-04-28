package sockets

import (
	"encoding/json"
	"errors"
	"my-texas-42-backend/games"
	"my-texas-42-backend/logger"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
	"regexp"
	"strconv"
)

const chatMessageMaxLength = 255

// invalid characters: ; ' " \n
const invalidCharactersRegex = `[;'"\n]`

func handleChatMessage(cm *ConnectionManager, username string, resultData string) {
	var data models.WSSendChatMessageAPIModel
	err := json.Unmarshal([]byte(resultData), &data)

	if err != nil {
		logger.Error("Failed to unmarshal chat message data: " + err.Error() + "\n" + resultData)
		return
	}

	game, err := games.GetGameManager().GetGameByUsername(username)
	if err != nil {
		sendErrorToPlayer(username, "Failed to get game: "+err.Error())
		cm.RemoveConnection(username)
		return
	}

	err = validateChatMessage(data.Message)
	if err != nil {
		sendErrorToPlayer(username, "Invalid chat message: "+err.Error())
		return
	}

	err = saveChatMessage(game.MatchId, username, data.Message)
	if err != nil {
		sendErrorToPlayer(username, err.Error())
		return
	}

	outgoingMessage := models.WSOutgoingMessageAPIModel{
		MessageType: models.WSMessageTypeChat,
		Message:     data.Message,
		Username:    username,
		GameData:    nil,
	}

	cm.SendMessageToGame(outgoingMessage, game)
}

func saveChatMessage(matchId int, username string, message string) error {
	query := sql_scripts.SaveChatMessage(matchId, username, message)
	err := services.Execute(query)
	if err != nil {
		return errors.New("failed to save chat message: " + err.Error())
	}
	return nil
}

func validateChatMessage(message string) error {
	if len(message) > chatMessageMaxLength {
		return errors.New("message must be shorter than " + strconv.Itoa(chatMessageMaxLength) + " characters")
	}

	// check for invalid characters with regex
	containsInvalidCharacters, err := regexp.MatchString(invalidCharactersRegex, message)
	if err != nil {
		return errors.New("failed to validate message with regex: " + err.Error())
	}

	if containsInvalidCharacters {
		return errors.New("message contains invalid characters")
	}

	return nil
}
