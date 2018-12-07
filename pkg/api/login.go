package api

import (
	"yuudidi.com/pkg/protocol/response"
	"github.com/gin-gonic/gin"
)

// @Summary 登录系统
// @Tags 管理后台 API
// @Description 用户（承兑商）登录系统
// @Accept  json
// @Produce  json
// @Param body body response.LoginArg true "Login argument"
// @Success 200 {object} response.LoginRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /login [post]
func WebLogin(c *gin.Context) {

}

// @Summary 承兑商登录APP
// @Tags 承兑商APP API
// @Description 用户（承兑商）登录系统
// @Accept  json
// @Produce  json
// @Param body body response.LoginArg true "输入参数"
// @Success 200 {object} response.LoginRet ""
// @Router /merchant/login [post]
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

// @Summary 承兑商注册APP
// @Tags 承兑商APP API
// @Description 用户（承兑商）注册
// @Accept  json
// @Produce  json
// @Param body body response.RegisterArg true "输入参数"
// @Success 200 {object} response.RegisterRet ""
// @Router /merchant/register [post]
func Register(c *gin.Context) {
	// TODO

	var ret response.RegisterRet
	ret.Status = "success"
	ret.Entity.Uid = 1
	c.JSON(200, ret)
}

// @Summary 获取随机验证码，承兑商“重置密码”时使用
// @Tags 承兑商APP API
// @Description 获取随机验证码，通知短信或者邮件发送。这个API在承兑商“重置密码”时使用。
// @Accept  json
// @Produce  json
// @Param v_code  query  string     true        "图形验证码"
// @Param account  query  string     true        "手机号或者邮箱"
// @Success 200 {object} response.GetRandomCodeRet ""
// @Router /merchant/randomcode [get]
func GetRandomCode(c *gin.Context) {
	// TODO

	var ret response.GetRandomCodeRet
	ret.Status = "success"
	c.JSON(200, ret)
}

// @Summary 重置承兑商密码
// @Tags 承兑商APP API
// @Description 重置承兑商密码
// @Accept  json
// @Produce  json
// @Param body body response.ResetPasswordArg true "输入参数"
// @Success 200 {object} response.ResetPasswordRet ""
// @Router /merchant/resetpassword [post]
func ResetPw(c *gin.Context) {
	// TODO

	var ret response.ResetPasswordRet
	ret.Status = "success"
	c.JSON(200, ret)
}

// @Summary 承兑商退出登录
// @Tags 承兑商APP API
// @Description 承兑商退出登录
// @Accept  json
// @Produce  json
// @Param body body response.AppLogoutArg true "输入参数"
// @Router /merchant/logout [post]
func AppLogout(c *gin.Context) {
	// TODO

	var ret response.AppLogoutRet
	ret.Status = "success"
	c.JSON(200, ret)
}
