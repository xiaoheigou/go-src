package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 上传文件
// @Tags 管理后台 API
// @Description 上传文件
// @Accept  json
// @Param distributorId query string false "平台商"
// @Param merchantId query string false "承兑商"
// @Param start_time query string false "筛选开始时间"
// @Param stop_time query string false "筛选截止时间"
// @Param time_field query string false "筛选字段"
// @Produce  json
// @Success 200 {object} response.GetDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/down [get]
func DownFile(c *gin.Context) {

	startTime := c.Query("start_time")
	stopTime := c.Query("stop_time")
	sort := c.DefaultQuery("sort", "desc")
	distributorId := c.Query("distributorId")
	timeFiled := c.DefaultQuery("time_field", "created_at")
	session := sessions.Default(c)
	distributor := session.Get("distributor")
	role := session.Get("userRole")

	status := c.Query("status")
	merchantId := c.Query("merchantId")
	originOrder := strings.TrimSpace(c.Query("origin_order"))
	direction := c.Query("direction")

	distributorIdTemp := distributor.(int64)
	if distributorIdTemp > 0 && role == 2 {
		data, fileName := service.GetOrdersByDistributorAndTimeSlot(utils.TransformTypeToString(distributor), startTime, stopTime, sort, timeFiled)

		c.Header("content-disposition", `attachment; filename=`+fileName)
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

		service.ExportExcel(data, c.Writer)
	}
	//客服可下载excel
	if role == 1 {
		data, fileName := service.GetOrdersExcel(distributorId, startTime, stopTime, sort, timeFiled, status, merchantId, originOrder, direction)

		c.Header("content-disposition", `attachment; filename=`+fileName)
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

		service.ExportExcel(data, c.Writer)
	}
}
