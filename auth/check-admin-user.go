package auth

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/logger"
	"my-texas-42-backend/request-util"
)

func CheckAdminUser(c *gin.Context) {
	user, err := request_util.GetRequestUser(c)
	if err != nil {
		logger.Error("Error getting request user: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !user.IsAdmin {
		c.JSON(401, gin.H{"error": "You need admin permissions to view this resource."})
		return
	}
}
