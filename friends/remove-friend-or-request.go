package friends

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
	"my-texas-42-backend/util"
)

func RemoveFriendOrRequest(c *gin.Context) {
	otherUserUsername := c.Param("username")

	currentUser, err := util.GetRequestUser(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Check if the users are the same
	if otherUserUsername == currentUser.Username {
		c.JSON(400, gin.H{"error": "You cannot be friends with yourself."})
		return
	}

	// Check if the user exists
	query := sql_scripts.GetUserProfileByUsername(otherUserUsername)
	userRows, err := services.Query[models.UserModel](query)
	if err != nil || len(userRows) == 0 {
		c.JSON(404, gin.H{"error": "User not found."})
		return
	}

	// Check if the friend request exists
	query = sql_scripts.CheckForExistingFriendRequest(otherUserUsername, currentUser.Username)
	friendRequestRows, err := services.Query[models.FriendRequestModel](query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(friendRequestRows) > 0 {
		query = sql_scripts.RemoveFriendRequest(otherUserUsername, currentUser.Username)
		err = services.Execute(query)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Friend request removed."})
		return
	}

	// Check if the users are already friends
	query = sql_scripts.CheckForExistingFriend(otherUserUsername, currentUser.Username)
	friendRows, err := services.Query[models.FriendModel](query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(friendRows) > 0 {
		query = sql_scripts.RemoveFriend(otherUserUsername, currentUser.Username)
		err = services.Execute(query)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Friend removed."})
		return
	}

	c.JSON(404, gin.H{"error": "Friend request or friend not found."})
}
