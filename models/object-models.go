package models

type UserID int

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
