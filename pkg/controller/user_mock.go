// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

func GetUser(c *gin.Context) {

	user := models.Merchant{
		Nickname: "test1",
	}

	utils.DB.Create(&user)
	c.JSON(200, user)
}

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

func CreateUser(c *gin.Context) {
	var user models.User


	c.JSON(200, user)
}