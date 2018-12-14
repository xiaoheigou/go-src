// +build !swagger

package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
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
// @Accept multipart/form-data
// @Produce json
// @Param uid  path  int  true  "用户id"
// @Param pay_type  query  string  true  "0:微信，1:支付宝，2:银行卡"
// @Param name  query  string  true  "收款人姓名"
// @Param amount  query  string  true  "微信或支付宝账号二维码对应的金额，为0时表示不固定金额"
// @Param account  query  string  true  "微信或支付宝账号，或者银行卡卡号"
// @Param bank  query  string  true  "银行名称"
// @Param bank_branch  query  string  true  "银行分行名称"
// @Success 200 {object} response.AddPaymentRet ""
// @Router /m/merchants/{uid}/settings/payments [post]
func AddPayment(c *gin.Context) {
	var uid int
	var err error
	if uid, err = strconv.Atoi(c.Param("uid")); err != nil {
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		var ret response.AddPaymentRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}

	var payType int
	if payType, err = strconv.Atoi(c.Param("pay_type")); err != nil {
		utils.Log.Errorf("pay_type [%v] is invalid, expect a integer", c.Param("pay_type"))
		var ret response.AddPaymentRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}
	if ! (payType == models.PaymentTypeWeixin || payType == models.PaymentTypeAlipay || payType == models.PaymentTypeBanck) {
		utils.Log.Errorf("pay_type [%v] is invalid", payType)
		var ret response.AddPaymentRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}
	name := c.Query("name")
	amount := c.Query("amount")
	account := c.Query("account")
	bank := c.Query("bank")
	bankBranch := c.Query("bank_branch")
	var amountFloat float64
	if amountFloat, err = strconv.ParseFloat(amount, 32); err != nil {
		utils.Log.Errorf("amount [%v] is invalid", amount)
		var ret response.AddPaymentRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	// Upload the file to specific dst.
	// c.SaveUploadedFile(file, dst)
	if err := c.SaveUploadedFile(file, file.Filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	c.JSON(200, service.AddPaymentInfo(uid, payType, name, amountFloat, account, bank, bankBranch))
	return
	//var ret response.AddPaymentRet
	//ret.Status = "success"
	//c.JSON(200, ret)
	//return
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
