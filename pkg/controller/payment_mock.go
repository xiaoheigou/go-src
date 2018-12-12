// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
)

func GetPayments(c *gin.Context) {
	// TODO

	var ret response.GetPaymentsRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	ret.Entity.Data = make([]response.Payment, 1, 1)
	ret.Entity.Data[0] = response.Payment{Id: 1, PayType: 2, Name: "sky", BankAccount: "", Bank: "", BankBranch: "", GrCode: "xxyy"}
	c.JSON(200, ret)
}

func AddPayment(c *gin.Context) {
	// TODO

	var ret response.AddPaymentRet
	ret.Status = "success"
	ret.Entity.Uid = 123

	c.JSON(200, ret)
}

func SetPayment(c *gin.Context) {
	// TODO

	var ret response.SetPaymentRet
	ret.Status = "success"
	ret.Entity.Uid = 123

	c.JSON(200, ret)
}

func DeletePayment(c *gin.Context) {
	// TODO

	var ret response.DeletePaymentRet
	ret.Status = "success"
	ret.Entity.Uid = 123

	c.JSON(200, ret)
}
