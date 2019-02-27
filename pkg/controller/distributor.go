// +build !swagger

package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 获取平台商列表
// @Tags 管理后台 API
// @Description 坐席获取平台商列表
// @Accept  json
// @Produce  json
// @Param page query int true "页数"
// @Param size query int true "每页数量"
// @Param status query string false "平台商状态 0/1/2 审核/正常/冻结"
// @Param start_time query string false "筛选开始时间 2006-01-02T15:04:05"
// @Param stop_time query string false "筛选截止时间 2006-01-02T15:04:05"
// @Param time_field query string false "筛选字段 create_at update_at"
// @Param sort query string false "排序方式 desc asc"
// @Param search query string false "搜索值"
// @Success 200 {object} response.GetDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/distributors [get]
func GetDistributors(c *gin.Context) {

	page := c.DefaultQuery("page", "1")
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
	var param response.CreateDistributorsArgs
	if err := c.ShouldBind(&param); err != nil {
		utils.Log.Debugf("request param is error,%v", err)
	}

	c.JSON(200, service.CreateDistributor(param))
}

// @Summary 修改平台商
// @Tags 管理后台 API
// @Description 坐席修改平台商信息
// @Accept  multipart/form-data
// @Produce  json
// @Param uid path int true "平台商id"
// @Param body body response.UpdateDistributorsArgs true "输入参数"
// @Success 200 {object} response.UpdateDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/distributors/{uid} [put]
func UpdateDistributors(c *gin.Context) {
	// TODO
	var param response.UpdateDistributorsArgs
	if err := c.ShouldBind(&param); err != nil {
		utils.Log.Debugf("request param is error,%v", err)
	}
	uid := c.Param("uid")

	c.JSON(200, service.UpdateDistributor(param, uid))
}

// @Summary 获取平台商
// @Tags 管理后台 API
// @Description 获取平台商
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Success 200 {object} response.GetDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/distributors/{uid} [get]
func GetDistributor(c *gin.Context) {
	session := sessions.Default(c)
	role := session.Get("userRole")
	if role == 2 {
		distributorId := session.Get("distributor")
		utils.Log.Debugf("distributor get distributor detail,%v", distributorId)
		if distributorId == nil {
			c.JSON(400, "bad request")
			return
		}

		c.JSON(200, service.GetDistributor(utils.TransformTypeToString(distributorId)))
		return
	}
	uid := c.Param("uid")
	c.JSON(200, service.GetDistributor(uid))
}

// @Summary 平台商上传证书
// @Tags 管理后台 API
// @Description 平台商上传证书
// @Accept  multipart/form-data
// @Produce  json
// @Param uid path int true "用户id"
// @Success 200 {object} response.GetDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/distributors/{uid}/upload [post]
func UploadCaPem(c *gin.Context) {
	c.JSON(200, service.UploadPem(c))
}

// @Summary 平台商提现
// @Tags 管理后台 API
// @Description 平台商提现
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Success 200 {object} response.GetDistributorsRet ""
// @Router /w/distributors/{uid}/withdraw [post]
func DistributorWithdraw(c *gin.Context) {
	session := sessions.Default(c)
	role := session.Get("userRole")
	if role == 2 {
		distributorId := session.Get("distributor")
		utils.Log.Debugf("distributor get distributor detail,%v", distributorId)
		if distributorId == nil {
			c.JSON(400, "bad request")
			return
		}

		uid := c.Param("uid")
		utils.Log.Debugf("func DistributorWithdraw, uid = %v", uid)

		// TODO 找对就关系

		//if distributorId != uid {
		//	// sessions中保存的id和这次请求中path中指定的id不匹配
		//	c.JSON(400, "bad request")
		//	return
		//}

		var param response.DistributorWithdrawArgs
		if err := c.ShouldBindJSON(&param); err != nil {
			utils.Log.Debugf("request param is error,%v", err)
			c.JSON(400, "bad request param")
			return
		}

		c.JSON(200, service.DistributorWithdraw(param, distributorId.(string), uid))
		return
	} else {
		// not a distributor user
		c.JSON(400, "bad request")
		return
	}
}
