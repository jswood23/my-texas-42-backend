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
	GamesPlayed          int `json:"games-played"`
	GamesWon             int `json:"games-won"`
	RoundsPlayed         int `json:"rounds-played"`
	RoundsWon            int `json:"rounds-won"`
	TotalPointsAsBidder  int `json:"total-points-as-bidder"`
	TotalRoundsAsBidder  int `json:"total-rounds-as-bidder"`
	TotalPointsAsSupport int `json:"total-points-as-support"`
	TotalRoundsAsSupport int `json:"total-rounds-as-support"`
	TotalPointsAsCounter int `json:"total-points-as-counter"`
	TotalRoundsAsCounter int `json:"total-rounds-as-counter"`
	TimesWinningBidTotal int `json:"times-winning-bid-total"`
	TimesCallingSuit     int `json:"times-calling-suit"`
	TimesCallingNil      int `json:"times-calling-nil"`
	TimesCallingSplash   int `json:"times-calling-splash"`
	TimesCallingPlunge   int `json:"times-calling-plunge"`
	TimesCallingSevens   int `json:"times-calling-sevens"`
	TimesCallingDelve    int `json:"times-calling-delve"`
}

type UserProfileAPIModel struct {
	UserID           int      `json:"id"`
	Username         string   `json:"username"`
	Email            string   `json:"email"`
	DisplayName      string   `json:"display-name"`
	Friends          []string `json:"friends"`
	IncomingRequests []string `json:"incoming-requests"`
	Stats            UserStatsAPIModel
}

type OtherUserProfileAPIModel struct {
	UserID        int    `json:"id"`
	Username      string `json:"username"`
	DisplayName   string `json:"display-name"`
	IsFriends     bool   `json:"is-friends"`
	IsRequestSent bool   `json:"is-request-sent"`
	Stats         UserStatsAPIModel
}
