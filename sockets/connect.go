package sockets

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"my-texas-42-backend/util"
	"net/http"
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
	defer conn.Close()

	user, err := util.GetRequestUser(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	manager.AddConnection(user.UserID, conn)
	defer manager.RemoveConnection(user.UserID)

	manager.HandleIncomingMessages(user.UserID)
}
