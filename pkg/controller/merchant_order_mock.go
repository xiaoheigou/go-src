// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
)

func GetOrder(c *gin.Context) {
	// TODO

	var ret response.GetOrderRet
	ret.Status = "success"
	//ret.Data = make([]response.MerchantOrder, 1, 1)
	//ret.Data[0] = response.MerchantOrder{OrderNum: 1, OrderStatus: 1, OrderType: 1, TotalPrice: "650"}
	c.JSON(200, ret)
}

func GetOrdersByMerchant(c *gin.Context) {
	// TODO

	var ret response.GetOrderRet
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10001,
		Price:             4,
		Quantity:          "",
		Amount:            100,
		PaymentRef:        "",
		Status:            0,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛大春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10002,
		Price:             5,
		Quantity:          "",
		Amount:            200,
		PaymentRef:        "",
		Status:            0,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛二春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10003,
		Price:             5,
		Quantity:          "",
		Amount:            100,
		PaymentRef:        "",
		Status:            0,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛三春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10004,
		Price:             5,
		Quantity:          "",
		Amount:            400,
		PaymentRef:        "",
		Status:            0,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛四春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10005,
		Price:             5,
		Quantity:          "",
		Amount:            500,
		PaymentRef:        "",
		Status:            1,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛五春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10006,
		Price:             5,
		Quantity:          "",
		Amount:            100,
		PaymentRef:        "",
		Status:            1,
		Direction:         1,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛六春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10007,
		Price:             5,
		Quantity:          "",
		Amount:            100,
		PaymentRef:        "",
		Status:            0,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛七春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10008,
		Price:             5,
		Quantity:          "",
		Amount:            100,
		PaymentRef:        "",
		Status:            0,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛八春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10009,
		Price:             5,
		Quantity:          "",
		Amount:            100,
		PaymentRef:        "",
		Status:            0,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛九春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10010,
		Price:             5,
		Quantity:          "",
		Amount:            1000,
		PaymentRef:        "",
		Status:            0,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛十春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp:         models.Timestamp{},
	})
	ret.PageCount = 100
	ret.PageNum = 1
	ret.PageSize = 10
	c.JSON(200, ret)
}

func GetOrderDetail(c *gin.Context) {

	var ret response.GetOrderDetailRet
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       10001,
		Price:             4,
		Quantity:          "",
		Amount:            100,
		PaymentRef:        "",
		Status:            0,
		Direction:         0,
		DistributorId:     0,
		MerchantId:        0,
		MerchantPaymentId: 0,
		AccountId:         "牛大春",
		CurrencyCrypto:    "BTUSD",
		CurrencyFiat:      "",
		PayType:           0,
		QrCode:            "",
		Name:              "",
		BankAccount:       "",
		Bank:              "",
		BankBranch:        "",
		Timestamp: models.Timestamp{
			CreatedAt: time.Date(2018, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	})
	c.JSON(200, ret)
}
