package models

import (
	"errors"
	"strconv"
	"strings"

	"my-texas-42-backend/util"
)

// Player numbering convention (matches legacy):
//   0 = team1[0], 1 = team2[0], 2 = team1[1], 3 = team2[1].

const historySep = "\\"
const dominoSep = "-"

// GetPlayerNumByUsername returns the interleaved player number (0..3) or -1 if not found.
func (game *GlobalGameState) GetPlayerNumByUsername(username string) int {
	if found, i := util.StringSliceContainsWithIndex(game.Team1UserNames, username); found {
		return i * 2
	}
	if found, i := util.StringSliceContainsWithIndex(game.Team2UserNames, username); found {
		return i*2 + 1
	}
	return -1
}

// GetUsernameByPosition returns the username at the given player number, or "" if out of range.
func (game *GlobalGameState) GetUsernameByPosition(position int) string {
	if position < 0 || position > 3 {
		return ""
	}
	idx := position / 2
	if position%2 == 0 {
		if idx < len(game.Team1UserNames) {
			return game.Team1UserNames[idx]
		}
	} else {
		if idx < len(game.Team2UserNames) {
			return game.Team2UserNames[idx]
		}
	}
	return ""
}

// getPlayerDominoesByPosition returns the hand of the player at the given position.
func (game *GlobalGameState) getPlayerDominoesByPosition(position int) []DominoName {
	return game.AllPlayerDominoes[position%2][position/2]
}

// removePlayerDomino removes the named domino from the hand at the given position.
// Returns true if removed.
func (game *GlobalGameState) removePlayerDomino(position int, domino DominoName) bool {
	dominoes := game.getPlayerDominoesByPosition(position)
	for i, d := range dominoes {
		if d == domino {
			newHand := make([]DominoName, 0, len(dominoes)-1)
			newHand = append(newHand, dominoes[:i]...)
			newHand = append(newHand, dominoes[i+1:]...)
			game.AllPlayerDominoes[position%2][position/2] = newHand
			return true
		}
	}
	return false
}

// isCalling: bidding is complete but trump has not been called.
func (game *GlobalGameState) isCalling() bool {
	return !game.IsInBidding && game.RoundRules.Trump == RuleUndecided
}

// isPlaying: the round is active (trump set).
func (game *GlobalGameState) isPlaying() bool {
	return !game.IsInBidding && game.RoundRules.Trump != RuleUndecided
}

// getPlayerPosition returns 0..3 — the order of the current player to act in the
// current trick (or bid round). Skips the nil-bidder's partner when applicable.
func (game *GlobalGameState) getPlayerPosition() int {
	if game.RoundRules.Variant != RuleNil {
		return (game.CurrentPlayerTurn + 4 - game.CurrentStartingPlayer) % 4
	}
	nilBidder := game.getNilBiddingPlayer()
	skippedPlayer := (nilBidder + 2) % 4

	newPosition := 0
	i := game.CurrentStartingPlayer
	for i != game.CurrentPlayerTurn {
		if i != skippedPlayer {
			newPosition++
		}
		i = (i + 1) % 4
	}
	return newPosition
}

// getNilBiddingPlayer returns the player number who placed the Nil call.
// Mirrors legacy getNilBiddingPlayer (defensive about Splash/Plunge call message at index 5).
func (game *GlobalGameState) getNilBiddingPlayer() int {
	if len(game.RoundHistory) < 6 {
		return 0
	}
	callMessage := game.RoundHistory[5]
	if _, _, move, ok := parseHistoryEntry(callMessage); ok {
		if move == RuleSplash || move == RulePlunge {
			if len(game.RoundHistory) >= 7 {
				callMessage = game.RoundHistory[6]
			}
		}
	}
	username, _, _, ok := parseHistoryEntry(callMessage)
	if !ok {
		return 0
	}
	n := game.GetPlayerNumByUsername(username)
	if n < 0 {
		return 0
	}
	return n
}

// formatHistoryEntry produces a "username\moveType\move" wire-format string.
func formatHistoryEntry(username, moveType, move string) string {
	return username + historySep + moveType + historySep + move
}

// parseHistoryEntry splits a "username\moveType\move" string. Returns ok=false if not 3 fields.
func parseHistoryEntry(entry string) (username, moveType, move string, ok bool) {
	parts := strings.Split(entry, historySep)
	if len(parts) != 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

// parseDomino parses "5-3" → (5, 3).
func parseDomino(d string) (int, int, error) {
	parts := strings.Split(d, dominoSep)
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid domino format")
	}
	a, err := strconv.Atoi(parts[0])
	if err != nil || a < 0 || a > 6 {
		return 0, 0, errors.New("invalid domino sides")
	}
	b, err := strconv.Atoi(parts[1])
	if err != nil || b < 0 || b > 6 {
		return 0, 0, errors.New("invalid domino sides")
	}
	return a, b, nil
}

// extractDominoFromHistory extracts the domino from a "user\play\6-4" history entry.
// Mirrors legacy `move.slice(-3).split('-')`.
func extractDominoFromHistory(entry string) (int, int, error) {
	if len(entry) < 3 {
		return 0, 0, errors.New("invalid history entry")
	}
	return parseDomino(entry[len(entry)-3:])
}

// containsDomino returns true if the slice contains target.
func containsDomino(dominoes []DominoName, target DominoName) bool {
	for _, d := range dominoes {
		if d == target {
			return true
		}
	}
	return false
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sliceContainsInt(s []int, v int) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

func tryParseSuit(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}
