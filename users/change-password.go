package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
)

func ChangePassword(c *gin.Context) {
	request, err := models.DecodeAPIModel[models.ChangePasswordAPIModel](c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body.",
			"reason":  err.Error(),
		})
		return
	}

	authToken := c.GetHeader("Authorization")

	err = services.ChangePasswordCognito(authToken, request.OldPassword, request.NewPassword)
	if err != nil {
		c.JSON(400, gin.H{
			"message": fmt.Sprintf("Change password failed. %v", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Password changed successfully.",
	})
}
