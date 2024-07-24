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
	UserID           int               `json:"id"`
	Username         string            `json:"username"`
	Email            string            `json:"email"`
	DisplayName      string            `json:"displayname"`
	Friends          []string          `json:"friends"`
	IncomingRequests []string          `json:"incomingrequests"`
	Stats            UserStatsAPIModel `json:"stats"`
}

type OtherUserProfileAPIModel struct {
	UserID        int               `json:"id"`
	Username      string            `json:"username"`
	DisplayName   string            `json:"displayname"`
	IsFriends     bool              `json:"isfriends"`
	IsRequestSent bool              `json:"isrequestsent"`
	Stats         UserStatsAPIModel `json:"stats"`
}
