// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
)

// @Summary 获取承兑商的订单列表
// @Tags 承兑商APP API
// @Description 获取承兑商的订单列表
// @Accept  json
// @Produce  json
// @Param uid  path  string  true  "承兑商用户id"
// @Param order_type  query  int  false  "订单类型。0/1/2分别表示：全部/买入/卖出，默认全部"
// @Param order_status  query  int  false  "订单状态。0/1/2分别表示：全部/未支付的/已支付的，默认全部"
// @Param page_num  query  int  false  "页号码，从0开始，默认为0"
// @Param page_size  query  int  false  "页大小，默认为10"
// @Success 200 {object} response.GetOrderRet ""
// @Router /m/merchants/{uid}/orders [get]
func GetOrder(c *gin.Context) {
	// TODO

	var ret response.GetOrderRet
	ret.Status = "success"
	ret.Data = make([]response.MerchantOrder, 1, 1)
	ret.Data[0] = response.MerchantOrder{OrderNum: 1, OrderStatus: 1, OrderType: 1, TotalPrice: "650"}
	c.JSON(200, ret)
}
