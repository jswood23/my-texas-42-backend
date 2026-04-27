package models

import (
	"errors"
	"math/rand"
	"my-texas-42-backend/util"
	"strconv"
	"strings"
)

const (
	RuleNoForcedBid    = "No forced bid"
	RuleForced31Bid    = "Forced 31 bid"
	RuleForcedNil      = "Forced Nil"
	RuleNil2Mark       = "2-mark Nil"
	RuleSplash2Mark    = "2-mark Splash"
	RulePlunge2Mark    = "2-mark Plunge"
	RuleSevens2Mark    = "2-mark Sevens"
	RuleDelve          = "Delve"
	RuleNil            = "Nil"
	RuleSplash         = "Splash"
	RulePlunge         = "Plunge"
	RuleSevens         = "Sevens"
	RuleDoublesHigh    = "Doubles-high"
	RuleDoublesLow     = "Doubles-low"
	RuleDoublesOwnSuit = "Doubles-own-suit"
	RuleDoublesTrump   = "Doubles-trump"
	RuleFollowMe       = "Follow-me"
	RuleUndecided      = "Undecided"
	RuleNoVariant      = ""
)

func (game *GlobalGameState) ContainsPlayer(username string) bool {
	isInTeam1, _ := util.StringSliceContainsWithIndex(game.Team1UserNames, username)
	isInTeam2, _ := util.StringSliceContainsWithIndex(game.Team2UserNames, username)
	return isInTeam1 || isInTeam2
}

func (game *GlobalGameState) ConnectDisconnectedPlayer(username string) {
	isInTeam1, _ := util.StringSliceContainsWithIndex(game.Team1UserNames, username)
	if isInTeam1 {
		for i, player := range game.Team1UserNames {
			if player == username {
				game.Team1Connected[i] = true
				break
			}
		}
	}

	isInTeam2, _ := util.StringSliceContainsWithIndex(game.Team2UserNames, username)
	if isInTeam2 {
		for i, player := range game.Team2UserNames {
			if player == username {
				game.Team2Connected[i] = true
				break
			}
		}
	}
}

func (game *GlobalGameState) SetPlayerAsDisconnected(username string) {
	isInTeam1, _ := util.StringSliceContainsWithIndex(game.Team1UserNames, username)
	if isInTeam1 {
		for i, player := range game.Team1UserNames {
			if player == username {
				game.Team1Connected[i] = false
				break
			}
		}
	}

	isInTeam2, _ := util.StringSliceContainsWithIndex(game.Team2UserNames, username)
	if isInTeam2 {
		for i, player := range game.Team2UserNames {
			if player == username {
				game.Team2Connected[i] = false
				break
			}
		}
	}
}

func (game *GlobalGameState) AddPlayer(username string, teamNumber int) error {
	if teamNumber == 1 {
		if len(game.Team1UserNames) >= 2 {
			return errors.New("team 1 is full")
		}
		game.Team1UserNames = append(game.Team1UserNames, username)
		game.Team1Connected = append(game.Team1Connected, true)
	} else if teamNumber == 2 {
		if len(game.Team2UserNames) >= 2 {
			return errors.New("team 2 is full")
		}
		game.Team2UserNames = append(game.Team2UserNames, username)
		game.Team2Connected = append(game.Team2Connected, true)
	} else {
		return errors.New("invalid team number")
	}

	return nil
}

func (game *GlobalGameState) GetPlayerGameState(username string) *PlayerGameState {
	isInTeam1, i := util.StringSliceContainsWithIndex(game.Team1UserNames, username)
	if isInTeam1 {
		return &PlayerGameState{
			GameState:      game.GameState,
			PlayerDominoes: game.AllPlayerDominoes[0][i],
			HasStarted:     game.HasStarted,
		}
	}

	isInTeam2, i := util.StringSliceContainsWithIndex(game.Team2UserNames, username)
	if isInTeam2 {
		return &PlayerGameState{
			GameState:      game.GameState,
			PlayerDominoes: game.AllPlayerDominoes[1][i],
			HasStarted:     game.HasStarted,
		}
	}

	return nil
}

func (game *GlobalGameState) GetAllUsernames() []string {
	allUsernames := make([]string, 0, len(game.Team1UserNames)+len(game.Team2UserNames))
	allUsernames = append(allUsernames, game.Team1UserNames...)
	allUsernames = append(allUsernames, game.Team2UserNames...)
	return allUsernames
}

func (game *GlobalGameState) GetAllConnectedUsernames() []string {
	allUsernames := make([]string, 0, len(game.Team1UserNames)+len(game.Team2UserNames))
	for i, player := range game.Team1UserNames {
		if game.Team1Connected[i] {
			allUsernames = append(allUsernames, player)
		}
	}
	for i, player := range game.Team2UserNames {
		if game.Team2Connected[i] {
			allUsernames = append(allUsernames, player)
		}
	}
	return allUsernames
}

func (game *GlobalGameState) SwitchPlayerTeam(username string) error {
	isInTeam1, i := util.StringSliceContainsWithIndex(game.Team1UserNames, username)
	if isInTeam1 {
		if len(game.Team2UserNames) >= 2 {
			return errors.New("team 2 is full")
		}

		// remove player from team 1
		game.Team1UserNames = append(game.Team1UserNames[:i], game.Team1UserNames[i+1:]...)
		game.Team1Connected = append(game.Team1Connected[:i], game.Team1Connected[i+1:]...)

		// add player to team 2
		game.Team2UserNames = append(game.Team2UserNames, username)
		game.Team2Connected = append(game.Team2Connected, true)

		return nil
	}

	isInTeam2, i := util.StringSliceContainsWithIndex(game.Team2UserNames, username)
	if !isInTeam2 {
		return errors.New("player not in any team")
	}

	if len(game.Team1UserNames) >= 2 {
		return errors.New("team 1 is full")
	}

	// remove player from team 2
	game.Team2UserNames = append(game.Team2UserNames[:i], game.Team2UserNames[i+1:]...)
	game.Team2Connected = append(game.Team2Connected[:i], game.Team2Connected[i+1:]...)

	// add player to team 1
	game.Team1UserNames = append(game.Team1UserNames, username)
	game.Team1Connected = append(game.Team1Connected, true)

	return nil
}

func (game *GlobalGameState) IsFull() bool {
	return len(game.Team1UserNames) == 2 && len(game.Team2UserNames) == 2
}

func (game *GlobalGameState) HasAllPlayersConnected() bool {
	for _, connected := range game.Team1Connected {
		if !connected {
			return false
		}
	}
	for _, connected := range game.Team2Connected {
		if !connected {
			return false
		}
	}
	return true
}

func (game *GlobalGameState) HasAllPlayersDisconnected() bool {
	for _, connected := range game.Team1Connected {
		if connected {
			return false
		}
	}
	for _, connected := range game.Team2Connected {
		if connected {
			return false
		}
	}
	return true
}

func (game *GlobalGameState) StartNextRound() {
	game.HasStarted = true

	gameState := &game.GameState
	gameState.RoundRules = RoundRules{
		Bid:         0,
		BiddingTeam: 0,
		Trump:       RuleUndecided,
		Variant:     RuleNoVariant,
	}

	var startingPlayer int
	if gameState.CurrentRound == 0 {
		startingPlayer = rand.Intn(4)
	} else {
		startingPlayer = (gameState.CurrentStartingBidder + 1) % 4
	}
	gameState.CurrentStartingBidder = startingPlayer
	gameState.CurrentStartingPlayer = startingPlayer
	gameState.CurrentPlayerTurn = startingPlayer

	gameState.IsInBidding = true
	gameState.CurrentRound++

	gameState.Team1RoundScore = 0
	gameState.Team2RoundScore = 0
	gameState.RoundHistory = []string{}

	game.AssignDominoes()
}

func (game *GlobalGameState) AssignDominoes() {
	// make a list of all dominoes in a double 6 set
	dominoes := make([]DominoName, 0, 28)
	for i := 0; i <= 6; i++ {
		for j := i; j <= 6; j++ {
			dominoes = append(dominoes, DominoName(strconv.Itoa(j)+":"+strconv.Itoa(i)))
		}
	}

	// shuffle the dominoes
	rand.Shuffle(len(dominoes), func(i, j int) {
		dominoes[i], dominoes[j] = dominoes[j], dominoes[i]
	})

	// assign dominoes to players
	game.AllPlayerDominoes = [2][2][]DominoName{
		{
			dominoes[0:7],
			dominoes[7:14],
		},
		{
			dominoes[14:21],
			dominoes[21:28],
		},
	}
}

func (game *GlobalGameState) ProcessMove(username string, moveStr string) error {
	err := game.validateTurn(username)
	if err != nil {
		return err
	}

	moveType, _, err := getMove(moveStr)
	if err != nil {
		return err
	}

	switch moveType {
	case MoveTypeBid:
		break
	case MoveTypePlay:
		break
	case MoveTypeCall:
		break
	}

	return nil
}

// getMove parses the move string and returns the move type, actual move, and an error if any, respectively
func getMove(moveStr string) (MoveType, ActualMove, error) {
	// validate that there is exactly one forward slash
	if strings.Count(moveStr, "/") != 1 {
		return "", "", errors.New("invalid move format")
	}

	// split the move string by forward slash
	parts := strings.Split(moveStr, "/")
	if len(parts) != 2 {
		return "", "", errors.New("invalid move format")
	}

	moveType := strings.TrimSpace(parts[0])
	actualMove := strings.TrimSpace(parts[1])
	if moveType == "" || actualMove == "" {
		return "", "", errors.New("invalid move format")
	}

	// validate that the move type is one of the allowed types
	allowedMoveTypes := []string{"bid", "play", "call"}
	if !util.StringSliceContains(allowedMoveTypes, moveType) {
		return "", "", errors.New("invalid move type")
	}

	return MoveType(moveType), ActualMove(actualMove), nil
}
