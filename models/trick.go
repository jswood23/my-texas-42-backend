package models

import "strconv"

// playersInRound returns 3 for Nil (partner sits out) or 4 for normal play.
func (game *GlobalGameState) playersInRound() int {
	if game.RoundRules.Variant == RuleNil {
		return 3
	}
	return 4
}

// getTrickScore returns points for the just-completed trick: 1 + sum of pip-totals
// of any dominoes whose pip-sum is divisible by 5.
func (game *GlobalGameState) getTrickScore() int {
	dominoes := game.trickDominoes()
	if dominoes == nil {
		return 0
	}
	score := 1
	for _, d := range dominoes {
		if (d[0]+d[1])%5 == 0 {
			score += d[0] + d[1]
		}
	}
	return score
}

// trickDominoes returns the parsed sides of the dominoes in the just-completed trick.
func (game *GlobalGameState) trickDominoes() [][2]int {
	playersInRound := game.playersInRound()
	n := len(game.RoundHistory)
	if n < playersInRound {
		return nil
	}
	out := make([][2]int, 0, playersInRound)
	for _, entry := range game.RoundHistory[n-playersInRound:] {
		a, b, err := extractDominoFromHistory(entry)
		if err != nil {
			return nil
		}
		out = append(out, [2]int{a, b})
	}
	return out
}

// adjustWinningPlayerOfTrick converts the winning trick-index (0..playersInRound-1)
// to the actual player number, accounting for the Nil-bidder's partner being skipped.
func (game *GlobalGameState) adjustWinningPlayerOfTrick(winningIndex int) int {
	if game.RoundRules.Variant != RuleNil {
		return (game.CurrentStartingPlayer + winningIndex) % 4
	}
	nilBidder := game.getNilBiddingPlayer()
	skippedPlayer := (nilBidder + 2) % 4
	newWinning := game.CurrentStartingPlayer
	for i := 0; i < winningIndex; i++ {
		newWinning = (newWinning + 1) % 4
		if newWinning == skippedPlayer {
			newWinning = (newWinning + 1) % 4
		}
	}
	return newWinning
}

// getWinningPlayerOfTrick returns the player number who took the just-completed trick,
// per the variant-specific rules. Mirrors legacy getWinningPlayerOfTrick.
func (game *GlobalGameState) getWinningPlayerOfTrick() int {
	rules := game.RoundRules
	dominoes := game.trickDominoes()
	if dominoes == nil {
		return -1
	}
	playersInRound := len(dominoes)

	// Default starting suit = max side of the lead domino.
	startingSuit := dominoes[0][0]
	if dominoes[0][1] > startingSuit {
		startingSuit = dominoes[0][1]
	}

	isNilDoublesHigh := rules.Variant == RuleNil && rules.Trump == RuleDoublesHigh
	isNilDoublesLow := rules.Variant == RuleNil && rules.Trump == RuleDoublesLow
	isNilDoublesOwnSuit := rules.Variant == RuleNil && rules.Trump == RuleDoublesOwnSuit

	if rules.Variant == RuleSevens {
		smallest := -1
		var smallestIdxs []int
		for i, d := range dominoes {
			diff := absInt(7 - (d[0] + d[1]))
			if smallest == -1 || diff < smallest {
				smallest = diff
				smallestIdxs = []int{i}
			} else if diff == smallest {
				smallestIdxs = append(smallestIdxs, i)
			}
		}
		// Bidding team plays at trick positions 0 and 2 (alternating teams).
		isBiddingTeamWinning := false
		for _, idx := range smallestIdxs {
			if idx == 0 || idx == 2 {
				isBiddingTeamWinning = true
			}
		}
		winner := 1
		if isBiddingTeamWinning {
			winner = 0
		}
		return (game.CurrentStartingPlayer + winner) % 4
	}

	if isNilDoublesLow {
		wdIdx := 0
		wdSides := dominoes[0]
		for i := 1; i < playersInRound; i++ {
			cd := dominoes[i]
			isOutOfSuit := cd[0] != startingSuit && cd[1] != startingSuit
			if cd[0] == cd[1] || isOutOfSuit {
				continue
			}
			winningHigh := wdSides[0]
			if wdSides[0] == startingSuit {
				winningHigh = wdSides[1]
			}
			currentHigh := cd[0]
			if cd[0] == startingSuit {
				currentHigh = cd[1]
			}
			if currentHigh > winningHigh || wdSides[0] == wdSides[1] {
				wdIdx = i
				wdSides = cd
			}
		}
		return game.adjustWinningPlayerOfTrick(wdIdx)
	}

	if isNilDoublesOwnSuit {
		wdIdx := 0
		wdSides := dominoes[0]
		isStartingDominoDouble := wdSides[0] == wdSides[1]
		for i := 1; i < playersInRound; i++ {
			cd := dominoes[i]
			isOutOfSuit := cd[0] != startingSuit && cd[1] != startingSuit
			if isStartingDominoDouble {
				isOutOfSuit = cd[0] != cd[1]
			}
			if isOutOfSuit {
				continue
			}
			if isStartingDominoDouble {
				if cd[0] > wdSides[0] {
					wdIdx = i
					wdSides = cd
				}
			} else {
				winningHigh := wdSides[0]
				if wdSides[0] == startingSuit {
					winningHigh = wdSides[1]
				}
				currentHigh := cd[0]
				if cd[0] == startingSuit {
					currentHigh = cd[1]
				}
				if currentHigh > winningHigh {
					wdIdx = i
					wdSides = cd
				}
			}
		}
		return game.adjustWinningPlayerOfTrick(wdIdx)
	}

	if trump, err := strconv.Atoi(rules.Trump); err == nil {
		if trump >= 0 && (dominoes[0][0] == trump || dominoes[0][1] == trump) {
			startingSuit = trump
		}
		wdIdx := 0
		wdSides := dominoes[0]
		for i := 1; i < 4; i++ {
			cd := dominoes[i]
			isOutOfSuit := !(cd[0] == startingSuit || cd[1] == startingSuit ||
				cd[0] == trump || cd[1] == trump)
			if isOutOfSuit {
				continue
			}
			winnerHasTrump := wdSides[0] == trump || wdSides[1] == trump
			if winnerHasTrump {
				if wdSides[0] == wdSides[1] {
					continue
				}
				if !(cd[0] == trump || cd[1] == trump) {
					continue
				}
				winningHigh := wdSides[0]
				if wdSides[0] == trump {
					winningHigh = wdSides[1]
				}
				currentHigh := cd[0]
				if cd[0] == trump {
					currentHigh = cd[1]
				}
				if currentHigh > winningHigh || cd[0] == cd[1] {
					wdIdx = i
					wdSides = cd
				}
			} else {
				if cd[0] == trump || cd[1] == trump {
					wdIdx = i
					wdSides = cd
					continue
				}
				if wdSides[0] == wdSides[1] {
					continue
				}
				winningHigh := wdSides[0]
				if wdSides[0] == startingSuit {
					winningHigh = wdSides[1]
				}
				currentHigh := cd[0]
				if cd[0] == startingSuit {
					currentHigh = cd[1]
				}
				if currentHigh > winningHigh || cd[0] == cd[1] {
					wdIdx = i
					wdSides = cd
				}
			}
		}
		return (game.CurrentStartingPlayer + wdIdx) % 4
	}

	if rules.Trump == RuleFollowMe || isNilDoublesHigh {
		wdIdx := 0
		wdSides := dominoes[0]
		for i := 1; i < playersInRound; i++ {
			cd := dominoes[i]
			isOutOfSuit := cd[0] != startingSuit && cd[1] != startingSuit
			if wdSides[0] == wdSides[1] || isOutOfSuit {
				continue
			}
			winningHigh := wdSides[0]
			if wdSides[0] == startingSuit {
				winningHigh = wdSides[1]
			}
			currentHigh := cd[0]
			if cd[0] == startingSuit {
				currentHigh = cd[1]
			}
			if currentHigh > winningHigh || cd[0] == cd[1] {
				wdIdx = i
				wdSides = cd
			}
		}
		if isNilDoublesHigh {
			return game.adjustWinningPlayerOfTrick(wdIdx)
		}
		return (game.CurrentStartingPlayer + wdIdx) % 4
	}

	if rules.Trump == RuleDoublesTrump {
		wdIdx := 0
		wdSides := dominoes[0]
		for i := 1; i < 4; i++ {
			cd := dominoes[i]
			if wdSides[0] == wdSides[1] {
				if wdSides[0] == 6 {
					continue
				}
				if cd[0] != cd[1] {
					continue
				}
				if cd[0] > wdSides[0] {
					wdIdx = i
					wdSides = cd
				}
			} else {
				if cd[0] == cd[1] {
					wdIdx = i
					wdSides = cd
					continue
				}
				isOutOfSuit := cd[0] != startingSuit && cd[1] != startingSuit
				winningHigh := wdSides[0]
				if wdSides[0] == startingSuit {
					winningHigh = wdSides[1]
				}
				currentHigh := cd[0]
				if cd[0] == startingSuit {
					currentHigh = cd[1]
				}
				if currentHigh > winningHigh && !isOutOfSuit {
					wdIdx = i
					wdSides = cd
				}
			}
		}
		return (game.CurrentStartingPlayer + wdIdx) % 4
	}

	return -1
}
