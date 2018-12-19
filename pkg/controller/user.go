// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 获取坐席
// @Tags 管理后台 API
// @Description 管理员查看坐席
// @Accept  json
// @Produce  json
// @Param uid path int true "坐席id"
// @Success 200 {object} response.UpdateDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/users/{uid} [get]
func GetUser(c *gin.Context) {
	uid := c.Param("uid")
	c.JSON(200, service.GetAgent(uid))
}

// @Summary 获取坐席列表
// @Tags 管理后台 API
// @Description 管理员查看坐席列表
// @Accept  json
// @Produce  json
// @Param page query int true "页数"
// @Param size query int true "每页数量"
// @Param status query string false "坐席状态0/1 正常/冻结"
// @Param start_time query string false "筛选开始时间 2006-01-02T15:04:05"
// @Param stop_time query string false "筛选截止时间 2006-01-02T15:04:05"
// @Param time_field query string false "筛选字段 created_at"
// @Param sort query string false "排序方式 desc asc"
// @Param search query string false "搜索值,只匹配username"
// @Success 200 {object} response.UpdateDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/users [get]
func GetUsers(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	size := c.DefaultQuery("size", "10")
	status := c.Query("status")
	startTime := c.Query("start_time")
	stopTime := c.Query("stop_time")
	sort := c.DefaultQuery("sort", "desc")
	timeFiled := c.DefaultQuery("time_field", "created_at")
	search := c.Query("search")
	c.JSON(200, service.GetUsers(page, size, status, startTime, stopTime, sort, timeFiled, search,"1"))
}

// @Summary 添加坐席
// @Tags 管理后台 API
// @Description 管理员添加坐席
// @Accept  json
// @Produce  json
// @Param body body response.UserArgs true "输入参数"
// @Success 200 {object} response.UpdateDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/users [post]
func CreateUser(c *gin.Context) {
	var param response.UserArgs

	if err := c.ShouldBind(&param); err != nil {
		utils.Log.Errorf("can't bind request body.err:%v",err)
	}
	param.Role = 1
	c.JSON(200, service.CreateUser(param,nil))
}

// @Summary 修改坐席
// @Tags 管理后台 API
// @Description 修改坐席
// @Accept  json
// @Produce  json
// @Param uid path int true "坐席id"
// @Param body body response.UserArgs true "输入参数"
// @Success 200 {object} response.UpdateDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/users/{uid} [put]
func UpdateUser(c *gin.Context) {
	var param response.UserArgs
	uid := c.Param("uid")

	if err := c.ShouldBind(&param); err != nil {
		utils.Log.Errorf("can't bind request body.err:%v",err)
	}
	param.Role = 1
	c.JSON(200, service.UpdateUser(param,uid))
}

// @Summary 修改密码
// @Tags 管理后台 API
// @Description 坐席修改密码
// @Accept  json
// @Produce  json
// @Param uid path int true "坐席id"
// @Param body body response.UserPasswordArgs true "输入参数"
// @Success 200 {object} response.UpdateDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/users/{uid}/password [put]
func UpdateUserPassword(c *gin.Context) {
	var param response.UserPasswordArgs
	uid := c.Param("uid")
	if err := c.ShouldBind(&param); err != nil {
		utils.Log.Errorf("can't bind request body.err:%v",err)
	}

	c.JSON(200, service.UpdateUserPassword(param,uid))
}

// @Summary 重置密码
// @Tags 管理后台 API
// @Description 管理员重置密码
// @Accept  json
// @Produce  json
// @Param uid path int true "坐席id"
// @Param body body response.UserArgs true "输入参数"
// @Success 200 {object} response.UpdateDistributorsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/users/{uid}/password/reset [put]
func ResetUserPassword(c *gin.Context) {
	var param response.UserArgs
	uid := c.Param("uid")
	if err := c.ShouldBind(&param); err != nil {
		utils.Log.Errorf("can't bind request body.err:%v",err)
	}
	param.Role = 1
	c.JSON(200, service.ResetUserPassword(param,uid))
}