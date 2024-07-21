package models

type CreateAccountAPIModel struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginAPIModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
