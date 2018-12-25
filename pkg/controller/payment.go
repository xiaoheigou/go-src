// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
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
// @Param uid  path  int  true  "用户id"
// @Param pay_type  query  string  false  "1:微信，2:支付宝，4:银行卡。不提供这个参数，或者为-1时表示查询全部收款账户信息"
// @Param page_size  query  string  false  "分页控制参数，页的大小。默认为10。不能超过50。"
// @Param page_num  query  string  false  "分页控制参数，第多少个页（从1开始）。默认为1"
// @Success 200 {object} response.GetPaymentsPageRet ""
// @Router /m/merchants/{uid}/settings/payments [get]
func GetPayments(c *gin.Context) {
	var uid int
	var err error
	if uid, err = strconv.Atoi(c.Param("uid")); err != nil {
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		var ret response.GetProfileRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}
	c.JSON(200, service.GetPaymentInfo(uid, c))
	return
}

// @Summary 增加承兑商的收款账户信息
// @Tags 承兑商APP API
// @Description 增加承兑商的收款账户信息
// @Accept multipart/form-data
// @Produce json
// @Param uid  path  int  true  "用户id"
// @Param pay_type  query  string  true  "1:微信，2:支付宝，4:银行卡"
// @Param name  query  string  true  "收款人姓名"
// @Param amount  query  string  false  "微信或支付宝账号二维码对应的金额，为0时表示不固定金额"
// @Param account  query  string  true  "微信或支付宝账号，或者银行卡卡号"
// @Param bank  query  string  false  "银行名称"
// @Param bank_branch  query  string  false  "银行分行名称"
// @Param account_default  query  string  false  "是否为默认银行卡，0：不是默认，1：默认"
// @Success 200 {object} response.CommonRet ""
// @Router /m/merchants/{uid}/settings/payments [post]
func AddPayment(c *gin.Context) {
	c.JSON(200, service.AddPaymentInfo(c))
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
// @Param uid  path  int  true  "用户id"
// @Param id  path  int  true "收款账号信息主键"
// @Param pay_type  query  string  true  "1:微信，2:支付宝，4:银行卡"
// @Param name  query  string  true  "收款人姓名"
// @Param amount  query  string  false  "微信或支付宝账号二维码对应的金额，为0时表示不固定金额"
// @Param account  query  string  true  "微信或支付宝账号，或者银行卡卡号"
// @Param bank  query  string  false  "银行名称"
// @Param bank_branch  query  string  false  "银行分行名称"
// @Param account_default  query  string  false  "是否为默认银行卡，0：不是默认，1：默认"
// @Success 200 {object} response.CommonRet ""
// @Router /m/merchants/{uid}/settings/payments/{id} [put]
func SetPayment(c *gin.Context) {
	c.JSON(200, service.UpdatePaymentInfo(c))
	return

	//var ret response.SetPaymentRet
	//ret.Status = "success"
	//c.JSON(200, ret)
}

// @Summary 删除承兑商的收款账户信息
// @Tags 承兑商APP API
// @Description 删除承兑商的收款账户信息。如果找不到相关的收款账户信息，status也会为success。
// @Accept  json
// @Produce  json
// @Param uid  path  int  true "用户id"
// @Param id  path  int  true "收款账号信息主键"
// @Success 200 {object} response.DeletePaymentRet ""
// @Router /m/merchants/{uid}/settings/payments/{id} [delete]
func DeletePayment(c *gin.Context) {
	var uid int
	var err error
	if uid, err = strconv.Atoi(c.Param("uid")); err != nil {
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		var ret response.DeletePaymentRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}

	var paymentId int
	if paymentId, err = strconv.Atoi(c.Param("id")); err != nil {
		utils.Log.Errorf("id [%v] is invalid, expect a integer", c.Param("id"))
		var ret response.DeletePaymentRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}

	c.JSON(200, service.DeletePaymentInfo(uid, paymentId))
	return

	//var ret response.DeletePaymentRet
	//ret.Status = "success"
	//c.JSON(200, ret)
}
