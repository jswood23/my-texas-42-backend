package models

import "github.com/gorilla/websocket"

type UserID int

type LogLevel int

type UserModel struct {
	UserID      UserID `db:"userid"`
	Username    string `db:"username"`
	Email       string `db:"email"`
	IsAdmin     bool   `db:"isadmin"`
	DisplayName string `db:"displayname"`
}

type UserStatsModel struct {
	UserStatsId          int    `db:"userstatsid"`
	GamesPlayed          int    `db:"gamesplayed"`
	GamesWon             int    `db:"gameswon"`
	RoundsPlayed         int    `db:"roundsplayed"`
	RoundsWon            int    `db:"roundswon"`
	TotalPointsAsBidder  int    `db:"totalpointsasbidder"`
	TotalRoundsAsBidder  int    `db:"totalroundsasbidder"`
	TotalPointsAsSupport int    `db:"totalpointsassupport"`
	TotalRoundsAsSupport int    `db:"totalroundsassupport"`
	TotalPointsAsCounter int    `db:"totalpointsascounter"`
	TotalRoundsAsCounter int    `db:"totalroundsascounter"`
	TimesWinningBidTotal int    `db:"timeswinningbidtotal"`
	TimesCallingSuit     int    `db:"timescallingsuit"`
	TimesCallingNil      int    `db:"timescallingnil"`
	TimesCallingSplash   int    `db:"timescallingsplash"`
	TimesCallingPlunge   int    `db:"timescallingplunge"`
	TimesCallingSevens   int    `db:"timescallingsevens"`
	TimesCallingDelve    int    `db:"timescallingdelve"`
	UserID               UserID `db:"userid"`
}

type FriendRequestModel struct {
	FriendRequestID UserID `db:"friendrequestid"`
	SenderUserID    UserID `db:"senderuserid"`
	ReceiverUserID  UserID `db:"receiveruserid"`
}

type FriendModel struct {
	FriendsID int    `db:"friendsid"`
	User1ID   UserID `db:"user1id"`
	User2ID   UserID `db:"user2id"`
}

type ConnectionMap map[UserID]*websocket.Conn

type RoundRules struct {
	Bid         int    `json:"bid"`
	BiddingTeam int    `json:"biddingTeam"`
	Trump       string `json:"trump"`
	Variant     string `json:"variant"`
}

type InviteCode string

type DominoName string

type GameState struct {
	MatchName              string      `json:"match_name"`
	MatchInviteCode        InviteCode  `json:"match_invite_code"`
	Rules                  []string    `json:"rules"`
	Team1UserNames         []string    `json:"team_1"`
	Team2UserNames         []string    `json:"team_2"`
	IsConnected            []bool      `json:"is_connected"`
	CurrentRound           int         `json:"current_round"`
	CurrentStartingBidder  int         `json:"current_starting_bidder"`
	CurrentStartingPlayer  int         `json:"current_starting_player"`
	CurrentIsBidding       bool        `json:"current_is_bidding"`
	CurrentPlayerTurn      int         `json:"current_player_turn"`
	CurrentRoundRules      interface{} `json:"current_round_rules"`
	CurrentTeam1RoundScore int         `json:"current_team_1_round_score"`
	CurrentTeam2RoundScore int         `json:"current_team_2_round_score"`
	CurrentTeam1TotalScore int         `json:"current_team_1_total_score"`
	CurrentTeam2TotalScore int         `json:"current_team_2_total_score"`
	CurrentRoundHistory    []string    `json:"current_round_history"`
	TotalRoundHistory      []string    `json:"total_round_history"`
}

type PlayerGameState struct {
	GameState
	PlayerDominoes []DominoName `json:"player_dominoes"`
}

type GlobalGameState struct {
	GameState
	AllPlayerDominoes []DominoName
	Team1PlayerIDs    []UserID
	Team2PlayerIDs    []UserID
}

type GameMap map[InviteCode]*GlobalGameState
