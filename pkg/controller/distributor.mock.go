// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/service"
)

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

func CreateDistributors(c *gin.Context) {
	// TODO
	var param response.CreateDistributorsArgs
	if err := c.ShouldBind(&param); err != nil {

	}

	c.JSON(200, "")
}

func UpdateDistributors(c *gin.Context) {
	// TODO

	c.JSON(200, "")
}
