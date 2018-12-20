package response

type  LoginArg struct {
	// 手机号，不支持邮箱登录
	Account  string `json:"account" binding:"required" example:"13112345678"`
	// 国家码
	NationCode int `json:"nation_code" example:86`
	// 用户设置的密码
	Password string `json:"password" binding:"required" example:"pwd123"`
}

type LoginData struct {
	// 用户id
	Uid int64 `json:"uid" example:123`
	// user_status可以为0/1/2，分别表示“待审核/正常/冻结”
	UserStatus int `json:"user_status" example:0`
	// user_cert可以为0/1，分别表示“未认证/已认证”
	UserCert int `json:"user_cert" example:0`
	// 用户昵称
	NickName string `json:"nickname" example:"老王"`
	// JWT
	Token string `json:"token"`
	// JWT过期时间
	TokenExpire int64 `json:"token_expire"`
}

type LoginRet struct {
	CommonRet
	Data []LoginData `json:"data"`
}

type RefreshTokenData struct {
	// JWT
	Token string `json:"token"`
	// JWT过期时间
	TokenExpire int64 `json:"token_expire"`
}

type RefreshTokenRet struct {
	CommonRet
	Data []LoginData `json:"data"`
}

type RegisterArg struct {
	Phone    string `json:"phone" binding:"required" example:"13112345678"`
	// 国家码
	NationCode int `json:"nation_code" binding:"exists" example:86`
	Email    string `json:"email" binding:"required" example:"xxx@sina.com"`
	Password string `json:"password" binding:"required" example:"pwd1234"`
	// 随机验证码（通过手机发送的）
	PhoneRandomCode string `json:"phone_random_code" binding:"required" example:"9823"`
	// 随机验证码（通过手机发送的）序号
	PhoneRandomCodeSeq int `json:"phone_random_code_seq" binding:"exists" example:12`
	// 随机验证码（通过邮件发送的）
	EmailRandomCode string `json:"email_random_code" binding:"required" example:"9823"`
	// 随机验证码（通过邮件发送的）序号
	EmailRandomCodeSeq int `json:"email_random_code_seq" binding:"exists" example:13`
}

type RegisterData struct {
	// 用户id
	Uid int64 `json:"uid" example:123`
}

type RegisterRet struct {
	CommonRet
	Data []RegisterData `json:"data"`
}


type SendRandomCodeArg struct {
	// 手机或邮箱
	Account string `json:"account" binding:"required" example:xxx@sina.com`
	// 国家码，当account为手机号时，需要提供
	NationCode int `json:"nation_code" example:86`
	// 获取随机码时指定的purpose，默认为register
	Purpose string `json:"purpose" example:"register"`
}

type SendRandomCodeData struct {
	// 验证码序号
	RandomCodeSeq int `json:"random_code_seq" example:123`
}

type SendRandomCodeRet struct {
	CommonRet
	Data []SendRandomCodeData `json:"data"`
}

type VerifyRandomCodeArg struct {
	// 手机和邮箱
	Account string `json:"account" binding:"required"  example:xxx@sina.com`
	// 国家码，当account为手机号时，需要提供
	NationCode int `json:"nation_code" example:86`
	// 随机验证码的内容
	RandomCode string `json:"random_code" binding:"required" example:"H3Q2A"`
	// 随机验证码的序号
	RandomCodeSeq int `json:"random_code_seq" binding:"exists" example:12`
	// 获取随机码时指定的purpose，默认为register
	Purpose string `json:"purpose" example:"register"`
}

type VerifyRandomCodeRet struct {
	CommonRet
}

type RegisterGeetestArg struct {
	// 手机或邮箱
	Account string `json:"account" binding:"required" example:xxx@sina.com`
	// 国家码，当account为手机号时，需要提供
	NationCode int `json:"nation_code" example:86`
	// 获取随机码时指定的purpose，默认为register
	Purpose string `json:"purpose" example:"register"`
}

type RegisterGeetestData struct {
	GeetestServerStatus string `json:"geetest_server_status"`
	GeetestChallenge string `json:"geetest_challenge"`
}

type RegisterGeetestRet struct {
	CommonRet
	Data []RegisterGeetestData `json:"data"`
}

type VerifyGeetestArg struct {
	// 手机或邮箱
	Account string `json:"account" binding:"required" example:xxx@sina.com`
	// 国家码，当account为手机号时，需要提供
	NationCode int `json:"nation_code" example:86`
	// 获取随机码时指定的purpose，默认为register
	Purpose string `json:"purpose" example:"register"`
	// 调用start-geetest时返回的geetest_server_status
	GeetestServerStatus string `json:"geetest_server_status"`
	// 调用start-geetest时返回的geetest_challenge
	GeetestChallenge string `json:"geetest_challenge"`
	//
	GeetestValidate string `json:"geetest_validate"`
	//
	GeetestSeccode string `json:"geetest_seccode"`
}

type ResetPasswordArg struct {
	// 所要重置密码的账号名
	Account string `json:"account" binding:"required" example:"13112345678"`
	// 随机验证码（通过手机或邮件发送的）
	RandomCode string `json:"random_code" binding:"required" example:"9823"`
	// 所设置的新密码
	Password string `json:"password" binding:"required" example:"pwd1234"`
}

type ResetPasswordRet struct {
	CommonRet
}

type ChangePasswordArg struct {
	// 旧密码
	OldPassword string `json:"old_password" binding:"required" example:"pwd1234"`
	// 所设置的新密码
	NewPassword string `json:"new_password" binding:"required" example:"pwd12345"`
}

type ChangePasswordRet struct {
	CommonRet
}

type AppLogoutArg struct {
	Uid int `json:"uid" example:123`
}

type AppLogoutRet struct {
	CommonRet
	Uid int `json:"uid" example:123`
}

type WebLoginRet struct {
	CommonRet
	Data []WebLoginResponse `json:"data"`
}

type WebLoginResponse struct {
	Uid      int64    `json:"uid"`
	Username string `json:"username"`
	//平台角色 0:管理员 1:坐席
	Role int `json:"role"`
}

type WebLoginArgs struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
