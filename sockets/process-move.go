package sockets

import (
	"my-texas-42-backend/games"
	"my-texas-42-backend/models"
)

func processMove(cm *ConnectionManager, username string, move string) {
	game, err := games.GetGameManager().GetGameByUsername(username)
	if err != nil {
		sendErrorToPlayer(username, "Failed to get game: "+err.Error())
		cm.RemoveConnection(username)
		return
	}

	err = game.ProcessMove(username, move)

	if err != nil {
		sendErrorToPlayer(username, "Invalid move: "+err.Error())
		return
	}

	messagePlayersInGame(game, models.WSMessageTypeGameUpdate, username+" made a move.")
}
