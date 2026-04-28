package models

import (
	"strconv"
	"strings"
	"testing"
)

// newTestGame builds a 4-player started match with deterministic seating.
// Players: t1p1=0=alice, t2p1=1=bob, t1p2=2=carol, t2p2=3=dave.
func newTestGame() *GlobalGameState {
	g := &GlobalGameState{
		GameState: GameState{
			MatchInviteCode: "TEST00",
			Team1UserNames:  []string{"alice", "carol"},
			Team2UserNames:  []string{"bob", "dave"},
			Team1Connected:  []bool{true, true},
			Team2Connected:  []bool{true, true},
			Rules:           []string{},
		},
	}
	g.StartNextRound()
	// Make starting player deterministic: alice (position 0).
	g.CurrentStartingBidder = 0
	g.CurrentStartingPlayer = 0
	g.CurrentPlayerTurn = 0
	return g
}

func usernameForPosition(g *GlobalGameState, pos int) string {
	return g.GetUsernameByPosition(pos)
}

// TestBasicBiddingThenAllPassResetsRound: 4 passes should restart the round.
func TestBasicBiddingThenAllPassResetsRound(t *testing.T) {
	g := newTestGame()
	g.Rules = []string{RuleNoForcedBid} // allow all-pass

	startingRound := g.CurrentRound

	for i := 0; i < 4; i++ {
		username := usernameForPosition(g, g.CurrentPlayerTurn)
		if err := g.ProcessMove(username, "bid/0"); err != nil {
			t.Fatalf("bid %d failed: %v", i, err)
		}
	}

	if g.CurrentRound != startingRound+1 {
		t.Fatalf("expected round to advance to %d after all-pass; got %d",
			startingRound+1, g.CurrentRound)
	}
	if !g.IsInBidding {
		t.Fatalf("expected to be in bidding phase after re-deal")
	}
}

// TestForcedBidRejectsPassFromLastPlayer: with default rules, last bidder must bid.
func TestForcedBidRejectsPassFromLastPlayer(t *testing.T) {
	g := newTestGame()

	for i := 0; i < 3; i++ {
		username := usernameForPosition(g, g.CurrentPlayerTurn)
		if err := g.ProcessMove(username, "bid/0"); err != nil {
			t.Fatalf("bid %d failed: %v", i, err)
		}
	}

	username := usernameForPosition(g, g.CurrentPlayerTurn)
	err := g.ProcessMove(username, "bid/0")
	if err == nil {
		t.Fatalf("expected last-player pass to be rejected when forced bid is on")
	}
	if !strings.Contains(err.Error(), "30") {
		t.Fatalf("expected forced-bid error to mention 30; got %v", err)
	}
}

// TestRejectMoveFromWrongPlayer
func TestRejectMoveFromWrongPlayer(t *testing.T) {
	g := newTestGame()
	if err := g.ProcessMove("bob", "bid/30"); err == nil {
		t.Fatalf("expected wrong-player bid to be rejected")
	}
}

// TestBiddingAdvancesToCalling: after 4 bids with one winner, transitions to call phase.
func TestBiddingAdvancesToCalling(t *testing.T) {
	g := newTestGame()

	bids := []string{"bid/30", "bid/0", "bid/0", "bid/0"}
	for i, b := range bids {
		username := usernameForPosition(g, g.CurrentPlayerTurn)
		if err := g.ProcessMove(username, b); err != nil {
			t.Fatalf("bid %d failed: %v", i, err)
		}
	}

	if g.IsInBidding {
		t.Fatalf("expected bidding to be complete")
	}
	if g.RoundRules.Bid != 30 {
		t.Fatalf("expected winning bid 30; got %d", g.RoundRules.Bid)
	}
	if g.RoundRules.BiddingTeam != 1 {
		t.Fatalf("expected bidding team 1 (alice); got %d", g.RoundRules.BiddingTeam)
	}
	if g.RoundRules.Trump != RuleUndecided {
		t.Fatalf("expected trump Undecided; got %s", g.RoundRules.Trump)
	}
	if g.CurrentPlayerTurn != 0 {
		t.Fatalf("expected alice (0) to call; got turn %d", g.CurrentPlayerTurn)
	}
}

// TestCallSetsTrump verifies that a numeric trump call moves the round to play phase.
func TestCallSetsTrump(t *testing.T) {
	g := newTestGame()
	for _, b := range []string{"bid/30", "bid/0", "bid/0", "bid/0"} {
		username := usernameForPosition(g, g.CurrentPlayerTurn)
		if err := g.ProcessMove(username, b); err != nil {
			t.Fatalf("bid failed: %v", err)
		}
	}
	username := usernameForPosition(g, g.CurrentPlayerTurn)
	if err := g.ProcessMove(username, "call/3"); err != nil {
		t.Fatalf("call failed: %v", err)
	}
	if g.RoundRules.Trump != "3" {
		t.Fatalf("expected trump=3; got %s", g.RoundRules.Trump)
	}
	if !g.isPlaying() {
		t.Fatalf("expected play phase after trump call")
	}
}

// TestRejectInvalidPlayDomino: caller can't play a domino they don't have.
func TestRejectInvalidPlayDomino(t *testing.T) {
	g := newTestGame()
	for _, b := range []string{"bid/30", "bid/0", "bid/0", "bid/0"} {
		username := usernameForPosition(g, g.CurrentPlayerTurn)
		if err := g.ProcessMove(username, b); err != nil {
			t.Fatalf("bid failed: %v", err)
		}
	}
	username := usernameForPosition(g, g.CurrentPlayerTurn)
	if err := g.ProcessMove(username, "call/3"); err != nil {
		t.Fatalf("call failed: %v", err)
	}

	// Try to play a domino we definitely don't have. Construct any 0-6 domino
	// and look for one not in alice's hand.
	aliceHand := g.getPlayerDominoesByPosition(0)
	have := map[DominoName]bool{}
	for _, d := range aliceHand {
		have[d] = true
	}
	var notInHand DominoName
	for j := 0; j <= 6; j++ {
		for k := 0; k <= j; k++ {
			d := DominoName(strconv.Itoa(j) + "-" + strconv.Itoa(k))
			if !have[d] {
				notInHand = d
				break
			}
		}
		if notInHand != "" {
			break
		}
	}
	if notInHand == "" {
		t.Fatalf("alice somehow has all 28 dominoes")
	}

	err := g.ProcessMove("alice", "play/"+string(notInHand))
	if err == nil {
		t.Fatalf("expected play of non-held domino to be rejected")
	}
}

// TestRejectMatchEnded: once MatchWinningTeam is set, ProcessMove rejects.
func TestRejectMatchEnded(t *testing.T) {
	g := newTestGame()
	g.MatchWinningTeam = 1
	if err := g.ProcessMove("alice", "bid/30"); err == nil {
		t.Fatalf("expected move after match end to be rejected")
	}
}

// TestPlayThroughOneTrick: 4 players each play one valid (lowest available) domino;
// after the 4th, processEndOfTrick must fire and update scores or starting player.
func TestPlayThroughOneTrick(t *testing.T) {
	g := newTestGame()
	for _, b := range []string{"bid/30", "bid/0", "bid/0", "bid/0"} {
		username := usernameForPosition(g, g.CurrentPlayerTurn)
		if err := g.ProcessMove(username, b); err != nil {
			t.Fatalf("bid failed: %v", err)
		}
	}
	username := usernameForPosition(g, g.CurrentPlayerTurn)
	if err := g.ProcessMove(username, "call/0"); err != nil {
		// trump = 0 to keep things simple
		t.Fatalf("call failed: %v", err)
	}

	startingRound := g.CurrentRound
	startingHistoryLen := len(g.RoundHistory)

	// Play whatever's first in each player's hand. Suit-following may make this
	// fail occasionally, but each player has 7 dominoes — odds are good.
	// To keep the test deterministic regardless of legal moves, we walk hand and
	// pick any domino that validates.
	for trickStep := 0; trickStep < 4; trickStep++ {
		turn := g.CurrentPlayerTurn
		username := usernameForPosition(g, turn)
		hand := append([]DominoName{}, g.getPlayerDominoesByPosition(turn)...)
		played := false
		for _, d := range hand {
			if err := g.ProcessMove(username, "play/"+string(d)); err == nil {
				played = true
				break
			}
		}
		if !played {
			t.Fatalf("step %d: %s couldn't legally play any domino", trickStep, username)
		}
	}

	// After 4 plays, the trick has resolved — RoundHistory should contain at
	// least one trick-end system message.
	if len(g.RoundHistory) <= startingHistoryLen+4 {
		t.Fatalf("expected trick-end message in history; len=%d, was %d",
			len(g.RoundHistory), startingHistoryLen)
	}
	if g.CurrentRound != startingRound {
		// Round shouldn't have ended with bid=30 after just one trick.
		t.Fatalf("round transitioned unexpectedly")
	}
}
