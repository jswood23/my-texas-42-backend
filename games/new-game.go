package games

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/util"
)

func NewGame(c *gin.Context) {
	user, err := util.GetRequestUser(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Verify that the user is not already in a game
	userCurrentGame := findUserCurrentGame(user.Username)
	println("User current game: " + userCurrentGame)
	if userCurrentGame != "" {
		c.JSON(400, gin.H{"error": "User is already in a game"})
		return
	}

	request, err := models.DecodeAPIModel[models.NewGameAPIModel](c.Request.Body)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	privacyLevel, err := util.ValidatePrivacyLevel(request.Privacy)

	newGame := GetGameManager().CreateNewGame(request.MatchName, privacyLevel, request.Rules, user.Username)

	respBody := models.GameAPIModel{
		MatchName:       newGame.GameState.MatchName,
		MatchInviteCode: newGame.GameState.MatchInviteCode,
		Rules:           newGame.GameState.Rules,
		Team1:           newGame.GameState.Team1UserNames,
		Team2:           newGame.GameState.Team2UserNames,
	}

	c.JSON(200, respBody)
}
