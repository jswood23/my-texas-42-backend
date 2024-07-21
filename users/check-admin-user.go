package users

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
)

func CheckAdminUser(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.AbortWithStatus(401)
		return
	}

	usernameStr, _ := username.(*string)

	query := sql_scripts.GetUserProfileByUsername(*usernameStr)

	result, err := services.Query[models.UserModel](query)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error querying database", "reason": err.Error()})
		return
	}

	if len(result) == 0 {
		c.AbortWithStatus(401)
		return
	}

	user := result[0]
	if !user.IsAdmin {
		c.AbortWithStatus(401)
		return
	}
}
