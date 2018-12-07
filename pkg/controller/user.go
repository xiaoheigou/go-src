package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

func GetUser(c *gin.Context) {

	user := models.Merchant{
		NickName:"test1",
	}

	utils.DB.Create(&user)
	c.JSON(200, user)
}
