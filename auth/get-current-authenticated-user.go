package auth

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/request-util"
)

func GetCurrentUser(c *gin.Context) {
	isAuthenticated := authenticateWithAuthToken(c, false)

	if !isAuthenticated {
		response := models.CurrentUserAPIModel{
			Exists: false,
			Attributes: models.AuthenticatedUserAttributes{
				Email:         "",
				EmailVerified: false,
				Sub:           "",
			},
			Username:    "",
			DisplayName: "",
		}

		c.JSON(200, response)
		return
	}

	user, err := request_util.GetRequestUser(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	verified, exists := c.Get("emailVerified")
	if !exists {
		c.JSON(500, gin.H{"error": "emailVerified not found in context"})
		return
	}
	emailVerifiedAttr := verified.(*cognitoidentityprovider.AttributeType)
	emailVerified := *emailVerifiedAttr.Value == "true"

	sub, exists := c.Get("sub")
	if !exists {
		c.JSON(500, gin.H{"error": "sub not found in context"})
		return
	}
	subAttr := sub.(*cognitoidentityprovider.AttributeType)
	subString := *subAttr.Value

	response := models.CurrentUserAPIModel{
		Exists: true,
		Attributes: models.AuthenticatedUserAttributes{
			Email:         user.Email,
			EmailVerified: emailVerified,
			Sub:           subString,
		},
		Username:    user.Username,
		DisplayName: user.DisplayName,
	}

	c.JSON(200, response)
}
