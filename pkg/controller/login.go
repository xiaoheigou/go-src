// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/service"
)

// @Summary 登录系统
// @Tags 管理后台 API
// @Description 用户（承兑商）登录系统
// @Accept  json
// @Produce  json
// @Param body body response.LoginArg true "Login argument"
// @Success 200 {object} response.LoginRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/login [post]
func WebLogin(c *gin.Context) {

}

// @Summary 承兑商登录APP
// @Tags 承兑商APP API
// @Description 用户（承兑商）登录系统
// @Accept  json
// @Produce  json
// @Param body body response.LoginArg true "输入参数"
// @Success 200 {object} response.LoginRet ""
// @Router /m/merchant/login [post]
func AppLogin(c *gin.Context) {
	//TODO

	var ret response.LoginRet
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.LoginData{
		Uid:        1,
		UserStatus: 0,
		UserCert:   0,
		NickName:   "老王"})
	c.JSON(200, ret)
}

// @Summary 承兑商注册APP
// @Tags 承兑商APP API
// @Description 用户（承兑商）注册
// @Accept  json
// @Produce  json
// @Param body body response.RegisterArg true "输入参数"
// @Success 200 {object} response.RegisterRet ""
// @Router /m/merchant/register [post]
func Register(c *gin.Context) {
	// TODO

	var ret response.RegisterRet
	var json response.RegisterArg
	if err := c.ShouldBindJSON(&json); err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode = 10400
		ret.ErrMsg = "参数输入错误"
		c.JSON(200, ret)
		return
	}

	uid := service.AddMerchant(json.Phone, json.Email)

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.RegisterData{Uid: uid})
	c.JSON(200, ret)
}

// @Summary 获取随机验证码
// @Tags 承兑商APP API
// @Description 获取随机验证码，通知短信或者邮件发送。这个API在承兑商注册用户时使用。
// @Accept  json
// @Produce  json
// @Param account  query  string  true  "手机号或者邮箱"
// @Param purpose  query  string  true  "表明获取随机验证码的用途，注册用户时获取随机码请填register"
// @Success 200 {object} response.GetRandomCodeRet ""
// @Router /m/merchant/random-code [get]
func GetRandomCode(c *gin.Context) {
	// TODO

	var ret response.GetRandomCodeRet
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.GetRandomCodeData{RandomCodeSeq: 123456})
	c.JSON(200, ret)
}


// @Summary 注册时验证身份（验证短信或者邮件发送的随机验证码）
// @Tags 承兑商APP API
// @Description 注册时验证身份（验证短信或者邮件发送的随机验证码）。注册时分为了几个步骤，APP端可以在前面步骤验证通过后再进行下一步操作。
// @Accept  json
// @Produce  json
// @Param body body response.VerifyRandomCodeArg true "输入参数"
// @Success 200 {object} response.VerifyRandomCodeRet ""
// @Router /m/merchant/verify-identity [post]
func VerifyRandomCode(c *gin.Context) {
	// TODO

	var ret response.VerifyRandomCodeRet
	ret.Status = response.StatusSucc
	c.JSON(200, ret)
}


// @Summary 承兑商重置密码
// @Tags 承兑商APP API
// @Description 承兑商重置密码
// @Accept  json
// @Produce  json
// @Param body body response.ResetPasswordArg true "输入参数"
// @Success 200 {object} response.ResetPasswordRet ""
// @Router /m/merchant/reset-password [post]
func ResetPw(c *gin.Context) {
	// TODO

	var ret response.ResetPasswordRet
	ret.Status = response.StatusSucc
	c.JSON(200, ret)
}


// @Summary 承兑商修改密码
// @Tags 承兑商APP API
// @Description 承兑商修改密码，需要发送手机随机码
// @Accept  json
// @Produce  json
// @Param body body response.ChangePasswordArg true "输入参数"
// @Success 200 {object} response.ChangePasswordRet ""
// @Router /m/merchant/change-password [post]
func ChangePw(c *gin.Context) {
	// TODO

	var ret response.ChangePasswordRet
	ret.Status = response.StatusSucc
	c.JSON(200, ret)
}


// @Summary 承兑商退出登录
// @Tags 承兑商APP API
// @Description 承兑商退出登录
// @Accept  json
// @Produce  json
// @Param body body response.AppLogoutArg true "输入参数"
// @Success 200 {object} response.AppLogoutRet ""
// @Router /m/merchant/logout [post]
func AppLogout(c *gin.Context) {
	// TODO

	var ret response.AppLogoutRet
	ret.Status = response.StatusSucc
	c.JSON(200, ret)
}
