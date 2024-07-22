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
