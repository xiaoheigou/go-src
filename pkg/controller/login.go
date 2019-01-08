// +build !swagger

package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 登录系统
// @Tags 管理后台 API
// @Description 用户（承兑商）登录系统
// @Accept  json
// @Produce  json
// @Param body body response.WebLoginArgs true "Login argument"
// @Success 200 {object} response.WebLoginRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/login [post]
func WebLogin(c *gin.Context) {
	var param response.WebLoginArgs

	if err := c.ShouldBind(&param); err != nil {
		utils.Log.Warnf("param is error,err:%v", err)
		ret := response.EntityResponse{}
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
	} else {
		session := sessions.Default(c)
		c.JSON(200, service.Login(param, session))
	}
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
	var json response.LoginArg
	if err := c.ShouldBindJSON(&json); err != nil {
		var retFail response.LoginRet
		retFail.Status = response.StatusFail
		retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, retFail)
		return
	}

	c.JSON(200, service.AppLogin(json))
	return
	//
	//var ret response.LoginRet
	//ret.Status = response.StatusSucc
	//ret.Data = append(ret.Data, response.LoginData{
	//	Uid:        1,
	//	UserStatus: 0,
	//	UserCert:   0,
	//	NickName:   "老王"})
	//c.JSON(200, ret)
}

// @Summary 承兑商获取新的jwt
// @Tags 承兑商APP API
// @Description 承兑商获取新的jwt，jwt的过期时间的固定的，为了保证良好的用户体验，App可以在token过期前申请新token。
// @Accept  json
// @Produce  json
// @Success 200 {object} response.RefreshTokenRet ""
// @Router /m/merchants/{uid}/refresh-token [get]
func RefreshToken(c *gin.Context) {
	var uid int
	var err error
	if uid, err = strconv.Atoi(c.Param("uid")); err != nil {
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		var ret response.RefreshTokenRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}

	c.JSON(200, service.RefreshToken(uid))
	return
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
	var json response.RegisterArg
	if err := c.ShouldBindJSON(&json); err != nil {
		utils.Log.Error(err)
		var retFail response.RegisterRet
		retFail.Status = response.StatusFail
		retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, retFail)
		return
	}

	c.JSON(200, service.AddMerchant(json))
	return
}

// @Summary 获取随机验证码
// @Tags 承兑商APP API
// @Description 获取随机验证码，通知短信或者邮件发送。这个API在承兑商注册用户时使用。
// @Accept  json
// @Produce  json
// @Param body body response.SendRandomCodeArg true "输入参数"
// @Success 200 {object} response.SendRandomCodeRet ""
// @Router /m/merchant/random-code [post]
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

// @Summary 注册时验证身份（验证短信或者邮件发送的随机验证码）
// @Tags 承兑商APP API
// @Description 注册时验证身份（验证短信或者邮件发送的随机验证码）。注册时分为了几个步骤，APP端可以在前面步骤验证通过后再进行下一步操作。
// @Accept  json
// @Produce  json
// @Param body body response.VerifyRandomCodeArg true "输入参数"
// @Success 200 {object} response.VerifyRandomCodeRet ""
// @Router /m/merchant/verify-identity [post]
func VerifyRandomCode(c *gin.Context) {
	var json response.VerifyRandomCodeArg
	if err := c.ShouldBindJSON(&json); err != nil {
		utils.Log.Error(err)
		var retFail response.VerifyRandomCodeRet
		retFail.Status = response.StatusFail
		retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, retFail)
		return
	}

	c.JSON(200, service.VerifyRandomCode(json))
	return
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
// @Param uid  path  int  true  "用户id"
// @Param body body response.ChangePasswordArg true "输入参数"
// @Success 200 {object} response.ChangePasswordRet ""
// @Router /m/merchants/{uid}/change-password [post]
func ChangePw(c *gin.Context) {
	var uid int
	var err error
	if uid, err = strconv.Atoi(c.Param("uid")); err != nil {
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		var ret response.ChangePasswordRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}

	var json response.ChangePasswordArg
	if err := c.ShouldBindJSON(&json); err != nil {
		var retFail response.ChangePasswordRet
		retFail.Status = response.StatusFail
		retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, retFail)
		return
	}

	c.JSON(200, service.ChangeMerchantPassword(uid, json))
	return

	//var ret response.ChangePasswordRet
	//ret.Status = response.StatusSucc
	//c.JSON(200, ret)
}

// @Summary 承兑商退出登录
// @Tags 承兑商APP API
// @Description 承兑商退出登录
// @Accept  json
// @Produce  json
// @Param body body response.AppLogoutArg true "输入参数"
// @Success 200 {object} response.AppLogoutRet ""
// @Router /m/merchants/{uid}/logout [post]
func AppLogout(c *gin.Context) {
	// TODO

	var ret response.AppLogoutRet
	ret.Status = response.StatusSucc
	c.JSON(200, ret)
}
