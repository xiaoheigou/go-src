package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"

	"github.com/zzh20/timewheel"
)

var (
	engine        *defaultEngine
	wheel         *timewheel.TimeWheel
	notifyWheel   *timewheel.TimeWheel
	confirmWheel  *timewheel.TimeWheel
	transferWheel *timewheel.TimeWheel
	awaitTimeout  int64
	retryTimeout  int64
	retries       int64
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
	Quantity float64 `json:"quantity"`
	//Price - rate between crypto and fiat
	Price float32 `json:"price"`
	//Amount of the order, in fiat currency
	Amount float64 `json:"amount"`
	//Payment type, chosen by trader
	PayType uint `json:"pay_type"`
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
		Price:          float32(price),
		Amount:         amount,
		PayType:        uint(payT),
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

func waitAcceptTimeout(data interface{}) {
	//ignore first fmanager object, add later if needed
	//key = order number
	//no one accept till timeout, re-fulfill it then

	orderNum := data.(string)
	//merchants := data.(map[string]interface{})["merchants"].([]int64)
	utils.Log.Debugf("Order %s not accepted by any merchant. Re-fulfill it...", orderNum)
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
	go reFulfillOrder(&orderToFulfill, 1)
}

func notifyPaidTimeout(data interface{}) {
	//ignore first fmanager object, add later if needed
	//key = order number
	//no one accept till timeout, re-fulfill it then
	orderNum := data.(string)
	utils.Log.Debugf("Order %s not notify paid timeout.", orderNum)

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

	if order.Status < models.NOTIFYPAID {
		if order.Direction == 0 {
			//释放币商的币
			asset := models.Assets{}
			if tx.Set("gorm:query_option", "FOR UPDATE").First(&asset, "merchant_id = ? AND currency_crypto = ? ", order.MerchantId, order.CurrencyCrypto).RecordNotFound() {
				utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
				tx.Rollback()
				utils.Log.Errorf("Can't find corresponding asset record of merchant_id %d, currency_crypto %s", order.MerchantId, order.CurrencyCrypto)
				utils.Log.Debugf("func doTransfer finished abnormally.")
				return
			}
			//释放冻结的币
			if err := tx.Model(&models.Assets{}).Where("merchant_id = ? AND currency_crypto = ? AND qty_frozen >= ?", order.MerchantId, order.CurrencyCrypto, order.Quantity).
				Updates(map[string]interface{}{"quantity": asset.Quantity + order.Quantity, "qty_frozen": asset.QtyFrozen - order.Quantity}).Error; err != nil {
				utils.Log.Errorf("notifyPaidTimeout release coin is failed,order number:%s,merchantId:%d", orderNum, order.MerchantId)
				utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
				tx.Rollback()
				return
			}
			//订单支付方式释放
			if err := tx.Model(&models.PaymentInfo{}).Where("uid = ? AND id = ?", order.MerchantId, order.MerchantPaymentId).Update("in_use", 0).Error; err != nil {
				utils.Log.Errorf("notifyPaidTimeout release payment info,merchantId:[%d],orderNUmber:[%s]", order.MerchantId, order.OrderNumber)
				utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
				tx.Rollback()
				return
			}
		} else if order.Direction == 1 {
			//用户提现单子,确认付款超时,不释放币
			//asset := models.Assets{}
			//if tx.Set("gorm:query_option", "FOR UPDATE").First(&asset, "distributor_id = ? AND currency_crypto = ? ", order.DistributorId, order.CurrencyCrypto).RecordNotFound() {
			//	utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
			//	tx.Rollback()
			//	utils.Log.Errorf("Can't find corresponding asset record of DistributorId %d, currency_crypto %s", order.DistributorId, order.CurrencyCrypto)
			//	return
			//}
			////释放冻结的币
			//if err := tx.Model(&models.Assets{}).Where("distributor_id = ? AND currency_crypto = ? AND qty_frozen >= ?", order.DistributorId, order.CurrencyCrypto, order.Quantity).
			//	Updates(map[string]interface{}{"quantity": asset.Quantity + order.Quantity, "qty_frozen": asset.QtyFrozen - order.Quantity}).Error; err != nil {
			//	utils.Log.Errorf("notifyPaidTimeout release coin is failed,order number:%s,DistributorId:%d", orderNum, order.DistributorId)
			//	utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
			//	tx.Rollback()
			//	return
			//}
		}

		//订单状态改为suspended
		//failed, highlight the order to set status to "SUSPENDED"
		if err := tx.Model(&order).Where("order_number = ? AND status < ?", order.OrderNumber, models.NOTIFYPAID).Update("status", models.SUSPENDED).Error; err != nil {
			utils.Log.Errorf("Update order %s status to SUSPENDED failed", order.OrderNumber)
			utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
			tx.Rollback()
			return
		}
		//fulfillment log 添加记录
		var fulfillment models.Fulfillment
		if err := tx.Order("seq_id desc").First(&fulfillment, "order_number = ?", orderNum).Error; err != nil {
			utils.Log.Errorf("get fulfillment order %s failed", order.OrderNumber)
			utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
			tx.Rollback()
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
			OriginStatus:  order.Status,
			UpdatedStatus: models.SUSPENDED,
		}
		if err := tx.Create(&fulfillmentLog).Error; err != nil {
			utils.Log.Errorf("notifyPaidTimeout create fulfillmentLog is failed,order number:%s", orderNum)
			utils.Log.Errorf("tx in func notifyPaidTimeout rollback, tx=[%v]", tx)
			tx.Rollback()
			return
		}

	}

	utils.Log.Debugf("tx in func notifyPaidTimeout commit, tx=[%v]", tx)
	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("error tx in func notifyPaidTimeout commit, err=[%v]", err)
	}

}

func confirmPaidTimeout(data interface{}) {
	//ignore first fmanager object, add later if needed
	//key = order number
	//no one accept till timeout, re-fulfill it then
	orderNum := data.(string)
	utils.Log.Debugf("Order %s confirm paid timeout", orderNum)
	order := models.Order{}
	if utils.DB.First(&order, "order_number = ?", orderNum).RecordNotFound() {
		utils.Log.Errorf("Order %s not found.", orderNum)
		return
	}
	message := models.Msg{
		MsgType:    models.ConfirmPaid,
		MerchantId: []int64{order.MerchantId},
		H5:         []string{orderNum},
		Timeout:    0,
		Data: []interface{}{
			map[string]interface{}{
				"order_number": order.OrderNumber,
				"direction":    order.Direction,
			},
		},
	}
	uponConfirmPaid(message)
}

func transferTimeout(data interface{}) {
	orderNum := data.(string)
	utils.Log.Debugf("transfer timeout begin,orderNum:%s", orderNum)
	doTransfer(orderNum)
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
	var merchants, alreadyAcceptNotify, alreadyFulfillMerchants []int64
	//去掉手动接单的并且已经接单的
	if order.Direction == 0 {
		//Buy, try to match all-automatic merchants firstly
		// 1. available merchants(online + in_work) + auto accept order/confirm payment + fix amount match
		merchants = GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, true, 0, 0)
		if len(merchants) == 0 { //no priority merchants with fix amount match found, another round call
			// 2. available merchants(online + in_work) + auto accept order/confirm payment
			merchants = GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, false, 0, 0)
			if len(merchants) == 0 { //no priority merchants with non-fix amount match found, then "manual operation" merchants
				// 3. available merchants(online + in_work) + manual accept order/confirm payment + has fix amount qrcode
				merchants = GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, true, 1, 0)
				if len(merchants) == 0 { //Sell, all should manually processed
					// 4. available merchants(online + in_work) + manual accept order/confirm payment + has arbitrary amount qrcode
					merchants = GetMerchantsQualified(order.Amount, order.Quantity, order.CurrencyCrypto, order.PayType, false, 1, 0)
				}
				//手动接单的,只允许同时接一个订单
				if err := utils.DB.Model(models.Order{}).Where("status <= ?", models.NOTIFYPAID).Pluck("merchant_id", &alreadyAcceptNotify).Error; err != nil {
					utils.Log.Errorf("func selectMerchantsToFulfillOrder error, the select order = [%+v]", order)
				}
				//if err := utils.DB.Model(models.Order{}).Where("status = ? AND direction = 1", models.CONFIRMPAID).Pluck("merchant_id", &alreadyAcceptConfirm).Error; err != nil {
				//	utils.Log.Errorf("func selectMerchantsToFulfillOrder error, the select order = [%+v]", order)
				//}
				utils.Log.Debugf("already accept order merchant,[%v],[%v]", alreadyAcceptNotify, alreadyAcceptNotify)
			}
		}
	} else {
		//Sell, any online + in_work could pickup order
		merchants = GetMerchantsQualified(0, 0, order.CurrencyCrypto, order.PayType, false, 1, 0)
	}

	//去掉重新派单时已经派过的币商
	if err := utils.DB.Model(&models.Fulfillment{}).Where("order_number = ?", order.OrderNumber).Pluck("distinct merchant_id", &alreadyFulfillMerchants).Error; err != nil {
		utils.Log.Errorf("selectMerchantsToFulfillOrder get fulfillment is failed,orderNumber:%s", order.OrderNumber)
	}

	//去掉已经派过单的币商
	var selectedMerchants []int64
	if data, err := utils.GetCacheSetMembers(utils.UniqueOrderSelectMerchantKey(order.OrderNumber)); err != nil {
		utils.Log.Errorf("func selectMerchantsToFulfillOrder error, the select order = [%+v]", order)
	} else if len(data) > 0 {
		utils.Log.Debugf("already fulfill merchant,[%v]", selectedMerchants)
		utils.ConvertStringToInt(data, &selectedMerchants)
	}

	merchants = utils.DiffSet(merchants, selectedMerchants, alreadyAcceptNotify, alreadyFulfillMerchants)

	utils.Log.Debugf("before sort by last order time, the merchants = [%+v]", merchants)
	merchants = sortMerchantsByLastOrderTime(merchants, order.Direction)
	utils.Log.Debugf(" after sort by last order time, the merchants = [%+v]", merchants)

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

		if err := SendSmsOrderAccepted(merchantID, orderNum); err != nil {
			utils.Log.Errorf("order [%v] is accepted by merchant [%v], send sms fail. error [%v]", orderNum, merchantID, err)
		}
	} else { //already accepted, reject the request
		utils.Log.Debugf("merchant %d accepted order is failed,order already by merchant %s accept.", merchantID, merchant)
		if err := NotifyThroughWebSocketTrigger(models.Picked, &[]int64{merchantID}, &[]string{}, 60, nil); err != nil {
			utils.Log.Errorf("Notify Picked through websocket ")
		}
	}
}

func (engine *defaultEngine) UpdateFulfillment(
	msg models.Msg,
) {
	utils.Log.Debugf("func UpdateFulfillment begin, msg = [%+v]", msg)
	utils.AddBackgroundJob(utils.UpdateFulfillmentTask, utils.NormalPriority, msg)
	utils.Log.Debugf("func UpdateFulfillment finished normally.")
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
	utils.Log.Debugf("fulfill for order [%+V]", order.OrderNumber)
	merchants := engine.selectMerchantsToFulfillOrder(&order)
	if len(*merchants) == 0 {
		utils.Log.Warnf("fulfillOrder function None merchant is available at moment, will re-fulfill later.")
		go reFulfillOrder(&order, 1)
		return nil
	}
	//send order to pick
	if err := sendOrder(&order, merchants); err != nil {
		utils.Log.Errorf("Send order failed: %v", err)
		return err
	}
	//push into timewheel to wait
	utils.Log.Debugf("await timeout wheel,%v", wheel)
	//if wheel == nil {
	//	utils.Log.Debugf("accept order timeout wheel init")
	//	wheel = timewheel.New(1*time.Second, int(timeout), waitAcceptTimeout) //process wheel per second
	//	wheel.Start()
	//}
	timeout := awaitTimeout + retries*retryTimeout + awaitTimeout
	selectedMerchantsToRedis(order.OrderNumber, timeout, merchants)
	wheel.Add(order.OrderNumber)
	return nil
}

func reFulfillOrder(order *OrderToFulfill, seq uint8) {
	time.Sleep(time.Duration(retryTimeout) * time.Second)
	//re-fulfill
	merchants := engine.selectMerchantsToFulfillOrder(order)
	utils.Log.Debugf("re-fulfill for order [%+V] and selectedMerchants [%v]", order.OrderNumber, merchants)
	if len(*merchants) == 0 {
		utils.Log.Warnf("None merchant is available at moment, will re-fulfill later.")
		if seq <= uint8(retries) {
			go reFulfillOrder(order, seq+1)
			return
		}

		tx := utils.DB.Begin()
		if tx.Error != nil {
			utils.Log.Debugf("tx in func reFulfillOrder begin fail, tx=[%v]", tx)
			utils.Log.Errorf("func reFulfillOrder finished abnormally.")
			return
		}

		//failed, highlight the order to set status to "SUSPENDED"
		suspendedOrder := models.Order{}
		if tx.Set("gorm:query_option", "FOR UPDATE").Find(&suspendedOrder, "order_number = ?  AND status < ?", order.OrderNumber, models.ACCEPTED).RecordNotFound() {
			utils.Log.Errorf("Unable to find order %s", order.OrderNumber)
		} else {
			utils.Log.Debugf("tx in func reFulfillOrder begin, tx=[%v]", tx)
			if err := tx.Model(&models.Order{}).Where("order_number = ? AND status < ?", order.OrderNumber, models.ACCEPTED).Update("status", models.SUSPENDED).Error; err != nil {
				utils.Log.Errorf("Update order %s status to SUSPENDED failed", order.OrderNumber)
				utils.Log.Errorf("tx in func reFulfillOrder rollback, tx=[%v]", tx)
				tx.Rollback()
				return
			}
			if suspendedOrder.Direction == 1 {
				asset := models.Assets{}
				if tx.Set("gorm:query_option", "FOR UPDATE").Find(&asset, "distributor_id = ? AND currency_crypto = ? ", suspendedOrder.DistributorId, order.CurrencyCrypto).RecordNotFound() {
					utils.Log.Errorf("Can't find corresponding asset record of distributor_id %d, currency_crypto %s", suspendedOrder.DistributorId, order.CurrencyCrypto)
					utils.Log.Errorf("tx in func reFulfillOrder rollback, tx=[%v]", tx)
					tx.Rollback()
					return
				}
				//释放冻结的币
				if err := tx.Model(&models.Assets{}).Where("distributor_id = ? AND currency_crypto = ? AND qty_frozen >= ?", suspendedOrder.DistributorId, order.CurrencyCrypto, order.Quantity).
					Updates(map[string]interface{}{"quantity": asset.Quantity + order.Quantity, "qty_frozen": asset.QtyFrozen - order.Quantity}).Error; err != nil {
					utils.Log.Errorf("notifyPaidTimeout release coin is failed,order number:%s,merchantId:%d", suspendedOrder.OrderNumber, suspendedOrder.MerchantId)
					utils.Log.Errorf("tx in func reFulfillOrder rollback, tx=[%v]", tx)
					tx.Rollback()
					return
				}
			}
		}
		utils.Log.Debugf("tx in func reFulfillOrder commit, tx=[%v]", tx)
		if err := tx.Commit().Error; err != nil {
			utils.Log.Errorf("error tx in func reFulfillOrder commit, err=[%v]", err)
		}
		return
	}
	//send order to pick
	if err := sendOrder(order, merchants); err != nil {
		utils.Log.Errorf("Send order failed: %v", err)
	}
	//push into timewheel
	timeout := awaitTimeout + retries*retryTimeout + awaitTimeout
	selectedMerchantsToRedis(order.OrderNumber, timeout, merchants)
	wheel.Add(order.OrderNumber)
	return
}

func selectedMerchantsToRedis(orderNumber string, timeout int64, merchants *[]int64) {
	utils.Log.Debugf("selectedMerchantsToRedis orderNumber:[%s],timeout:[%d],merchants:[%v]", orderNumber, timeout, *merchants)
	key := utils.UniqueOrderSelectMerchantKey(orderNumber)
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
	if err := NotifyThroughWebSocketTrigger(models.SendOrder, merchants, &h5, uint(timeout), []OrderToFulfill{*order}); err != nil {
		utils.Log.Errorf("Send order through websocket trigger API failed: %v", err)
		utils.Log.Debugf("func sendOrder finished abnormally.")
		return err
	}
	utils.Log.Debugf("func sendOrder finished normally.")
	return nil
}

func acceptOrder(queue string, args ...interface{}) error {
	//book keeping of all merchants who accept the order
	//recover OrderToFulfill and merchants ID map from args
	utils.Log.Debugf("acceptOrder begin")
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
	wheel.Remove(order.OrderNumber)
	var fulfillment *OrderFulfillment
	var err error
	if fulfillment, err = FulfillOrderByMerchant(order, merchantID, 0); err != nil {
		wheel.Add(order.OrderNumber)
		return fmt.Errorf("Unable to connect order with merchant: %v", err)
	}
	notifyFulfillment(fulfillment)
	utils.Log.Debugf("acceptOrder end")
	return nil
}

func notifyFulfillment(fulfillment *OrderFulfillment) error {
	merchantID := fulfillment.MerchantID
	orderNumber := fulfillment.OrderNumber
	timeoutStr := utils.Config.GetString("fulfillment.timeout.notifypaid")
	timeout, _ := strconv.ParseInt(timeoutStr, 10, 32)
	utils.Log.Debugf("notifyFulfillment start, merchantID %v", merchantID)
	if err := NotifyThroughWebSocketTrigger(models.FulfillOrder, &[]int64{merchantID}, &[]string{orderNumber}, uint(timeout), []OrderFulfillment{*fulfillment}); err != nil {
		wheel.Add(fulfillment.OrderNumber)
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
		uponNotifyPaid(msg)
	case models.ConfirmPaid:
		uponConfirmPaid(msg)
	case models.Transferred:
		utils.Log.Warnf("msg with type Transferred should not occur in redis queue, it processed directly after confirm paid")
	case models.AutoConfirmPaid:
		uponAutoConfirmPaid(msg)
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

	tx := utils.DB.Begin()
	if tx.Error != nil {
		utils.Log.Debugf("tx in func uponNotifyPaid begin fail, tx=[%v]", tx)
		utils.Log.Errorf("func uponNotifyPaid finished abnormally.")
		return
	}
	utils.Log.Debugf("tx in func uponNotifyPaid begin, tx=[%v]", tx)

	//Trader buy, update order status, fulfillment
	order := models.Order{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&order, "order_number = ?", ordNum).RecordNotFound() {
		utils.Log.Errorf("Record not found: order with number %s.", ordNum)
		utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
		tx.Rollback()
		return
	}
	if direction == 0 {
		fulfillment := models.Fulfillment{}
		if tx.Set("gorm:query_option", "FOR UPDATE").Where("order_number = ?", ordNum).Order("seq_id DESC").First(&fulfillment).RecordNotFound() {
			utils.Log.Errorf("No fulfillment with order number %s found.", ordNum)
			utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
			tx.Rollback()
			return
		}

		//update order status
		if err := tx.Model(&order).Update("status", models.NOTIFYPAID).Error; err != nil {
			utils.Log.Errorf("Can't update order %s status to %s. %v", ordNum, "NOTIFYPAID", err)
			utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
			tx.Rollback()
			return
		}
		if err := tx.Model(&fulfillment).Updates(models.Fulfillment{Status: models.NOTIFYPAID, PaidAt: time.Now()}).Error; err != nil {
			utils.Log.Errorf("Can't update order %s fulfillment info. %v", ordNum, err)
			utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
			tx.Rollback()
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
			utils.Log.Errorf("Can't create order %s fulfillment log. %v", ordNum, err)
			utils.Log.Errorf("tx in func uponNotifyPaid rollback, tx=[%v]", tx)
			tx.Rollback()
			return
		}
		utils.Log.Debugf("tx in func uponNotifyPaid commit, tx=[%v]", tx)
		if err := tx.Commit().Error; err != nil {
			utils.Log.Errorf("error tx in func uponNotifyPaid commit, err=[%v]", err)
		}

		timeoutStr := utils.Config.GetString("fulfillment.timeout.notifypaymentconfirmed")
		timeout, _ := strconv.ParseInt(timeoutStr, 10, 32)
		//then notify partner the same message - only direction = 0, Trader Buy
		if err := NotifyThroughWebSocketTrigger(models.NotifyPaid, &msg.MerchantId, &msg.H5, uint(timeout), msg.Data); err != nil {
			utils.Log.Errorf("Notify partner notify paid messaged failed.")
		}
		//if confirmWheel == nil {
		//	//confirm paid timeout
		//	confirmWheel = timewheel.New(1*time.Second, int(timeout), confirmPaidTimeout) //process wheel per second
		//	confirmWheel.Start()
		//}
		notifyWheel.Remove(order.OrderNumber)
		confirmWheel.Add(order.OrderNumber)

		if err := SendSmsOrderPaid(fulfillment.MerchantID, ordNum); err != nil {
			utils.Log.Errorf("order [%v] is marked as paid by user, send sms to merchant [%v] fail. error [%v]", ordNum, fulfillment.MerchantID, err)
		}
	} else { //Trader Sell, trigger confirm paid automaticaly
		message := models.Msg{
			MsgType:    models.ConfirmPaid,
			MerchantId: msg.MerchantId,
			H5:         msg.H5,
			Timeout:    0,
			Data: []interface{}{
				map[string]interface{}{
					"order_number": order.OrderNumber,
					"direction":    order.Direction,
				},
			},
		}
		utils.Log.Debugf("tx in func uponNotifyPaid commit, tx=[%v]", tx)
		if err := tx.Commit().Error; err != nil {
			utils.Log.Errorf("error tx in func uponNotifyPaid commit, err=[%v]", err)
		}

		//as if we got confirm paid message from APP
		uponConfirmPaid(message)
	}
}

func uponConfirmPaid(msg models.Msg) {
	utils.Log.Debugf("func uponConfirmPaid begin, msg = [%+v]", msg)
	ordNum, _ := getOrderNumberAndDirectionFromMessage(msg)

	tx := utils.DB.Begin()
	if tx.Error != nil {
		utils.Log.Debugf("tx in func uponConfirmPaid begin fail, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return
	}
	utils.Log.Debugf("tx in func uponConfirmPaid begin, tx=[%v]", tx)

	order := models.Order{}
	if tx.Set("gorm:query_option", "FOR UPDATE").Where("order_number = ?", ordNum).First(&order).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("Record not found: order with number %s.", ordNum)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return
	}

	fulfillment := models.Fulfillment{}
	if tx.Set("gorm:query_option", "FOR UPDATE").Where("order_number = ?", ordNum).Order("seq_id DESC").First(&fulfillment).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("No fulfillment with order number %s found.", ordNum)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return
	}

	// check current status
	if fulfillment.Status == models.CONFIRMPAID {
		tx.Rollback()
		utils.Log.Errorf("order number %s is already with status %d (CONFIRMPAID), do nothing.", ordNum, models.CONFIRMPAID)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return
	} else if fulfillment.Status == models.TRANSFERRED {
		tx.Rollback()
		utils.Log.Errorf("order number %s has status %d (TRANSFERRED), cannot change it to %d (CONFIRMPAID).", ordNum, models.TRANSFERRED, models.CONFIRMPAID)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
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
		OriginStatus:  fulfillment.Status,
		UpdatedStatus: models.CONFIRMPAID,
	}
	if err := tx.Create(&fulfillmentLog).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Can't create order %s fulfillment log. %v", ordNum, err)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return
	}

	// update order status
	if err := tx.Model(&order).Update("status", models.CONFIRMPAID).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Can't update order %s status to %s. %v", ordNum, "CONFIRMPAID", err)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return
	}
	if err := tx.Model(&fulfillment).Updates(models.Fulfillment{Status: models.CONFIRMPAID, PaymentConfirmedAt: time.Now()}).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Can't update order %s fulfillment info. %v", ordNum, err)
		utils.Log.Errorf("tx in func uponConfirmPaid rollback, tx=[%v]", tx)
		utils.Log.Errorf("func uponConfirmPaid finished abnormally.")
		return
	}
	utils.Log.Debugf("tx in func uponConfirmPaid commit, tx=[%v]", tx)
	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("error tx in func uponConfirmPaid commit, err=[%v]", err)
	}

	notifyMerchant := []int64{fulfillment.MerchantID}

	// notify partner
	if err := NotifyThroughWebSocketTrigger(models.ConfirmPaid, &notifyMerchant, &[]string{}, 0, msg.Data); err != nil {
		utils.Log.Errorf("Notify partner notify paid messaged failed.")
	}

	if order.Direction == 0 {
		doTransfer(ordNum)
	} else {
		transferWheel.Add(ordNum)
	}

	confirmWheel.Remove(order.OrderNumber)
	utils.Log.Debugf("func uponConfirmPaid finished normally.")
}

func doTransfer(ordNum string) {
	utils.Log.Debugf("func doTransfer begin, OrderNumber = [%+v]", ordNum)

	tx := utils.DB.Begin()
	if tx.Error != nil {
		utils.Log.Debugf("tx in func doTransfer begin fail, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally.")
		return
	}
	utils.Log.Debugf("tx in func doTransfer begin, tx=[%v]", tx)

	//Trader buy, update order status, fulfillment
	order := models.Order{}
	if tx.Set("gorm:query_option", "FOR UPDATE").Where("order_number = ?", ordNum).First(&order).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("Record not found: order with number %s.", ordNum)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally.")
		return
	}
	fulfillment := models.Fulfillment{}
	if tx.Set("gorm:query_option", "FOR UPDATE").Where("order_number = ?", ordNum).Order("seq_id DESC").First(&fulfillment).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("No fulfillment with order number %s found.", ordNum)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally.")
		return
	}

	if fulfillment.Status == models.TRANSFERRED {
		tx.Rollback()
		utils.Log.Errorf("order number %s is already with status %d (TRANSFERRED), do nothing.", ordNum, models.TRANSFERRED)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally.")
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
		OriginStatus:  fulfillment.Status,
		UpdatedStatus: models.TRANSFERRED,
	}
	if err := tx.Create(&fulfillmentLog).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Can't create order %s fulfillment log. %v", ordNum, err)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally.")
		return
	}

	//update order
	if err := tx.Model(&order).Update("status", models.TRANSFERRED, "merchant_payment_id", fulfillment.MerchantPaymentID).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Can't update order %s status to %s. %v", ordNum, "TRANSFERRED", err)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally.")
		return
	}
	transferredAt := time.Now()
	if err := tx.Model(&fulfillment).Updates(models.Fulfillment{Status: models.TRANSFERRED, TransferredAt: transferredAt}).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Can't update order %s fulfillment info. %v", ordNum, err)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally.")
		return
	}

	asset := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&asset, "merchant_id = ? AND currency_crypto = ? ", order.MerchantId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("Can't find corresponding asset record of merchant_id %d, currency_crypto %s", order.MerchantId, order.CurrencyCrypto)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally.")
		return
	}

	// 找到平台商记录
	assetForDist := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForDist, "distributor_id = ? AND currency_crypto = ? ", order.DistributorId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		utils.Log.Errorf("Can't find corresponding asset record of distributor_id %d, currency_crypto %s", order.DistributorId, order.CurrencyCrypto)
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("func doTransfer finished abnormally.")
		return
	}

	if order.Direction == 0 {
		// Trader Buy
		utils.Log.Debugf("Freeze [%v] %v for merchant (uid=[%v])", order.Quantity, order.CurrencyCrypto, fulfillment.MerchantPaymentID)
		if err := tx.Table("assets").Where("id = ? and qty_frozen >= ?", asset.Id, order.Quantity).Update("qty_frozen", asset.QtyFrozen-order.Quantity).Error; err != nil {
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("Can't freeze asset for merchant (uid=[%v]). err: %v", asset.MerchantId, err)
			utils.Log.Errorf("func doTransfer finished abnormally.")
			tx.Rollback()
			return
		}

		// 转币给平台商
		if err := tx.Table("assets").Where("id = ? ", assetForDist.Id).Update("quantity", assetForDist.Quantity+order.Quantity).Error; err != nil {
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("Can't transfer asset to distributor (distributor_id=[%v]). err: %v", assetForDist.DistributorId, err)
			utils.Log.Errorf("func doTransfer finished abnormally.")
			tx.Rollback()
			return
		}

		if err := tx.Table("payment_infos").Where("id = ?", fulfillment.MerchantPaymentID).Update("in_use", 0).Error; err != nil {
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("Can't change in_use to 0, record id=[%v], err=[%v]", fulfillment.MerchantPaymentID, err)
			utils.Log.Errorf("func doTransfer finished abnormally.")
			tx.Rollback()
			return
		}
	} else {
		// Trader Sell
		utils.Log.Debugf("Add [%v] %v for merchant (uid=[%v])", order.Quantity, order.CurrencyCrypto, fulfillment.MerchantPaymentID)
		if err := tx.Table("assets").Where("id = ?", asset.Id).Update("quantity", asset.Quantity+order.Quantity).Error; err != nil {
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("Can't add [%s] for merchant (uid=[%v]): %v", order.CurrencyCrypto, asset.MerchantId, err)
			utils.Log.Errorf("func doTransfer finished abnormally.")
			tx.Rollback()
			return
		}

		// 从平台商扣币
		if err := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForDist.Id, order.Quantity).Update("qty_frozen", assetForDist.QtyFrozen-order.Quantity).Error; err != nil {
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("Can't reduce [%s] for distributor (distributor_id=[%v]). err: %v", order.CurrencyCrypto, assetForDist.DistributorId, err)
			utils.Log.Errorf("func doTransfer finished abnormally.")
			tx.Rollback()
			return
		}
	}

	// Add asset history
	assetHistory := models.AssetHistory{
		Currency:    order.CurrencyCrypto,
		Direction:   order.Direction,
		MerchantId:  order.MerchantId,
		Quantity:    order.Quantity,
		IsOrder:     1,
		OrderNumber: ordNum,
	}
	if err := tx.Model(&models.AssetHistory{}).Create(&assetHistory).Error; err != nil {
		utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
		utils.Log.Errorf("create asset history for merchant (uid=[%v]) failed. err:[%v]", order.MerchantId, err)
		utils.Log.Errorf("func doTransfer finished abnormally.")
		tx.Rollback()
		return
	}
	if order.Direction == 0 {
		utils.Log.Infof("merchant (uid=[%v]) pay out %v %v", order.MerchantId, order.Quantity, order.CurrencyCrypto)
	} else if order.Direction == 1 {
		utils.Log.Infof("merchant (uid=[%v]) receive %v %v", order.MerchantId, order.Quantity, order.CurrencyCrypto)
	}

	utils.Log.Debugf("tx in func doTransfer commit, tx=[%v]", tx)
	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("error tx in func doTransfer commit, err=[%v]", err)
	}

	if err := utils.UpdateMerchantLastOrderTime(strconv.FormatInt(order.MerchantId, 10), order.Direction, transferredAt); err != nil {
		utils.Log.Warnf("func doTransfer call UpdateMerchantLastOrderTime fail [%+v].", err)
	}

	AsynchronousNotifyDistributor(order)

	utils.Log.Debugf("func doTransfer finished normally.")
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
	uponConfirmPaid(message)
}

//RegisterFulfillmentFunctions - register fulfillment functions, called by server
func RegisterFulfillmentFunctions() {
	//register worker function
	utils.RegisterWorkerFunc(utils.FulfillOrderTask, fulfillOrder)
	utils.RegisterWorkerFunc(utils.AcceptOrderTask, acceptOrder)
	utils.RegisterWorkerFunc(utils.UpdateFulfillmentTask, updateFulfillment)
}

func init() {
	timeoutStr := utils.Config.GetString("fulfillment.timeout.awaitaccept")
	awaitTimeout, _ = strconv.ParseInt(timeoutStr, 10, 64)
	utils.Log.Debugf("wheel init,timeout:%d", awaitTimeout)
	wheel = timewheel.New(1*time.Second, int(awaitTimeout), waitAcceptTimeout) //process wheel per second
	wheel.Start()

	timeoutStr = utils.Config.GetString("fulfillment.timeout.notifypaid")
	timeout, _ := strconv.ParseInt(timeoutStr, 10, 64)
	utils.Log.Debugf("notify wheel init,timeout:%d", timeout)
	notifyWheel = timewheel.New(1*time.Second, int(timeout), notifyPaidTimeout) //process wheel per second
	notifyWheel.Start()

	//confirm paid timeout
	timeoutStr = utils.Config.GetString("fulfillment.timeout.notifypaymentconfirmed")
	timeout, _ = strconv.ParseInt(timeoutStr, 10, 64)
	utils.Log.Debugf("confirm wheel init,timeout:%d", timeout)
	confirmWheel = timewheel.New(1*time.Second, int(timeout), confirmPaidTimeout) //process wheel per second
	confirmWheel.Start()

	timeoutStr = utils.Config.GetString("fulfillment.timeout.transfer")
	timeout, _ = strconv.ParseInt(timeoutStr, 10, 64)
	utils.Log.Debugf("transfer wheel init,timeout:%d", timeout)
	confirmWheel = timewheel.New(1*time.Second, int(timeout), transferTimeout) //process wheel per second
	confirmWheel.Start()

	timeoutStr = utils.Config.GetString("fulfillment.timeout.retry")
	retryTimeout, _ = strconv.ParseInt(timeoutStr, 10, 64)
	utils.Log.Debugf("retry timeout:%d", retryTimeout)

	retryStr := utils.Config.GetString("fulfillment.retries")
	retries, _ = strconv.ParseInt(retryStr, 10, 64)
	utils.Log.Debugf("retries:%d", retries)
}
