package models

import (
	"errors"
	"fmt"
	"strconv"
)

func (game *GlobalGameState) validateTurn(username string) error {
	if !game.ContainsPlayer(username) {
		return fmt.Errorf("user %s is not part of the game", username)
	}

	if game.CurrentPlayerTurn%2 == 0 {
		if game.Team1UserNames[game.CurrentPlayerTurn/2] != username {
			return fmt.Errorf("it is not %s's turn", username)
		}
	} else {
		if game.Team2UserNames[game.CurrentPlayerTurn/2] != username {
			return fmt.Errorf("it is not %s's turn", username)
		}
	}

	return nil
}

func (game *GlobalGameState) validateBid(bid ActualMove) error {
	if !game.IsInBidding {
		return errors.New("game is not in bidding phase")
	}

	bidValue, err := strconv.Atoi(string(bid))

	if err != nil {
		return errors.New("invalid bid format")
	}

	currentHighestBid := game.getCurrentHighestBid()

	if bidValue <= currentHighestBid {
		return fmt.Errorf("bid must be higher than the current highest bid (%s)", strconv.Itoa(currentHighestBid))
	}
	return nil
}

func (game *GlobalGameState) getCurrentHighestBid() int {
	return 0
}
