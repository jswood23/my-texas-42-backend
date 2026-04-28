package models

import (
	"errors"
	"strconv"
	"strings"

	"my-texas-42-backend/util"
)

var twoMarkBidVariants = []string{RuleNil, RuleSplash, RulePlunge, RuleSevens}
var splashBidVariants = []string{RuleSplash, RulePlunge}
var nilDoublesCalls = []string{RuleDoublesHigh, RuleDoublesLow, RuleDoublesOwnSuit}
var trumpCallNames = []string{RuleFollowMe, RuleDoublesTrump}

// validateCall returns nil if `callStr` is a valid trump/variant call.
// Mirrors the calling branch of legacy checkValidity().
func (game *GlobalGameState) validateCall(callStr string) error {
	move := strings.Split(callStr, " ")[0]
	if move == "" {
		return errors.New("invalid call")
	}

	rules := game.RoundRules

	isForcedNil := false
	if util.StringSliceContains(game.Rules, RuleForcedNil) && len(game.RoundHistory) >= 4 {
		_, _, m0, ok0 := parseHistoryEntry(game.RoundHistory[0])
		_, _, m1, ok1 := parseHistoryEntry(game.RoundHistory[1])
		_, _, m2, ok2 := parseHistoryEntry(game.RoundHistory[2])
		_, _, m3, ok3 := parseHistoryEntry(game.RoundHistory[3])
		if ok0 && ok1 && ok2 && ok3 && m0 == "0" && m1 == "0" && m2 == "0" && m3 == "42" {
			isForcedNil = true
		}
	}

	is2MarkBid := util.StringSliceContains(twoMarkBidVariants, move)
	isSplashBid := util.StringSliceContains(splashBidVariants, move)

	if is2MarkBid && !util.StringSliceContains(game.Rules, move) {
		return errors.New(move + " is not allowed by the current match rules")
	}

	if is2MarkBid && rules.Bid < 84 && !(move == RuleNil && isForcedNil) {
		return errors.New(move + " can only be called as a 2-mark bid")
	}

	// Splash/Plunge partner is choosing trump — can't make another 2-mark call.
	if is2MarkBid && util.StringSliceContains(splashBidVariants, rules.Variant) {
		return errors.New("you cannot make this call for splash or plunge. Please choose a trump")
	}

	if isSplashBid {
		doublesInHand := 0
		for _, d := range game.getPlayerDominoesByPosition(game.CurrentPlayerTurn) {
			a, b, err := parseDomino(string(d))
			if err == nil && a == b {
				doublesInHand++
			}
		}
		if move == RuleSplash && doublesInHand < 3 {
			return errors.New("you must have 3 or more doubles to bid splash")
		}
		if move == RulePlunge && doublesInHand < 4 {
			return errors.New("you must have 4 or more doubles to bid plunge")
		}
	}

	if move == RuleNil {
		nilCalls := strings.Split(callStr, " ")
		if len(nilCalls) < 2 || !util.StringSliceContains(nilDoublesCalls, nilCalls[1]) {
			return errors.New("you must call doubles high, low, or a suit of their own for nil")
		}
	}

	if !is2MarkBid {
		// Must be a numeric trump (0..6) or one of the named trump calls.
		if _, err := strconv.Atoi(move); err != nil {
			if !util.StringSliceContains(trumpCallNames, move) {
				return errors.New("you must call the trump for the round")
			}
		} else {
			n, _ := strconv.Atoi(move)
			if n < 0 || n > 6 {
				return errors.New("you must call the trump for the round")
			}
		}
	}

	return nil
}

// applyCall mutates RoundRules per the player's call. Mirrors legacy setRoundRules.
func (game *GlobalGameState) applyCall(callStr string) {
	currentRules := game.RoundRules

	// numeric trump
	if _, err := strconv.Atoi(callStr); err == nil {
		currentRules.Trump = callStr
		game.RoundRules = currentRules
		game.CurrentPlayerTurn = game.CurrentStartingPlayer
		return
	}

	// Splash / Plunge: set variant and pass trump pick to partner.
	if util.StringSliceContains(splashBidVariants, callStr) {
		currentRules.Variant = callStr
		if callStr == RulePlunge {
			currentRules.Bid += 42
		}
		game.RoundRules = currentRules
		game.CurrentStartingPlayer = (game.CurrentStartingPlayer + 2) % 4
		game.CurrentPlayerTurn = game.CurrentStartingPlayer
		return
	}

	if callStr == RuleSevens {
		currentRules.Variant = RuleSevens
		currentRules.Trump = ""
	}

	if strings.Contains(callStr, RuleNil) {
		currentRules.Variant = RuleNil
		parts := strings.Split(callStr, " ")
		if len(parts) >= 2 {
			currentRules.Trump = parts[1]
		}
	}

	if strings.Contains(callStr, RuleFollowMe) {
		currentRules.Trump = RuleFollowMe
	}
	if strings.Contains(callStr, RuleDoublesTrump) {
		currentRules.Trump = RuleDoublesTrump
	}

	game.RoundRules = currentRules
	game.CurrentPlayerTurn = game.CurrentStartingPlayer
}
