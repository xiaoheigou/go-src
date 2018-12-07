package api

import (
	"github.com/gin-gonic/gin"
	"YuuPay_core-service/pkg/models"
	"YuuPay_core-service/pkg/utils"
)

func GetUser(c *gin.Context) {

	user := models.Merchant{
		NickName:"test1",
	}

	utils.DB.Create(&user)
	c.JSON(200, user)
}
