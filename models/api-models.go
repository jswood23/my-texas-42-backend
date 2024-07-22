package models

type SignupAPIModel struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ConfirmSignupAPIModel struct {
	Username         string `json:"username"`
	VerificationCode string `json:"verificationcode"`
}

type LoginAPIModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
