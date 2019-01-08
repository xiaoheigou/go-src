// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/errcode"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

func WebLogin(c *gin.Context) {
	var ret response.EntityResponse
	webLoginRet := []response.WebLoginResponse{{}}
	webLoginRet[0].Uid = 1
	webLoginRet[0].Role = 1
	ret.Status = response.StatusSucc
	ret.Data = webLoginRet

	c.JSON(200, ret)
}

func AppLogin(c *gin.Context) {
	//TODO

	var ret response.LoginRet
	ret.Status = "success"
	data := response.LoginData{
		Uid:        1,
		UserStatus: 1,
		NickName:   "1",
		UserCert:   1,
	}
	ret.Data = []response.LoginData{data}

	c.JSON(200, ret)
}

func Register(c *gin.Context) {
	// TODO

	var ret response.RegisterRet
	ret.Status = "success"
	var data response.RegisterData
	data.Uid = 1
	ret.Data = []response.RegisterData{data}
	c.JSON(200, ret)
}

func GetRandomCode(c *gin.Context) {
	// TODO

	var ret response.SendRandomCodeRet
	var data response.SendRandomCodeData
	data.RandomCodeSeq = 123
	ret.Status = "success"
	ret.Data = []response.SendRandomCodeData{data}
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

func SendRandomCode(c *gin.Context) {
	var json response.SendRandomCodeArg
	if err := c.ShouldBindJSON(&json); err != nil {
		utils.Log.Error(err)
		var retFail response.SendRandomCodeRet
		retFail.Status = response.StatusFail
		retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, retFail)
		return
	}

	c.JSON(200, service.GetRandomCode(json))
	return
}

func OrderComplaint(c *gin.Context) {
	// TODO

	var ret response.OrderComplaintRet
	ret.Status = "success"
	c.JSON(200, ret)
}

func AppLogout(c *gin.Context) {
	// TODO

	var ret response.AppLogoutRet
	ret.Status = "success"
	c.JSON(200, ret)
}

func ChangePw(c *gin.Context) {
	// TODO

	var ret response.ChangePasswordRet
	ret.Status = response.StatusSucc
	c.JSON(200, ret)
}
