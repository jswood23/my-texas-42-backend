package users

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"net/http"
)

func CreateAccount(c *gin.Context) {
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	var request models.CreateAccountAPIModel
	err := dec.Decode(&request)
	if err != nil {
		responseDatabaseError := make(map[string]string)
		responseDatabaseError["message"] = "Invalid request body."

		c.JSON(http.StatusBadRequest, responseDatabaseError)
		return
	}

	println("request:")
	println(request.Username)
	println(request.Email)
	println(request.Password)
}
