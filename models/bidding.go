package models

import (
	"errors"
	"strconv"

	"my-texas-42-backend/util"
)

// validateBid returns nil if `bidStr` is a valid bid for the current player.
// Mirrors the bidding branch of the legacy checkValidity().
func (game *GlobalGameState) validateBid(bidStr string) error {
	bid, err := strconv.Atoi(bidStr)
	if err != nil {
		return errors.New("invalid bid input")
	}

	playerPosition := game.getPlayerPosition()

	highestBid := 0
	if playerPosition > 0 {
		n := len(game.RoundHistory)
		if n >= playerPosition {
			previousBids := game.RoundHistory[n-playerPosition:]
			for _, entry := range previousBids {
				_, _, move, ok := parseHistoryEntry(entry)
				if !ok {
					continue
				}
				bidNum, err := strconv.Atoi(move)
				if err == nil && bidNum > highestBid {
					highestBid = bidNum
				}
			}
		}
	}

	if bid != 0 && bid < 30 {
		return errors.New("the minimum bid is 30")
	}

	if bid > 42 && bid%42 != 0 {
		return errors.New("invalid bid amount")
	}

	if bid > 84 && bid/42 != highestBid/42+1 {
		return errors.New("you cannot bid more than 1 mark above the current bid")
	}

	if bid != 0 && bid <= highestBid {
		return errors.New("you must either bid higher than the current highest bid (" + strconv.Itoa(highestBid) + ") or pass")
	}

	if playerPosition == 3 &&
		!util.StringSliceContains(game.Rules, RuleNoForcedBid) &&
		highestBid == 0 {
		if util.StringSliceContains(game.Rules, RuleForced31Bid) && bid < 31 {
			return errors.New("you must bid 31 or higher")
		}
		if !util.StringSliceContains(game.Rules, RuleForced31Bid) && bid < 30 {
			return errors.New("you must bid 30 or higher")
		}
	}

	return nil
}

// processBids transitions from bidding to calling phase. Called after the 4th bid.
// If everyone passed, starts the next round.
func (game *GlobalGameState) processBids() {
	n := len(game.RoundHistory)
	if n < 4 {
		return
	}
	allBids := game.RoundHistory[n-4:]

	highestBid := 0
	bidWinner := 0
	for i, entry := range allBids {
		_, _, move, ok := parseHistoryEntry(entry)
		if !ok {
			continue
		}
		bidNum, err := strconv.Atoi(move)
		if err == nil && bidNum > highestBid {
			highestBid = bidNum
			bidWinner = (game.CurrentStartingBidder + i) % 4
		}
	}

	if highestBid == 0 {
		game.StartNextRound()
		return
	}

	username := game.GetUsernameByPosition(bidWinner)
	game.RoundHistory = append(game.RoundHistory, username+" has won the bid.")

	biddingTeam := 1
	if bidWinner%2 != 0 {
		biddingTeam = 2
	}
	game.RoundRules = RoundRules{
		Bid:         highestBid,
		BiddingTeam: biddingTeam,
		Trump:       RuleUndecided,
		Variant:     RuleNoVariant,
	}
	game.IsInBidding = false
	game.CurrentStartingPlayer = bidWinner
	game.CurrentPlayerTurn = bidWinner
}
