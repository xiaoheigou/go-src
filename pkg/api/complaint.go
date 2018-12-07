package api

import "github.com/gin-gonic/gin"


// @Summary 获取申诉列表
// @Tags 管理后台 API
// @Description 坐席获取申诉列表
// @Accept  json
// @Produce  json
// @Param page query int true "页数"
// @Param size query int true "每页数量"
// @Param status query string false "申诉处理状态"
// @Param start_time query string false "筛选开始时间"
// @Param stop_time query string false "筛选截止时间"
// @Param time_field query string false "筛选字段"
// @Param search query string false "搜索值"
// @Success 200 {object} response.GetComplaintsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /complaints [get]
func GetComplaints(c *gin.Context) {

}

// @Summary 处理申诉
// @Tags 管理后台 API
// @Description 坐席处理申诉
// @Accept  json
// @Produce  json
// @Param id path int true "申诉信息id"
// @Param body body response.HandleComplaintsArgs true "输入参数"
// @Success 200 {object} response.HandleComplaintsRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /complaints/{id} [put]
func HandleComplaints(c *gin.Context) {

}
