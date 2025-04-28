package sockets

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"my-texas-42-backend/games"
	"my-texas-42-backend/logger"
	"my-texas-42-backend/models"
	"my-texas-42-backend/util"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default
		return true
	},
}

func Connect(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to websocket: %v", err)
		c.JSON(500, gin.H{"error": "Failed to upgrade connection to websocket"})
		return
	}
	defer closeConnection(conn)

	user, err := util.GetRequestUser(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	manager.AddConnection(user.Username, conn)
	defer disconnectPlayer(user)

	game, err := addPlayerToGame(c)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if game.IsFull() && game.HasAllPlayersConnected() && !game.HasStarted {
		startGame(game)
	}

	manager.handleIncomingMessages(user.Username)
}

func closeConnection(conn *websocket.Conn) {
	err := conn.Close()
	if err != nil {
		logger.Error("Failed to close ws connection: " + err.Error())
	}
}

func addPlayerToGame(c *gin.Context) (*models.GlobalGameState, error) {
	var matchInviteCode = models.InviteCode(c.Query("match_invite_code"))
	username := c.Query("username")
	teamNumber, err := strconv.Atoi(c.Query("team_number"))
	if err != nil {
		return nil, errors.New("invalid team number")
	}

	if matchInviteCode == "" {
		return nil, errors.New("match invite code is required")
	}

	game := games.GetGameManager().GetGameByInviteCode(matchInviteCode)

	if game.ContainsPlayer(username) {
		messagePlayersInGame(game, models.WSMessageTypeChat, username+" reconnected.")
		game.ConnectDisconnectedPlayer(username)
		return game, nil
	}

	err = game.AddPlayer(username, teamNumber)
	if err != nil {
		return game, err
	}

	messagePlayersInGame(game, models.WSMessageTypeChat, username+" joined the game.")

	return game, nil
}

func disconnectPlayer(user *models.UserModel) {
	manager.RemoveConnection(user.Username)

	game, err := games.GetGameManager().GetGameByUsername(user.Username)
	if err != nil {
		logger.Error("Failed to disconnect player " + user.Username + " from game: " + err.Error())
		return
	}

	game.SetPlayerAsDisconnected(user.Username)

	if game.HasAllPlayersDisconnected() {
		logger.Info("All players have disconnected from game ID " + strconv.Itoa(game.MatchId) + ". Removing game.")
		games.GetGameManager().RemoveGame(game.MatchInviteCode)
	}

	messagePlayersInGame(game, models.WSMessageTypeChat, user.Username+" disconnected.")
}

func startGame(game *models.GlobalGameState) {
	if game.HasStarted {
		return
	}

	game.StartNextRound()
	messagePlayersInGame(game, models.WSMessageTypeGameUpdate, "The game is starting.")
}
