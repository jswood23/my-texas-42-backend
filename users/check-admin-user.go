package users

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/util"
)

func CheckAdminUser(c *gin.Context) {
	user, err := util.GetRequestUser(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !user.IsAdmin {
		c.JSON(401, gin.H{"error": "You need admin permissions to view this resource."})
		return
	}
}
