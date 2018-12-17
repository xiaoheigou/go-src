// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
)

func GetPayments(c *gin.Context) {
	// TODO

	var ret response.GetPaymentsPageRet
	ret.Status = "success"
	c.JSON(200, ret)
}

func AddPayment(c *gin.Context) {
	// TODO

	var ret response.CommonRet
	ret.Status = "success"
	c.JSON(200, ret)
}

func SetPayment(c *gin.Context) {
	// TODO

	var ret response.CommonRet
	ret.Status = "success"
	c.JSON(200, ret)
}

func DeletePayment(c *gin.Context) {
	// TODO

	var ret response.DeletePaymentRet
	ret.Status = "success"
	c.JSON(200, ret)
}
