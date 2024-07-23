package friends

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
	"my-texas-42-backend/util"
)

func AcceptFriendRequest(c *gin.Context) {
	senderUsername := c.Param("username")

	username, err := util.GetRequestUsername(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Check if the users are already friends
	query := sql_scripts.CheckForExistingFriend(senderUsername, *username)
	friendRows, err := services.Query[models.FriendModel](query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(friendRows) > 0 {
		c.JSON(400, gin.H{"error": "The user is already your friend."})
		return
	}

	// Check if the friend request exists
	query = sql_scripts.CheckForExistingFriendRequest(senderUsername, *username)
	friendRequestRows, err := services.Query[models.FriendRequestModel](query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(friendRequestRows) == 0 {
		c.JSON(400, gin.H{"error": "Friend request does not exist."})
		return
	}

	query = sql_scripts.RemoveFriendRequest(senderUsername, *username)
	err = services.Execute(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	query = sql_scripts.NewFriend(senderUsername, *username)
	err = services.Execute(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(200)
}
