package users

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/services"
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

	c.Set("authResult", authResult)
}
