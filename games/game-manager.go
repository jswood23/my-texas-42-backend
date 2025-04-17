package games

import (
	"errors"
	"my-texas-42-backend/models"
	"my-texas-42-backend/util"
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

func (gm *GameManager) GetAllGames() models.GameMap {
	return gm.games
}

func (gm *GameManager) CreateNewGame(matchName string, matchPrivacy models.PrivacyLevel, rules []string, ownerUsername string) *models.GlobalGameState {
	game := &models.GlobalGameState{
		GameState: models.GameState{
			MatchInviteCode:        util.GenerateInviteCode(),
			MatchName:              matchName,
			MatchPrivacy:           matchPrivacy,
			Rules:                  rules,
			OwnerUsername:          ownerUsername,
			Team1UserNames:         []string{ownerUsername},
			Team2UserNames:         make([]string, 0),
			IsConnected:            []bool{false},
			CurrentRound:           0,
			CurrentStartingBidder:  0,
			CurrentStartingPlayer:  0,
			CurrentIsBidding:       false,
			CurrentPlayerTurn:      0,
			CurrentRoundRules:      nil,
			CurrentTeam1RoundScore: 0,
			CurrentTeam2RoundScore: 0,
			CurrentTeam1TotalScore: 0,
			CurrentTeam2TotalScore: 0,
			CurrentRoundHistory:    make([]string, 0),
			TotalRoundHistory:      make([]string, 0),
		},
		HasStarted:        false,
		AllPlayerDominoes: make([]models.DominoName, 0),
		Team1PlayerIDs:    make([]models.UserID, 0),
		Team2PlayerIDs:    make([]models.UserID, 0),
	}

	gm.addGame(game)

	return game
}

func (gm *GameManager) addGame(game *models.GlobalGameState) {
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
