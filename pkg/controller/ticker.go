package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/service"
)

// @Summary 获取工单详情
// @Tags 管理后台 API
// @Description 管理员查看充值申请
// @Accept  json
// @Produce  json
// @Param orderNumber path string true "订单号"
// @Success 200 {object} response.GetMerchantAssetHistoryRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/orders/ticket/{orderNumber} [get]
func GetTicket(c *gin.Context) {
	orderNumber := c.Param("orderNumber")
	c.JSON(200, service.GetTicket(orderNumber))
}

// @Summary 获取工单变动详情
// @Tags 管理后台 API
// @Description 查看工单变动详情
// @Accept  json
// @Produce  json
// @Param page query int true "页数"
// @Param size query int true "每页数量"
// @Param start_time query string false "筛选开始时间 2006-01-02T15:04:05"
// @Param stop_time query string false "筛选截止时间 2006-01-02T15:04:05"
// @Param time_field query string false "筛选字段 create_at update_at"
// @Param sort query string false "排序方式 desc asc"
// @Param search query string false "搜索值"
// @Param ticketId path int true "工单id"
// @Success 200 {object} response.EntityResponse "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/tickets/{ticketId} [get]
func GetTicketUpdates(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")
	startTime := c.Query("start_time")
	stopTime := c.Query("stop_time")
	sort := c.DefaultQuery("sort", "desc")
	timeFiled := c.DefaultQuery("time_field", "created_at")
	search := c.Query("search")
	ticketId := c.Param("ticketId")
	c.JSON(200, service.GetTicketUpdates(page, size, startTime, stopTime, sort, timeFiled, search, ticketId))
}
