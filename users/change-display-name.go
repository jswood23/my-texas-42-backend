package users

import (
	"github.com/gin-gonic/gin"
	"my-texas-42-backend/models"
	"my-texas-42-backend/services"
	"my-texas-42-backend/sql_scripts"
	"my-texas-42-backend/util"
)

func ChangeDisplayName(c *gin.Context) {
	request, err := models.DecodeAPIModel[models.ChangeDisplayNameAPIModel](c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	username, err := util.GetRequestUsername(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	query := sql_scripts.ChangeDisplayName(request.NewDisplayName, *username)
	err = services.Execute(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(200)
}
