package request_util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
)

var userNotFoundError = errors.New("user not found in context")

func GetRequestUser(c *gin.Context) (*models.UserModel, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, userNotFoundError
	}

	userModel, ok := user.(models.UserModel)
	if !ok {
		return nil, userNotFoundError
	}

	return &userModel, nil
}
