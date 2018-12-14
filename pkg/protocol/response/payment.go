package response

type Payment struct {
	// 主键，无实际意义
	Id int `json:"id" example:1`
	// 支付类型 1:银行卡 2:微信 3:支付宝
	PayType int `json:"pay_type" example:3`
	// 账户名
	Name string `json:"name" example:"sky"`
	// 银行卡账号
	BankAccount string `json:"bank_account" example:"88888888"`
	// 所属银行
	Bank string `json:"bank" example:"工商银行"`
	// 所属银行分行
	BankBranch string `json:"bank_branch" example:"工商银行日本分行"`
	// 二维码信息
	GrCode string `json:"gr_code" example:"xxxx"`
}

type GetPaymentsRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid  int       `json:"uid" example:123`
		Data []Payment `json:"data"`
	}
}

type AddPaymentRet struct {
	CommonRet
}

type SetPaymentArg struct {
	// 收款账号信息主键
	Id int `json:"id" example:1`
	// 支付类型 1:银行卡 2:微信 3:支付宝
	PayType int `json:"pay_type" example:3`
	// 账户名
	Name string `json:"name" example:"sky"`
	// 银行卡账号
	BankAccount string `json:"bank_account" example:"88888888"`
	// 所属银行
	Bank string `json:"bank" example:"工商银行"`
	// 所属银行分行
	BankBranch string `json:"bank_branch" example:"工商银行日本分行"`
	// 二维码信息
	GrCode string `json:"gr_code" example:"xxxx"`
}

type SetPaymentRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid int `json:"uid" example:123`
	}
}

type DeletePaymentArg struct {
	// 承兑商用户id
	Uid int `json:"uid" example:1`
	// 收款账号信息主键
	Id int `json:"id" example:1`
}

type DeletePaymentRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid int `json:"uid" example:123`
	}
}
