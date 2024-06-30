package models

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ParseAPIRequestToModel(c *gin.Context, v any) error {
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(v)
	if err != nil {
		responseDatabaseError := make(map[string]string)
		responseDatabaseError["message"] = "Invalid request body."

		c.JSON(http.StatusBadRequest, responseDatabaseError)
		return err
	}

	return nil
}
