package models

import "errors"

// validatePlay returns nil if the player can legally play `playStr` (e.g. "5-3").
// Mirrors the play branch of legacy checkValidity().
func (game *GlobalGameState) validatePlay(playStr string) error {
	a, b, err := parseDomino(playStr)
	if err != nil {
		return errors.New("invalid domino input")
	}

	playerHand := game.getPlayerDominoesByPosition(game.CurrentPlayerTurn)
	if !containsDomino(playerHand, DominoName(playStr)) {
		return errors.New("you do not have this domino")
	}

	rules := game.RoundRules

	if rules.Variant == RuleSevens {
		thisDiff := absInt(7 - (a + b))
		smallestDiff := thisDiff
		for _, d := range playerHand {
			x, y, err := parseDomino(string(d))
			if err != nil {
				continue
			}
			if diff := absInt(7 - (x + y)); diff < smallestDiff {
				smallestDiff = diff
			}
		}
		if thisDiff > smallestDiff {
			return errors.New("you must play the domino in your hand that is closest to seven")
		}
		return nil
	}

	playerPosition := game.getPlayerPosition()
	if playerPosition == 0 {
		return nil
	}

	n := len(game.RoundHistory)
	if n < playerPosition {
		return nil
	}
	previousMoves := game.RoundHistory[n-playerPosition:]
	if len(previousMoves) == 0 {
		return nil
	}

	firstA, firstB, err := extractDominoFromHistory(previousMoves[0])
	if err != nil {
		return nil
	}

	trump := rules.Trump
	isDoublesOwnSuit := trump == RuleDoublesTrump || trump == RuleDoublesOwnSuit
	isStartingDominoDouble := firstA == firstB

	trumpInt, trumpIsNumeric := tryParseSuit(trump)

	var startingSuit int
	startingSuitIsTrump := false
	isDominoNotStartingSuit := false

	if isDoublesOwnSuit && isStartingDominoDouble {
		isDominoNotStartingSuit = a != b
		startingSuit = firstA
	} else if isDoublesOwnSuit && !isStartingDominoDouble {
		startingSuit = firstA
		if a == b {
			isDominoNotStartingSuit = true
		} else {
			isDominoNotStartingSuit = !(a == startingSuit || b == startingSuit)
		}
	} else {
		// trump on the lead domino → trump is the suit; otherwise the higher side leads
		if trumpIsNumeric && (firstA == trumpInt || firstB == trumpInt) {
			startingSuit = trumpInt
			startingSuitIsTrump = true
		} else {
			startingSuit = firstA
			startingSuitIsTrump = trumpIsNumeric && trumpInt == firstA
		}

		hasStartingSuit := a == startingSuit || b == startingSuit
		hasTrump := trumpIsNumeric && (a == trumpInt || b == trumpInt)

		// If the led suit is non-trump and this domino has both led suit and trump,
		// it must be played as trump (not following suit).
		isDominoNotStartingSuit = !hasStartingSuit ||
			(!startingSuitIsTrump && hasStartingSuit && hasTrump)
	}

	if !isDominoNotStartingSuit {
		return nil
	}

	// Player can only play a non-led-suit domino if they have no led-suit-following domino.
	hasViable := false
	for _, d := range playerHand {
		x, y, err := parseDomino(string(d))
		if err != nil {
			continue
		}

		normalStart := false
		if !isDoublesOwnSuit {
			hasSS := x == startingSuit || y == startingSuit
			hasTr := trumpIsNumeric && (x == trumpInt || y == trumpInt)
			ssIsTrump := trumpIsNumeric && trumpInt == startingSuit
			normalStart = hasSS && !(!ssIsTrump && hasSS && hasTr)
		}
		doublesStartDouble := false
		if isDoublesOwnSuit && isStartingDominoDouble && x == y {
			doublesStartDouble = true
		}
		doublesStartNotDouble := false
		if isDoublesOwnSuit && !isStartingDominoDouble {
			hasSS := x == startingSuit || y == startingSuit
			if hasSS && x != y {
				doublesStartNotDouble = true
			}
		}
		if normalStart || doublesStartDouble || doublesStartNotDouble {
			hasViable = true
			break
		}
	}

	if hasViable {
		return errors.New("you need to play a domino that matches the starting suit")
	}
	return nil
}

// applyPlay removes the played domino from the player's hand. Mirrors legacy playDomino.
func (game *GlobalGameState) applyPlay(playStr string) {
	game.removePlayerDomino(game.CurrentPlayerTurn, DominoName(playStr))
}
