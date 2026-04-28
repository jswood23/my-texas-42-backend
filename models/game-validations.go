package models

import "errors"

func (game *GlobalGameState) validateTurn(username string) error {
	var currentTurnUsername string
	if game.CurrentPlayerTurn%2 == 0 {
		currentTurnUsername = game.Team1UserNames[game.CurrentPlayerTurn/2]
	} else {
		// dividing by 2 truncates toward zero, therefore 1/2 = 0 and 3/2 = 1
		currentTurnUsername = game.Team2UserNames[game.CurrentPlayerTurn/2]
	}

	if username != currentTurnUsername {
		return errors.New("it is not your turn")
	}

	return nil
}

