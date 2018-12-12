// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
)

func GetOrder(c *gin.Context) {
	// TODO

	var ret response.GetOrderRet
	ret.Status = "success"
	ret.Entity.Data = make([]response.MerchantOrder, 1, 1)
	ret.Entity.Data[0] = response.MerchantOrder{OrderNum: 1, OrderStatus: 1, OrderType: 1, TotalPrice: "650"}
	c.JSON(200, ret)
}
