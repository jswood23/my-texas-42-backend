package games

import (
	"errors"
	"my-texas-42-backend/models"
)

type GameManager struct {
	games       models.GameMap
	inviteCodes []models.InviteCode
}

var manager = &GameManager{
	games: make(models.GameMap),
}

func GetGameManager() *GameManager {
	return manager
}

func (gm *GameManager) GetGameCount() int {
	return len(gm.games)
}

func (gm *GameManager) AddGame(game *models.GlobalGameState) {
	gm.games[game.MatchInviteCode] = game
	gm.inviteCodes = append(gm.inviteCodes, game.MatchInviteCode)
}

func (gm *GameManager) RemoveGame(inviteCode models.InviteCode) {
	delete(gm.games, inviteCode)
	// remove the invite code from the list of invite codes
	for i, code := range gm.inviteCodes {
		if code == inviteCode {
			gm.inviteCodes = append(gm.inviteCodes[:i], gm.inviteCodes[i+1:]...)
			break
		}
	}
}

func (gm *GameManager) GetGameByInviteCode(inviteCode models.InviteCode) *models.GlobalGameState {
	return gm.games[inviteCode]
}

func (gm *GameManager) GetGameByUserID(userID models.UserID) (*models.GlobalGameState, error) {
	for _, game := range gm.games {
		for _, playerID := range game.Team1PlayerIDs {
			if playerID == userID {
				return game, nil
			}
		}
		for _, playerID := range game.Team2PlayerIDs {
			if playerID == userID {
				return game, nil
			}
		}
	}
	return nil, errors.New("game not found")
}
