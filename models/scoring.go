package models

import "strconv"

// DefaultMatchTargetMarks is the number of marks needed to win the match.
// Legacy had no match-end condition; the new server defaults to 7.
const DefaultMatchTargetMarks = 7

// matchTargetMarks returns the number of marks needed to win the match.
// (Configurable via Rules in a future change; currently returns the default.)
func (game *GlobalGameState) matchTargetMarks() int {
	return DefaultMatchTargetMarks
}

// processRoundWinner awards marks to the winning team and appends end-of-round
// messages to the round and total histories. Mirrors legacy processRoundWinner.
func (game *GlobalGameState) processRoundWinner(winningTeam int) {
	roundMarks := (game.RoundRules.Bid + 41) / 42 // ceil(bid / 42)
	if winningTeam == 1 {
		game.CurrentTeam1TotalScore += roundMarks
	} else {
		game.CurrentTeam2TotalScore += roundMarks
	}

	endOfRoundMessage := "Team " + strconv.Itoa(winningTeam) +
		" wins round worth " + strconv.Itoa(roundMarks) + " marks."
	game.RoundHistory = append(game.RoundHistory, endOfRoundMessage)
	game.TotalRoundHistory = append(game.TotalRoundHistory, endOfRoundMessage)
}

// isMatchOver returns the winning team (1 or 2) if a team has reached the
// match target, or 0 if the match continues. processRoundWinner only increments
// one team's score per round, so simultaneous match-end isn't possible.
func (game *GlobalGameState) isMatchOver() int {
	target := game.matchTargetMarks()
	if game.CurrentTeam1TotalScore >= target {
		return 1
	}
	if game.CurrentTeam2TotalScore >= target {
		return 2
	}
	return 0
}

// processEndOfTrick handles the moment when a trick (or bidding round) completes.
// Either advances bidding → calling, or scores the trick and possibly transitions
// the round / match. Mirrors legacy processEndOfTrick.
func (game *GlobalGameState) processEndOfTrick() {
	if game.IsInBidding {
		game.processBids()
		return
	}

	roundRules := game.RoundRules
	winningPlayer := game.getWinningPlayerOfTrick()
	if winningPlayer < 0 {
		return
	}
	winningPlayerUsername := game.GetUsernameByPosition(winningPlayer)
	winningTeam := 1
	if winningPlayer%2 != 0 {
		winningTeam = 2
	}
	trickScore := game.getTrickScore()

	endOfTrickMessage := "Team " + strconv.Itoa(winningTeam) +
		" (" + winningPlayerUsername + ") wins trick worth " +
		strconv.Itoa(trickScore) + " points."
	game.RoundHistory = append(game.RoundHistory, endOfTrickMessage)

	isEndOfRound := false
	if winningTeam == 1 {
		game.Team1RoundScore += trickScore
		if roundRules.Variant == RuleNil {
			// In Nil, team 1's hand is exhausted when both team-1 players are out.
			// (Nil bidder's partner sits out — they keep their dealt hand.)
			isRoundOver := game.nilBidderHandEmpty()
			if roundRules.BiddingTeam == 2 && isRoundOver {
				isEndOfRound = true
				game.processRoundWinner(2)
			} else if roundRules.BiddingTeam == 1 {
				// Nil bidder's team won a trick → bid broken immediately.
				isEndOfRound = true
				game.processRoundWinner(2)
			}
		} else if (roundRules.BiddingTeam == 1 && game.Team1RoundScore >= roundRules.Bid) ||
			(roundRules.BiddingTeam == 2 && game.Team1RoundScore > 42-roundRules.Bid) {
			isEndOfRound = true
			game.processRoundWinner(1)
		}
	} else {
		game.Team2RoundScore += trickScore
		if roundRules.Variant == RuleNil {
			isRoundOver := game.nilBidderHandEmpty()
			if roundRules.BiddingTeam == 1 && isRoundOver {
				isEndOfRound = true
				game.processRoundWinner(1)
			} else if roundRules.BiddingTeam == 2 {
				isEndOfRound = true
				game.processRoundWinner(1)
			}
		} else if (roundRules.BiddingTeam == 2 && game.Team2RoundScore >= roundRules.Bid) ||
			(roundRules.BiddingTeam == 1 && game.Team2RoundScore > 42-roundRules.Bid) {
			isEndOfRound = true
			game.processRoundWinner(2)
		}
	}

	if isEndOfRound {
		if winner := game.isMatchOver(); winner != 0 {
			game.MatchWinningTeam = winner
			game.RoundHistory = append(game.RoundHistory,
				"Team "+strconv.Itoa(winner)+" wins the match!")
			return
		}
		game.StartNextRound()
		return
	}

	game.CurrentStartingPlayer = winningPlayer
	game.CurrentPlayerTurn = winningPlayer
}

// nilBidderHandEmpty returns true when the Nil bidder has played all their dominoes.
// In Nil, the bidder's partner does not play, so we only check the bidder's hand.
func (game *GlobalGameState) nilBidderHandEmpty() bool {
	bidder := game.getNilBiddingPlayer()
	return len(game.getPlayerDominoesByPosition(bidder)) == 0
}
