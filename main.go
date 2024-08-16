package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"my-texas-42-backend/admin"
	"my-texas-42-backend/auth"
	"my-texas-42-backend/friends"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sockets"
	"my-texas-42-backend/system"
	"my-texas-42-backend/users"
	"net/http"
	"os"
	"time"
)

func main() {
	err := system.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	allowCors(r)

	r.GET("/health", getAppHealth)

	r.GET("/app", auth.Authenticate, auth.CheckAdminUser, admin.GetAppStats)

	r.POST("/users", users.Signup)
	r.PUT("/users/confirm", users.ConfirmSignup)
	r.POST("/users/login", users.Login)
	r.GET("/users/current", auth.GetCurrentUser)
	r.PUT("/users/change-password", auth.Authenticate, users.ChangePassword)
	r.PUT("/users/change-display-name", auth.Authenticate, users.ChangeDisplayName)
	r.GET("/users/:username", auth.Authenticate, users.GetUserProfile)

	r.POST("/friends/:username", auth.Authenticate, friends.AddFriend)
	r.POST("/friends/:username/accept", auth.Authenticate, friends.AcceptFriendRequest)
	r.DELETE("/friends/:username", auth.Authenticate, friends.RemoveFriendOrRequest)

	r.GET("/ws", auth.Authenticate, sockets.Connect)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func allowCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://mytexas42.com", "https://www.mytexas42.com", "https://staging-app.mytexas42.com", "https://www.staging-app.mytexas42.com"}, // Replace with your allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
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
