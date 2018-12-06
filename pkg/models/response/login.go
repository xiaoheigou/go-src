package response

type LoginArg struct {
	Account  string `json:"account" binding:"required" example:"13112345678"`
	Password string `json:"password" binding:"required" example:"pwd123"`
}

type LoginRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid int `json:"uid" example:123`
		// user_status可以为0/1/2，分别表示“正常/待审核/冻结”
		UserStatus int `json:"user_status" example:0`
		// user_cert可以为0/1，分别表示“已认证/未认证”
		UserCert int `json:"user_cert" example:0`
		// 用户昵称
		NickName string `json:"nickname" example:"老王"`
	}
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
	Entity struct {
		// 用户id
		Uid int `json:"uid" example:123`
	}
}

type GetRandomCodeRet struct {
	CommonRet
	Entity struct {
	}
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
	Entity struct {
	}
}

type AppLogoutArg struct {
	Uid int `json:"uid" example:123`
}

type AppLogoutRet struct {
	CommonRet
	Entity struct {
	}
}
