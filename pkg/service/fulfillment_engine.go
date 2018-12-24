package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

var engine *defaultEngine

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
	Quantity float64 `json:"quantity"`
	//Price - rate between crypto and fiat
	Price float64 `json:"price"`
	//Amount of the order, in fiat currency
	Amount float64 `json:"amount"`
	//Payment type, chosen by trader
	PayType int `json:"pay_type"`
	//微信或支付宝二维码地址
	QrCode string `gorm:"type:varchar(255)" json:"qr_code"`
	//微信或支付宝账号
	Name string `gorm:"type:varchar(100)" json:"name"`
	//银行账号
	BankAccount string `gorm:"" json:"bank_account"`
	//所属银行
	Bank string `gorm:"" json:"bank"`
	//所属银行分行
	BankBranch string `gorm:"" json:"bank_branch"`
}

func getOrderToFulfillFromMapStrings(values map[string]interface{}) OrderToFulfill {
	var distributorID, direct, payT int64
	var quantity, price, amount float64

	if distributorN, ok := values["distributor"].(json.Number); ok {
		distributorID, _ = distributorN.Int64()
	}
	if directN, ok := values["direction"].(json.Number); ok {
		direct, _ = directN.Int64()
	}
	if payTN, ok := values["pay_type"].(json.Number); ok {
		payT, _ = payTN.Int64()
	}
	if quantityN, ok := values["quantity"].(json.Number); ok {
		quantity, _ = quantityN.Float64()
	}
	if priceN, ok := values["price"].(json.Number); ok {
		price, _ = priceN.Float64()
	}
	if amountN, ok := values["amount"].(json.Number); ok {
		amount, _ = amountN.Float64()
	}
	var qrCode, name, bank, bankAccount, bankBranch string
	if values["qr_code"] != nil {
		qrCode = values["qr_code"].(string)
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

	return OrderToFulfill{
		OrderNumber:    values["order_number"].(string),
		Direction:      int(direct),
		AccountID:      values["account"].(string),
		OriginOrder:    values["origin_order"].(string),
		DistributorID:  distributorID,
		CurrencyCrypto: values["currency_crypto"].(string),
		CurrencyFiat:   values["currency_fiat"].(string),
		Quantity:       quantity,
		Price:          price,
		Amount:         amount,
		PayType:        int(payT),
		QrCode:         qrCode,
		Name:           name,
		Bank:           bank,
		BankAccount:    bankAccount,
		BankBranch:     bankBranch,
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

func getPaymentInfoFromMapStrings(data []interface{}) []models.PaymentInfo {
	//get payment_info firstly, []interface{}, only 1 item included
	if len(data) < 1 {
		utils.Log.Errorf("No data presented in the parameters list")
		return []models.PaymentInfo{}
	}
	datum, ok := data[0].(map[string]interface{})
	if !ok {
		utils.Log.Errorf("Invalid data{} object presented in parameters list")
		return []models.PaymentInfo{}
	}
	result := []models.PaymentInfo{}
	var eAmount float64
	//datum => map[string]interface{}
	var uid, payT int64
	if uidN, ok := datum["uid"].(json.Number); ok {
		uid, _ = uidN.Int64()
	}
	if payTN, ok := datum["pay_type"].(json.Number); ok {
		payT, _ = payTN.Int64()
	}
	if eAmountN, ok := datum["e_amount"].(json.Number); ok {
		eAmount, _ = eAmountN.Float64()
	}
	var pi models.PaymentInfo
	switch payT {
	case 1:
		fallthrough
	case 2:
		pi = models.PaymentInfo{
			Uid:       uid,
			PayType:   int(payT),
			EAccount:  datum["e_account"].(string),
			QrCode:    datum["qr_code"].(string),
			QrCodeTxt: datum["qr_code_txt"].(string),
			EAmount:   eAmount,
		}
	case 4:
		pi = models.PaymentInfo{
			Uid:         uid,
			PayType:     int(payT),
			Name:        datum["name"].(string),
			Bank:        datum["bank"].(string),
			BankAccount: datum["bank_account"].(string),
			BankBranch:  datum["bank_branch"].(string),
		}
	}
	result = append(result, pi)
	return result
}

func getFulfillmentInfoFromMapStrings(values map[string]interface{}) OrderFulfillment {
	var merchantID int64
	if merchantIDN, ok := values["merchant_id"].(json.Number); ok {
		merchantID, _ = merchantIDN.Int64()
	}
	orderToFulfill := getOrderToFulfillFromMapStrings(values)
	data, ok := values["payment_info"].([]interface{})
	if !ok {
		utils.Log.Errorf("Wrong msg.data.payment_info format")
		return OrderFulfillment{}
	}
	paymentInfo := getPaymentInfoFromMapStrings(data)
	return OrderFulfillment{
		OrderToFulfill:    orderToFulfill,
		MerchantID:        merchantID,
		MerchantNickName:  values["merchant_nickname"].(string),
		MerchantAvatarURI: values["merchant_avartar_uri"].(string),
		PaymentInfo:       paymentInfo,
	}
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
	// ReFulfillOrder - Rerun fulfillment logic upon receiving NO "pick-order" response
	// from last round SendOrder. The last round fulfillment options would be stored
	// in the database (maybe also in the cache), with a "sequence" number indicator.
	// Every time of the re-fulfill, the "sequence" number increases.
	ReFulfillOrder(
		orderNumber string, // Order number to be re-fulfilled.
	)
	// SendOrder - notify merchants to accept order.
	// Order is being set at SENT status after SendOrder.
	SendOrder(
		order *OrderToFulfill, // order to be fulfilled
		merchants *[]int64, // a list of merchants ID to pick-up the order
	)
	// AcceptOrder - receive merchants' response on accept order.
	// it then pick up the winner of all responded merchants and call
	// NotifyFulfillment to inform the winner
	AcceptOrder(
		order OrderToFulfill, //order number to accept
		merchantID int64, //merchant id
	)
	// NotifyFulfillment - notify trader/merchant about the fulfillment.
	// Before notification, order is set to ACCEPTED
	NotifyFulfillment(
		fulfillment *OrderFulfillment, //the fulfillment choice decided by engine
	)
	// UpdateFulfillment - update fulfillment processing like payment notified, confirm payment, etc..
	// Upon receiving these message, fulfillment engine should update order/fulfillment status + appended extra message
	UpdateFulfillment(
		msg models.Msg, // Order number
		//operation int, // fulfilment operation such as notify_paid, payment_confirmed, etc..
		//data interface{}, // arbitrary notification data according to different operation
	)
}

// NewOrderFulfillmentEngine - return a new fulfillment engine.
// Factory method of fulfillment engine to adopt to future engine extension.
func NewOrderFulfillmentEngine(_ /*config*/ interface{}) OrderFulfillmentEngine {
	//Singleton, may init engine by config, now ignore it
	if engine == nil {
		utils.SetSettings()
		engine = new(defaultEngine)
	}
	return engine
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

func (engine *defaultEngine) ReFulfillOrder(
	orderNumber string,
) {
	//get corresponding fulfillment object, update it then re-run FulfillOrder
	var lastFulfillment *OrderFulfillment
	lastFulfillment = getFufillmentByOrderNumber(orderNumber)
	engine.FulfillOrder(
		&OrderToFulfill{
			OrderNumber:    lastFulfillment.OrderNumber,
			Direction:      lastFulfillment.Direction,
			CurrencyCrypto: lastFulfillment.CurrencyCrypto,
			CurrencyFiat:   lastFulfillment.CurrencyFiat,
			Quantity:       lastFulfillment.Quantity,
			Price:          lastFulfillment.Price,
			Amount:         lastFulfillment.Amount,
			PayType:        lastFulfillment.PayType,
		})
}

func (engine *defaultEngine) SendOrder(
	order *OrderToFulfill,
	merchants *[]int64,
) {
	//send "order to fulfill" to selected merchants
	utils.AddBackgroundJob(utils.SendOrderTask, utils.NormalPriority, *order, *merchants)
}

func (engine *defaultEngine) selectMerchantsToFulfillOrder(order *OrderToFulfill) *[]int64 {
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
	var merchants []int64
	if order.Direction == 0 {
		//Buy, try to match all-automatic merchants firstly
		// 1. available merchants(online + in_work) + auto accept order/confirm payment + fix amount match
		merchants = GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, true, 0, 0)
		if len(merchants) == 0 { //no priority merchants with fix amount match found, another round call
			// 2. available merchants(online + in_work) + auto accept order/confirm payment
			merchants = GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, false, 0, 0)
			if len(merchants) == 0 { //no priority merchants with non-fix amount match found, then "manual operation" merchants
				// 3. available merchants(online + in_work) + manual accept order/confirm payment
				merchants = GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, true, 1, 0)
			}
		} else { //Sell, all should manually processed
			merchants = GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, false, 1, 0)
		}
	} else {
		//Sell, any online + in_work could pickup order
		merchants = GetMerchantsQualified(0, 0, order.CurrencyCrypto, order.PayType, false, 1, 0)
	}
	return &merchants
}

//GetMerchantsQualified - return mock data
func GetMerchantsQualified(amount, quantity float64, currencyCrypto string, payType int, fix bool, group uint8, limit uint8) []int64 {
	var key string
	var merchantIds []int64
	var assetMerchantIds []int64
	var paymentMerchantIds []int64
	result := []int64{}
	//获取承兑商在线列表
	if group == 0 {
		key = utils.UniqueMerchantOnlineAutoKey()
	} else if group == 1 {
		key = utils.UniqueMerchantOnlineAcceptKey()
	} else {
		return result
	}
	//将承兑商在线列表从string数组转为int64数组
	if tempIds, err := utils.GetCacheSetMembers(key); err != nil {
		return result
	} else {
		if merchantIds, err = convertStringToInt(tempIds); err != nil {
			return result
		}
	}

	//查询资产符合情况的币商列表
	db := utils.DB.Model(&models.Assets{}).Where("currency_crypto = ? AND quantity >= ?", currencyCrypto, quantity)
	if err := db.Pluck("merchant_id", &assetMerchantIds).Error; err != nil {
		utils.Log.Errorf("Gets a list of asset conformance is failed.")
		return result
	}
	//通过支付方式过滤
	db = utils.DB.Model(&models.PaymentInfo{}).Where("in_use = ? ", 0)

	//fix - 是否只查询具有固定支付金额对应（支付宝，微信）二维码的币商
	//true - 只查询固定支付金额二维码（支付宝，微信）
	//false - 查询所有支付方式（即只要支付方式满足即可）
	if fix {
		db = db.Where("e_amount = ?", amount)
	} else {
		//0表示非固定金额
		db = db.Where("e_amount = ?", 0)
	}
	//pay_type - 支付类型混合值，示例： 1 - 微信， 2 - 支付宝, 4 - 银行， 3 - 银行+支付宝， 5 - 银行+微信，6 - 微信+支付宝， 7 - 所有
	switch payType {
	case 1:
		db = db.Where("pay_type = ?", 1)
	case 2:
		db = db.Where("pay_type = ?", 2)
	case 3:
		db = db.Where("pay_type = ? AND pay_type= ?", 1, 2)
	case 4:
		db = db.Where("pay_type = ?", 4)
	case 5:
		db = db.Where("pay_type = ? AND pay_type= ?", 1, 4)
	case 6:
		db = db.Where("pay_type = ? AND pay_type= ?", 2, 4)
	case 7:
		//所有的支付方式，不过滤
	default:
		return result
	}

	if err := db.Pluck("uid", &paymentMerchantIds).Error; err != nil {
		utils.Log.Errorf("Gets a list of payment conformance is failed.")
		return result
	}
	merchantIds = mergeList(merchantIds, assetMerchantIds, paymentMerchantIds)

	//限制返回条数 0 代表全部返回
	if limit == 0 {
		return merchantIds
	} else if limit > 0 {
		return merchantIds[0:limit]
	}
	return result
}

func (engine *defaultEngine) AcceptOrder(
	order OrderToFulfill,
	merchantID int64,
) {
	utils.AddBackgroundJob(utils.AcceptOrderTask, utils.HighPriority, order, merchantID)
}

func (engine *defaultEngine) NotifyFulfillment(
	fulfillment *OrderFulfillment,
) {
	//notify fulfillment information to merchant.
	utils.AddBackgroundJob(utils.NotifyFulfillmentTask, utils.HighPriority, fulfillment)
}

func (engine *defaultEngine) UpdateFulfillment(
	msg models.Msg,
) {
	utils.AddBackgroundJob(utils.UpdateFulfillmentTask, utils.NormalPriority, msg)
}

// waitWinner - wait till winner comes.
func waitWinner(
	orderNumer string,
	winner *models.Merchant,
) {
	//per each orderNumber, there will be a timer to wait till some one response to "accept order".
	//if no one accept till timeout, then ReFulfillOrder will be called.
}

func getFufillmentByOrderNumber(orderNumber string) *OrderFulfillment {
	//get current fulfillment by order number, search from cache,
	//then persistency if not found
	return &OrderFulfillment{}
}

//wrapper methods complies to goworker func.
func fulfillOrder(queue string, args ...interface{}) error {
	//recover OrderToFulfill from args
	var order OrderToFulfill
	if orderArg, ok := args[0].(map[string]interface{}); ok {
		order = getOrderToFulfillFromMapStrings(orderArg)
	} else {
		return fmt.Errorf("Wrong order arg: %v", args[0])
	}
	merchants := engine.selectMerchantsToFulfillOrder(&order)
	if len(*merchants) == 0 {
		//TODO: no merchants found, will re-fulfill order later
		return nil
	}
	//send order to pick
	engine.SendOrder(&order, merchants)
	return nil
}

func sendOrder(queue string, args ...interface{}) error {
	//recover OrderToFulfill and merchants ID map from args
	var order OrderToFulfill
	if orderArg, ok := args[0].(map[string]interface{}); ok {
		order = getOrderToFulfillFromMapStrings(orderArg)
	} else {
		return fmt.Errorf("Wrong order arg: %v", args[0])
	}
	var merchants []int64
	if merchangtsArg, ok := args[1].([]interface{}); ok {
		for _, id := range merchangtsArg {
			if mid, ok := id.(json.Number); ok {
				n, _ := mid.Int64()
				merchants = append(merchants, n)
			}
		}
	} else {
		return fmt.Errorf("Wrong merchant IDs: %v", args[1])
	}
	utils.Log.Debugf("Order %v sent to: %v", order, merchants)
	h5 := []string{order.OrderNumber}
	if err := NotifyThroughWebSocketTrigger(models.SendOrder, &merchants, &h5, 600, []OrderToFulfill{order}); err != nil {
		utils.Log.Errorf("Send order through websocket trigger API failed: %v", err)
	}

	return nil
}

func acceptOrder(queue string, args ...interface{}) error {
	//book keeping of all merchants who accept the order
	//recover OrderToFulfill and merchants ID map from args
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
	//now just choose the first responder as winner TODO: decide the winner

	var fulfillment *OrderFulfillment
	var err error
	if fulfillment, err = FulfillOrderByMerchant(order, merchantID, 0); err != nil {
		return fmt.Errorf("Unable to connect order with merchant: %v", err)
	}

	//notify fulfillment
	eng := NewOrderFulfillmentEngine(nil)
	eng.NotifyFulfillment(fulfillment)
	return nil
}

func notifyFulfillment(queue string, args ...interface{}) error {
	//recover order fulfillment information from args...
	//args:
	// fulfillment - OrderFulfillment which keeps both OrderToFulfill and Merchant information
	var fulfillment OrderFulfillment
	if fulfillmentArg, ok := args[0].(map[string]interface{}); ok {
		fulfillment = getFulfillmentInfoFromMapStrings(fulfillmentArg)
	} else {
		return fmt.Errorf("Wrong format of OrderFulfillment arg: %v", args[0])
	}
	utils.Log.Debugf("Fulfillment: %v", fulfillment)
	merchantID := fulfillment.MerchantID
	orderNumber := fulfillment.OrderNumber
	if err := NotifyThroughWebSocketTrigger(models.FulfillOrder, &[]int64{merchantID}, &[]string{orderNumber}, 600, []OrderFulfillment{fulfillment}); err != nil {
		utils.Log.Errorf("Send fulfillment through websocket trigger API failed: %v", err)
		return err
	}
	return nil
}

var msgTypes = map[string]models.MsgType{
	"send_order":    models.SendOrder,
	"accept":        models.Accept,
	"fulfill_order": models.FulfillOrder,
	"notify_paid":   models.NotifyPaid,
	"confirm_paid":  models.NotifyPaid,
	"transferred":   models.Transferred,
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
		uponNotifyPaid(msg)
	case models.ConfirmPaid:
		uponConfirmPaid(msg)
	case models.Transferred:
		uponTransferred(msg)
	}
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

func uponNotifyPaid(msg models.Msg) {
	//update order-fulfillment information
	ordNum, direction := getOrderNumberAndDirectionFromMessage(msg)
	if direction == 0 {
		//Trader buy, update order status, fulfillment
		order := models.Order{}
		if err := utils.DB.First(&order, "order_number = ?", ordNum).Error; err != nil {
			utils.Log.Errorf("Unable to find order with number %s. %v", ordNum, err)
			return
		}
		fulfillment := models.Fulfillment{}
		if err := utils.DB.Where("order_number = ?", ordNum).Order("seq_id DESC").First(&fulfillment).Error; err != nil {
			utils.Log.Errorf("No fulfillment with order number %s found. %v", ordNum, err)
			return
		}
		tx := utils.DB.Begin()
		//update order status
		if err := tx.Model(&order).Update("status", models.NOTIFYPAID).Error; err != nil {
			tx.Rollback()
			utils.Log.Errorf("Can't update order %s status to %s. %v", ordNum, "NOTIFYPAID", err)
			return
		}
		timeoutStr := utils.Config.GetString("fulfillment.timeout.notifypaymentconfirmed")
		timeout, _ := strconv.ParseInt(timeoutStr, 10, 32)
		if err := tx.Model(&fulfillment).Updates(models.Fulfillment{Status: models.NOTIFYPAID, PaidAt: time.Now(), NotifyPaymentConfirmedBefore: time.Now().Add(time.Duration(timeout) * time.Second)}).Error; err != nil {
			tx.Rollback()
			utils.Log.Errorf("Can't update order %s fulfillment info. %v", ordNum, err)
			return
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
			OriginStatus:  models.ACCEPTED,
			UpdatedStatus: models.NOTIFYPAID,
		}
		if err := tx.Create(&fulfillmentLog).Error; err != nil {
			tx.Rollback()
			utils.Log.Errorf("Can't create order %s fulfillment log. %v", ordNum, err)
			return
		}
		tx.Commit()
	}
	//then notify partner the same message - only direction = 0, Trader Buy
	NotifyThroughWebSocketTrigger(models.NotifyPaid, &msg.MerchantId, &msg.H5, 600, msg.Data)
}

func uponConfirmPaid(msg models.Msg) {
	//update order-fulfillment information
	//no need to notify the partner as he/she already exits
}

func uponTransferred(models.Msg) {
	//TODO: currently automatically transfer crypto coin to buyer after payment confirmed
}

//RegisterFulfillmentFunctions - register fulfillment functions, called by server
func RegisterFulfillmentFunctions() {
	//register worker function
	utils.RegisterWorkerFunc(utils.FulfillOrderTask, fulfillOrder)
	utils.RegisterWorkerFunc(utils.SendOrderTask, sendOrder)
	utils.RegisterWorkerFunc(utils.AcceptOrderTask, acceptOrder)
	utils.RegisterWorkerFunc(utils.NotifyFulfillmentTask, notifyFulfillment)
	utils.RegisterWorkerFunc(utils.UpdateFulfillmentTask, updateFulfillment)
}
