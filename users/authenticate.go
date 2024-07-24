package users

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
)

func Authenticate(c *gin.Context) {
	authenticateWithCognito(c, true)
}

func authenticateWithCognito(c *gin.Context, abortIfNotAuthenticated bool) bool {
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		if abortIfNotAuthenticated {
			c.JSON(401, gin.H{"error": "No authorization token provided."})
		}
		return false
	}
	authResult, err := services.AuthenticateRequest(authToken)
	if err != nil {
		if !abortIfNotAuthenticated {
			c.JSON(401, gin.H{"error": "Invalid authorization token."})
		}
		return false
	}

	query := sql_scripts.GetUserProfileByUsername(*authResult.Username)
	result, err := services.Query[models.UserModel](query)
	if err != nil || len(result) == 0 {
		if abortIfNotAuthenticated {
			c.JSON(500, gin.H{"error": "User data was not found."})
		}
		return false
	}

	c.Set("user", result[0])
	c.Set("sub", authResult.UserAttributes[0])
	c.Set("emailVerified", authResult.UserAttributes[1])
	c.Set("email", authResult.UserAttributes[2])

	return true
}
