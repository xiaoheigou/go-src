// +build !swagger

package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 获取承兑商的订单列表
// @Tags 承兑商APP API
// @Description 获取承兑商的订单列表
// @Accept  json
// @Produce  json
// @Param uid  path  int  true  "承兑商用户id"
// @Param direction  query  int  false  "订单类型。0/1表示平台商用户买入/卖出。不传或者传入-1表示全部。"
// @Param in_progress  query  int  false  "订单是否在进行中。0表示已经结束的订单，1表示进行中的订单。不传或者传入-1表示全部。"
// @Param page_num  query  int  false  "页号码，从1开始，默认为1"
// @Param page_size  query  int  false  "页大小，默认为10。不能超过50。"
// @Success 200 {object} response.OrdersRet ""
// @Router /m/merchants/{uid}/orders [get]
func GetOrdersByMerchant(c *gin.Context) {
	var uid int
	var err error
	if uid, err = strconv.Atoi(c.Param("uid")); err != nil {
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		var ret response.OrdersRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}

	var pageNum int
	if c.Query("page_num") == "" {
		utils.Log.Infof("page_num is missing, use 1 as default")
		pageNum = 1
	} else if pageNum, err = strconv.Atoi(c.Query("page_num")); err != nil {
		utils.Log.Errorf("page_num [%v] is invalid, expect a integer", c.Query("page_num"))
		var ret response.OrdersRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}

	var pageSize int
	if c.Query("page_size") == "" {
		utils.Log.Infof("page_size is missing, use 10 as default")
		pageSize = 10
	} else if pageSize, err = strconv.Atoi(c.Query("page_size")); err != nil {
		utils.Log.Errorf("page_size [%v] is invalid, expect a integer", c.Query("page_size"))
		var ret response.OrdersRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}
	if pageSize > 50 {
		utils.Log.Errorf("page_size [%v] is too large, must <= 50", pageSize)
		var ret response.OrdersRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrPageSizeTooLarge.Data()
		c.JSON(200, ret)
		return
	}

	directionStr := c.Query("direction")
	var direction int
	if directionStr == "" {
		utils.Log.Infof("direction is missing, use -1 as default")
		direction = -1
	} else if directionStr == "0" || directionStr == "1" || directionStr == "-1" {
		direction, _ = strconv.Atoi(directionStr)
	} else {
		utils.Log.Errorf("direction [%v] is invalid, expect 0/1/-1", directionStr)
		var ret response.OrdersRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}

	inProgressStr := c.Query("in_progress")
	var inProgress int
	if inProgressStr == "" {
		utils.Log.Infof("in_progress is missing, use -1 as default")
		inProgress = -1
	} else if inProgressStr == "0" || inProgressStr == "1" || inProgressStr == "-1" {
		inProgress, _ = strconv.Atoi(inProgressStr)
	} else {
		utils.Log.Errorf("in_progress [%v] is invalid, expect 0/1/-1", directionStr)
		var ret response.OrdersRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}

	c.JSON(200, service.GetOrdersByMerchant(pageNum, pageSize, direction, inProgress, int64(uid)))
	return

	//var ret response.GetOrderRet
	//ret.Status = response.StatusSucc
	//ret.Data = append(ret.Data, models.Order{
	//	OrderNumber:       "10001",
	//	Price:             4,
	//	Quantity:          "",
	//	Amount:            100,
	//	PaymentRef:        "",
	//	Status:            0,
	//	Direction:         0,
	//	DistributorId:     0,
	//	MerchantId:        0,
	//	MerchantPaymentId: 0,
	//	AccountId:         "牛大春",
	//	CurrencyCrypto:    "BTUSD",
	//	CurrencyFiat:      "",
	//	PayType:           0,
	//	QrCode:            "",
	//	Name:              "",
	//	BankAccount:       "",
	//	Bank:              "",
	//	BankBranch:        "",
	//	Timestamp:         models.Timestamp{},
	//})
	//ret.Data = append(ret.Data, models.Order{
	//	OrderNumber:       "10002",
	//	Price:             5,
	//	Quantity:          "",
	//	Amount:            200,
	//	PaymentRef:        "",
	//	Status:            0,
	//	Direction:         0,
	//	DistributorId:     0,
	//	MerchantId:        0,
	//	MerchantPaymentId: 0,
	//	AccountId:         "牛二春",
	//	CurrencyCrypto:    "BTUSD",
	//	CurrencyFiat:      "",
	//	PayType:           0,
	//	QrCode:            "",
	//	Name:              "",
	//	BankAccount:       "",
	//	Bank:              "",
	//	BankBranch:        "",
	//	Timestamp:         models.Timestamp{},
	//})
	//ret.Data = append(ret.Data, models.Order{
	//	OrderNumber:       "10003",
	//	Price:             5,
	//	Quantity:          "",
	//	Amount:            100,
	//	PaymentRef:        "",
	//	Status:            0,
	//	Direction:         0,
	//	DistributorId:     0,
	//	MerchantId:        0,
	//	MerchantPaymentId: 0,
	//	AccountId:         "牛三春",
	//	CurrencyCrypto:    "BTUSD",
	//	CurrencyFiat:      "",
	//	PayType:           0,
	//	QrCode:            "",
	//	Name:              "",
	//	BankAccount:       "",
	//	Bank:              "",
	//	BankBranch:        "",
	//	Timestamp:         models.Timestamp{},
	//})
	//ret.Data = append(ret.Data, models.Order{
	//	OrderNumber:       "10004",
	//	Price:             5,
	//	Quantity:          "",
	//	Amount:            400,
	//	PaymentRef:        "",
	//	Status:            0,
	//	Direction:         0,
	//	DistributorId:     0,
	//	MerchantId:        0,
	//	MerchantPaymentId: 0,
	//	AccountId:         "牛四春",
	//	CurrencyCrypto:    "BTUSD",
	//	CurrencyFiat:      "",
	//	PayType:           0,
	//	QrCode:            "",
	//	Name:              "",
	//	BankAccount:       "",
	//	Bank:              "",
	//	BankBranch:        "",
	//	Timestamp:         models.Timestamp{},
	//})
	//ret.Data = append(ret.Data, models.Order{
	//	OrderNumber:       "10005",
	//	Price:             5,
	//	Quantity:          "",
	//	Amount:            500,
	//	PaymentRef:        "",
	//	Status:            1,
	//	Direction:         0,
	//	DistributorId:     0,
	//	MerchantId:        0,
	//	MerchantPaymentId: 0,
	//	AccountId:         "牛五春",
	//	CurrencyCrypto:    "BTUSD",
	//	CurrencyFiat:      "",
	//	PayType:           0,
	//	QrCode:            "",
	//	Name:              "",
	//	BankAccount:       "",
	//	Bank:              "",
	//	BankBranch:        "",
	//	Timestamp:         models.Timestamp{},
	//})
	//ret.Data = append(ret.Data, models.Order{
	//	OrderNumber:       "10006",
	//	Price:             5,
	//	Quantity:          "",
	//	Amount:            100,
	//	PaymentRef:        "",
	//	Status:            1,
	//	Direction:         1,
	//	DistributorId:     0,
	//	MerchantId:        0,
	//	MerchantPaymentId: 0,
	//	AccountId:         "牛六春",
	//	CurrencyCrypto:    "BTUSD",
	//	CurrencyFiat:      "",
	//	PayType:           0,
	//	QrCode:            "",
	//	Name:              "",
	//	BankAccount:       "",
	//	Bank:              "",
	//	BankBranch:        "",
	//	Timestamp:         models.Timestamp{},
	//})
	//ret.Data = append(ret.Data, models.Order{
	//	OrderNumber:       "10007",
	//	Price:             5,
	//	Quantity:          "",
	//	Amount:            100,
	//	PaymentRef:        "",
	//	Status:            0,
	//	Direction:         0,
	//	DistributorId:     0,
	//	MerchantId:        0,
	//	MerchantPaymentId: 0,
	//	AccountId:         "牛七春",
	//	CurrencyCrypto:    "BTUSD",
	//	CurrencyFiat:      "",
	//	PayType:           0,
	//	QrCode:            "",
	//	Name:              "",
	//	BankAccount:       "",
	//	Bank:              "",
	//	BankBranch:        "",
	//	Timestamp:         models.Timestamp{},
	//})
	//ret.Data = append(ret.Data, models.Order{
	//	OrderNumber:       "10008",
	//	Price:             5,
	//	Quantity:          "",
	//	Amount:            100,
	//	PaymentRef:        "",
	//	Status:            0,
	//	Direction:         0,
	//	DistributorId:     0,
	//	MerchantId:        0,
	//	MerchantPaymentId: 0,
	//	AccountId:         "牛八春",
	//	CurrencyCrypto:    "BTUSD",
	//	CurrencyFiat:      "",
	//	PayType:           0,
	//	QrCode:            "",
	//	Name:              "",
	//	BankAccount:       "",
	//	Bank:              "",
	//	BankBranch:        "",
	//	Timestamp:         models.Timestamp{},
	//})
	//ret.Data = append(ret.Data, models.Order{
	//	OrderNumber:       "10009",
	//	Price:             5,
	//	Quantity:          "",
	//	Amount:            100,
	//	PaymentRef:        "",
	//	Status:            0,
	//	Direction:         0,
	//	DistributorId:     0,
	//	MerchantId:        0,
	//	MerchantPaymentId: 0,
	//	AccountId:         "牛九春",
	//	CurrencyCrypto:    "BTUSD",
	//	CurrencyFiat:      "",
	//	PayType:           0,
	//	QrCode:            "",
	//	Name:              "",
	//	BankAccount:       "",
	//	Bank:              "",
	//	BankBranch:        "",
	//	Timestamp:         models.Timestamp{},
	//})
	//ret.Data = append(ret.Data, models.Order{
	//	OrderNumber:       "10010",
	//	Price:             5,
	//	Quantity:          "",
	//	Amount:            1000,
	//	PaymentRef:        "",
	//	Status:            0,
	//	Direction:         0,
	//	DistributorId:     0,
	//	MerchantId:        0,
	//	MerchantPaymentId: 0,
	//	AccountId:         "牛十春",
	//	CurrencyCrypto:    "BTUSD",
	//	CurrencyFiat:      "",
	//	PayType:           0,
	//	QrCode:            "",
	//	Name:              "",
	//	BankAccount:       "",
	//	Bank:              "",
	//	BankBranch:        "",
	//	Timestamp:         models.Timestamp{},
	//})
	//ret.PageCount = 100
	//ret.PageNum = 1
	//ret.PageSize = 10
	//c.JSON(200, ret)
}

// @Summary 承兑商获取某一条订单的详情
// @Tags 承兑商APP API
// @Description 承兑商获取某一条订单的详情
// @Accept  json
// @Produce  json
// @Param uid  path  int  true  "承兑商用户id"
// @Param order_number  path  string  true  "订单id"
// @Success 200 {object} response.OrdersRet ""
// @Router /m/merchants/{uid}/orders/{order_number} [get]
func GetOrderDetail(c *gin.Context) {
	var uid int
	var err error
	if uid, err = strconv.Atoi(c.Param("uid")); err != nil {
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		var ret response.OrdersRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}

	orderNumber := c.Param("order_number")
	if orderNumber == "" {
		utils.Log.Errorf("order_number is empty")
		var ret response.OrdersRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return
	}

	c.JSON(200, service.GetOrderByMerchantIdAndOrderNumber(int64(uid), orderNumber))
	return

	//orderNumber := c.Param("order-id")
	//var ret response.OrdersRet
	//ret.Status = response.StatusSucc
	//if orderNumber == "10001" {
	//	ret.Data = append(ret.Data, models.Order{
	//		OrderNumber:       "10001",
	//		Price:             4.4,
	//		Quantity:          "25",
	//		Amount:            110,
	//		PaymentRef:        "",
	//		Status:            0,
	//		Direction:         1,
	//		DistributorId:     0,
	//		MerchantId:        0,
	//		MerchantPaymentId: 0,
	//		AccountId:         "牛大春",
	//		CurrencyCrypto:    "BTUSD",
	//		CurrencyFiat:      "RMB",
	//		PayType:           0,
	//		QrCode:            "https://image.baidu.com/search/down?tn=download&word=download&ie=utf8&fr=detail&url=https%3A%2F%2Ftimgsa.baidu.com%2Ftimg%3Fimage%26quality%3D80%26size%3Db9999_10000%26sec%3D1545215853444%26di%3Dd9eae39078073fa50b6c75cfb4f6b4ce%26imgtype%3D0%26src%3Dhttp%253A%252F%252Fpic.chinaz.com%252F2018%252F0409%252F18040918011940118.jpg&thumburl=https%3A%2F%2Fss1.bdstatic.com%2F70cFvXSh_Q1YnxGkpoWK1HF6hhy%2Fit%2Fu%3D3021223247%2C92123690%26fm%3D26%26gp%3D0.jpg",
	//		Name:              "",
	//		BankAccount:       "",
	//		Bank:              "",
	//		BankBranch:        "",
	//		Timestamp: models.Timestamp{
	//			CreatedAt: time.Date(2018, 11, 17, 20, 34, 58, 651387237, time.UTC),
	//		},
	//	})
	//} else if orderNumber == "10002" {
	//	ret.Data = append(ret.Data, models.Order{
	//		OrderNumber:       "10002",
	//		Price:             6.5,
	//		Quantity:          "20",
	//		Amount:            130,
	//		PaymentRef:        "",
	//		Status:            0,
	//		Direction:         0,
	//		DistributorId:     0,
	//		MerchantId:        0,
	//		MerchantPaymentId: 0,
	//		AccountId:         "牛二春",
	//		CurrencyCrypto:    "BTUSD",
	//		CurrencyFiat:      "RMB",
	//		PayType:           1,
	//		QrCode:            "http://13.250.12.109:8086/1_0_1_favicon.png",
	//		Name:              "",
	//		BankAccount:       "",
	//		Bank:              "",
	//		BankBranch:        "",
	//		Timestamp: models.Timestamp{
	//			CreatedAt: time.Date(2018, 11, 17, 20, 34, 58, 651387237, time.UTC),
	//		},
	//	})
	//}
	//
	//c.JSON(200, ret)
}

// @Summary 订单操作
// @Tags 承兑商APP API
// @Description 承兑商对订单进行操作，如接单、确认收款、确认付款
// @Accept  json
// @Produce  json
// @Param uid path int true "承兑商用户id"
// @Param body body models.Msg true "消息"
// @Success 200 {object} response.OrdersRet ""
// @Router /m/merchants/{uid}/orders/fulfill [put]
func OrderFulfill(c *gin.Context) {
	var ret response.EntityResponse
	var msg models.Msg

	uid := c.Param("uid")
	if err := c.ShouldBind(&msg); err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
	}
	var orderToFulfill service.OrderToFulfill
	engine := service.NewOrderFulfillmentEngine(nil)
	if msg.MsgType == models.Accept {
		data := msg.Data
		if len(data) > 0 {
			if id, err := strconv.ParseInt(uid, 10, 64); err == nil {
				if b, err := json.Marshal(data[0]); err == nil {
					if err := json.Unmarshal(b, &orderToFulfill); err == nil {
						utils.Log.Debugf("accept msg,%v", orderToFulfill)
						engine.AcceptOrder(orderToFulfill, id)
					}
				}
			}
		}
	} else {
		engine.UpdateFulfillment(msg)
	}

	ret.Status = response.StatusSucc
	c.JSON(200, ret)
}
