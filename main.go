package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/admin"
	"my-texas-42-backend/auth"
	"my-texas-42-backend/friends"
	"my-texas-42-backend/games"
	"my-texas-42-backend/logger"
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
		logger.Critical("Failed to start server: " + err.Error())
		return
	}
	logger.Info("Starting server")
	defer logger.Info("Stopping server") // this doesn't work right now

	r := getRouter()

	r.GET("/health", getAppHealth)

	r.GET("/app", auth.Authenticate, auth.CheckAdminUser, admin.GetAppStats)

	r.POST("/users", users.Signup)
	r.PUT("/users/confirm", users.ConfirmSignup)
	r.POST("/users/resend-confirmation", users.ResendConfirmation)
	r.POST("/users/login", users.Login)
	r.GET("/users/current", auth.GetCurrentUser)
	r.PUT("/users/change-password", auth.Authenticate, users.ChangePassword)
	r.PUT("/users/change-display-name", auth.Authenticate, users.ChangeDisplayName)
	r.GET("/users/:username", auth.Authenticate, users.GetUserProfile)

	r.POST("/friends/:username", auth.Authenticate, friends.AddFriend)
	r.POST("/friends/:username/accept", auth.Authenticate, friends.AcceptFriendRequest)
	r.DELETE("/friends/:username", auth.Authenticate, friends.RemoveFriendOrRequest)

	r.GET("/games", auth.Authenticate, games.ListGames)
	r.POST("/games", auth.Authenticate, games.NewGame)

	r.GET("/ws", auth.Authenticate, sockets.Connect)

	err = r.Run(":8080")
	if err != nil {
		logger.Critical("Failed to start server: " + err.Error())
		return
	}
}

func getRouter() *gin.Engine {
	r := gin.Default()

	origins := []string{
		"https://mytexas42.com",
		"https://www.mytexas42.com",
	}

	if system.GetEnv() == "staging" {
		origins = []string{
			"http://localhost:3000",
			"https://staging-app.mytexas42.com",
			"https://www.staging-app.mytexas42.com",
		}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return r
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
