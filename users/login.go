package users

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
)

func Login(c *gin.Context) {
	request, err := models.DecodeAPIModel[models.LoginAPIModel](c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body.",
			"reason":  err.Error(),
		})
		return
	}

	err, authResult := services.LoginCognito(request.Username, request.Password)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Login failed.",
			"reason":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successful.",
		"token":   authResult.AccessToken,
	})
}
