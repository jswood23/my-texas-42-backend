package admin

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/sockets"
)

func GetAppStats(c *gin.Context) {
	connectionManager := sockets.GetConnectionManager()
	connectionCount := connectionManager.GetConnectionCount()

	c.JSON(200, gin.H{
		"connectionCount": connectionCount,
	})
}
