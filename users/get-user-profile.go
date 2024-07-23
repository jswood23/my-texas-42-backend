package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
	"my-texas-42-backend/util"
)

type friendRow struct {
	Username string `db:"username"`
}

var userDoesNotExistError = errors.New("user does not exist")

func GetUserProfile(c *gin.Context) {
	username := c.Param("username")

	currentUser, err := util.GetRequestUser(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error getting username from request", "reason": err.Error()})
		return
	}

	if username == currentUser.Username {
		userProfile, err := GetCurrentUserProfile(currentUser.Username)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error getting user profile", "reason": err.Error()})
			return
		}
		c.JSON(200, userProfile)
		return
	} else {
		otherUserProfile, err := GetOtherUserProfile(username, currentUser.UserID, currentUser.Username)
		if err != nil {
			if errors.Is(err, userDoesNotExistError) {
				c.JSON(404, gin.H{"error": err.Error()})
				return
			}
			c.JSON(500, gin.H{"error": "Error getting user profile", "reason": err.Error()})
			return
		}
		c.JSON(200, otherUserProfile)
		return
	}
}

func GetCurrentUserProfile(username string) (*models.UserProfileAPIModel, error) {
	query := sql_scripts.GetUserProfileByUsername(username)
	userResult, err := services.Query[models.UserModel](query)
	if err != nil {
		return nil, err
	}

	if len(userResult) == 0 {
		return nil, userDoesNotExistError
	}
	user := userResult[0]

	// Get list of usernames that the user is friends with
	query = sql_scripts.GetUserFriends(username)
	friendResult, err := services.Query[friendRow](query)
	if err != nil {
		return nil, err
	}
	var friendUsernames = []string{}
	for _, friend := range friendResult {
		friendUsernames = append(friendUsernames, friend.Username)
	}

	// Get list of usernames sending friend requests to the user
	query = sql_scripts.GetIncomingFriendRequests(username)
	friendRequestsResult, err := services.Query[models.FriendRequestModel](query)
	if err != nil {
		return nil, err
	}
	var friendRequests = []string{}
	for _, friendRequest := range friendRequestsResult {
		query = sql_scripts.GetUser(friendRequest.SenderUserID)
		senderResult, err := services.Query[models.UserModel](query)
		if err != nil {
			return nil, err
		}
		if len(senderResult) == 0 {
			return nil, errors.New("friend request sender user not found")
		}
		friendRequests = append(friendRequests, senderResult[0].Username)
	}

	userStats, err := getUserStats(username)
	if err != nil {
		return nil, err
	}

	return &models.UserProfileAPIModel{
		UserID:           user.UserID,
		Username:         user.Username,
		Email:            user.Email,
		DisplayName:      user.DisplayName,
		Friends:          friendUsernames,
		IncomingRequests: friendRequests,
		Stats:            *userStats,
	}, nil
}

func GetOtherUserProfile(username string, currentUserID int, currentUserName string) (*models.OtherUserProfileAPIModel, error) {
	query := sql_scripts.GetUserProfileByUsername(username)
	userResult, err := services.Query[models.UserModel](query)
	if err != nil {
		return nil, err
	}

	if len(userResult) == 0 {
		return nil, userDoesNotExistError
	}
	user := userResult[0]

	// See if the current user has sent a friend request to the user
	query = sql_scripts.GetIncomingFriendRequests(username)
	friendRequestsResult, err := services.Query[models.FriendRequestModel](query)
	if err != nil {
		return nil, err
	}
	isRequestSent := false
	for _, friendRequest := range friendRequestsResult {
		if friendRequest.SenderUserID == currentUserID {
			isRequestSent = true
			break
		}
	}

	// See if the current user is friends with the user
	query = sql_scripts.GetUserFriends(user.Username)
	friendResult, err := services.Query[friendRow](query)
	if err != nil {
		return nil, err
	}
	isFriend := false
	for _, friend := range friendResult {
		if friend.Username == currentUserName {
			isFriend = true
			break
		}
	}

	userstats, err := getUserStats(username)
	if err != nil {
		return nil, err
	}

	return &models.OtherUserProfileAPIModel{
		UserID:        user.UserID,
		Username:      user.Username,
		DisplayName:   user.DisplayName,
		IsFriends:     isFriend,
		IsRequestSent: isRequestSent,
		Stats:         *userstats,
	}, nil
}

func getUserStats(username string) (*models.UserStatsAPIModel, error) {
	query := sql_scripts.GetUserStats(username)
	userStatsResult, err := services.Query[models.UserStatsModel](query)
	if err != nil {
		return nil, err
	}

	if len(userStatsResult) == 0 {
		return nil, errors.New("user stats not found")
	}
	userStats := userStatsResult[0]

	return &models.UserStatsAPIModel{
		GamesPlayed:          userStats.GamesPlayed,
		GamesWon:             userStats.GamesWon,
		RoundsPlayed:         userStats.RoundsPlayed,
		RoundsWon:            userStats.RoundsWon,
		TotalPointsAsBidder:  userStats.TotalPointsAsBidder,
		TotalRoundsAsBidder:  userStats.TotalRoundsAsBidder,
		TotalPointsAsSupport: userStats.TotalPointsAsSupport,
		TotalRoundsAsSupport: userStats.TotalRoundsAsSupport,
		TotalPointsAsCounter: userStats.TotalPointsAsCounter,
		TotalRoundsAsCounter: userStats.TotalRoundsAsCounter,
		TimesWinningBidTotal: userStats.TimesWinningBidTotal,
		TimesCallingSuit:     userStats.TimesCallingSuit,
		TimesCallingNil:      userStats.TimesCallingNil,
		TimesCallingSplash:   userStats.TimesCallingSplash,
		TimesCallingPlunge:   userStats.TimesCallingPlunge,
		TimesCallingSevens:   userStats.TimesCallingSevens,
		TimesCallingDelve:    userStats.TimesCallingDelve,
	}, nil
}
