// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
)

// @Summary 获取平台商列表
// @Tags C端相关 API
// @Description 失败订单再处理api
// @Accept  json
// @Produce  json
// @Param body body response.ReprocessOrderRequest true "请求体"
// @Success 200 {object} response.ReprocessOrderResponse "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/reprocess-order [post]
func ReprocessOrder(c *gin.Context){
	var req response.ReprocessOrderRequest
	c.ShouldBind(&req)
	var ret response.ReprocessOrderResponse
	ret.Status=response.StatusSucc
	ret.ErrCode=123
	ret.ErrMsg="reprecess success"
	ret.Data=[]response.ReprocessOrderEntity{
		{
			Url:          "www.otc.com",
			OrderSuccess: "Notify Order Created",
			TotalCount:   "12",
			OrderNo:      "12332",
			OrderType:    "2",
		},
	}
	c.JSON(200,ret)
}
