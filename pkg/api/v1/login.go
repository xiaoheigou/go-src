package v1

import (
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

}


type LoginArg struct {
	Account     string `json:"account" binding:"required" example:"13112345678"`
	Password string `json:"password" binding:"required" example:"pwd123"`
}

type LoginRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      int `json:"uid" example:123`
		// user_status可以为0/1/2，分别表示“正常/待审核/冻结”
		UserStatus int `json:"user_status" example:0`
		// user_cert可以为0/1，分别表示“已认证/未认证”
		UserCert int `json:"user_cert" example:0`
		// 用户昵称
		NickName string `json:"nickname" example:"老王"`
	}
}

// @Summary 承兑商登录APP
// @Tags 承兑商APP API
// @Description 用户（承兑商）登录系统
// @Accept  json
// @Produce  json
// @Param body body v1.LoginArg true "输入参数"
// @Success 200 {object} v1.LoginRet ""
// @Router /merchant/login [post]
func AppLogin(c *gin.Context) {
	//TODO

	var ret LoginRet
	ret.Status = "success"
	ret.Entity.Uid = 1
	ret.Entity.UserStatus = 0
	ret.Entity.UserCert = 0
	ret.Entity.NickName = "老王"
	c.JSON(200, ret)
}

type RegisterArg struct {
	Phone     string `json:"phone" binding:"required" example:"13112345678"`
	Email string `json:"email" binding:"required" example:"xxx@sina.com"`
	// 图形验证码
	PicCode string `json:"pic_code" binding:"required" example:"E87A"`
	// 随机验证码（通过手机或邮件发送的）
	RandomCode string `json:"random_code" binding:"required" example:"9823"`
	Password string `json:"password" binding:"required" example:"pwd1234"`
}

type RegisterRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      int `json:"uid" example:123`
	}
}

// @Summary 承兑商注册APP
// @Tags 承兑商APP API
// @Description 用户（承兑商）注册
// @Accept  json
// @Produce  json
// @Param body body v1.RegisterArg true "输入参数"
// @Success 200 {object} v1.RegisterRet ""
// @Router /merchant/register [post]
func Register(c *gin.Context) {
	// TODO

	var ret RegisterRet
	ret.Status = "success"
	ret.Entity.Uid = 1
	c.JSON(200, ret)
}

type SendCodeRet struct {
	CommonRet
	Entity struct {
	}
}

// @Summary 获取随机验证码，承兑商“重置密码”时使用
// @Tags 承兑商APP API
// @Description 获取随机验证码，通知短信或者邮件发送。这个API在承兑商“重置密码”时使用。
// @Accept  json
// @Produce  json
// @Param v_code  path  string     true        "图形验证码"
// @Param account  path  string     true        "手机号或者邮箱"
// @Success 200 {object} v1.SendCodeRet ""
// @Router /merchant/sendcode [get]
func SendCode(c *gin.Context) {
	// TODO

	var ret SendCodeRet
	ret.Status = "success"
	c.JSON(200, ret)
}

type ResetPasswordArg struct {
	// 所要重置密码的账号名
	Account     string `json:"account" binding:"required" example:"13112345678"`
	// 图形验证码
	PicCode string `json:"pic_code" binding:"required" example:"E87A"`
	// 随机验证码（通过手机或邮件发送的）
	RandomCode string `json:"random_code" binding:"required" example:"9823"`
	// 所设置的新密码
	Password string `json:"password" binding:"required" example:"pwd1234"`
}

type ResetPasswordRet struct {
	CommonRet
	Entity struct {
	}
}

// @Summary 重置承兑商密码
// @Tags 承兑商APP API
// @Description 重置承兑商密码
// @Accept  json
// @Produce  json
// @Param body body v1.ResetPasswordArg true "输入参数"
// @Success 200 {object} v1.ResetPasswordRet ""
// @Router /merchant/resetpassword [post]
func ResetPw(c *gin.Context) {
	// TODO

	var ret ResetPasswordRet
	ret.Status = "success"
	c.JSON(200, ret)
}


type AppLogoutRet struct {
	CommonRet
	Entity struct {
	}
}


// @Summary 承兑商退出登录
// @Tags 承兑商APP API
// @Description 承兑商退出登录
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Success 200 {object} v1.AppLogoutRet ""
// @Router /merchant/logout [post]
func AppLogout(c *gin.Context) {
	// TODO

	var ret AppLogoutRet
	ret.Status = "success"
	c.JSON(200, ret)
}