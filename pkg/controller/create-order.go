// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
)

// @Summary c端客户下单
// @Tags C端相关 API
// @Description c端客户下单api
// @Accept  json
// @Produce  json
// @Param body body response.CreateOrderRequest true "请求体"
// @Success 200 {object} response.CreateOrderRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/create-order [post]
func CreateOrder(c *gin.Context) {
	var ret response.CreateOrderRet
	var req response.CreateOrderRequest
	c.ShouldBind(&req)
	ret.ErrCode = 0
	ret.Status = "success"
	ret.ErrMsg = "order create ok"
	ret.Data = []models.CreateOrderResult{
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
