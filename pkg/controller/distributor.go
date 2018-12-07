package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/service"
)

// @Summary 获取平台商列表
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
// @Success 200 {object} response.GetDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/distributors [get]
func GetDistributors(c *gin.Context) {

	page := c.DefaultQuery("page", "0")
	size := c.DefaultQuery("size", "10")
	status := c.Query("status")
	startTime := c.Query("start_time")
	stopTime := c.Query("stop_time")
	timefield := c.DefaultQuery("time_field", "createAt")
	//search only match distributorId and name
	search := c.Query("search")

	data := service.GetDistributors(page, size, status, startTime, stopTime, timefield, search)

	obj := response.GetDistributorsRet{}

	obj.Status = "success"
	obj.ErrCode = 123
	obj.ErrMsg = "test"
	obj.Entity.Data = data

	c.JSON(200, obj)
}

// @Summary 创建平台商
// @Tags 管理后台 API
// @Description 坐席创建平台商
// @Accept  json
// @Produce  json
// @Param body body response.CreateDistributorsArgs true "输入参数"
// @Success 200 {object} response.CreateDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/distributors [post]
func CreateDistributors(c *gin.Context) {
	// TODO
	var param response.CreateDistributorsArgs
	if err := c.ShouldBind(&param); err != nil {

	}

	c.JSON(200, "")
}

// @Summary 修改平台商
// @Tags 管理后台 API
// @Description 坐席修改平台商信息
// @Accept  json
// @Produce  json
// @Param body body response.UpdateDistributorsArgs true "输入参数"
// @Success 200 {object} response.UpdateDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/distributors [put]
func UpdateDistributors(c *gin.Context) {
	// TODO

	c.JSON(200, "")
}
