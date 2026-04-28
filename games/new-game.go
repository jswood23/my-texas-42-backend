package games

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/request-util"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
	"my-texas-42-backend/util"
)

func NewGame(c *gin.Context) {
	user, err := request_util.GetRequestUser(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Verify that the user is not already in a game
	userCurrentGame := findUserCurrentGame(user.Username)
	if userCurrentGame != "" {
		c.JSON(400, gin.H{"error": "User is already in a game"})
		return
	}

	request, err := models.DecodeAPIModel[models.NewGameAPIModel](c.Request.Body) // todo: match name and privacy are not populating correctly

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	privacyLevelStr, err := util.ValidatePrivacyLevel(request.Privacy)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	privacyLevel := models.PrivacyLevel(privacyLevelStr)

	type matchIdResponse struct {
		MatchId int `db:"matchid"`
	}

	rulesString := ""
	for i, rule := range request.Rules {
		rulesString += rule
		if i < len(request.Rules)-1 {
			rulesString += ", "
		}
	}
	query := sql_scripts.NewMatch(request.MatchName, string(privacyLevel), rulesString, user.Username)
	response, err := services.Query[matchIdResponse](query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(response) == 0 {
		c.JSON(500, gin.H{"error": "Failed to create match"})
		return
	}

	newGame := GetGameManager().CreateNewGame(response[0].MatchId, request.MatchName, privacyLevel, request.Rules, user.Username)

	respBody := models.GameAPIModel{
		MatchName:       newGame.GameState.MatchName,
		MatchInviteCode: newGame.GameState.MatchInviteCode,
		Rules:           newGame.GameState.Rules,
		Team1:           newGame.GameState.Team1UserNames,
		Team2:           newGame.GameState.Team2UserNames,
	}

	c.JSON(200, respBody)
}
