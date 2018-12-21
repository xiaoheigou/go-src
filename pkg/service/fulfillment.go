package service

import (
	"encoding/json"
	"fmt"

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
}

func getOrderToFulfillFromMapStrings(values map[string]interface{}) OrderToFulfill {
	var direct, payT int64

	if directN, ok := values["direction"].(json.Number); ok {
		direct, _ = directN.Int64()
	}
	if payTN, ok := values["pay_type"].(json.Number); ok {
		payT, _ = payTN.Int64()
	}
	return OrderToFulfill{
		OrderNumber:    values["order_number"].(string),
		Direction:      int(direct),
		CurrencyCrypto: values["currency_crypto"].(string),
		CurrencyFiat:   values["currency_fiat"].(string),
		Quantity:       values["quantity"].(float64),
		Price:          values["price"].(float64),
		Amount:         values["amount"].(float64),
		PayType:        int(payT),
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

func getPaymentInfoFromMapStrings(values map[string]interface{}) []models.PaymentInfo {
	//get payment_info firstly, []interface{}
	var paymentInfo = values["payment_info"]
	var result []models.PaymentInfo
	if paymentList, ok := paymentInfo.([]map[string]interface{}); ok {
		for _, payment := range paymentList {
			var uid, payT int64
			if uidN, ok := payment["uid"].(json.Number); ok {
				uid, _ = uidN.Int64()
			}
			if payTN, ok := payment["pay_type"].(json.Number); ok {
				payT, _ = payTN.Int64()
			}
			var pi models.PaymentInfo
			switch payT {
			case 1:
			case 2:
				pi = models.PaymentInfo{
					Uid:       uid,
					PayType:   int(payT),
					EAccount:  payment["e_account"].(string),
					QrCode:    payment["qr_code"].(string),
					QrCodeTxt: payment["qr_code_txt"].(string),
					EAmount:   payment["e_amount"].(float64),
				}
			case 4:
				pi = models.PaymentInfo{
					Uid:         uid,
					PayType:     int(payT),
					Name:        payment["name"].(string),
					Bank:        payment["bank"].(string),
					BankAccount: payment["bank_account"].(string),
					BankBranch:  payment["bank_branch"].(string),
				}
			}
			result = append(result, pi)
		}
		return result
	}
	//error occured
	utils.Log.Errorf("Parsing payment_info from queue message failed")
	return []models.PaymentInfo{}
}

func getFulfillmentInfoFromMapStrings(values map[string]interface{}) OrderFulfillment {
	var merchantID int64
	if merchantIDN, ok := values["merchant_id"].(json.Number); ok {
		merchantID, _ = merchantIDN.Int64()
	}
	orderToFulfill := getOrderToFulfillFromMapStrings(values)
	paymentInfo := getPaymentInfoFromMapStrings(values)
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
	fmt.Printf("merchants contents: %d\n", len(merchants))
	return &merchants
}

//GetMerchantsQualified - return mock data
func GetMerchantsQualified(amount, quantity float64, currencyCrypto string, payType int, fix bool, group uint8, limit uint8) []int64 {
	return []int64{1, 2, 3}
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
	//according to different operation + data, update order/fulfillment accordingly.
	//in additon, send notification to impacted partner of the operation
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
		//no merchants found, will re-fulfill order later
		return nil
	}
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
	if err := NotifyThroughWebSocketTrigger(models.FulfillOrder, &[]int64{merchantID}, &[]string{orderNumber}, 600, fulfillment); err != nil {
		utils.Log.Errorf("Send fulfillment through websocket trigger API failed: %v", err)
		return err
	}
	return nil
}

//RegisterFulfillmentFunctions - register fulfillment functions, called by server
func RegisterFulfillmentFunctions() {
	//register worker function
	utils.RegisterWorkerFunc(utils.FulfillOrderTask, fulfillOrder)
	utils.RegisterWorkerFunc(utils.SendOrderTask, sendOrder)
	utils.RegisterWorkerFunc(utils.NotifyFulfillmentTask, notifyFulfillment)
}
