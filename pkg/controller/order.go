// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 获取订单列表
// @Tags 管理后台 API
// @Description 坐席获取订单列表
// @Accept  json
// @Produce  json
// @Param page query int true "页数"
// @Param size query int true "每页数量"
// @Param status query string false "订单状态"
// @Param distributor_id query string false "平台商id"
// @Param merchant_id query string false "承兑商id"
// @Param start_time query string false "筛选开始时间"
// @Param stop_time query string false "筛选截止时间"
// @Param time_field query string false "筛选字段"
// @Param search query string false "搜索值"
// @Success 200 {object} response.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/orders [get]
func GetOrders(c *gin.Context) {
	var ret response.OrdersRet
	ret.Status = "success"
	ret.Data = []models.Order{
		{
			OrderNumber:   2,
			MerchantId:    1,
			DistributorId: 1,
			Price:         1,
			Amount:        6.666,
		},
	}
	c.JSON(200, ret)
}

// @Summary 获取订单
// @Tags 管理后台 API
// @Description 坐席获取订单列表
// @Accept  json
// @Produce  json
// @Param orderNumber path int true "订单id"
// @Success 200 {object} response.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/order/{orderNumber} [get]
func GetOrderByOrderNumber(c *gin.Context) {

	id, error := strconv.ParseInt(c.Param("id"), 10, 64)
	var ret response.OrdersRet

	if error != nil {
		utils.Log.Error(error)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoOrderNumberErr.Data()
		c.JSON(200, ret)
	}

	ret = service.GetOrderByOrderNumber(id)

	c.JSON(200, ret)
}

// @Summary 获取订单列表
// @Tags C端相关 API
// @Description 根据accountId及distributorId获取订单列表
// @Accept  json
// @Produce  json
// @Param page query int false "页数"
// @Param size query int false "每页数量"
// @Param status query string false "订单状态"
// @Param distributor_id query string true "平台商id"
//@Param  account_id query string true "账户id"
// @Param start_time query string false "筛选开始时间"
// @Param stop_time query string false "筛选截止时间"
// @Success 200 {object} response.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/order/list [get]
func GetOrderList(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	size := c.DefaultQuery("size", "10")
	accountId := c.Query("accountId")
	distributorId := c.Query("distributorId")

	var ret response.PageResponse
	ret = service.GetOrderList(page, size, accountId, distributorId)

	c.JSON(200, ret)

}
