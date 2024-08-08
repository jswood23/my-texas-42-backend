package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
	"my-texas-42-backend/util"
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

	if util.IsEmailValid(request.Username) {
		username, err := getUsernameByEmail(request.Username)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Login failed.",
				"reason":  err.Error(),
			})
			return
		}
		request.Username = username
	}

	authResult, err := services.LoginCognito(request.Username, request.Password)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Login failed.",
			"reason":  err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Login successful.",
		"token":   authResult.AccessToken,
	})
}

func getUsernameByEmail(email string) (string, error) {
	type resultRow struct {
		Username string `db:"username"`
	}

	query := sql_scripts.GetUsernameByEmail(email)
	result, err := services.Query[resultRow](query)
	if err != nil {
		return "", err
	}

	if len(result) == 0 {
		return "", errors.New("no user found with this email")
	}

	return result[0].Username, nil
}
