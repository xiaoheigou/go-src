package order

import (
	"YuuPay_core-service/pkg/models"
	"github.com/gin-gonic/gin"
)

type CommonRet struct {
	// status可以为success或者fail
	Status   string `json:"status" binding:"required" example:"success"`
	// err_msg仅在失败时设置
	ErrMsg  string `json:"err_msg" example:"由于xx原因，导致操作失败"`
	// err_code仅在失败时设置
	ErrCode int `json:"err_code" example:1001`
}

type OrdersRet struct {
	CommonRet

	Entity struct {
		Data []models.Order
	}
}

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
// @Success 200 {object} order.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /orders [get]
func GetOrders(c *gin.Context) {

	c.JSON(200, "")
}
