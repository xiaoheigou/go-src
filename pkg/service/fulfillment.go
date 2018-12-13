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
	//Account ID
	AccountID string
	//Distributor ID
	DistributorID int64
	//Order number to fulfill.
	OrderNumber string
	//Buy or sell
	Direction int
	//Crypto currency
	CurrencyCrypto string
	//Fiat currency
	CurrencyFiat string
	//Quantity, in crypto currency
	Quantity string
	//Price - rate between crypto and fiat
	Price string
	//Amount of the order, in fiat currency
	Amount string
	//Payment type, chosen by trader
	PayType int
}

func getOrderToFulfillFromMapStrings(values map[string]interface{}) OrderToFulfill {
	var distrID, direct, payT int64

	if distrN, ok := values["DistributorID"].(json.Number); ok {
		distrID, _ = distrN.Int64()
	}
	if directN, ok := values["Direction"].(json.Number); ok {
		direct, _ = directN.Int64()
	}
	if payTN, ok := values["PayType"].(json.Number); ok {
		payT, _ = payTN.Int64()
	}
	return OrderToFulfill{
		AccountID:      values["AccountID"].(string),
		DistributorID:  distrID,
		OrderNumber:    values["OrderNumber"].(string),
		Direction:      int(direct),
		CurrencyCrypto: values["CurrencyCrypto"].(string),
		CurrencyFiat:   values["CurrencyFiat"].(string),
		Quantity:       values["Quantity"].(string),
		Price:          values["Price"].(string),
		Amount:         values["Amount"].(string),
		PayType:        int(payT),
	}
}

// OrderFulfillment - Order fulfillment result.
type OrderFulfillment struct {
	OrderToFulfill
	//Merchant ID
	MerchantID int64
	//Merchant nickname
	MerchantNickName string
	//Merchant avatar URI
	MerchantAvatarURI string
	//Payment information, by reference
	models.PaymentInfo
}

func getPaymentInfoFromMapStrings(values map[string]interface{}) models.PaymentInfo {
	var uid, payT int64
	if uidN, ok := values["uid"].(json.Number); ok {
		uid, _ = uidN.Int64()
	}
	if payTN, ok := values["pay_type"].(json.Number); ok {
		payT, _ = payTN.Int64()
	}

	return models.PaymentInfo{
		Uid:         uid,
		PayType:     int(payT),
		Name:        values["name"].(string),
		Bank:        values["bank"].(string),
		BankAccount: values["bank_account"].(string),
		BankBranch:  values["bank_branch"].(string),
	}
}

func getFulfillmentInfoFromMapStrings(values map[string]interface{}) OrderFulfillment {
	var merchantID int64
	if merchantIDN, ok := values["MerchantID"].(json.Number); ok {
		merchantID, _ = merchantIDN.Int64()
	}
	orderToFulfill := getOrderToFulfillFromMapStrings(values)
	paymentInfo := getPaymentInfoFromMapStrings(values)
	return OrderFulfillment{
		OrderToFulfill:    orderToFulfill,
		MerchantID:        merchantID,
		MerchantNickName:  values["MerchantNickName"].(string),
		MerchantAvatarURI: values["MerchantAvatarURI"].(string),
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
	selectMerchantsToFulfillOrder(order *OrderToFulfill) *[]models.Merchant
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
			AccountID:      lastFulfillment.AccountID,
			DistributorID:  lastFulfillment.DistributorID,
			OrderNumber:    lastFulfillment.OrderNumber,
			Direction:      lastFulfillment.Direction,
			CurrencyCrypto: lastFulfillment.CurrencyCrypto,
			CurrencyFiat:   lastFulfillment.CurrencyFiat,
			Quantity:       lastFulfillment.Quantity,
			Price:          lastFulfillment.Price,
			Amount:         lastFulfillment.Amount,
			PayType:        lastFulfillment.PaymentInfo.PayType,
		})
}

func (engine *defaultEngine) SendOrder(
	order *OrderToFulfill,
	merchants *[]int64,
) {
	//send "order to fulfill" to selected merchants
	utils.AddBackgroundJob(utils.SendOrderTask, utils.NormalPriority, *order, *merchants)
}

func (engine *defaultEngine) selectMerchantsToFulfillOrder(order *OrderToFulfill) *[]models.Merchant {
	//search logic starts here!
	//1. filter out merchants currently not in "accept order" mode;
	//2. filter out merchants who don't have enough "coins" to take Trader-Buy order;
	//3. filter out merchants who don't support Trader specified payment type;
	//4. prioritize those merchants who do have "fix-amount" payment QR code matching demand;
	//5. constraints: merchant's payment info can serve one order at same time (locked if already matched previous order)
	//6. constraints: merchant can only take one same "amount" order at same time;
	//7. constraints: risk-control concerns which may reject to assign order to some merchant (TODO: to be added later)
	return &[]models.Merchant{}
}

func (engine *defaultEngine) NotifyFulfillment(
	fulfillment *OrderFulfillment,
) {
	//notify fulfillment information to merchant.
	utils.AddBackgroundJob(utils.NotifyFulfillmentTask, utils.HighPriority, fulfillment)
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
	engine.selectMerchantsToFulfillOrder(&order)
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
	return nil
}

//RegisterFulfillmentFunctions - register fulfillment functions, called by server
func RegisterFulfillmentFunctions() {
	//register worker function
	utils.RegisterWorkerFunc(utils.FulfillOrderTask, fulfillOrder)
	utils.RegisterWorkerFunc(utils.SendOrderTask, sendOrder)
	utils.RegisterWorkerFunc(utils.NotifyFulfillmentTask, notifyFulfillment)
}
