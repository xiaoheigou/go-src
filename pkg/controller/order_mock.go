// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
)

func GetOrders(c *gin.Context) {
	var ret response.OrdersRet
	ret.Status = "success"
	ret.Data = []models.Order{
		{
			OrderNumber: 2,
			MerchantId:  1,
			DistributorId: 1,
			Price: 1,
			Amount: 6.666,
		},
	}
	c.JSON(200, ret)
}

func GetOrderByOrderNumber(c *gin.Context) {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var ret response.OrdersRet
	ret.Status = "success"
	ret.Data = []models.Order{
		{
			OrderNumber: id,
			MerchantId:  1,
			DistributorId: 1,
			Price: 1,
			Amount: 6.666,
		},
	}
	c.JSON(200, ret)
}

func GetOrderList (c *gin.Context) {
	var ret response.OrdersRet
	ret.Status=response.StatusSucc
	ret.ErrCode=123
	ret.ErrMsg="get orderList success"
	ret.Data = []models.Order{
		{

			MerchantId:  1,
			DistributorId: 1,
			Price: 1,
			Amount: 6.666,
		},
	}
	c.JSON(200, ret)

}