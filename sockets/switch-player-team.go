package sockets

import (
	"my-texas-42-backend/games"
	"my-texas-42-backend/models"
)

func switchPlayerTeam(cm *ConnectionManager, username string) {
	// lock the connection manager to avoid a race condition for switching teams
	cm.mu.Lock()

	game, err := games.GetGameManager().GetGameByUsername(username)
	if err != nil {
		sendErrorToPlayer(username, "Failed to get game: "+err.Error())
		cm.RemoveConnection(username)
		return
	}

	err = game.SwitchPlayerTeam(username)

	if err != nil {
		sendErrorToPlayer(username, "Failed to switch teams: "+err.Error())
		return
	}

	cm.mu.Unlock()

	messagePlayersInGame(game, models.WSMessageTypeGameUpdate, username+" switched teams.")
}
