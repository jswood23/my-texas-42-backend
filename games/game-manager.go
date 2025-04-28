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

func (gm *GameManager) CreateNewGame(matchId int, matchName string, matchPrivacy models.PrivacyLevel, rules []string, ownerUsername string) *models.GlobalGameState {
	// Generate a match invite code that doesn't already exist
	var matchInviteCode models.InviteCode
	for {
		matchInviteCode = util.GenerateInviteCode()
		if _, exists := gm.games[matchInviteCode]; !exists {
			break
		}
	}

	game := &models.GlobalGameState{
		GameState: models.GameState{
			MatchInviteCode:       matchInviteCode,
			MatchName:             matchName,
			MatchPrivacy:          matchPrivacy,
			Rules:                 rules,
			OwnerUsername:         ownerUsername,
			Team1UserNames:        []string{ownerUsername},
			Team2UserNames:        make([]string, 0),
			Team1Connected:        []bool{false},
			Team2Connected:        make([]bool, 0),
			CurrentRound:          0,
			CurrentStartingBidder: 0,
			CurrentStartingPlayer: 0,
			IsInBidding:           false,
			CurrentPlayerTurn:     0,
			RoundRules: models.RoundRules{
				Bid:         0,
				BiddingTeam: 0,
				Trump:       models.RuleUndecided,
				Variant:     models.RuleNoVariant,
			},
			Team1RoundScore:        0,
			Team2RoundScore:        0,
			CurrentTeam1TotalScore: 0,
			CurrentTeam2TotalScore: 0,
			RoundHistory:           make([]string, 0),
			TotalRoundHistory:      make([]string, 0),
		},
		HasStarted:        false,
		AllPlayerDominoes: [2][2][]models.DominoName{{{}, {}}, {{}, {}}},
		MatchId:           matchId,
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

func (gm *GameManager) GetGameByUsername(username string) (*models.GlobalGameState, error) {
	for _, game := range gm.games {
		for _, playerUsername := range game.GameState.Team1UserNames {
			if playerUsername == username {
				return game, nil
			}
		}
		for _, playerUsername := range game.GameState.Team2UserNames {
			if playerUsername == username {
				return game, nil
			}
		}
	}
	return nil, errors.New("game not found")
}
