package response

type LoginArg struct {
	Account  string `json:"account" binding:"required" example:"13112345678"`
	Password string `json:"password" binding:"required" example:"pwd123"`
}

type LoginRet struct {
	CommonRet
	// 用户id
	Uid int `json:"uid" example:123`
	// user_status可以为0/1/2，分别表示“待审核/正常/冻结”
	UserStatus int `json:"user_status" example:0`
	// user_cert可以为0/1，分别表示“未认证/已认证”
	UserCert int `json:"user_cert" example:0`
	// 用户昵称
	NickName string `json:"nickname" example:"老王"`
}

type RegisterArg struct {
	Phone string `json:"phone" binding:"required" example:"13112345678"`
	Email string `json:"email" binding:"required" example:"xxx@sina.com"`
	// 图形验证码
	PicCode string `json:"pic_code" binding:"required" example:"E87A"`
	// 随机验证码（通过手机或邮件发送的）
	RandomCode string `json:"random_code" binding:"required" example:"9823"`
	Password   string `json:"password" binding:"required" example:"pwd1234"`
}

type RegisterRet struct {
	CommonRet
	// 用户id
	Uid int `json:"uid" example:123`
}

type GetRandomCodeRet struct {
	CommonRet
}

type ResetPasswordArg struct {
	// 所要重置密码的账号名
	Account string `json:"account" binding:"required" example:"13112345678"`
	// 图形验证码
	PicCode string `json:"pic_code" binding:"required" example:"E87A"`
	// 随机验证码（通过手机或邮件发送的）
	RandomCode string `json:"random_code" binding:"required" example:"9823"`
	// 所设置的新密码
	Password string `json:"password" binding:"required" example:"pwd1234"`
}

type ResetPasswordRet struct {
	CommonRet
}

type AppLogoutArg struct {
	Uid int `json:"uid" example:123`
}

type AppLogoutRet struct {
	CommonRet
	Uid int `json:"uid" example:123`
}
