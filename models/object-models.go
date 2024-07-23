package models

type UserModel struct {
	UserID      int    `db:"userid"`
	Username    string `db:"username"`
	Email       string `db:"email"`
	IsAdmin     bool   `db:"isadmin"`
	DisplayName string `db:"displayname"`
}

type FriendRequestModel struct {
	FriendRequestID int `db:"friendrequestid"`
	SenderUserID    int `db:"senderuserid"`
	ReceiverUserID  int `db:"receiveruserid"`
}

type FriendModel struct {
	FriendsID int `db:"friendsid"`
	User1ID   int `db:"user1id"`
	User2ID   int `db:"user2id"`
}
