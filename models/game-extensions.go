package models

import "errors"

func (game *GlobalGameState) ContainsPlayer(username string) bool {
	isInTeam1, _ := teamContains(game.Team1UserNames, username)
	isInTeam2, _ := teamContains(game.Team2UserNames, username)
	return isInTeam1 || isInTeam2
}

func teamContains(slice []string, item string) (bool, int) {
	for i, v := range slice {
		if v == item {
			return true, i
		}
	}
	return false, -1
}

func (game *GlobalGameState) ConnectDisconnectedPlayer(username string) {
	isInTeam1, _ := teamContains(game.Team1UserNames, username)
	if isInTeam1 {
		for i, player := range game.Team1UserNames {
			if player == username {
				game.Team1Connected[i] = true
				break
			}
		}
	}

	isInTeam2, _ := teamContains(game.Team2UserNames, username)
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
	isInTeam1, _ := teamContains(game.Team1UserNames, username)
	if isInTeam1 {
		for i, player := range game.Team1UserNames {
			if player == username {
				game.Team1Connected[i] = false
				break
			}
		}
	}

	isInTeam2, _ := teamContains(game.Team2UserNames, username)
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
	isInTeam1, i := teamContains(game.Team1UserNames, username)
	if isInTeam1 {
		return &PlayerGameState{
			GameState:      game.GameState,
			PlayerDominoes: game.AllPlayerDominoes[0][i],
		}
	}

	isInTeam2, i := teamContains(game.Team1UserNames, username)
	if isInTeam2 {
		return &PlayerGameState{
			GameState:      game.GameState,
			PlayerDominoes: game.AllPlayerDominoes[1][i],
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
