package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"my-texas-42-backend/friends"
	"my-texas-42-backend/services"
	"my-texas-42-backend/system"
	"my-texas-42-backend/users"
	"net/http"
	"os"
)

func main() {
	err := system.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/health", getAppHealth)

	r.POST("/users", users.Signup)
	r.PUT("/users/confirm", users.ConfirmSignup)
	r.POST("/users/login", users.Login)
	r.PUT("/users/change-password", users.Authenticate, users.ChangePassword)
	r.PUT("/users/change-display-name", users.Authenticate, users.ChangeDisplayName)
	r.GET("/users/:username", users.Authenticate, users.GetUserProfile)

	r.POST("/friends", users.Authenticate, friends.AddFriend)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func getAppHealth(c *gin.Context) {
	err := services.CheckDBConnection()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"environment": os.Getenv("ENVIRONMENT"),
			"status":      "unhealthy",
			"reason":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"environment": os.Getenv("ENVIRONMENT"),
		"status":      "great",
	})
}
