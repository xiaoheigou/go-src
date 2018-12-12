// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
)

func WebLogin(c *gin.Context) {
	var ret response.EntityResponse
	webLoginRet := []response.WebLoginRet{{Uid:"1",Role:0}}
	ret.Status = response.StatusSucc
	ret.Data = webLoginRet

	c.JSON(200,ret)
}

func AppLogin(c *gin.Context) {
	//TODO

	var ret response.LoginRet
	ret.Status = "success"
	ret.Entity.Uid = 1
	ret.Entity.UserStatus = 0
	ret.Entity.UserCert = 0
	ret.Entity.NickName = "老王"
	c.JSON(200, ret)
}

func Register(c *gin.Context) {
	// TODO

	var ret response.RegisterRet
	ret.Status = "success"
	ret.Entity.Uid = 1
	c.JSON(200, ret)
}

func GetRandomCode(c *gin.Context) {
	// TODO

	var ret response.SendRandomCodeRet
	ret.Status = "success"
	ret.Seq = 113456
	c.JSON(200, ret)
}

func VerifyRandomCode(c *gin.Context) {
	// TODO

	var ret response.VerifyRandomCodeRet
	ret.Status = "success"
	c.JSON(200, ret)
}


func ResetPw(c *gin.Context) {
	// TODO

	var ret response.ResetPasswordRet
	ret.Status = "success"
	c.JSON(200, ret)
}

func AppLogout(c *gin.Context) {
	// TODO

	var ret response.AppLogoutRet
	ret.Status = "success"
	c.JSON(200, ret)
}
