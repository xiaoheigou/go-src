package merchant

import "github.com/gin-gonic/gin"


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
		Uid      int `json:"uid" example:123`
		data []Payment `json:"data"`
	}
}

// @Summary 获取承兑商的收款账户信息
// @Tags 承兑商APP API
// @Description 获取承兑商的收款账户信息
// @Accept  json
// @Produce  json
// @Param uid  query  string     true        "用户id"
// @Success 200 {object} merchant.GetPaymentsRet ""
// @Router /merchant/settings/payments [get]
func GetPayments(c *gin.Context) {
	// TODO

	var ret GetPaymentsRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	ret.Entity.data = make([]Payment, 3, 5)
	ret.Entity.data = append(ret.Entity.data,
		Payment{Id: 1, PayType: 2, Name: "sky", BankAccount: "", Bank: "", BankBranch:"", GrCode: "xxyy"})
	c.JSON(200, ret)
}


type AddPaymentArg struct {
	// 承兑商用户id
	Uid int `json:"uid" example:1`
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

type AddPaymentRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      int `json:"uid" example:123`
	}
}

// @Summary 增加承兑商的收款账户信息
// @Tags 承兑商APP API
// @Description 增加承兑商的收款账户信息
// @Accept  json
// @Produce  json
// @Param body body merchant.AddPaymentArg true "输入参数"
// @Success 200 {object} merchant.AddPaymentRet ""
// @Router /merchant/settings/payments [post]
func AddPayment(c *gin.Context) {
	// TODO

	var ret AddPaymentRet
	ret.Status = "success"
	ret.Entity.Uid = 123

	c.JSON(200, ret)
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
		Uid      int `json:"uid" example:123`
	}
}

// @Summary 修改承兑商的收款账户信息
// @Tags 承兑商APP API
// @Description 修改承兑商的收款账户信息
// @Accept  json
// @Produce  json
// @Param body body merchant.SetPaymentArg true "输入参数"
// @Success 200 {object} merchant.SetPaymentRet ""
// @Router /merchant/settings/payments [put]
func SetPayment(c *gin.Context) {
	// TODO

	var ret SetPaymentRet
	ret.Status = "success"
	ret.Entity.Uid = 123

	c.JSON(200, ret)
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
		Uid      int `json:"uid" example:123`
	}
}

// @Summary 删除承兑商的收款账户信息
// @Tags 承兑商APP API
// @Description 删除承兑商的收款账户信息。需要先确定是否有正在进行的订单，如果有不允许删除。
// @Accept  json
// @Produce  json
// @Param body body merchant.DeletePaymentArg true "输入参数"
// @Success 200 {object} merchant.DeletePaymentRet ""
// @Router /merchant/settings/payments [delete]
func DeletePayment(c *gin.Context) {
	// TODO

	var ret DeletePaymentRet
	ret.Status = "success"
	ret.Entity.Uid = 123

	c.JSON(200, ret)
}