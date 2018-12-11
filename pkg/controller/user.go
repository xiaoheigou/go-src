// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

func GetUser(c *gin.Context) {

	user := models.Merchant{
		NickName: "test1",
	}

	utils.DB.Create(&user)
	c.JSON(200, user)
}



// @Summary 添加坐席
// @Tags 管理后台 API
// @Description 管理员添加坐席
// @Accept  json
// @Produce  json
// @Param body body response.UserArgs true "输入参数"
// @Success 200 {object} response.UpdateDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/users [post]
func CreateUser(c *gin.Context) {
	var param response.UserArgs

	if err := c.ShouldBind(&param); err != nil {
		utils.Log.Errorf("can't bind request body.err:%v",err)
	}
	c.JSON(200, service.CreateUser(param))
}
