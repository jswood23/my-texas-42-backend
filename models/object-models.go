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
	UserSub     string `db:"usersub"`
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

type ConnectionMap map[string]*websocket.Conn // the string is a username

type RoundRules struct {
	Bid         int    `json:"bid"`
	BiddingTeam int    `json:"biddingTeam"`
	Trump       string `json:"trump"`
	Variant     string `json:"variant"`
}

type InviteCode string

type DominoName string

type PrivacyLevel string

const (
	PrivacyPublic  PrivacyLevel = "public"
	PrivacyPrivate PrivacyLevel = "private"
	PrivacyFriends PrivacyLevel = "friends"
)

type GameState struct {
	MatchName              string       `json:"match_name"`
	MatchInviteCode        InviteCode   `json:"match_invite_code"`
	MatchPrivacy           PrivacyLevel `json:"match_privacy"`
	Rules                  []string     `json:"rules"`
	OwnerUsername          string       `json:"owner_username"`
	Team1UserNames         []string     `json:"team_1"`
	Team2UserNames         []string     `json:"team_2"`
	Team1Connected         []bool       `json:"team_1_connected"`
	Team2Connected         []bool       `json:"team_2_connected"`
	CurrentRound           int          `json:"current_round"`
	CurrentStartingBidder  int          `json:"current_starting_bidder"`
	CurrentStartingPlayer  int          `json:"current_starting_player"`
	IsInBidding            bool         `json:"current_is_bidding"`
	CurrentPlayerTurn      int          `json:"current_player_turn"`
	RoundRules             RoundRules   `json:"current_round_rules"`
	Team1RoundScore        int          `json:"current_team_1_round_score"`
	Team2RoundScore        int          `json:"current_team_2_round_score"`
	CurrentTeam1TotalScore int          `json:"current_team_1_total_score"`
	CurrentTeam2TotalScore int          `json:"current_team_2_total_score"`
	RoundHistory           []string     `json:"current_round_history"`
	TotalRoundHistory      []string     `json:"total_round_history"`
}

type PlayerGameState struct {
	GameState
	PlayerDominoes []DominoName `json:"player_dominoes"`
	HasStarted     bool         `json:"has_started"`
}

type GlobalGameState struct {
	GameState
	HasStarted        bool
	AllPlayerDominoes [2][2][]DominoName
	MatchId           int
}

type GameMap map[InviteCode]*GlobalGameState
