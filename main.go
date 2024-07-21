package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"my-texas-42-backend/data_access"
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

	//r.GET("/", auth.Authenticate, testRoot)
	r.POST("/users/new", users.CreateAccount)
	r.GET("/health", getAppHealth)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func getAppHealth(c *gin.Context) {
	err := data_access.CheckDBConnection()

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

func testRoot(c *gin.Context) {
	testObj := make(map[string]any)
	testObj["one"] = 1
	testObj["two"] = 2
	testObj["three"] = "four"

	c.JSON(http.StatusOK, testObj)
}

//func testNewSession(c *gin.Context) {
//	err := auth.NewSession(c)
//	if err != nil {
//		println(err.Error())
//		return
//	}
//
//	testObj := make(map[string]any)
//	testObj["message"] = "New session successfully added."
//
//	c.JSON(http.StatusOK, testObj)
//}
