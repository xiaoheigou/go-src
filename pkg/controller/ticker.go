package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 获取工单详情
// @Tags 管理后台 API
// @Description 管理员查看充值申请
// @Accept  json
// @Produce  json
// @Param orderNumber path string true "订单号"
// @Success 200 {object} models.Tickets "成功（status为success）失败（status为fail）都会返回200"
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
// @Success 200 {object} models.TicketUpdate "成功（status为success）失败（status为fail）都会返回200"
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
// @Summary 创建工单
// @Tags C端相关 API
// @Description 创建工单Api
// @Accept  json
// @Produce  json
// @Param body body response.CreateDistributorsArgs true "输入参数"
// @Success 200 {object} response.CommonRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/c/ticket [post]
func CreateTicket(c *gin.Context) {
	var ret response.CommonRet

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		utils.Log.Errorf("err is :[%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
		return
	}
	utils.Log.Debugf("the ticket requestBody is :[%v]", body)

	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	token := utils.Config.Get("tickettoken.token")
	token1 := fmt.Sprintf("%v", token)
	sign := c.Query("signature")
	str := service.SortString(token1, timestamp, nonce, string(body))
	sign1 := service.Sha1(str)

	if sign != sign1 {
		utils.Log.Error("sign is not right,sign=[%v]", sign1)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.IllegalSignErr.Data()
		c.JSON(200, ret)
		return
	}

	ret = service.DealTicket(body)
	c.JSON(200, ret)

}
