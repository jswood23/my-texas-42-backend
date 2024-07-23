package friends

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
	"my-texas-42-backend/util"
)

func AddFriend(c *gin.Context) {
	request, err := models.DecodeAPIModel[models.AddFriendAPIModel](c.Request.Body)
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

	// Check if the friend request already exists
	query := sql_scripts.CheckForExistingFriendRequest(*username, request.ReceiverUsername)
	friendRequestRows, err := services.Query[models.FriendRequestModel](query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(friendRequestRows) > 0 {
		c.JSON(400, gin.H{"error": "Friend request already exists."})
		return
	}

	// Check if the users are already friends
	query = sql_scripts.CheckForExistingFriend(*username, request.ReceiverUsername)
	friendRows, err := services.Query[models.FriendModel](query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(friendRows) > 0 {
		c.JSON(400, gin.H{"error": fmt.Sprintf("The user %s is already your friend.", request.ReceiverUsername)})
		return
	}

	query = sql_scripts.NewFriendRequest(*username, request.ReceiverUsername)
	err = services.Execute(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(200)
}
