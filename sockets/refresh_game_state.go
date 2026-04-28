package sockets

import (
	"my-texas-42-backend/games"
	"my-texas-42-backend/models"
)

func refreshGameStateForPlayer(cm *ConnectionManager, username string) {
	game, err := games.GetGameManager().GetGameByUsername(username)

	if err != nil {
		sendErrorToPlayer(username, "Failed to get game: "+err.Error())
		cm.RemoveConnection(username)
		return
	}

	messagePlayersInGame(game, models.WSMessageTypeGameUpdate, "Game state updated")
}
