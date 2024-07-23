package users

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
)

func Signup(c *gin.Context) {
	request, err := models.DecodeAPIModel[models.SignupAPIModel](c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body.",
			"reason":  err.Error(),
		})
		return
	}

	checkExistingUserRequest := sql_scripts.CheckForExistingUser(request.Username, request.Email)
	result, err := services.Query[models.UserModel](checkExistingUserRequest)
	if err != nil || len(result) > 0 {
		c.JSON(400, gin.H{
			"message": "User already exists.",
		})
		return
	}

	err = services.SignUpCognito(request.Email, request.Username, request.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to create user auth.",
		})
		return
	}

	query := sql_scripts.NewUser(request.Email, request.Username)
	err = services.Execute(query)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to create user data.",
		})
		return
	}

	c.AbortWithStatus(201)
}

func ConfirmSignup(c *gin.Context) {
	request, err := models.DecodeAPIModel[models.ConfirmSignupAPIModel](c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body.",
			"reason":  err.Error(),
		})
		return
	}

	_, err = services.ConfirmSignUpCognito(request.Username, request.VerificationCode)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to confirm user.",
		})
		return
	}

	c.AbortWithStatus(200)
}
