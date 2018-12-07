// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
)

// @Summary 获取承兑商的收款账户信息
// @Tags 承兑商APP API
// @Description 获取承兑商的收款账户信息
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Success 200 {object} response.GetPaymentsRet ""
// @Router /m/merchants/{uid}/settings/payments [get]
func GetPayments(c *gin.Context) {
	// TODO

	var ret response.GetPaymentsRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	ret.Entity.Data = make([]response.Payment, 1, 1)
	ret.Entity.Data[0] = response.Payment{Id: 1, PayType: 2, Name: "sky", BankAccount: "", Bank: "", BankBranch: "", GrCode: "xxyy"}
	c.JSON(200, ret)
}

// @Summary 增加承兑商的收款账户信息
// @Tags 承兑商APP API
// @Description 增加承兑商的收款账户信息
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Param body body response.AddPaymentArg true "输入参数"
// @Success 200 {object} response.AddPaymentRet ""
// @Router /m/merchants/{uid}/settings/payments [post]
func AddPayment(c *gin.Context) {
	// TODO

	var ret response.AddPaymentRet
	ret.Status = "success"
	ret.Entity.Uid = 123

	c.JSON(200, ret)
}

// @Summary 修改承兑商的收款账户信息
// @Tags 承兑商APP API
// @Description 修改承兑商的收款账户信息
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Param body body response.SetPaymentArg true "输入参数"
// @Success 200 {object} response.SetPaymentRet ""
// @Router /m/merchants/{uid}/settings/payments [put]
func SetPayment(c *gin.Context) {
	// TODO

	var ret response.SetPaymentRet
	ret.Status = "success"
	ret.Entity.Uid = 123

	c.JSON(200, ret)
}

// @Summary 删除承兑商的收款账户信息
// @Tags 承兑商APP API
// @Description 删除承兑商的收款账户信息。需要先确定是否有正在进行的订单，如果有不允许删除。
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Param body body response.DeletePaymentArg true "输入参数"
// @Success 200 {object} response.DeletePaymentRet ""
// @Router /m/merchants/{uid}/settings/payments [delete]
func DeletePayment(c *gin.Context) {
	// TODO

	var ret response.DeletePaymentRet
	ret.Status = "success"
	ret.Entity.Uid = 123

	c.JSON(200, ret)
}
