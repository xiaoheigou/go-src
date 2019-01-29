// +build !swagger

package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)


// @Summary 查询失败的回调消息
// @Tags 管理后台 API
// @Description 查询失败的回调消息
// @Accept  json
// @Produce  json
// @Param page query int false "页数"
// @Param size query int false "每页数量"
// @Success 200 {object} response.NotifyRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/orders/notify/list [get]
func GetFailNotifyList(c *gin.Context) {
	var ret response.PageResponse
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")
	ret = service.GetNotifyListBySendStatus(page, size)
	c.JSON(200, ret)

}

// @Summary 手动批量发送回调消息
// @Tags 管理后台 API
// @Description 手动批量发送回调消息
// @Accept  json
// @Produce  json
// @Param body body  response.NotifyListReq true "通知id"
// @Success 200 {object} response.NotifyRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/orders/notify/batch/manual [post]
func ManualSendMessage(c *gin.Context) {
	var ret response.CommonRet
	var req response.NotifyListReq

	body, _ := ioutil.ReadAll(c.Request.Body)
	utils.Log.Debugf("the method ManualSendMessage's requestbody is :[%s] ", body)
	err := json.Unmarshal(body, &req)
	if err != nil {
		utils.Log.Error("err,%v", err)
	}
	for _, orderNumber := range req.OrderNumber {
		ret = service.ManualPushMessage(orderNumber)
		if ret.Status == response.StatusFail {
			utils.Log.Errorf("there is something wrong,when send notify by hand ,orderNumber=[%s]", orderNumber)
			break
		}
	}

	c.JSON(200, ret)

}
