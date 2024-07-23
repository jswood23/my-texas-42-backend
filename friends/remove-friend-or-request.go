package friends

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
	"my-texas-42-backend/util"
)

func RemoveFriendOrRequest(c *gin.Context) {
	request, err := models.DecodeAPIModel[models.RemoveFriendOrRequestAPIModel](c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body.",
			"reason":  err.Error(),
		})
		return
	}

	username, err := util.GetRequestUsername(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Check if the friend request exists
	query := sql_scripts.CheckForExistingFriendRequest(request.SenderUsername, *username)
	friendRequestRows, err := services.Query[models.FriendRequestModel](query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(friendRequestRows) > 0 {
		query = sql_scripts.RemoveFriendRequest(request.SenderUsername, *username)
		err = services.Execute(query)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		return
	}

	// Check if the users are already friends
	query = sql_scripts.CheckForExistingFriend(request.SenderUsername, *username)
	friendRows, err := services.Query[models.FriendModel](query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(friendRows) > 0 {
		query = sql_scripts.RemoveFriend(request.SenderUsername, *username)
		err = services.Execute(query)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		return
	}

	c.JSON(404, gin.H{"error": "Friend request or friend not found."})
}
