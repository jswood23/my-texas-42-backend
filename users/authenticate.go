package users

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
)

func Authenticate(c *gin.Context) {
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		c.AbortWithStatus(401)
		return
	}
	authResult, err := services.AuthenticateRequest(authToken)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	query := sql_scripts.GetUserProfileByUsername(*authResult.Username)
	result, err := services.Query[models.UserModel](query)
	if err != nil || len(result) == 0 {
		c.JSON(500, gin.H{"error": "User data was not found."})
		return
	}

	c.Set("user", result[0])
	c.Set("sub", authResult.UserAttributes[0])
	c.Set("emailVerified", authResult.UserAttributes[1])
	c.Set("email", authResult.UserAttributes[2])
}
