package auth

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/logger"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
	"regexp"
)

func Authenticate(c *gin.Context) {
	requestPath := c.Request.RequestURI

	pattern := `^/ws\?`
	matched, err := regexp.MatchString(pattern, requestPath)
	if err != nil {
		logger.Error("Error checking request path pattern: " + err.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "Error checking request path pattern"})
		return
	}

	if matched {
		authenticateWithUserSub(c, true)
	} else {
		authenticateWithAuthToken(c, true)
	}
}

func authenticateWithUserSub(c *gin.Context, abortIfNotAuthenticated bool) {
	username := c.Query("username")
	userSub := c.Query("sub")

	if username == "" || userSub == "" {
		if abortIfNotAuthenticated {
			c.AbortWithStatusJSON(401, gin.H{"error": "No username or sub provided."})
		}
		return
	}

	query := sql_scripts.GetUserProfileByUserSub(username, userSub)
	result, err := services.Query[models.UserModel](query)
	if err != nil || len(result) == 0 {
		if abortIfNotAuthenticated {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid username or sub."})
		}
		return
	}

	c.Set("user", result[0])
	c.Set("sub", userSub)
	c.Set("emailVerified", true)
	c.Set("email", result[0].Email)

	return
}

func authenticateWithAuthToken(c *gin.Context, abortIfNotAuthenticated bool) bool {
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		if abortIfNotAuthenticated {
			c.AbortWithStatusJSON(401, gin.H{"error": "No authorization token provided."})
		}
		return false
	}
	authResult, err := services.AuthenticateRequest(authToken)
	if err != nil {
		if abortIfNotAuthenticated {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization token."})
		}
		return false
	}

	var subAttr, emailVerifiedAttr, emailAttr *cognitoidentityprovider.AttributeType
	for _, attr := range authResult.UserAttributes {
		switch *attr.Name {
		case "sub":
			subAttr = attr
		case "email_verified":
			emailVerifiedAttr = attr
		case "email":
			emailAttr = attr
		}
	}

	if subAttr == nil {
		if abortIfNotAuthenticated {
			logger.Error("User sub not found in token")
			c.AbortWithStatusJSON(500, gin.H{"error": "User sub not found in token."})
		}
		return false
	}

	query := sql_scripts.GetUserByUserSub(*subAttr.Value)
	result, err := services.Query[models.UserModel](query)
	if err != nil || len(result) == 0 {
		if abortIfNotAuthenticated {
			errMsg := "User data was not found"
			if err != nil {
				errMsg += ": " + err.Error()
			}
			logger.Error(errMsg)
			c.AbortWithStatusJSON(500, gin.H{"error": "User data was not found."})
		}
		return false
	}

	c.Set("user", result[0])
	c.Set("sub", subAttr)
	c.Set("emailVerified", emailVerifiedAttr)
	c.Set("email", emailAttr)

	return true
}
