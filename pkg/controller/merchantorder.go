// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/models"
)

// @Summary 获取承兑商的订单列表
// @Tags 承兑商APP API
// @Description 获取承兑商的订单列表
// @Accept  json
// @Produce  json
// @Param uid  path  string  true  "承兑商用户id"
// @Param direction  query  int  false  "订单类型。0/1表示平台商用户买入/卖出，默认为全部。"
// @Param status  query  int  false  "订单状态。0/1分别表示：未支付的/已支付的，默认为全部。"
// @Param page_num  query  int  false  "页号码，从0开始，默认为0"
// @Param page_size  query  int  false  "页大小，默认为10"
// @Success 200 {object} response.GetOrderRet ""
// @Router /m/merchants/{uid}/orders [get]
func GetOrdersByMerchant(c *gin.Context) {
	// TODO

	var ret response.GetOrderRet
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10001,
		Price:             4,
		Quantity:          "",
		Amount:            100,
		PaymentRef:        "",
		Status:            0,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛大春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10001,
		Price:             4,
		Quantity:          "",
		Amount:            100,
		PaymentRef:        "",
		Status:            0,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛二春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	ret.PageCount = 100
	ret.PageNum = 1
	ret.PageSize = 10
	c.JSON(200, ret)
}


// @Summary 承兑商获取某一条订单的详情
// @Tags 承兑商APP API
// @Description 承兑商获取某一条订单的详情
// @Accept  json
// @Produce  json
// @Param uid  path  string  true  "承兑商用户id"
// @Param order-id  path  string  true  "订单id"
// @Success 200 {object} response.GetOrderDetailRet ""
// @Router /m/merchants/{uid}/orders/{order-id} [get]
func GetOrderDetail(c *gin.Context) {

	var ret response.GetOrderDetailRet
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10001,
		Price:             4,
		Quantity:          "",
		Amount:            100,
		PaymentRef:        "",
		Status:            0,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛大春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	c.JSON(200, ret)
}
