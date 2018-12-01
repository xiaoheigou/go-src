package user

import (
	"github.com/gin-gonic/gin"
	"otc-project/pkg/models"
	"otc-project/pkg/utils"
)

func GetUser(c *gin.Context) {

	user := models.User{
		NickName:"test1",
	}
	user.Asset = models.Assets{
		QtyFrozen:float32(123),
	}

	user.Identify = models.Identities{
		Name:"sky",
	}
	user.Payments = []models.PaymentInfo{
		{
			Account:"111111",
		},
	}

	utils.DB.Create(&user)
	c.JSON(200, user)
}
