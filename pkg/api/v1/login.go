package v1

import (
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

}


type LoginArg struct {
	User     string `json:"user" binding:"required" example:"13112345678"`
	Password string `json:"password" binding:"required" example:"pwd123"`
}

type LoginRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      string `json:"uid" example:"123"`
		// user_status可以为0/1/2，分别表示“正常/待审核/冻结”
		UserStatus int `json:"user_status" example:0`
		// user_cert可以为0/1，分别表示“已认证/未认证”
		UserCert int `json:"user_cert" example:0`
		// nickname是用户昵称
		NickName string `json:"nickname" example:"老王"`
	}
}

// @Summary 登录系统
// @Tags 承兑商APP API
// @Description 用户（承兑商）登录系统
// @Accept  json
// @Produce  json
// @Param body body v1.LoginArg true "Login argument"
// @Success 200 {object} v1.LoginRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /user/login [post]
func AppLogin(c *gin.Context) {
	//TODO

	c.JSON(200, gin.H{
		"token": 123,
		"uid":   123,
	})
}


