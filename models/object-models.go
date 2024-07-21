package models

type UserModel struct {
	UserID      int    `db:"userid"`
	Username    string `db:"username"`
	Email       string `db:"email"`
	IsAdmin     bool   `db:"isadmin"`
	DisplayName string `db:"displayname"`
}
