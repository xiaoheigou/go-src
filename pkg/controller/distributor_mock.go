// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
)

func GetDistributors(c *gin.Context) {


	data := []models.Distributor{
		{
			Id:1,
			Name:"test",
		},
	}

	obj := response.GetDistributorsRet{}

	obj.Status = "success"
	obj.ErrCode = 123
	obj.ErrMsg = "test"
	obj.Data = data

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