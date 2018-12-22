// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
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
	var req response.CreateOrderRequest
	var orderNumber string

	if err := c.ShouldBind(&req); err != nil {
		utils.Log.Debugf("request param is error,%v", err)
	}

	orderNumber = service.PlaceOrder(req)

	//c.Redirect(301, redirectUrl)
	c.JSON(200,orderNumber)

}
