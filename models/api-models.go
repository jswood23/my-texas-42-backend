package models

type AuthenticatedUserAttributes struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailverified"`
	Sub           string `json:"sub"`
}

type CurrentUserAPIModel struct {
	Exists      bool                        `json:"exists"`
	Attributes  AuthenticatedUserAttributes `json:"attributes"`
	Username    string                      `json:"username"`
	DisplayName string                      `json:"displayname"`
}

type SignupAPIModel struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ConfirmSignupAPIModel struct {
	Username         string `json:"username"`
	VerificationCode string `json:"verificationcode"`
}

type ChangePasswordAPIModel struct {
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}

type ChangeDisplayNameAPIModel struct {
	NewDisplayName string `json:"displayname"`
}

type LoginAPIModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserStatsAPIModel struct {
	GamesPlayed          int `json:"gamesplayed"`
	GamesWon             int `json:"gameswon"`
	RoundsPlayed         int `json:"roundsplayed"`
	RoundsWon            int `json:"roundswon"`
	TotalPointsAsBidder  int `json:"totalpointsasbidder"`
	TotalRoundsAsBidder  int `json:"totalroundsasbidder"`
	TotalPointsAsSupport int `json:"totalpointsassupport"`
	TotalRoundsAsSupport int `json:"totalroundsassupport"`
	TotalPointsAsCounter int `json:"totalpointsascounter"`
	TotalRoundsAsCounter int `json:"totalroundsascounter"`
	TimesWinningBidTotal int `json:"timeswinningbidtotal"`
	TimesCallingSuit     int `json:"timescallingsuit"`
	TimesCallingNil      int `json:"timescallingnil"`
	TimesCallingSplash   int `json:"timescallingsplash"`
	TimesCallingPlunge   int `json:"timescallingplunge"`
	TimesCallingSevens   int `json:"timescallingsevens"`
	TimesCallingDelve    int `json:"timescallingdelve"`
}

type UserProfileAPIModel struct {
	UserID           UserID            `json:"id"`
	Username         string            `json:"username"`
	Email            string            `json:"email"`
	DisplayName      string            `json:"displayname"`
	Friends          []string          `json:"friends"`
	IncomingRequests []string          `json:"incomingrequests"`
	Stats            UserStatsAPIModel `json:"stats"`
}

type OtherUserProfileAPIModel struct {
	UserID        UserID            `json:"id"`
	Username      string            `json:"username"`
	DisplayName   string            `json:"displayname"`
	IsFriends     bool              `json:"isfriends"`
	IsRequestSent bool              `json:"isrequestsent"`
	Stats         UserStatsAPIModel `json:"stats"`
}

type WSOutgoingMessageAPIModel struct {
	MessageType string           `json:"message_type"`
	Message     string           `json:"message"`
	Username    string           `json:"username"`
	GameData    *PlayerGameState `json:"game_data"`
}

type WSIncomingMessageAPIModel struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

type WSSendChatMessageAPIModel struct {
	Message string `json:"message"`
}

const (
	WSMessageTypeChat       = "chat"
	WSMessageTypeGameUpdate = "game-update"
	WSMessageTypeGameError  = "game-error"
)

type GameAPIModel struct {
	MatchName       string     `json:"match_name"`
	MatchInviteCode InviteCode `json:"match_invite_code"`
	Rules           []string   `json:"rules"`
	Team1           []string   `json:"team_1"`
	Team2           []string   `json:"team_2"`
}

type ListGamesAPIModel struct {
	InGame       InviteCode     `json:"in_game"`
	PublicGames  []GameAPIModel `json:"public_games"`
	PrivateGames []GameAPIModel `json:"private_games"`
}

type NewGameAPIModel struct {
	MatchName string   `json:"match_name"`
	Privacy   string   `json:"privacy"`
	Rules     []string `json:"rules"`
}
