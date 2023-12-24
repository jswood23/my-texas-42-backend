package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"my-texas-42-backend/auth"
	"my-texas-42-backend/data_access"
	"net/http"
)

func testRoot(c *gin.Context) {
	testObj := make(map[string]any)
	testObj["one"] = 1
	testObj["two"] = 2
	testObj["three"] = "four"

	c.JSON(http.StatusOK, testObj)
}

func testNewSession(c *gin.Context) {
	err := auth.NewSession(c)
	if err != nil {
		println(err.Error())
		return
	}

	testObj := make(map[string]any)
	testObj["message"] = "New session successfully added."

	c.JSON(http.StatusOK, testObj)
}

func main() {
	err := data_access.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/", auth.Authenticate, testRoot)
	r.POST("/login", testNewSession)

	err = r.Run(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
