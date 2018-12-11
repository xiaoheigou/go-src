package response

type LoginArg struct {
	Account  string `json:"account" binding:"required" example:"13112345678"`
	Password string `json:"password" binding:"required" example:"pwd123"`
}

type LoginData struct {
	// 用户id
	Uid int `json:"uid" example:123`
	// user_status可以为0/1/2，分别表示“待审核/正常/冻结”
	UserStatus int `json:"user_status" example:0`
	// user_cert可以为0/1，分别表示“未认证/已认证”
	UserCert int `json:"user_cert" example:0`
	// 用户昵称
	NickName string `json:"nickname" example:"老王"`
}

type LoginRet struct {
	CommonRet
	Data []LoginData `json:"data"`
}

type RegisterArg struct {
	Phone    string `json:"phone" binding:"required" example:"13112345678"`
	Email    string `json:"email" binding:"required" example:"xxx@sina.com"`
	Password string `json:"password" binding:"required" example:"pwd1234"`
	// 随机验证码（通过手机发送的）
	PhoneRandomCode string `json:"phone_random_code" binding:"required" example:"9823"`
	// 随机验证码（通过手机发送的）序号
	PhoneRandomCodeSeq int `json:"phone_random_code_seq" binding:"required" example:12`
	// 随机验证码（通过邮件发送的）
	EmailRandomCode string `json:"email_random_code" binding:"required" example:"9823"`
	// 随机验证码（通过邮件发送的）序号
	EmailRandomCodeSeq int `json:"email_random_code_seq" binding:"required" example:13`
}

type RegisterData struct {
	// 用户id
	Uid int `json:"uid" example:123`
}

type RegisterRet struct {
	CommonRet
	Data []RegisterData `json:"data"`
}

type GetRandomCodeData struct {
	// 验证码序号
	RandomCodeSeq int `json:"random_code_seq" example:123`
}

type GetRandomCodeRet struct {
	CommonRet
	Data []GetRandomCodeData `json:"data"`
}

type VerifyRandomCodeArg struct {
	// 手机和邮箱
	Account int `json:"account" example:xxx@sina.com`
	// 随机验证码的内容
	RandomCode string `json:"random_code" binding:"required" example:"H3Q2A"`
	// 随机验证码的序号
	RandomCodeSeq int `json:"random_code_seq" binding:"required" example:12`
}

type VerifyRandomCodeRet struct {
	CommonRet
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
	// 承兑商uid
	Uid int `json:"uid" example:123`
	// 手机收到的随机验证码
	RandomCode string `json:"random_code" binding:"required" example:"9823"`
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
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	//平台角色 0:管理员 1:坐席
	Role int `json:"role"`
}

type WebLoginArgs struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
