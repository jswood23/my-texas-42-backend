package util

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func GetRequestUsername(c *gin.Context) (*string, error) {
	username, exists := c.Get("username")
	if !exists {
		return nil, errors.New("username not found in context")
	}

	usernameStr, _ := username.(*string)

	return usernameStr, nil
}
