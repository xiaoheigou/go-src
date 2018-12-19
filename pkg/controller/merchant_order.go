// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
)

// @Summary 获取承兑商的订单列表
// @Tags 承兑商APP API
// @Description 获取承兑商的订单列表
// @Accept  json
// @Produce  json
// @Param uid  path  int  true  "承兑商用户id"
// @Param direction  query  int  false  "订单类型。0/1表示平台商用户买入/卖出。不传或者传入-1表示全部。"
// @Param status  query  int  false  "订单状态。0/1分别表示：未支付的/已支付的。不传或者传入-1表示全部。"
// @Param page_num  query  int  false  "页号码，从0开始，默认为0"
// @Param page_size  query  int  false  "页大小，默认为10"
// @Success 200 {object} response.GetOrderRet ""
// @Router /m/merchants/{uid}/orders [get]
func GetOrdersByMerchant(c *gin.Context) {
	// TODO

	var ret response.GetOrderRet
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, models.Order{
		OrderNumber:       "10001",
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
		OrderNumber:       "10002",
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
		OrderNumber:       "10003",
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
		OrderNumber:       "10004",
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
		OrderNumber:       "10005",
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
		OrderNumber:       "10006",
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
		OrderNumber:       "10007",
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
		OrderNumber:       "10008",
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
		OrderNumber:       "10009",
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
		OrderNumber:       "10010",
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


// @Summary 承兑商获取某一条订单的详情
// @Tags 承兑商APP API
// @Description 承兑商获取某一条订单的详情
// @Accept  json
// @Produce  json
// @Param uid  path  int  true  "承兑商用户id"
// @Param order-id  path  string  true  "订单id"
// @Success 200 {object} response.GetOrderDetailRet ""
// @Router /m/merchants/{uid}/orders/{order-id} [get]
func GetOrderDetail(c *gin.Context) {
	orderNumber := c.Query("order-id")
	var ret response.GetOrderDetailRet
	ret.Status = response.StatusSucc
	if orderNumber == "10001" {
		ret.Data = append(ret.Data, models.Order{
			OrderNumber:       "10001",
			Price:             4.4,
			Quantity:          "25",
			Amount:            110,
			PaymentRef:        "",
			Status:            0,
			Direction:         1,
			DistributorId:     0,
			MerchantId:        0,
			MerchantPaymentId: 0,
			AccountId:         "牛大春",
			CurrencyCrypto:    "BTUSD",
			CurrencyFiat:      "RMB",
			PayType:           0,
			QrCode:            "https://image.baidu.com/search/down?tn=download&word=download&ie=utf8&fr=detail&url=https%3A%2F%2Ftimgsa.baidu.com%2Ftimg%3Fimage%26quality%3D80%26size%3Db9999_10000%26sec%3D1545215853444%26di%3Dd9eae39078073fa50b6c75cfb4f6b4ce%26imgtype%3D0%26src%3Dhttp%253A%252F%252Fpic.chinaz.com%252F2018%252F0409%252F18040918011940118.jpg&thumburl=https%3A%2F%2Fss1.bdstatic.com%2F70cFvXSh_Q1YnxGkpoWK1HF6hhy%2Fit%2Fu%3D3021223247%2C92123690%26fm%3D26%26gp%3D0.jpg",
			Name:              "",
			BankAccount:       "",
			Bank:              "",
			BankBranch:        "",
			Timestamp:         models.Timestamp{
				CreatedAt: time.Date(2018, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
		})
	} else if orderNumber == "10002" {
		ret.Data = append(ret.Data, models.Order{
			OrderNumber:       "10002",
			Price:             6.5,
			Quantity:          "20",
			Amount:            130,
			PaymentRef:        "",
			Status:            0,
			Direction:         0,
			DistributorId:     0,
			MerchantId:        0,
			MerchantPaymentId: 0,
			AccountId:         "牛二春",
			CurrencyCrypto:    "BTUSD",
			CurrencyFiat:      "RMB",
			PayType:           1,
			QrCode:            "http://13.250.12.109:8086/1_0_1_favicon.png",
			Name:              "",
			BankAccount:       "",
			Bank:              "",
			BankBranch:        "",
			Timestamp:         models.Timestamp{
				CreatedAt: time.Date(2018, 11, 17, 20, 34, 58, 651387237, time.UTC),
			},
		})
	}

	c.JSON(200, ret)
}
