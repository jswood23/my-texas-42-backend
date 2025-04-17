package games

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
	"my-texas-42-backend/util"
)

type friendRow struct {
	Username string `db:"username"`
}

func ListGames(c *gin.Context) {
	user, err := util.GetRequestUser(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	respBody := models.ListGamesAPIModel{
		InGame:       "",
		PublicGames:  make([]models.GameAPIModel, 0),
		PrivateGames: make([]models.GameAPIModel, 0),
	}

	userCurrentGame := findUserCurrentGame(user.Username)
	if userCurrentGame != "" {
		respBody.InGame = userCurrentGame
		c.JSON(200, respBody)
		return
	}

	// Get the list of usernames that the user is friends with
	query := sql_scripts.GetUserFriends(user.Username)
	friendResult, err := services.Query[friendRow](query)
	if err != nil {
		c.JSON(500, gin.H{"error": `Error getting user's friends: ` + err.Error()})
		return
	}
	var friendUsernames = []string{}
	for _, friend := range friendResult {
		friendUsernames = append(friendUsernames, friend.Username)
	}

	games := GetGameManager().GetAllGames()

	for _, game := range games {
		if game.HasStarted {
			continue
		}

		if game.GameState.MatchPrivacy == models.PrivacyPublic {
			addToGameList(&respBody.PublicGames, game)
		} else if game.GameState.MatchPrivacy == models.PrivacyFriends && util.SliceContains(friendUsernames, game.GameState.OwnerUsername) {
			addToGameList(&respBody.PrivateGames, game)
		}
	}

	c.JSON(200, respBody)
}

func findUserCurrentGame(username string) models.InviteCode {
	games := GetGameManager().GetAllGames()

	for _, game := range games {
		if util.SliceContains(game.GameState.Team1UserNames, username) {
			return game.GameState.MatchInviteCode
		}

		if util.SliceContains(game.GameState.Team2UserNames, username) {
			return game.GameState.MatchInviteCode
		}
	}

	return ""
}

func addToGameList(gameList *[]models.GameAPIModel, game *models.GlobalGameState) {
	*gameList = append(*gameList, models.GameAPIModel{
		MatchName:       game.GameState.MatchName,
		MatchInviteCode: game.GameState.MatchInviteCode,
		Rules:           game.GameState.Rules,
		Team1:           game.GameState.Team1UserNames,
		Team2:           game.GameState.Team2UserNames,
	})
}
