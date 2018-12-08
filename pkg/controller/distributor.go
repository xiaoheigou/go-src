// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
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
// @Param time_field query string false "筛选字段 create_at update_at"
// @Param search query string false "搜索值"
// @Success 200 {object} response.GetDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/distributors [get]
func GetDistributors(c *gin.Context) {

	page := c.DefaultQuery("page", "0")
	size := c.DefaultQuery("size", "10")
	status := c.Query("status")
	startTime := c.Query("start_time")
	stopTime := c.Query("stop_time")
	sort := c.DefaultQuery("sort", "desc")
	timeFiled := c.DefaultQuery("time_field", "created_at")
	//search only match distributorId and name
	search := c.Query("search")

	c.JSON(200, service.GetDistributors(page, size, status, startTime, stopTime, sort, timeFiled, search))
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
	var param models.Distributor
	if err := c.ShouldBind(&param); err != nil {
		utils.Log.Debugf("request param is error,%v",err)
	}

	c.JSON(200, service.CreateDistributor(param))
}

// @Summary 修改平台商
// @Tags 管理后台 API
// @Description 坐席修改平台商信息
// @Accept  json
// @Produce  json
// @Param uid path int true "平台商id"
// @Param body body response.UpdateDistributorsArgs true "输入参数"
// @Success 200 {object} response.UpdateDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/distributors/{uid} [put]
func UpdateDistributors(c *gin.Context) {
	// TODO

	c.JSON(200, "")
}

// @Summary 获取承兑商
// @Tags 管理后台 API
// @Description 审核承冻结或者解冻
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Success 200 {object} response.GetDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/distributors/{uid} [get]
func GetDistributor(c *gin.Context) {
	var ret response.GetDistributorsRet

	ret.Status = "success"
	ret.Data = []models.Distributor{{
		Id:1,
		Name:"test",
		Phone:"13112345678",
	}}

	c.JSON(200, ret)
}