package users

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
)

func GetUserProfile(c *gin.Context) {
	username := c.Param("username")

	query := sql_scripts.GetUserProfileByUsername(username)

	result, err := services.Query[models.UserModel](query)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error querying database", "reason": err.Error()})
		return
	}

	if len(result) == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	user := result[0]

	c.JSON(200, gin.H{
		"id":           user.UserID,
		"username":     user.Username,
		"email":        user.Email,
		"display-name": user.DisplayName,
	})
}
