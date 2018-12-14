// +build !swagger

package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 承兑商获取资金变动历史
// @Tags 管理后台 API
// @Description 查看资金变动历史
// @Accept  json
// @Produce  json
// @Param uid path int true "承兑商id"
// @Param page query int true "页数"
// @Param size query int true "每页数量"
// @Param start_time query string false "筛选开始时间 2006-01-02T15:04:05"
// @Param stop_time query string false "筛选截止时间 2006-01-02T15:04:05"
// @Param time_field query string false "筛选字段 create_at update_at"
// @Param sort query string false "排序方式 desc asc"
// @Param search query string false "搜索值"
// @Success 200 {object} response.GetMerchantAssetHistoryRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/merchants/{uid}/assets/history [get]
func GetMerchantAssetHistory(c *gin.Context) {
	merchantId := c.Param("uid")
	page := c.DefaultQuery("page", "0")
	size := c.DefaultQuery("size", "10")
	startTime := c.Query("start_time")
	stopTime := c.Query("stop_time")
	sort := c.DefaultQuery("sort", "desc")
	timeFiled := c.DefaultQuery("time_field", "created_at")
	//search only match distributorId and name
	search := c.Query("search")
	c.JSON(200, service.GetAssetHistories(page, size, startTime, stopTime, sort, timeFiled, search, merchantId, true))
}

// @Summary 平台商获取资金变动历史
// @Tags 管理后台 API
// @Description 平台商查看资金变动历史
// @Accept  json
// @Produce  json
// @Param uid path int true "平台商id"
// @Param page query int true "页数"
// @Param size query int true "每页数量"
// @Param start_time query string false "筛选开始时间 2006-01-02T15:04:05"
// @Param stop_time query string false "筛选截止时间 2006-01-02T15:04:05"
// @Param time_field query string false "筛选字段 create_at update_at"
// @Param sort query string false "排序方式 desc asc"
// @Param search query string false "搜索值"
// @Success 200 {object} response.GetMerchantAssetHistoryRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/distributors/{uid}/assets/history [get]
func GetDistributorAssetHistory(c *gin.Context) {
	merchantId := c.Param("uid")
	page := c.DefaultQuery("page", "0")
	size := c.DefaultQuery("size", "10")
	startTime := c.Query("start_time")
	stopTime := c.Query("stop_time")
	sort := c.DefaultQuery("sort", "desc")
	timeFiled := c.DefaultQuery("time_field", "created_at")
	//search only match distributorId and name
	search := c.Query("search")
	c.JSON(200, service.GetAssetHistories(page, size, startTime, stopTime, sort, timeFiled, search, merchantId, false))
}

// @Summary 获取充值申请列表
// @Tags 管理后台 API
// @Description 管理员查看充值申请
// @Accept  json
// @Produce  json
// @Param page query int true "页数"
// @Param size query int true "每页数量"
// @Param status query string false "申请审核状态 0/1 未审核/已审核"
// @Param start_time query string false "筛选开始时间 2006-01-02T15:04:05"
// @Param stop_time query string false "筛选截止时间 2006-01-02T15:04:05"
// @Param time_field query string false "筛选字段 create_at update_at"
// @Param sort query string false "排序方式 desc asc"
// @Param search query string false "搜索值"
// @Success 200 {object} response.GetMerchantAssetHistoryRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/recharge/applies [get]
func GetRechargeApplies(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	size := c.DefaultQuery("size", "10")
	status := c.Query("status")
	startTime := c.Query("start_time")
	stopTime := c.Query("stop_time")
	sort := c.DefaultQuery("sort", "desc")
	timeFiled := c.DefaultQuery("time_field", "created_at")
	//search only match distributorId and name
	search := c.Query("search")
	c.JSON(200, service.GetAssetApplies(page, size, status, startTime, stopTime, sort, timeFiled, search))
}

// @Summary 充值确认
// @Tags 管理后台 API
// @Description 查看资金变动历史
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Param assetId path int true "资产id"
// @Success 200 {object} response.EntityResponse "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/merchants/{uid}/assets/apply/{assetApplyId} [put]
func RechargeConfirm(c *gin.Context) {
	session := sessions.Default(c)
	userId := utils.TransformTypeToString(session.Get("userId"))
	uid := c.Param("uid")
	assetApplyId := c.Param("applyId")

	c.JSON(200,service.RechargeConfirm(uid,assetApplyId,userId))
}

// @Summary 充值申请
// @Tags 管理后台 API
// @Description 给承兑商充值
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Param body body response.RechargeArgs true "充值"
// @Success 200 {object} response.RechargeRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/merchants/{uid}/assets [post]
func Recharge(c *gin.Context) {
	var param response.RechargeArgs
	merchantId := c.Param("uid")
	if err := c.ShouldBind(&param); err != nil {
		utils.Log.Warnf("param is error,err:%v", err)
		ret := response.EntityResponse{}
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
	}
	c.JSON(200, service.RechargeApply(merchantId, param))
}
