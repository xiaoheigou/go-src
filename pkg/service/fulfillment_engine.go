package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"reflect"
	"strconv"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/utils"
	"yuudidi.com/pkg/utils/timewheel"
)

var (
	engine                      *defaultEngine
	wheel                       *timewheel.TimeWheel // 如果币商不接单，则超时后，调用函数waitAcceptTimeout实现重新派单
	autoOrderAcceptWheel        *timewheel.TimeWheel // 对于自动订单，如果币商不接单，则超时后，调用函数waitAcceptTimeout实现重新派单
	officialMerchantAcceptWheel *timewheel.TimeWheel // 如果官方币商不接单，则超时后，调用函数waitOfficialMerchantAcceptTimeout实现重新派单
	notifyWheel                 *timewheel.TimeWheel // 如果一直不点击"我已付款"，则超时后（如900秒）会把订单状态改为5
	confirmWheel                *timewheel.TimeWheel // 如果一直没有确认收到对方的付款，则超时后（如900秒）会把订单状态改为5
	transferWheel               *timewheel.TimeWheel // 用户提现订单，冻结1小时（生产环境的时间配置）才放币
	suspendedWheel              *timewheel.TimeWheel // 订单方法异常，将订单修改为5，1
	unfreezeWheel               *timewheel.TimeWheel // 订单超时,45分钟后自动解冻
	awaitTimeout                int64
	retryTimeout                int64
	retries                     int64
	officialMerchantRetries     int64
	forbidNewOrderIfUnfinished  bool
)

// OrderToFulfill - order information for merchants to pick-up
type OrderToFulfill struct {
	//Order number to fulfill.
	OrderNumber string `json:"order_number"`
	//Trader Buy = 0, Trader Sell = 1
	Direction int `json:"direction"`
	//OriginOrder number
	OriginOrder string `json:"origin_order"`
	//AccountID
	AccountID string `json:"account"`
	//Distributor ID
	DistributorID int64 `json:"distributor"`
	//Crypto currency
	CurrencyCrypto string `json:"currency_crypto"`
	//Fiat currency
	CurrencyFiat string `json:"currency_fiat"`
	//Quantity, in crypto currency
	Quantity decimal.Decimal `json:"quantity"`
	//Price - rate between crypto and fiat
	Price float32 `json:"price"`
	//Amount of the order, in fiat currency
	Amount float64 `json:"amount"`
	//Payment type, chosen by trader
	PayType uint `json:"pay_type"`
	//微信或支付宝二维码地址
	QrCode string `gorm:"type:varchar(255)" json:"qr_code"`
	//微信或支付宝二维码所编码的字符串。对于自动订单，币商接单时，要在App端把这个值填上传回来（微信支付类型必须要传回来）
	QrCodeTxt string `json:"qr_code_txt"`
	//收款人姓名
	Name string `gorm:"type:varchar(100)" json:"name"`
	//银行账号
	BankAccount string `gorm:"" json:"bank_account"`
	//所属银行
	Bank string `gorm:"" json:"bank"`
	//所属银行分行
	BankBranch string `gorm:"" json:"bank_branch"`
	// 订单的接单类型，0表示手动接单订单，1表示自动接单订单。
	AcceptType int `json:"accept_type"`
	// 支付宝或微信的用户支付Id。对于自动订单，币商接单时，要在App端把把这个值填上传回来
	UserPayId string `json:"user_pay_id"`
	// 前端App生成二维码时，所需要的备注信息
	QrCodeMark string `json:"qr_code_mark"`
}

func getOrderToFulfillFromMapStrings(values map[string]interface{}) OrderToFulfill {
	var distributorID, direct, payT, acceptType int64
	var price, amount float64
	var quantity decimal.Decimal

	if distributorN, ok := values["distributor"].(json.Number); ok {
		distributorID, _ = distributorN.Int64()
	} else {
		utils.Log.Errorf("Type assertion fail, type of values[distributor] is %s", reflect.TypeOf(values["distributor"]).Kind())
	}
	if directN, ok := values["direction"].(json.Number); ok {
		direct, _ = directN.Int64()
	} else {
		utils.Log.Errorf("Type assertion fail, type of values[direction] is %s", reflect.TypeOf(values["direction"]).Kind())
	}
	if payTN, ok := values["pay_type"].(json.Number); ok {
		payT, _ = payTN.Int64()
	} else {
		utils.Log.Errorf("Type assertion fail, type of values[pay_type] is %s", reflect.TypeOf(values["pay_type"]).Kind())
	}
	if acceptTypeN, ok := values["accept_type"].(json.Number); ok {
		acceptType, _ = acceptTypeN.Int64()
	} else {
		utils.Log.Errorf("Type assertion fail, type of values[accept_type] is %s", reflect.TypeOf(values["accept_type"]).Kind())
	}
	if quantityS, ok := values["quantity"].(string); ok {
		quantity, _ = decimal.NewFromString(quantityS)
	} else {
		utils.Log.Errorf("Type assertion fail, type of values[quantity] is %s", reflect.TypeOf(values["quantity"]).Kind())
	}
	if priceN, ok := values["price"].(json.Number); ok {
		price, _ = priceN.Float64()
	} else {
		utils.Log.Errorf("Type assertion fail, type of values[price] is %s", reflect.TypeOf(values["price"]).Kind())
	}
	if amountN, ok := values["amount"].(json.Number); ok {
		amount, _ = amountN.Float64()
	} else {
		utils.Log.Errorf("Type assertion fail, type of values[amount] is %s", reflect.TypeOf(values["amount"]).Kind())
	}
	var qrCode, qrCodeTxt, name, bank, bankAccount, bankBranch, userPayId string
	if values["qr_code"] != nil {
		qrCode = values["qr_code"].(string)
	}
	if values["qr_code_txt"] != nil {
		qrCode = values["qr_code_txt"].(string)
	}
	if values["name"] != nil {
		name = values["name"].(string)
	}
	if values["bank"] != nil {
		bank = values["bank"].(string)
	}
	if values["bank_branch"] != nil {
		bankBranch = values["bank_branch"].(string)
	}
	if values["bank_account"] != nil {
		bankAccount = values["bank_account"].(string)
	}
	if values["user_pay_id"] != nil {
		userPayId = values["user_pay_id"].(string)
	}

	return OrderToFulfill{
		OrderNumber:    values["order_number"].(string),
		Direction:      int(direct),
		AccountID:      values["account"].(string),
		OriginOrder:    values["origin_order"].(string),
		DistributorID:  distributorID,
		CurrencyCrypto: values["currency_crypto"].(string),
		CurrencyFiat:   values["currency_fiat"].(string),
		Quantity:       quantity,
		Price:          float32(price),
		Amount:         amount,
		PayType:        uint(payT),
		QrCode:         qrCode,
		QrCodeTxt:      qrCodeTxt,
		Name:           name,
		Bank:           bank,
		BankAccount:    bankAccount,
		BankBranch:     bankBranch,
		AcceptType:     int(acceptType),
		UserPayId:      userPayId,
	}
}

// OrderFulfillment - Order fulfillment result.
type OrderFulfillment struct {
	OrderToFulfill
	//Merchant ID
	MerchantID int64 `json:"merchant_id"`
	//Merchant nickname
	MerchantNickName string `json:"merchant_nickname"`
	//Merchant avatar URI
	MerchantAvatarURI string `json:"merchant_avartar_uri"`
	//Paytype - 1 wechat, 2 zhifubao, 4 bank, support combination
	PayType int
	//Payment information, by reference
	PaymentInfo []models.PaymentInfo `json:"payment_info"`
}

// OrderFulfillmentEngine - engine interface of order fulfillment.
// The platform may change to new engine according to fulfillment rules changing.
type OrderFulfillmentEngine interface {
	// FulfillOrder - Async request to fulfill an order.
	FulfillOrder(
		order *OrderToFulfill, //Order to fulfill (demands information)
	)
	// selectMerchantsToFulfillOrder - Select merchants to fulfill the specified orders.
	// The returned merchant(s) would receive the OrderToFulfill object through notification channel.
	// When there's only one merchant returned in the result, it might be exhausted matching result
	// or the first automatic processing merchant selected. No matter which situation, just send OrderToFulfill
	// message to the selected merchant. [no different process logic needed by caller]
	selectMerchantsToFulfillOrder(order *OrderToFulfill) *[]int64
	// AcceptOrder - receive merchants' response on accept order.
	// it then pick up the winner of all responded merchants and notify the fulfillment
	AcceptOrder(
		order OrderToFulfill, //order number to accept
		merchantID int64, //merchant id
	)
	// UpdateFulfillment - update fulfillment processing like payment notified, confirm payment, etc..
	// Upon receiving these message, fulfillment engine should update order/fulfillment status + appended extra message
	UpdateFulfillment(
		msg models.Msg, // Order number
		//operation int, // fulfilment operation such as notify_paid, payment_confirmed, etc..
		//data interface{}, // arbitrary notification data according to different operation
	)
	DeleteWheel(orderNumber string)
}

// NewOrderFulfillmentEngine - return a new fulfillment engine.
// Factory method of fulfillment engine to adopt to future engine extension.
func NewOrderFulfillmentEngine(_ /*config*/ interface{}) OrderFulfillmentEngine {
	//Singleton, may init engine by config, now ignore it
	if engine == nil {
		utils.SetSettings()
		engine = new(defaultEngine)
		//init timewheel
		//never stop till process killed!
	}
	return engine
}

// 一轮派单后，如果超时没有接，这个函数就会启动
func waitAcceptTimeout(data interface{}) {
	//no one accept till timeout, re-fulfill it then

	orderNum := data.(string)
	utils.Log.Infof("func waitAcceptTimeout, order %s not accepted by any merchant. Re-fulfill it...", orderNum)
	order := models.Order{}
	if utils.DB.First(&order, "order_number = ?", orderNum).RecordNotFound() {
		utils.Log.Errorf("Order %s not found.", orderNum)
		return
	}
	orderToFulfill := OrderToFulfill{
		OrderNumber:    order.OrderNumber,
		Direction:      order.Direction,
		OriginOrder:    order.OriginOrder,
		AccountID:      order.AccountId,
		DistributorID:  order.DistributorId,
		CurrencyCrypto: order.CurrencyCrypto,
		CurrencyFiat:   order.CurrencyFiat,
		Quantity:       order.Quantity,
		Price:          order.Price,
		Amount:         order.Amount,
		PayType:        order.PayType,
		QrCode:         order.QrCode,
		Name:           order.Name,
		Bank:           order.Bank,
		BankAccount:    order.BankAccount,
		BankBranch:     order.BankBranch,
	}
	// 发送给了候选币商，但他们都没有接单，重新派单
	go reFulfillOrder(&orderToFulfill, 1)
}

// 超时没点“我已付款”时，下面函数会被调用
func notifyPaidTimeout(data interface{}) {
	//key = order number
	//no one accept till timeout, re-fulfill it then
	orderNum := data.(string)
	utils.Log.Infof("Order %s paid timeout.", orderNum)

	tx := utils.DB.Begin()
	if tx.Error != nil {
		utils.Log.Debugf("tx in func notifyPaidTimeout begin fail, tx=[%v]", tx)
		utils.Log.Errorf("func notifyPaidTimeout finished abnormally.")
		return
	}

	utils.Log.Debugf("tx in func notifyPaidTimeout begin, tx=[%v]", tx)

	order := models.Order{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&order, "order_number = ?", orderNum).RecordNotFound() {
		utils.Log.Errorf("Order %s not found.", orderNum)
		return
	}
	originStatus := order.Status
	if order.Status < models.NOTIFYPAID {
		if order.Direction == 0 {
			//订单支付方式释放
			if err := tx.Model(&models.PaymentInfo{}).Where("uid = ? AND id = ?", order.MerchantId, order.MerchantPaymentId).Update("in_use", 0).Error; err != nil {
				utils.Log.Errorf("notifyPaidTimeout release payment info,merchantId:[%d],orderNUmber:[%s]", order.MerchantId, order.OrderNumber)
				utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
				tx.Rollback()
				//超时更新失败，修改订单状态suspended
				suspendedWheel.Add(orderNum)
				return
			}

			if err := SendSmsOrderPaidTimeout(order.MerchantId, orderNum); err != nil {
				utils.Log.Errorf("order [%v] is not paid, and timeout, send sms fail. error [%v]", orderNum, order.MerchantId, err)
			}

		} else if order.Direction == 1 {
			//用户提现单子,确认付款超时,不释放币
		}

		//订单状态改为suspended
		//failed, highlight the order to set status to "SUSPENDED"
		if err := tx.Model(&order).Where("order_number = ? AND status < ?", order.OrderNumber, models.NOTIFYPAID).Updates(models.Order{Status: models.SUSPENDED, StatusReason: models.PAIDTIMEOUT}).Error; err != nil {
			utils.Log.Errorf("Update order %s status to SUSPENDED failed", order.OrderNumber)
			utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
			tx.Rollback()
			//超时更新失败，修改订单状态suspended
			suspendedWheel.Add(orderNum)
			return
		}
		//fulfillment log 添加记录
		var fulfillment models.Fulfillment
		if err := tx.Order("seq_id desc").First(&fulfillment, "order_number = ?", orderNum).Error; err != nil {
			utils.Log.Errorf("get fulfillment order %s failed", order.OrderNumber)
			utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
			tx.Rollback()
			//超时更新失败，修改订单状态suspended
			suspendedWheel.Add(orderNum)
			return
		}
		fulfillmentLog := models.FulfillmentLog{
			FulfillmentID: fulfillment.Id,
			OrderNumber:   order.OrderNumber,
			SeqID:         fulfillment.SeqID,
			IsSystem:      true,
			MerchantID:    order.MerchantId,
			AccountID:     order.AccountId,
			DistributorID: order.DistributorId,
			OriginStatus:  originStatus,
			UpdatedStatus: models.SUSPENDED,
		}
		if err := tx.Create(&fulfillmentLog).Error; err != nil {
			utils.Log.Errorf("notifyPaidTimeout create fulfillmentLog is failed,order number:%s", orderNum)
			utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
			tx.Rollback()
			//超时更新失败，修改订单状态suspended
			suspendedWheel.Add(orderNum)
			return
		}

	}

	utils.Log.Debugf("tx in func notifyPaidTimeout commit, tx=[%v]", tx)
	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("error tx in func notifyPaidTimeout commit, err=[%v]", err)
		tx.Rollback()
		//超时更新失败，修改订单状态suspended
		suspendedWheel.Add(orderNum)
	}

	//充值订单，如果h5端没有点击我已付款，开始记时,超时后自动解冻
	if order.Direction == 0 {
		unfreezeWheel.Add(orderNum)
	}

}

//充值订单付款超时,超时会自动解冻
func autoUnfreeze(data interface{}) {
	orderNumber := data.(string)
	utils.Log.Debugf("autoUnfreeze is begin,orderNumber:%s", orderNumber)

	//调用解冻方法,username和userId 传空是为了将来区分自动解冻和客服人工解冻
	ret := UnFreezeCoin(orderNumber, "", -1)
	//解冻失败,将订单
	if ret.Status != response.StatusSucc {
		suspendedWheel.Add(orderNumber)
	}

}

// 超时没点“已收到对方付款”时，下面函数会被调用
func confirmPaidTimeout(data interface{}) {
	//ignore first fmanager object, add later if needed
	//key = order number
	//no one accept till timeout, re-fulfill it then
	orderNum := data.(string)
	utils.Log.Infof("Order %s confirm paid timeout", orderNum)

	// 确认付款超时，不放币
	//order := models.Order{}
	//if utils.DB.First(&order, "order_number = ?", orderNum).RecordNotFound() {
	//	utils.Log.Errorf("Order %s not found.", orderNum)
	//	return
	//}
	//message := models.Msg{
	//	MsgType:    models.ConfirmPaid,
	//	MerchantId: []int64{order.MerchantId},
	//	H5:         []string{orderNum},
	//	Timeout:    0,
	//	Data: []interface{}{
	//		map[string]interface{}{
	//			"order_number": order.OrderNumber,
	//			"direction":    order.Direction,
	//		},
	//	},
	//}
	//if _, err := uponConfirmPaid(message); err != nil {
	//	utils.Log.Errorf("confirmPaidTimeout to uponConfirmPaid is failed,OrderNumber:%s", orderNum)
	//	suspendedWheel.Add(orderNum)
	//}

	tx := utils.DB.Begin()
	if tx.Error != nil {
		utils.Log.Debugf("tx in func confirmPaidTimeout begin fail, tx=[%v]", tx)
		utils.Log.Errorf("func confirmPaidTimeout finished abnormally.")
		return
	}

	utils.Log.Debugf("tx in func confirmPaidTimeout begin, tx=[%v]", tx)

	order := models.Order{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&order, "order_number = ?", orderNum).RecordNotFound() {
		utils.Log.Errorf("Order %s not found.", orderNum)
		return
	}
	originStatus := order.Status
	if order.Status < models.CONFIRMPAID {
		if order.Direction == 0 {
			//订单支付方式释放
			if err := tx.Model(&models.PaymentInfo{}).Where("uid = ? AND id = ?", order.MerchantId, order.MerchantPaymentId).Update("in_use", 0).Error; err != nil {
				utils.Log.Errorf("confirmPaidTimeout release payment info,merchantId:[%d],orderNUmber:[%s]", order.MerchantId, order.OrderNumber)
				utils.Log.Errorf("tx in func confirmPaidTimeout rollback, tx=[%v]", tx)
				tx.Rollback()
				//超时更新失败，修改订单状态suspended
				suspendedWheel.Add(orderNum)
				return
			}

			// 币商超时为确认收款，发送短信
			//if err := SendSmsOrderPaidTimeout(order.MerchantId, orderNum); err != nil {
			//	utils.Log.Errorf("order [%v] is not paid, and timeout, send sms fail. error [%v]", orderNum, order.MerchantId, err)
			//}

		} else if order.Direction == 1 {
			//用户提现单子,没有确认收款超时
		}

		//订单状态改为suspended
		//failed, highlight the order to set status to "SUSPENDED"
		if err := tx.Model(&order).Where("order_number = ? AND status < ?", order.OrderNumber, models.CONFIRMPAID).Updates(models.Order{Status: models.SUSPENDED, StatusReason: models.CONFIRMTIMEOUT}).Error; err != nil {
			utils.Log.Errorf("Update order %s status to SUSPENDED failed", order.OrderNumber)
			utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
			tx.Rollback()
			//超时更新失败，修改订单状态suspended
			suspendedWheel.Add(orderNum)
			return
		}
		//fulfillment log 添加记录
		var fulfillment models.Fulfillment
		if err := tx.Order("seq_id desc").First(&fulfillment, "order_number = ?", orderNum).Error; err != nil {
			utils.Log.Errorf("get fulfillment order %s failed", order.OrderNumber)
			utils.Log.Errorf("tx in func confirmPaidTimeout rollback, tx=[%v]", tx)
			tx.Rollback()
			//超时更新失败，修改订单状态suspended
			suspendedWheel.Add(orderNum)
			return
		}
		fulfillmentLog := models.FulfillmentLog{
			FulfillmentID: fulfillment.Id,
			OrderNumber:   order.OrderNumber,
			SeqID:         fulfillment.SeqID,
			IsSystem:      true,
			MerchantID:    order.MerchantId,
			AccountID:     order.AccountId,
			DistributorID: order.DistributorId,
			OriginStatus:  originStatus,
			UpdatedStatus: models.SUSPENDED,
		}
		if err := tx.Create(&fulfillmentLog).Error; err != nil {
			utils.Log.Errorf("confirmPaidTimeout create fulfillmentLog is failed,order number:%s", orderNum)
			utils.Log.Errorf("tx in func confirmPaidTimeout rollback, tx=[%v]", tx)
			tx.Rollback()
			//超时更新失败，修改订单状态suspended
			suspendedWheel.Add(orderNum)
			return
		}

	}

	utils.Log.Debugf("tx in func confirmPaidTimeout commit, tx=[%v]", tx)
	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("error tx in func confirmPaidTimeout commit, err=[%v]", err)
		tx.Rollback()
		//超时更新失败，修改订单状态suspended
		suspendedWheel.Add(orderNum)
	}
}

//转账时间到
func transferTimeout(data interface{}) {
	orderNum := data.(string)
	utils.Log.Debugf("transfer timeout begin,orderNum:%s", orderNum)
	if err := doTransfer(orderNum); err != nil {
		utils.Log.Errorf("transferTimeout to doTransfer is failed,orderNumber:%s", orderNum)
		suspendedWheel.Add(orderNum)
	}
}

//由于系统内部错误（如数据库异常等）导致操作不成功后会调用这个方法，把订单状态修改为5（suspended），保证该订单有机会在管理后台进行进一步操作（管理后台目前只能修改状态为5的订单）
func updateOrderStatusAsSuspended(data interface{}) {
	orderNum := data.(string)
	tx := utils.DB.Begin()
	order := models.Order{}
	if tx.Set("gorm:query_option", "FOR UPDATE").Where("order_number = ?", orderNum).First(&order).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("Record not found: order with number %s.", orderNum)
		utils.Log.Errorf("tx in func updateOrderStatusAsSuspended rollback, tx=[%v]", tx)
		utils.Log.Errorf("func updateOrderStatusAsSuspended finished abnormally.")
		suspendedWheel.Add(orderNum)
		return
	}

	originStatus := order.Status
	// 不能把已经完成的订单标记为系统异常状态 5.1 (SUSPENDED.SYSTEMUPDATEFAIL)，否则转币会出现问题：
	// 完成订单时进行了转币操作，标记为系统异常后，在管理后台还可以再次进行转币操作。
	if originStatus == models.TRANSFERRED {
		tx.Rollback()
		utils.Log.Errorf("order %s has status TRANSFERRED, cannot change to SUSPENDED.SYSTEMUPDATEFAIL", orderNum)
		utils.Log.Errorf("tx in func updateOrderStatusAsSuspended rollback, tx=[%v]", tx)
		utils.Log.Errorf("func updateOrderStatusAsSuspended finished abnormally. order_number = %s", orderNum)
		return
	}
	if originStatus == models.SUSPENDED && (order.StatusReason == models.MARKCOMPLETED || order.StatusReason == models.CANCEL) {
		tx.Rollback()
		utils.Log.Errorf("order %s has status MARKCOMPLETED or CANCEL, cannot change to SUSPENDED.SYSTEMUPDATEFAIL", orderNum)
		utils.Log.Errorf("tx in func updateOrderStatusAsSuspended rollback, tx=[%v]", tx)
		utils.Log.Errorf("func updateOrderStatusAsSuspended finished abnormally. order_number = %s", orderNum)
		return
	}

	if err := tx.Model(&models.Order{}).Where("order_number = ?", orderNum).Updates(models.Order{Status: models.SUSPENDED, StatusReason: models.SYSTEMUPDATEFAIL}).Error; err != nil {
		utils.Log.Errorf("update order status as suspended,is fail ,will retry,orderNumber:%s", orderNum)
		tx.Rollback()
		suspendedWheel.Add(orderNum)
		return
	}
	//fulfillment log 添加记录
	var fulfillment models.Fulfillment
	if err := tx.Order("seq_id desc").First(&fulfillment, "order_number = ?", orderNum).Error; err != nil {
		utils.Log.Errorf("get fulfillment order %s failed", order.OrderNumber)
		utils.Log.Errorf("tx in func updateOrderStatusAsSuspended rollback, tx=[%v]", tx)
		tx.Rollback()
		suspendedWheel.Add(orderNum)
		return
	}
	fulfillmentLog := models.FulfillmentLog{
		FulfillmentID: fulfillment.Id,
		OrderNumber:   order.OrderNumber,
		SeqID:         fulfillment.SeqID,
		IsSystem:      true,
		MerchantID:    order.MerchantId,
		AccountID:     order.AccountId,
		DistributorID: order.DistributorId,
		OriginStatus:  originStatus,
		StatusReason:  models.SYSTEMUPDATEFAIL,
		UpdatedStatus: models.SUSPENDED,
	}
	if err := tx.Create(&fulfillmentLog).Error; err != nil {
		utils.Log.Errorf("updateOrderStatusAsSuspended create fulfillmentLog is failed,order number:%s", orderNum)
		utils.Log.Errorf("tx in func updateOrderStatusAsSuspended rollback, tx=[%v]", tx)
		tx.Rollback()
		suspendedWheel.Add(orderNum)
		return
	}
	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("updateOrderStatusAsSuspended commit is failed,order number:%s", orderNum)
		utils.Log.Errorf("tx in func updateOrderStatusAsSuspended rollback, tx=[%v]", tx)
		tx.Rollback()
		suspendedWheel.Add(orderNum)
	}
}

//defaultEngine - hidden default OrderFulfillmentEngine
type defaultEngine struct {
}

func (engine *defaultEngine) FulfillOrder(
	order *OrderToFulfill,
) {
	utils.AddBackgroundJob(
		utils.FulfillOrderTask,
		utils.NormalPriority,
		order)
}

func (engine *defaultEngine) selectMerchantsToFulfillOrder(order *OrderToFulfill) *[]int64 {
	utils.Log.Debugf("func selectMerchantsToFulfillOrder begin, order = [%+v] and selected merchants [%v]", order)
	//search logic(in business prospective):
	//0. prioritize those run in "automatically comfirm payment" && "accept order" mode merchant, verify to see if anyone meets the demands
	//   (coin, payment type, fix-amount payment QR). If none matches, then:

	//1. filter out merchants currently not in "accept order" mode;
	//2. filter out merchants who don't have enough "coins" to take Trader-Buy order;
	//3. filter out merchants who don't support Trader specified payment type;
	//5. prioritize those merchants who do have "fix-amount" payment QR code matching demand;
	//6. constraints: merchant's payment info can serve one order at same time (locked if already matched previous order)
	//7. constraints: merchant can only take one same "amount" order at same time;
	//8. constraints: risk-control concerns which may reject to assign order to some merchant (TODO: to be added later)

	//implementation:
	//call service.GetMerchantsQualified(quote string, currencyCrypto string, pay_type uint8, fix bool, group uint8, limit uin8) []int64
	// with parameters copied from order set, in order:
	var merchants, merchantsUnfinished, alreadyFulfillMerchants, selectedMerchants []int64
	//去掉已经派过单的币商
	if data, err := utils.GetCacheSetMembers(utils.RedisKeyMerchantSelected(order.OrderNumber)); err != nil {
		utils.Log.Errorf("func selectMerchantsToFulfillOrder error, the select order = [%+v]", order)
	} else if len(data) > 0 {
		utils.Log.Infof("order %s had sent to merchants [%v] before, filter out them in this round", order.OrderNumber, selectedMerchants)
		utils.ConvertStringToInt(data, &selectedMerchants)
	}

	var isAutoOrder = false
	if order.Direction == 0 {
		//如果是银行卡,先优先匹配相同银行,在匹配不同银行,通过固定金额的参数fix进行区分,并且银行卡只有手动
		if order.PayType >= 4 {
			//1. fix 为true 只查询银行相同的币商
			merchants = utils.DiffSet(GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, true, false, 0, 0), selectedMerchants)
			if len(merchants) == 0 {
				// 2. available merchants(online + in_work) + manual accept order/confirm payment + has arbitrary amount qrcode
				merchants = utils.DiffSet(GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, false, false, 0, 0), selectedMerchants)
			}
		} else if order.PayType > 0 {
			//Buy, try to match all-automatic merchants firstly
			// 1. available merchants(online + in_work) + auto accept order/confirm payment
			merchants = utils.DiffSet(GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, false, true, 0, 0), selectedMerchants)
			if len(merchants) > 0 { // 找到了可以接收自动订单的币商
				isAutoOrder = true
			} else if len(merchants) == 0 { //no priority merchants with non-fix amount match found, then "manual operation" merchants
				// 2. available merchants(online + in_work) + manual accept order/confirm payment + has fix amount qrcode
				merchants = utils.DiffSet(GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, true, false, 0, 0), selectedMerchants)
				if len(merchants) == 0 { //Sell, all should manually processed
					// 3. available merchants(online + in_work) + manual accept order/confirm payment + has arbitrary amount qrcode
					merchants = utils.DiffSet(GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, false, false, 0, 0), selectedMerchants)
				}
			}
		}

		if forbidNewOrderIfUnfinished {
			// 只允许同时接一个订单
			if err := utils.DB.Model(models.Order{}).Where("status <= ? AND merchant_id > 0", models.NOTIFYPAID).Pluck("merchant_id", &merchantsUnfinished).Error; err != nil {
				utils.Log.Errorf("func selectMerchantsToFulfillOrder error, the select order = [%+v]", order)
			}
			utils.Log.Debugf("merchants [%v] have unfinished orders, filter out them in this round.", merchantsUnfinished)
		}

	} else {
		//Sell, any online + in_work could pickup order
		merchants = utils.DiffSet(GetMerchantsQualified(0, decimal.Zero, order.CurrencyCrypto, order.PayType, true, false, 0, 1), selectedMerchants)
		if len(merchants) == 0 {
			merchants = GetMerchantsQualified(0, decimal.Zero, order.CurrencyCrypto, order.PayType, false, false, 0, 1)
		}

		// 对于用户提现单，正常派单时，不派给官方币商
		// 官方币商仅当没有接单时，才会派给他们
		officialMerchants := getOfficialMerchants()
		if len(officialMerchants) > 0 {
			utils.Log.Debugf("filter out official merchants %v in normal fulfillment", officialMerchants)
			merchants = utils.DiffSet(merchants, officialMerchants)
		}
	}

	//重新派单时，去除已经接过这个订单的币商
	if err := utils.DB.Model(&models.Fulfillment{}).Where("order_number = ?", order.OrderNumber).Pluck("distinct merchant_id", &alreadyFulfillMerchants).Error; err != nil {
		utils.Log.Errorf("selectMerchantsToFulfillOrder get fulfillment is failed,orderNumber:%s", order.OrderNumber)
	}
	merchants = utils.DiffSet(merchants, selectedMerchants, merchantsUnfinished, alreadyFulfillMerchants)

	utils.Log.Debugf("before sort by last order time, the merchants = [%+v]", merchants)
	merchants = sortMerchantsByLastOrderTime(merchants, order.Direction)
	utils.Log.Debugf(" after sort by last order time, the merchants = [%+v]", merchants)

	if len(merchants) > 0 {
		if isAutoOrder {
			order.AcceptType = 1

			if len(merchants) > 1 {
				// 对于自动订单，只发订单给一个币商
				merchants = merchants[0:1]
			}
		} else {
			order.AcceptType = 0

			// 限制一轮最多给oneRoundSize个币商派单
			var oneRoundSize int64
			var err error
			if oneRoundSize, err = strconv.ParseInt(utils.Config.GetString("fulfillment.oneroundsize"), 10, 64); err != nil {
				utils.Log.Warnf("invalid configuration fulfillment.oneroundsize [%s], use 10 as default", utils.Config.GetString("fulfillment.oneroundsize"))
				oneRoundSize = 10
			}
			if len(merchants) > int(oneRoundSize) {
				utils.Log.Debugf("the candidate num [%d] is more than max size [%d] in one round, pick first [%d] merchants in this round", len(merchants), oneRoundSize, oneRoundSize)
				// 只选前oneRoundSize个币商
				merchants = merchants[0:oneRoundSize]
			}
		}
	}

	utils.Log.Debugf("func selectMerchantsToFulfillOrder finished, the select merchants = [%+v]", merchants)
	return &merchants
}

// 按merchants接单时间排序
func sortMerchantsByLastOrderTime(merchants []int64, direction int) []int64 {
	var redisSorted []string
	var redisSortedInt64 []int64
	var err error

	// 按redis中保存的merchants的接单时间，对merchants进行排序（接单早的排在前面）
	// 如果merchant还没有接过单，则在redis中没有记录，它也不会出现在结果集redisSorted中
	if redisSorted, err = utils.GetMerchantsSortedByLastOrderTime(direction); err != nil {
		utils.Log.Error("func sortMerchantsByLastOrderTime fail, call GetMerchantsSortedByLastOrderTime fail [%v]", err)
		return merchants
	}
	if err := utils.ConvertStringToInt(redisSorted, &redisSortedInt64); err != nil {
		utils.Log.Error("func sortMerchantsByLastOrderTime fail, call ConvertStringToInt fail [%v]", err)
		return merchants
	}

	var merchantsWithoutSuccOrder = utils.DiffSet(merchants, redisSortedInt64) // 从未接过单的merchants

	return append(merchantsWithoutSuccOrder, utils.InterSetInt64(redisSortedInt64, merchants)...)
}

func (engine *defaultEngine) AcceptOrder(
	order OrderToFulfill,
	merchantID int64,
) {
	utils.Log.Debugf("func AcceptOrder begin, order = [%+v], merchant = %d", order, merchantID)
	//check cache to see if anyone already accepted this order
	orderNum := order.OrderNumber
	key := utils.Config.GetString("cache.redis.prefix") + ":" + utils.Config.GetString("cache.key.acceptorder") + ":" + orderNum
	if merchant, err := utils.RedisClient.Get(key).Result(); err == redis.Nil {
		//book merchant
		utils.Log.Debugf("Order %s already accepted by %d", orderNum, merchantID)
		periodStr := utils.Config.GetString("fulfillment.timeout.accept")
		period, _ := strconv.ParseInt(periodStr, 10, 0)
		utils.RedisClient.Set(key, merchantID, time.Duration(2*period)*time.Second)
		//remove it from wheel
		//wheel.Remove(order.OrderNumber)
		utils.AddBackgroundJob(utils.AcceptOrderTask, utils.HighPriority, order, merchantID)

	} else { //already accepted, reject the request
		utils.Log.Debugf("merchant %d accepted order is failed,order already by merchant %s accept.", merchantID, merchant)
		data := []OrderToFulfill{{
			OrderNumber: orderNum,
		}}
		if err := NotifyThroughWebSocketTrigger(models.Picked, &[]int64{merchantID}, &[]string{}, 60, data); err != nil {
			utils.Log.Errorf("Notify Picked through websocket ")
		}
	}

	utils.Log.Debugf("func AcceptOrder finished finished, order_number, merchant = %d", order.OrderNumber, merchantID)
}

func (engine *defaultEngine) UpdateFulfillment(
	msg models.Msg,
) {
	utils.Log.Debugf("func UpdateFulfillment begin, msg = [%+v]", msg)
	utils.AddBackgroundJob(utils.UpdateFulfillmentTask, utils.NormalPriority, msg)
	utils.Log.Debugf("func UpdateFulfillment finished finished.")
}

func (engine *defaultEngine) DeleteWheel(orderNumber string) {
	utils.Log.Debugf("func DeleteWheel begin, msg = [%+v]", orderNumber)
	utils.AddBackgroundJob(utils.DeleteWheel, utils.HighPriority, orderNumber)
	utils.Log.Debugf("func DeleteWheel finished finished.")
}

//wrapper methods complies to goworker func.
func fulfillOrder(queue string, args ...interface{}) error {
	utils.Log.Debugf("func fulfillOrder begin.")
	//recover OrderToFulfill from args
	var order OrderToFulfill
	if orderArg, ok := args[0].(map[string]interface{}); ok {
		order = getOrderToFulfillFromMapStrings(orderArg)
	} else {
		return fmt.Errorf("Wrong order arg: %v", args[0])
	}
	utils.Log.Debugf("fulfill for order [%+V]", order.OrderNumber)
	merchants := engine.selectMerchantsToFulfillOrder(&order)
	if len(*merchants) == 0 {
		utils.Log.Warnf("func fulfillOrder, no merchant is available at moment, re-fulfill order %s later.", order.OrderNumber)
		// 第一轮派单，没找到候选币商，下面开始重新派单
		go reFulfillOrder(&order, 1)
		return nil
	}
	//send order to pick
	if err := sendOrder(&order, merchants); err != nil {
		utils.Log.Errorf("Send order %s to merchants failed: %v", order.OrderNumber, err)
		return err
	}
	// 等待币商接单
	if order.AcceptType == 0 {
		wheel.Add(order.OrderNumber)
	} else {
		// 自动订单的接单超时时间比一般订单的接单超时时间更短
		autoOrderAcceptWheel.Add(order.OrderNumber)
	}

	utils.Log.Debugf("func fulfillOrder finished normally. order_number = %s", order.OrderNumber)
	return nil
}

func getOfficialMerchants() []int64 {
	officialMerchants := []int64{}

	// 先从redis读取
	if officialMerchantsStr, err := utils.GetCacheSetMembers(utils.RedisKeyMerchantRole1()); err != nil {
		utils.ConvertStringToInt(officialMerchantsStr, &officialMerchants)
	}

	// 读不到，则从db中读取
	if len(officialMerchants) == 0 {
		db := utils.DB.Model(&models.Merchant{}).Where("role = 1")
		if err := db.Pluck("id", &officialMerchants).Error; err != nil {
			utils.Log.Errorf("getOfficialMerchants from db failed.")
		}

		// 保存到redis中
		for _, officialMerchant := range officialMerchants {
			expireTimeInSecond := 600 // 10分钟过期，过期后重新从数据库读取
			if err := utils.SetCacheSetMember(utils.RedisKeyMerchantRole1(), expireTimeInSecond, officialMerchant); err != nil {
				utils.Log.Errorf("add official Merchant %s to redis fail, err", officialMerchant, err)
			}
		}
	}

	utils.Log.Debugf("official merchants :%v", officialMerchants)
	return officialMerchants
}

// 派单给官方币商后，如果超时没有接，这个函数就会启动，重新派单
func waitOfficialMerchantAcceptTimeout(data interface{}) {
	orderNum := data.(string)
	utils.Log.Infof("func waitOfficialMerchantAcceptTimeout, order %s not accepted by any official merchant. Re-fulfill it...", orderNum)
	order := models.Order{}
	if utils.DB.First(&order, "order_number = ?", orderNum).RecordNotFound() {
		utils.Log.Errorf("Order %s not found.", orderNum)
		return
	}
	if order.Status == models.TRANSFERRED ||
		(order.Status == models.SUSPENDED && (order.StatusReason == models.MARKCOMPLETED || order.StatusReason == models.CANCEL)) {
		utils.Log.Warnf("Order %s has final status, cannot reFulfill", orderNum)
		return
	}
	orderToFulfill := OrderToFulfill{
		OrderNumber:    order.OrderNumber,
		Direction:      order.Direction,
		OriginOrder:    order.OriginOrder,
		AccountID:      order.AccountId,
		DistributorID:  order.DistributorId,
		CurrencyCrypto: order.CurrencyCrypto,
		CurrencyFiat:   order.CurrencyFiat,
		Quantity:       order.Quantity,
		Price:          order.Price,
		Amount:         order.Amount,
		PayType:        order.PayType,
		QrCode:         order.QrCode,
		Name:           order.Name,
		Bank:           order.Bank,
		BankAccount:    order.BankAccount,
		BankBranch:     order.BankBranch,
	}

	// 发送给了官方币商，但他们都没有接单，接着重新派单
	go reFulfillOrderToOfficialMerchants(&orderToFulfill)
}

func reFulfillOrderToOfficialMerchants(order *OrderToFulfill) {
	utils.Log.Debugf("func reFulfillOrderToOfficialMerchants begin, order_number = %s", order.OrderNumber)
	if order.Direction == 0 {
		utils.Log.Warnf("func reFulfillOrderToOfficialMerchants is not applicable for order with direction = 0")
		return
	} else if order.Direction == 1 {

		time.Sleep(time.Duration(retryTimeout) * time.Second)

		seq := utils.RedisGetRefulfillTimesToOfficialMerchants(order.OrderNumber)

		if seq < officialMerchantRetries {
			merchants := getOfficialMerchants()
			if len(merchants) == 0 {
				utils.Log.Errorf("func reFulfillOrderToOfficialMerchants, can not find any official merchants")
			} else {
				utils.Log.Debugf("func reFulfillOrderToOfficialMerchants, send order to official merchant %v", merchants)
				if err := sendOrder(order, &merchants); err != nil {
					utils.Log.Errorf("func reFulfillOrderToOfficialMerchants, send order failed: %v", err)
				}
			}

			// 等待官方币商接单
			officialMerchantAcceptWheel.Add(order.OrderNumber)

			utils.RedisIncreaseRefulfillTimesToOfficialMerchants(order.OrderNumber)
			return
		}

		utils.Log.Infof("func reFulfillOrderToOfficialMerchants, reach max trytimes %d", officialMerchantRetries)
		utils.Log.Infof("func reFulfillOrderToOfficialMerchants, order % not accepted by any official merchants, try change it to status 8 (ACCEPTTIMEOUT)", order.OrderNumber)
		// 超过了最大次数限制，修改订单为AcceptTimeout
		if err := utils.DB.Model(models.Order{}).Where("order_number = ? AND status < ?", order.OrderNumber, models.ACCEPTED).
			Update("status", models.ACCEPTTIMEOUT).Error; err != nil {
			utils.Log.Errorf("func reFulfillOrderToOfficialMerchants, update order %s to status SUSPENDED failed", order.OrderNumber)
			return
		}
	}
}

func reFulfillOrder(order *OrderToFulfill, seq uint8) {
	utils.Log.Infof("func reFulfillOrder begin. order_number = %s, seq = %d", order.OrderNumber, seq)

	time.Sleep(time.Duration(retryTimeout) * time.Second)
	//re-fulfill
	merchants := engine.selectMerchantsToFulfillOrder(order)
	utils.Log.Debugf("re-fulfill for order %s, candidate merchants: [%v]", order.OrderNumber, merchants)
	if len(*merchants) > 0 {
		//send order to pick
		if err := sendOrder(order, merchants); err != nil {
			utils.Log.Errorf("Send order failed: %v", err)
		}
		// 等待币商接单
		if order.AcceptType == 0 {
			wheel.Add(order.OrderNumber)
		} else {
			// 自动订单的接单超时时间比一般订单的接单超时时间更短
			autoOrderAcceptWheel.Add(order.OrderNumber)
		}

		utils.Log.Debugf("func reFulfillOrder finished normally. order_number = %s", order.OrderNumber)
		return
	}

	utils.Log.Warnf("func reFulfillOrder, no merchant is available at moment, re-fulfill order %s later.", order.OrderNumber)

	// 没找到合适币商，且少于重派次数，接着重派
	if seq <= uint8(retries) {
		go reFulfillOrder(order, seq+1)
		return
	}

	utils.Log.Warnf("func reFulfillOrder, order %s reach max fulfill times [%d].", order.OrderNumber, retries)

	// 用户提现订单，多次都没人接单，派单给具有“官方客服”身份的币商，以尽最大努力完成订单
	if order.Direction == 1 {
		go reFulfillOrderToOfficialMerchants(order)
		return
	}

	tx := utils.DB.Begin()
	if tx.Error != nil {
		utils.Log.Debugf("tx in func reFulfillOrder begin fail, tx=[%v]", tx)
		utils.Log.Errorf("func reFulfillOrder finished abnormally.")
		return
	}
	utils.Log.Debugf("tx in func reFulfillOrder begin, tx=[%v]", tx)

	//failed, highlight the order to set status to "ACCEPTTIMEOUT"
	suspendedOrder := models.Order{}
	if tx.Set("gorm:query_option", "FOR UPDATE").Find(&suspendedOrder, "order_number = ?  AND status < ?", order.OrderNumber, models.ACCEPTED).RecordNotFound() {
		utils.Log.Errorf("Unable to find order %s", order.OrderNumber)
	} else {
		if suspendedOrder.Direction == 0 { // 平台用户充值，找不到币商时，把订单改为ACCEPTTIMEOUT，这个订单不会再处理
			// 通知h5，没币商接单
			h5 := []string{order.OrderNumber}
			if err := NotifyThroughWebSocketTrigger(models.AcceptTimeout, &[]int64{}, &h5, 0, []OrderToFulfill{*order}); err != nil {
				utils.Log.Errorf("Notify accept timeout through websocket fail [%s]", err)
			}

			if err := tx.Model(&models.Order{}).Where("order_number = ? AND status < ?", order.OrderNumber, models.ACCEPTED).Update("status", models.ACCEPTTIMEOUT).Error; err != nil {
				utils.Log.Errorf("Update order %s status to SUSPENDED failed", order.OrderNumber)
				utils.Log.Errorf("tx in func reFulfillOrder rollback, tx=[%v]", tx)
				tx.Rollback()
				return
			}

			utils.Log.Debugf("call AsynchronousNotifyDistributor for %s, order status is 8 (ACCEPTTIMEOUT)", order.OrderNumber)
			AsynchronousNotifyDistributor(suspendedOrder)

		} else if suspendedOrder.Direction == 1 { // 平台用户提现，找不到币商时，把订单改为SUSPENDED，以后再处理
			if err := tx.Model(&models.Order{}).Where("order_number = ? AND status < ?", order.OrderNumber, models.ACCEPTED).Update("status", models.ACCEPTTIMEOUT).Error; err != nil {
				utils.Log.Errorf("Update order %s status to SUSPENDED failed", order.OrderNumber)
				utils.Log.Errorf("tx in func reFulfillOrder rollback, tx=[%v]", tx)
				tx.Rollback()
				return
			}

		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("error tx in func reFulfillOrder commit, err=[%v]", err)
	}
	utils.Log.Debugf("tx in func reFulfillOrder commit, tx=[%v]", tx)
	return
}

func selectedMerchantsToRedis(orderNumber string, timeout int64, merchants *[]int64) {
	utils.Log.Debugf("selectedMerchantsToRedis orderNumber:[%s],timeout:[%d],merchants:[%v]", orderNumber, timeout, *merchants)
	key := utils.RedisKeyMerchantSelected(orderNumber)
	var temp []interface{}
	for _, v := range *merchants {
		temp = append(temp, v)
	}
	if err := utils.SetCacheSetMember(key, int(2*timeout), temp...); err != nil {
		utils.Log.Warnf("order %v", orderNumber)
	}
}

func sendOrder(order *OrderToFulfill, merchants *[]int64) error {
	utils.Log.Debugf("func sendOrder begin, order = [%+v], merchants = [%+v]", order, merchants)
	timeoutStr := utils.Config.GetString("fulfillment.timeout.accept")
	timeout, _ := strconv.ParseInt(timeoutStr, 10, 32)
	h5 := []string{order.OrderNumber}

	// Android App接受微信收款方式的自动订单时，需要生成收款二维码并返回，这里把生成二维码所需要的备注提供给Android App
	if order.AcceptType == 1 && order.PayType == models.PaymentTypeWeixin {
		order.QrCodeMark = utils.GenQrCodeMark(order.OrderNumber)
	}

	if err := NotifyThroughWebSocketTrigger(models.SendOrder, merchants, &h5, uint(timeout), []OrderToFulfill{*order}); err != nil {
		utils.Log.Errorf("Send order through websocket trigger API failed: %v", err)
		utils.Log.Debugf("func sendOrder finished abnormally.")
		return err
	}

	// 把发送过订单的币商保存到redis中，某个币商抢到订单后，会通知其它币商
	timeout = awaitTimeout + retries*retryTimeout + awaitTimeout
	selectedMerchantsToRedis(order.OrderNumber, timeout, merchants)
	utils.Log.Debugf("func sendOrder finished normally. order_number = %s", order.OrderNumber)
	return nil
}

// 币商点击"抢单"后，下面函数会被调用
func acceptOrder(queue string, args ...interface{}) error {
	//book keeping of all merchants who accept the order
	//recover OrderToFulfill and merchants ID map from args
	utils.Log.Debugf("func acceptOrder begin")
	var order OrderToFulfill
	if orderArg, ok := args[0].(map[string]interface{}); ok {
		order = getOrderToFulfillFromMapStrings(orderArg)
	} else {
		return fmt.Errorf("Wrong order arg: %v", args[0])
	}
	var merchantID int64
	if mid, ok := args[1].(json.Number); ok {
		merchantID, _ = mid.Int64()
	} else {
		return fmt.Errorf("Wrong merchant IDs: %v", args[1])
	}
	var fulfillment *OrderFulfillment
	var err error
	if fulfillment, err = FulfillOrderByMerchant(order, merchantID, 0); err != nil {
		if err.Error() == "already accepted by others" {
			wheel.Remove(order.OrderNumber)
			autoOrderAcceptWheel.Remove(order.OrderNumber)
			officialMerchantAcceptWheel.Remove(order.OrderNumber) // 已经被其它官方币商接单，不再派单了
			utils.RedisDelRefulfillTimesToOfficialMerchants(order.OrderNumber)
			return nil
		}

		return fmt.Errorf("Unable to connect order with merchant: %v", err)
	}

	if order.AcceptType == 0 {
		utils.Log.Debugf("merchant %d accept order %s success", merchantID, order.OrderNumber)
	} else if order.AcceptType == 1 {
		utils.Log.Debugf("merchant %d accept auto order %s success", merchantID, order.OrderNumber)
	}

	// 更新币商接单时间（这个时间会影响币商的下次派单优先级）
	if err := utils.UpdateMerchantLastOrderTime(merchantID, order.Direction, time.Now()); err != nil {
		utils.Log.Warnf("func acceptOrder call UpdateMerchantLastOrderTime fail [%+v].", err)
	}

	notifyFulfillment(fulfillment)

	// 发短信通币商，抢单成功
	//if err := SendSmsOrderAccepted(merchantID, order.OrderNumber); err != nil {
	//	utils.Log.Errorf("order [%v] is accepted by merchant [%v], send sms fail. error [%v]", order.OrderNumber, merchantID, err)
	//}

	wheel.Remove(order.OrderNumber)
	autoOrderAcceptWheel.Remove(order.OrderNumber)
	officialMerchantAcceptWheel.Remove(order.OrderNumber) // 已经被其它官方币商接单，不在派单了
	utils.RedisDelRefulfillTimesToOfficialMerchants(order.OrderNumber)

	//币商已接单,推送其他没有接单的人说接单失败
	data := []OrderToFulfill{{
		OrderNumber: order.OrderNumber,
	}}
	var selectedMerchants []int64
	if data, err := utils.GetCacheSetMembers(utils.RedisKeyMerchantSelected(order.OrderNumber)); err != nil {
		utils.Log.Errorf("func accept selectMerchantsToFulfillOrder error, the select order = [%+v]", order)
	} else if len(data) > 0 {
		utils.Log.Debugf("order %s had sent to merchants [%v] before, only %d accepted, send others picked.", order.OrderNumber, selectedMerchants, merchantID)
		utils.ConvertStringToInt(data, &selectedMerchants)
	}
	//未抢到订单的币商
	notAccept := utils.RemoveElement(selectedMerchants, merchantID)
	utils.Log.Debugf("send pick msg to not accept merchant=[%v]", notAccept)
	if err := NotifyThroughWebSocketTrigger(models.Picked, &notAccept, &[]string{}, 60, data); err != nil {
		utils.Log.Errorf("Notify Picked through websocket ")
	}

	utils.Log.Debugf("func acceptOrder finished normally. order_number = %s", order.OrderNumber)
	return nil
}

func notifyFulfillment(fulfillment *OrderFulfillment) error {
	utils.Log.Debugf("func notifyFulfillment, arg fulfillment = %+v", fulfillment)

	merchantID := fulfillment.MerchantID
	orderNumber := fulfillment.OrderNumber
	timeoutStr := utils.Config.GetString("fulfillment.timeout.notifypaid")
	timeout, _ := strconv.ParseInt(timeoutStr, 10, 32)
	utils.Log.Debugf("notifyFulfillment start, merchantID %v", merchantID)
	if err := NotifyThroughWebSocketTrigger(models.FulfillOrder, &[]int64{merchantID}, &[]string{orderNumber}, uint(timeout), []OrderFulfillment{*fulfillment}); err != nil {
		utils.Log.Errorf("Send fulfillment through websocket trigger API failed: %v", err)
		return err
	}
	//if notifyWheel == nil {
	//	//notify paid timeout
	//	timeoutStr := utils.Config.GetString("fulfillment.timeout.notifypaid")
	//	timeout, _ := strconv.ParseInt(timeoutStr, 10, 0)
	//	notifyWheel = timewheel.New(1*time.Second, int(timeout), notifyPaidTimeout) //process wheel per second
	//	notifyWheel.Start()
	//}
	notifyWheel.Add(fulfillment.OrderNumber)
	return nil
}

var msgTypes = map[string]models.MsgType{
	"send_order":        models.SendOrder,
	"accept":            models.Accept,
	"fulfill_order":     models.FulfillOrder,
	"notify_paid":       models.NotifyPaid,
	"confirm_paid":      models.ConfirmPaid,
	"transferred":       models.Transferred,
	"auto_confirm_paid": models.AutoConfirmPaid,
}

func getMessageFromMapStrings(values map[string]interface{}) models.Msg {
	var result models.Msg
	msgType := msgTypes[values["msg_type"].(string)]
	//get Merchant id list
	var merchants []int64
	if ms, ok := values["merchant_id"].([]interface{}); ok {
		for _, mid := range ms {
			if number, ok := mid.(json.Number); ok {
				n64, _ := number.Int64()
				merchants = append(merchants, n64)
			}
		}
	}
	//get H5 string array
	var h5 []string
	if h5s, ok := values["h5"].([]interface{}); ok {
		for _, h5c := range h5s {
			h5 = append(h5, h5c.(string))
		}
	}
	//get timeout
	var timeout int
	if tn, ok := values["timeout"].(json.Number); ok {
		if t64, err := tn.Int64(); err == nil {
			timeout = int(t64)
		} else {
			utils.Log.Errorf("Error timeout in args: %v\n", err)
		}
	}
	if data, ok := values["data"].([]interface{}); ok {
		result = models.Msg{
			MsgType:    msgType,
			MerchantId: merchants,
			H5:         h5,
			Timeout:    timeout,
			Data:       data,
		}
	} else {
		utils.Log.Errorf("Error parsing websocket message")
		result = models.Msg{}
	}
	return result
}

func updateFulfillment(queue string, args ...interface{}) error {
	//according to different operation + data, update order/fulfillment accordingly.
	//in additon, send notification to impacted partner of the operation
	var msg models.Msg
	if msgArg, ok := args[0].(map[string]interface{}); ok {
		msg = getMessageFromMapStrings(msgArg)
	} else {
		return fmt.Errorf("Wrong format of Msg arg: %v", args[0])
	}
	//switch to different condition
	switch msg.MsgType {
	case models.NotifyPaid:
		if orderNum, err := uponNotifyPaid(msg); err != nil {
			utils.Log.Errorf("uponNotifyPaid is failed,orderNumber:%s", orderNum)
			//suspendedWheel.Add(orderNum)
		}
	case models.ConfirmPaid:
		if orderNum, err := uponConfirmPaid(msg); err != nil {
			utils.Log.Errorf("uponConfirmPaid is failed,orderNumber:%s", orderNum)
			//suspendedWheel.Add(orderNum)
		}
	case models.Transferred:
		utils.Log.Warnf("msg with type Transferred should not occur in redis queue, it processed directly after confirm paid")
	case models.AutoConfirmPaid:
		uponAutoConfirmPaid(msg)
	}
	return nil
}

func deleteWheel(queue string, args ...interface{}) error {
	utils.Log.Debugf("func deleteWheel begin,order:%v", args)
	orderNumber := args[0].(string)
	wheel.Remove(orderNumber)
	autoOrderAcceptWheel.Remove(orderNumber)
	officialMerchantAcceptWheel.Remove(orderNumber)
	notifyWheel.Remove(orderNumber)
	confirmWheel.Remove(orderNumber)
	transferWheel.Remove(orderNumber)
	unfreezeWheel.Remove(orderNumber)
	utils.Log.Debugf("func deleteWheel end.")
	return nil
}

func getOrderNumberAndDirectionFromMessage(msg models.Msg) (orderNumber string, direction int) {
	//get order number from msg.data.order_number
	if d, ok := msg.Data[0].(map[string]interface{}); ok {
		orderNumber = d["order_number"].(string)
		if dn, ok := d["direction"].(json.Number); ok {
			d64, _ := dn.Int64()
			direction = int(d64)
		}
	}
	return orderNumber, direction
}

func uponNotifyPaid(msg models.Msg) (string, error) {
	//update order-fulfillment information
	ordNum, direction := getOrderNumberAndDirectionFromMessage(msg)

	tx := utils.DB.Begin()
	if tx.Error != nil {
		utils.Log.Debugf("tx in func uponNotifyPaid begin fail, tx=[%v]", tx)
		utils.Log.Errorf("func uponNotifyPaid finished abnormally.")
		return ordNum, tx.Error
	}
	utils.Log.Debugf("tx in func uponNotifyPaid begin, tx=[%v]", tx)

	//Trader buy, update order status, fulfillment
	order := models.Order{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&order, "order_number = ? and status < ?", ordNum, models.NOTIFYPAID).RecordNotFound() {
		utils.Log.Errorf("Record not found: order with number %s.", ordNum)
		utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
		tx.Rollback()
		return ordNum, errors.New("record not found")
	}
	originStatus := order.Status

	if originStatus != models.ACCEPTED {
		tx.Rollback()
		utils.Log.Errorf("uponNotifyPaid order status is error,orderNumber:%s,status:%d", ordNum, originStatus)
		utils.Log.Errorf("func uponNotifyPaid finished abnormally.")
		return ordNum, nil
	}

	fulfillment := models.Fulfillment{}
	if tx.Set("gorm:query_option", "FOR UPDATE").Where("order_number = ?", ordNum).Order("seq_id DESC").First(&fulfillment).RecordNotFound() {
		utils.Log.Errorf("No fulfillment with order number %s found.", ordNum)
		utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
		tx.Rollback()
		return ordNum, errors.New("record not found")
	}

	// check current status
	if fulfillment.Status == models.NOTIFYPAID {
		tx.Rollback()
		utils.Log.Errorf("order number %s is already with status %d (NOTIFYPAID), do nothing.", ordNum, models.NOTIFYPAID)
		utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponNotifyPaid finished abnormally.")
		return ordNum, errors.New("uponNotifyPaid fulfillment order status is = notifypaid")
	} else if fulfillment.Status == models.CONFIRMPAID {
		tx.Rollback()
		utils.Log.Errorf("order number %s has status %d (CONFIRMPAID), cannot change it to %d (NOTIFYPAID)", ordNum, models.CONFIRMPAID)
		utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponNotifyPaid finished abnormally.")
		return ordNum, errors.New("uponNotifyPaid fulfillment order status is = confirmpaid")
	} else if fulfillment.Status == models.TRANSFERRED {
		tx.Rollback()
		utils.Log.Errorf("order number %s has status %d (TRANSFERRED), cannot change it to %d (NOTIFYPAID).", ordNum, models.TRANSFERRED, models.NOTIFYPAID)
		utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponNotifyPaid finished abnormally.")
		return ordNum, errors.New("uponNotifyPaid fulfillment order status is = transfered")
	}

	//update order
	if direction == 0 {
		if err := tx.Model(&order).Update("status", models.NOTIFYPAID).Error; err != nil {
			utils.Log.Errorf("Can't update order %s status to %s. %v", ordNum, "NOTIFYPAID", err)
			utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
			tx.Rollback()
			return ordNum, err
		}
	} else {
		if err := tx.Model(&order).Updates(models.Order{Status: models.NOTIFYPAID, BTUSDFlowStatus: models.BTUSDFlowD1TraderFrozenToMerchantFrozen}).Error; err != nil {
			utils.Log.Errorf("Can't update order %s status to %s. %v", ordNum, "NOTIFYPAID", err)
			utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
			tx.Rollback()
			return ordNum, err
		}
	}
	if err := tx.Model(&fulfillment).Updates(models.Fulfillment{Status: models.NOTIFYPAID, PaidAt: time.Now()}).Error; err != nil {
		utils.Log.Errorf("Can't update order %s fulfillment info. %v", ordNum, err)
		utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
		tx.Rollback()
		return ordNum, err
	}
	//update fulfillment with one new log
	fulfillmentLog := models.FulfillmentLog{
		FulfillmentID: fulfillment.Id,
		OrderNumber:   ordNum,
		SeqID:         fulfillment.SeqID,
		IsSystem:      false,
		MerchantID:    fulfillment.MerchantID,
		AccountID:     order.AccountId,
		DistributorID: order.DistributorId,
		OriginStatus:  originStatus,
		UpdatedStatus: models.NOTIFYPAID,
	}
	if err := tx.Create(&fulfillmentLog).Error; err != nil {
		utils.Log.Errorf("Can't create order %s fulfillment log. %v", ordNum, err)
		utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
		tx.Rollback()
		return ordNum, err
	}

	if direction == 0 {
		if err := tx.Commit().Error; err != nil {
			utils.Log.Errorf("error tx in func uponNotifyPaid commit, err=[%v]", err)
			return ordNum, err
		}
		utils.Log.Debugf("tx in func uponNotifyPaid commit, tx=[%v]", tx)

		timeoutStr := utils.Config.GetString("fulfillment.timeout.notifypaymentconfirmed")
		timeout, _ := strconv.ParseInt(timeoutStr, 10, 32)
		//then notify partner the same message - only direction = 0, Trader Buy
		if err := NotifyThroughWebSocketTrigger(models.NotifyPaid, &msg.MerchantId, &msg.H5, uint(timeout), msg.Data); err != nil {
			utils.Log.Errorf("Notify partner notify paid messaged failed.")
		}
		confirmWheel.Add(order.OrderNumber)

		if err := SendSmsOrderPaid(fulfillment.MerchantID, ordNum); err != nil {
			utils.Log.Errorf("order [%v] is marked as paid by user, send sms to merchant [%v] fail. error [%v]", ordNum, fulfillment.MerchantID, err)
		}
	} else { //Trader Sell
		// 币商点击"我已完成付款"，进行下面操作：
		// 增加部分币到币商的冻结账号中
		// 增加部分币到金融滴滴平台的冻结账号中
		// 把币从平台的冻结账号中扣除

		// 找到平台asset记录
		assetForDist := models.Assets{}
		if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForDist, "distributor_id = ? AND currency_crypto = ? ", order.DistributorId, order.CurrencyCrypto).RecordNotFound() {
			tx.Rollback()
			utils.Log.Errorf("Can't find corresponding asset record of distributor_id %d, currency_crypto %s", order.DistributorId, order.CurrencyCrypto)
			utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
			utils.Log.Errorf("func uponNotifyPaid finished abnormally.")
			return ordNum, errors.New("distributor asset uponNotifyPaid record not found")
		}

		// 找到币商asset记录
		asset := models.Assets{}
		if tx.Set("gorm:query_option", "FOR UPDATE").First(&asset, "merchant_id = ? AND currency_crypto = ? ", order.MerchantId, order.CurrencyCrypto).RecordNotFound() {
			tx.Rollback()
			utils.Log.Errorf("Can't find corresponding asset record of merchant_id %d, currency_crypto %s", order.MerchantId, order.CurrencyCrypto)
			utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
			utils.Log.Errorf("func uponNotifyPaid finished abnormally.")
			return ordNum, errors.New("merchant asset uponNotifyPaid record not found")
		}

		// 找到金融滴滴平台asset记录
		assetForPlatform := models.Assets{}
		platformDistId := 1 // 金融滴滴平台的distributor_id为1
		if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForPlatform, "distributor_id = ? AND currency_crypto = ? ",
			platformDistId, order.CurrencyCrypto).RecordNotFound() {
			tx.Rollback()
			utils.Log.Errorf("Can't find corresponding asset record for platform, currency_crypto %s", order.MerchantId, order.CurrencyCrypto)
			utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
			utils.Log.Errorf("func uponNotifyPaid finished abnormally.")
			return ordNum, errors.New("assetForPlatform uponNotifyPaid record not found")
		}

		if err := TransferCoinFromTraderFrozenToMerchantFrozen(tx, &assetForDist, &asset, &assetForPlatform, &order); err != nil {
			tx.Rollback()
			utils.Log.Errorf("func TransferCoinFromTraderFrozenToMerchantFrozen fail %v", err)
			utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
			utils.Log.Errorf("func uponNotifyPaid finished abnormally.")
			return ordNum, errors.New("TransferCoinFromTraderFrozenToMerchantFrozen fail" + err.Error())
		}

		utils.Log.Debugf("tx in func uponNotifyPaid commit, tx=[%v]", tx)
		if err := tx.Commit().Error; err != nil {
			utils.Log.Errorf("error tx in func uponNotifyPaid commit, err=[%v]", err)
			return ordNum, errors.New("uponNotifyPaid tx commit is failed")
		}

		message := models.Msg{
			MsgType:    models.ConfirmPaid,
			MerchantId: msg.MerchantId,
			H5:         msg.H5,
			Timeout:    0,
			Data: []interface{}{
				map[string]interface{}{
					"order_number": order.OrderNumber,
					"direction":    1,
				},
			},
		}
		//as if we got confirm paid message from APP
		if _, err := uponConfirmPaid(message); err != nil {
			utils.Log.Errorf("uponNotifyPaid to uponConfirmPaid is failed,OrderNumber:%s", ordNum)
			//超时更新失败，修改订单状态suspended
			suspendedWheel.Add(order.OrderNumber)
		}
	}
	notifyWheel.Remove(order.OrderNumber)
	return ordNum, nil
}

// 下面函数当确认"对方已付款"时，会被调用。
func uponConfirmPaid(msg models.Msg) (string, error) {
	utils.Log.Debugf("func uponConfirmPaid begin, msg = [%+v]", msg)
	ordNum, _ := getOrderNumberAndDirectionFromMessage(msg)

	tx := utils.DB.Begin()
	if tx.Error != nil {
		utils.Log.Debugf("tx in func uponConfirmPaid begin fail, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return ordNum, tx.Error
	}
	utils.Log.Debugf("tx in func uponConfirmPaid begin, tx=[%v]", tx)

	order := models.Order{}
	if tx.Set("gorm:query_option", "FOR UPDATE").Where("order_number = ?", ordNum).First(&order).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("Record not found: order with number %s.", ordNum)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return ordNum, errors.New("uponConfirmPaid order Record not found,orderNumber:" + ordNum)
	}
	originStatus := order.Status

	//因为充值单app增加了业务逻辑为：只要用户接单就可以点击确认付款，因此增加用户已接单状态可以点击确认收款按钮状态的判断
	if originStatus != models.ACCEPTED && originStatus != models.NOTIFYPAID {
		//如果订单状态是付款超时异常,即status和statusReason为5和2时，允许其点击确认收款
		if !(originStatus == models.SUSPENDED && order.StatusReason == models.PAIDTIMEOUT) {
			tx.Rollback()
			utils.Log.Errorf("uponConfirmPaid order status is error,orderNumber:%s,status:%d", ordNum, originStatus)
			utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
			return ordNum, nil
		}
	}

	fulfillment := models.Fulfillment{}
	if tx.Set("gorm:query_option", "FOR UPDATE").Where("order_number = ?", ordNum).Order("seq_id DESC").First(&fulfillment).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("No fulfillment with order number %s found.", ordNum)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return ordNum, errors.New("uponConfirmPaid fulfillment Record not found,orderNumber:" + ordNum)
	}

	// check current status
	if fulfillment.Status == models.CONFIRMPAID {
		tx.Rollback()
		utils.Log.Errorf("order number %s is already with status %d (CONFIRMPAID), do nothing.", ordNum, models.CONFIRMPAID)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return ordNum, errors.New("uponConfirmPaid fulfillment status confirmPaid is error,orderNumber:" + ordNum)
	} else if fulfillment.Status == models.TRANSFERRED {
		tx.Rollback()
		utils.Log.Errorf("order number %s has status %d (TRANSFERRED), cannot change it to %d (CONFIRMPAID).", ordNum, models.TRANSFERRED, models.CONFIRMPAID)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return ordNum, errors.New("uponConfirmPaid fulfillment status transferred is error,orderNumber:" + ordNum)
	}

	//update fulfillment with one new log
	fulfillmentLog := models.FulfillmentLog{
		FulfillmentID: fulfillment.Id,
		OrderNumber:   ordNum,
		SeqID:         fulfillment.SeqID,
		IsSystem:      false,
		MerchantID:    fulfillment.MerchantID,
		AccountID:     order.AccountId,
		DistributorID: order.DistributorId,
		OriginStatus:  originStatus,
		UpdatedStatus: models.CONFIRMPAID,
	}
	if err := tx.Create(&fulfillmentLog).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Can't create order %s fulfillment log. %v", ordNum, err)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return ordNum, err
	}

	// update order status
	if err := tx.Model(&order).Updates(map[string]interface{}{"status": models.CONFIRMPAID, "status_reason": 0}).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Can't update order %s status to %s. %v", ordNum, "CONFIRMPAID", err)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return ordNum, err
	}
	if err := tx.Model(&fulfillment).Updates(models.Fulfillment{Status: models.CONFIRMPAID, PaymentConfirmedAt: time.Now()}).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Can't update order %s fulfillment info. %v", ordNum, err)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return ordNum, err
	}
	utils.Log.Debugf("tx in func uponConfirmPaid commit, tx=[%v]", tx)
	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("error tx in func uponConfirmPaid commit,orderNumber:%s, err=[%v]", ordNum, err)
		return ordNum, err
	}

	notifyMerchant := []int64{fulfillment.MerchantID}

	// notify partner
	if err := NotifyThroughWebSocketTrigger(models.ConfirmPaid, &notifyMerchant, &[]string{order.OrderNumber}, 0, msg.Data); err != nil {
		utils.Log.Errorf("Notify partner notify paid messaged failed.")
	}

	if order.Direction == 0 {
		if err := doTransfer(ordNum); err != nil {
			//转账失败，修改订单状态suspended
			suspendedWheel.Add(ordNum)
		}
	} else {
		// 等待一定时间后，释放冻结的币
		transferWheel.Add(order.OrderNumber)
	}

	//付款超时的也允许确认收款，要将解冻的时间轮任务移除掉
	unfreezeWheel.Remove(order.OrderNumber)
	confirmWheel.Remove(order.OrderNumber)
	utils.Log.Debugf("func uponConfirmPaid finished normally. order_number = %s", order.OrderNumber)
	return ordNum, nil
}

func doTransfer(ordNum string) error {
	utils.Log.Debugf("func doTransfer begin, OrderNumber = [%+v]", ordNum)

	tx := utils.DB.Begin()
	if tx.Error != nil {
		utils.Log.Debugf("tx in func doTransfer begin fail, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
		return tx.Error
	}
	utils.Log.Debugf("tx in func doTransfer begin, tx=[%v]", tx)

	//Trader buy, update order status, fulfillment
	order := models.Order{}
	if tx.Set("gorm:query_option", "FOR UPDATE").Where("order_number = ?", ordNum).First(&order).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("Record not found: order with number %s.", ordNum)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
		return errors.New("not found order record,orderNumber:" + ordNum)
	}
	originStatus := order.Status

	if originStatus != models.CONFIRMPAID {
		tx.Rollback()
		utils.Log.Errorf("doTransfer order status is error,orderNumber:%s,status=[%v]", ordNum, originStatus)
		utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
		return nil
	}

	fulfillment := models.Fulfillment{}
	if tx.Set("gorm:query_option", "FOR UPDATE").Where("order_number = ?", ordNum).Order("seq_id DESC").First(&fulfillment).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("No fulfillment with order number %s found.", ordNum)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
		return errors.New("not found fulfillment record,orderNumber:" + ordNum)
	}

	if fulfillment.Status == models.TRANSFERRED {
		tx.Rollback()
		utils.Log.Errorf("order number %s is already with status %d (TRANSFERRED), do nothing.", ordNum, models.TRANSFERRED)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
		return errors.New(" fulfillment order status is equal transferred,orderNumber:" + ordNum)
	}

	//update fulfillment with one new log
	fulfillmentLog := models.FulfillmentLog{
		FulfillmentID: fulfillment.Id,
		OrderNumber:   ordNum,
		SeqID:         fulfillment.SeqID,
		IsSystem:      false,
		MerchantID:    fulfillment.MerchantID,
		AccountID:     order.AccountId,
		DistributorID: order.DistributorId,
		OriginStatus:  originStatus,
		UpdatedStatus: models.TRANSFERRED,
	}
	if err := tx.Create(&fulfillmentLog).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Can't create order %s fulfillment log. %v", ordNum, err)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
		return err
	}

	//update order
	if err := tx.Model(&order).Updates(map[string]interface{}{"status": models.TRANSFERRED, "merchant_payment_id": fulfillment.MerchantPaymentID}).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Can't update order %s status to %s. %v", ordNum, "TRANSFERRED", err)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
		return err
	}
	transferredAt := time.Now()
	if err := tx.Model(&fulfillment).Updates(models.Fulfillment{Status: models.TRANSFERRED, TransferredAt: transferredAt}).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Can't update order %s fulfillment info. %v", ordNum, err)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
		return err
	}

	// 找到平台商记录
	assetForDist := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForDist, "distributor_id = ? AND currency_crypto = ? ", order.DistributorId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("Can't find corresponding asset record of distributor_id %d, currency_crypto %s", order.DistributorId, order.CurrencyCrypto)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
		return errors.New(fmt.Sprintf("not found distributor asset,orderNuber:%s,distributorId:%d", ordNum, order.DistributorId))
	}

	asset := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&asset, "merchant_id = ? AND currency_crypto = ? ", order.MerchantId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("Can't find corresponding asset record of merchant_id %d, currency_crypto %s", order.MerchantId, order.CurrencyCrypto)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
		return errors.New(fmt.Sprintf("not found merchant asset,orderNumber:%s,merchantId:%d", ordNum, order.MerchantId))
	}

	if order.Direction == 0 {
		// Trader Buy
		utils.Log.Debugf("Freeze [%v] %v for merchant (uid=[%v])", order.Quantity, order.CurrencyCrypto, fulfillment.MerchantPaymentID)
		if asset.QtyFrozen.GreaterThanOrEqual(order.Quantity) { // 避免 qty_frozen 出现负数
			if err := tx.Table("assets").Where("id = ?", asset.Id).
				Update("qty_frozen", asset.QtyFrozen.Sub(order.Quantity)).Error; err != nil {
				utils.Log.Errorf("update asset record for merchant fail, order_number = %s", order.OrderNumber)
				tx.Rollback()
				utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
				return errors.New("update asset record for merchant fail")
			}
		} else {
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("Can't deduct %s %s frozen asset for merchant (uid=[%v]). asset for merchant = [%+v], order_number = %s",
				order.Quantity, order.CurrencyCrypto, asset.MerchantId, asset, order.OrderNumber)
			utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
			tx.Rollback()
			return errors.New(fmt.Sprintf("Can't deduct %s %s frozen asset for merchant (uid=[%v])", order.Quantity, order.CurrencyCrypto, asset.MerchantId))
		}

		// 转币给平台商
		if err := tx.Table("assets").Where("id = ? ", assetForDist.Id).
			Update("quantity", assetForDist.Quantity.Add(order.Quantity)).Error; err != nil {
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("Can't transfer asset to distributor (distributor_id=[%v]). err: %v", assetForDist.DistributorId, err)
			utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
			tx.Rollback()
			return err
		}

		if err := tx.Table("payment_infos").Where("id = ?", fulfillment.MerchantPaymentID).Update("in_use", 0).Error; err != nil {
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("Can't change in_use to 0, record id=[%v], err=[%v]", fulfillment.MerchantPaymentID, err)
			utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
			tx.Rollback()
			return err
		}

		// Add asset history for distributor
		assetHistory := models.AssetHistory{
			Currency:      order.CurrencyCrypto,
			Direction:     order.Direction,
			DistributorId: order.DistributorId,
			Quantity:      order.Quantity,
			IsOrder:       1,
			OrderNumber:   ordNum,
		}
		if err := tx.Model(&models.AssetHistory{}).Create(&assetHistory).Error; err != nil {
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("create asset history for distributor (uid=[%v]) failed. err:[%v]", order.MerchantId, err)
			utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
			tx.Rollback()
			return err
		}

		// Add asset history for merchant
		assetMerchantHistory := models.AssetHistory{
			Currency:    order.CurrencyCrypto,
			Direction:   order.Direction,
			MerchantId:  order.MerchantId,
			Quantity:    order.Quantity.Neg(),
			IsOrder:     1,
			OrderNumber: ordNum,
		}
		if err := tx.Model(&models.AssetHistory{}).Create(&assetMerchantHistory).Error; err != nil {
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("create asset history for merchant (uid=[%v]) failed. err:[%v]", order.MerchantId, err)
			utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
			tx.Rollback()
			return err
		}
	} else {
		// Trader Sell
		utils.Log.Debugf("Add [%v] %v for merchant (uid=[%v])", order.Quantity, order.CurrencyCrypto, fulfillment.MerchantPaymentID)

		// 找到金融滴滴平台记录
		assetForPlatform := models.Assets{}
		platformDistId := 1 // 金融滴滴平台的distributor_id为1
		if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForPlatform, "distributor_id = ? AND currency_crypto = ? ",
			platformDistId, order.CurrencyCrypto).RecordNotFound() {
			tx.Rollback()
			utils.Log.Errorf("Can't find corresponding asset record for platform, currency_crypto %s", order.MerchantId, order.CurrencyCrypto)
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
			return errors.New(fmt.Sprintf("not found jrdidi asset record,orderNumber:%s", ordNum))
		}

		if err := TransferNormally(tx, &assetForDist, &asset, &assetForPlatform, &order, nil); err != nil {
			tx.Rollback()
			utils.Log.Errorf("func TransferNormally fail %v", err)
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
			return errors.New("TransferNormally fail")
		}
	}

	utils.Log.Debugf("tx in func doTransfer commit, tx=[%v]", tx)
	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("error tx in func doTransfer commit, err=[%v]", err)
		utils.Log.Errorf("func doTransfer finished abnormally. order_number = %s", ordNum)
		return err
	}

	utils.Log.Debugf("call AsynchronousNotifyDistributor for %s, order status is 7 (TRANSFERRED)", order.OrderNumber)
	AsynchronousNotifyDistributor(order)

	utils.Log.Debugf("func doTransfer finished normally. order_number = %s", ordNum)
	return nil
}

func getAutoConfirmPaidFromMessage(msg models.Msg) (merchant int64, amount float64) {
	//get merchant, amount, ts from msg.data
	if d, ok := msg.Data[0].(map[string]interface{}); ok {
		mn, ok := d["merchant_id"].(json.Number)
		if ok {
			merchant, _ = mn.Int64()
		}
		an, ok := d["amount"].(json.Number)
		if ok {
			amount, _ = an.Float64()
		}
	}
	return merchant, amount
}

func uponAutoConfirmPaid(msg models.Msg) {
	//check to get merchant_id, amount, timestamp, compare them with all ongoing processing orders to match
	merchantID, amount := getAutoConfirmPaidFromMessage(msg)
	ts := time.Now()
	//record the event
	record := models.AutoConfirmLog{
		Uid:       merchantID,
		Amount:    amount,
		Timestamp: ts,
	}
	if err := utils.DB.Create(&record).Error; err != nil {
		utils.Log.Errorf("Unable to record auto_confirm_paid message: %v", msg)
	}
	order := models.Order{}
	timeoutStr := utils.Config.GetString("fulfillment.timeout.autoconfirmpaid")
	timeout, _ := strconv.ParseInt(timeoutStr, 10, 0)
	if utils.DB.First(&order, "merchant_id = ? and amount = ? and updated_at > ?", merchantID, amount, ts.UTC().Add(time.Duration(-1*timeout)*time.Second).Format("2006-01-02T15:04:05")).RecordNotFound() {
		utils.Log.Debugf("Auto confirm paid information doesn't match any ongoing order.")
		return
	}
	//found, send server_confirm_paid message
	type Data struct {
		OrderNumber string    `json:"order_number"`
		Direction   int       `json:"direction"`
		MerchantID  int64     `json:"merchant_id"`
		Timestamp   time.Time `json:"timestamp"`
	}
	data := Data{
		OrderNumber: order.OrderNumber,
		MerchantID:  merchantID,
		Direction:   order.Direction,
		Timestamp:   ts,
	}
	if err := NotifyThroughWebSocketTrigger(models.ServerConfirmPaid, &msg.MerchantId, &msg.H5, 0, data); err != nil {
		utils.Log.Errorf("Notify merchant server_confirm_paid messaged failed.")
	}
	message := models.Msg{
		MsgType:    models.ConfirmPaid,
		MerchantId: msg.MerchantId,
		H5:         msg.H5,
		Timeout:    0,
		Data: []interface{}{
			map[string]interface{}{
				"order_number": data.OrderNumber,
				"direction":    data.Direction,
			},
		},
	}
	//as if we got confirm paid message from APP
	if orderNum, err := uponConfirmPaid(message); err != nil {
		utils.Log.Errorf("autoConfirm to uponConfirmPaid is failed,OrderNumber:%s", orderNum)
		suspendedWheel.Add(orderNum)
	}
}

//RegisterFulfillmentFunctions - register fulfillment functions, called by server
func RegisterFulfillmentFunctions() {
	//register worker function
	utils.RegisterWorkerFunc(utils.FulfillOrderTask, fulfillOrder)
	utils.RegisterWorkerFunc(utils.AcceptOrderTask, acceptOrder)
	utils.RegisterWorkerFunc(utils.UpdateFulfillmentTask, updateFulfillment)
	utils.RegisterWorkerFunc(utils.DeleteWheel, deleteWheel)
}

func InitWheel() {
	timeoutStr := utils.Config.GetString("fulfillment.timeout.awaitaccept")
	key := utils.UniqueTimeWheelKey("awaitaccept")
	awaitTimeout, _ = strconv.ParseInt(timeoutStr, 10, 64)
	utils.Log.Debugf("wheel init,timeout:%d", awaitTimeout)
	wheel = timewheel.New(1*time.Second, int(awaitTimeout), key, waitAcceptTimeout) //process wheel per second
	wheel.Start()

	timeoutStr = utils.Config.GetString("fulfillment.timeout.awaitautoorderaccept")
	key = utils.UniqueTimeWheelKey("awaitautoorderaccept")
	awaitTimeout, _ = strconv.ParseInt(timeoutStr, 10, 64)
	utils.Log.Debugf("autoOrderAcceptWheel init,timeout:%d", awaitTimeout)
	autoOrderAcceptWheel = timewheel.New(1*time.Second, int(awaitTimeout), key, waitAcceptTimeout)
	autoOrderAcceptWheel.Start()

	key = utils.UniqueTimeWheelKey("awaitacceptofficialmerchant")
	utils.Log.Debugf("officialMerchantAcceptWheel init,timeout:%d", awaitTimeout)
	officialMerchantAcceptWheel = timewheel.New(1*time.Second, int(awaitTimeout), key, waitOfficialMerchantAcceptTimeout)
	officialMerchantAcceptWheel.Start()

	timeoutStr = utils.Config.GetString("fulfillment.timeout.notifypaid")
	timeout, _ := strconv.ParseInt(timeoutStr, 10, 64)
	key = utils.UniqueTimeWheelKey("notifypaid")
	utils.Log.Debugf("notify wheel init,timeout:%d", timeout)
	notifyWheel = timewheel.New(1*time.Second, int(timeout), key, notifyPaidTimeout) //process wheel per second
	notifyWheel.Start()

	//confirm paid timeout
	timeoutStr = utils.Config.GetString("fulfillment.timeout.notifypaymentconfirmed")
	timeout, _ = strconv.ParseInt(timeoutStr, 10, 64)
	key = utils.UniqueTimeWheelKey("confirmed")
	utils.Log.Debugf("confirm wheel init,timeout:%d", timeout)
	confirmWheel = timewheel.New(1*time.Second, int(timeout), key, confirmPaidTimeout) //process wheel per second
	confirmWheel.Start()

	timeoutStr = utils.Config.GetString("fulfillment.timeout.transfer")
	timeout, _ = strconv.ParseInt(timeoutStr, 10, 64)
	key = utils.UniqueTimeWheelKey("transfer")
	utils.Log.Debugf("transfer wheel init,timeout:%d", timeout)
	transferWheel = timewheel.New(1*time.Second, int(timeout), key, transferTimeout) //process wheel per second
	transferWheel.Start()

	//update order suspended retry time
	timeout = utils.Config.GetInt64("fulfillment.timeout.retrytime")
	key = utils.UniqueTimeWheelKey("retrytime")
	utils.Log.Debugf("suspendedWheel wheel init,timeout:%d", timeout)
	suspendedWheel = timewheel.New(1*time.Second, int(timeout), key, updateOrderStatusAsSuspended) //process wheel per second
	suspendedWheel.Start()

	timeout = utils.Config.GetInt64("fulfillment.timeout.autounfreeze")
	key = utils.UniqueTimeWheelKey("autounfreeze")
	utils.Log.Debugf("suspendedWheel wheel init,timeout:%d", timeout)
	unfreezeWheel = timewheel.New(1*time.Second, int(timeout), key, autoUnfreeze) //process wheel per second
	unfreezeWheel.Start()

	timeoutStr = utils.Config.GetString("fulfillment.timeout.retry")
	retryTimeout, _ = strconv.ParseInt(timeoutStr, 10, 64)
	utils.Log.Debugf("retry timeout:%d", retryTimeout)

	retryStr := utils.Config.GetString("fulfillment.retries")
	retries, _ = strconv.ParseInt(retryStr, 10, 64)
	utils.Log.Debugf("retries:%d", retries)

	officialMerchantRetriesStr := utils.Config.GetString("fulfillment.officialmerchant.retries")
	officialMerchantRetries, _ = strconv.ParseInt(officialMerchantRetriesStr, 10, 64)
	utils.Log.Debugf("officialMerchantRetries:%d", officialMerchantRetries)

	forbidNewOrderIfUnfinishedStr := utils.Config.GetString("fulfillment.forbidneworderifunfinished")
	var err error
	if forbidNewOrderIfUnfinished, err = strconv.ParseBool(forbidNewOrderIfUnfinishedStr); err != nil {
		utils.Log.Errorf("Wrong configuration: fulfillment.forbidneworderifunfinished, should be boolean. Set to default true.")
		forbidNewOrderIfUnfinished = true
	}
	utils.Log.Debugf("forbidNewOrderIfUnfinished:%s", forbidNewOrderIfUnfinished)
}
