// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
)

func CreateOrder(c *gin.Context) {
	var ret response.CreateOrderRet
	var req response.CreateOrderRequest
	c.ShouldBind(&req)
	ret.ErrCode = 0
	ret.Status = "success"
	ret.ErrMsg = "order create ok"
	ret.Data = []response.CreateOrderResult{
		{
			Url:          "www.otc.com",
			OrderSuccess: "Notify Order Created",
			TotalCount:   "12",
			OrderNo:      "12332",
			OrderType:    "2",
		},
	}

	c.JSON(200, ret)

}
