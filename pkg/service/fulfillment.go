package service

import (
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
	"github.com/benmanns/goworker"
)

// OrderToFulfill - order information for merchants to pick-up
type OrderToFulfill struct {
	//Account ID
	AccountID string
	//Distributor ID
	DistributorID int
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

// OrderFulfillment - Order fulfillment result.
type OrderFulfillment struct {
	*OrderToFulfill
	//Merchant ID
	MerchantID int
	//Merchant nickname
	MerchantNickName string
	//Merchant avatar URI
	MerchantAvatarURI string
	//Payment information, by reference
	PaymentInfo *models.PaymentInfo
}

// OrderFulfillmentEngine - engine interface of order fulfillment.
// The platform may change to new engine according to fulfillment rules changing.
type OrderFulfillmentEngine interface {
	// FulfillOrder - Async request to fulfill an order.
	FulfillOrder(
		accountID string, // Account identifier, may be blank
		distributorID int, // Distributor id, must not be blank
		orderNumber string, // Order number to be fulfilled.
		direction int, //Buy or Sell order, 0-buy, 1-sell, from trader's perspective.
		currencyCrypto string, //crypto-coin currency
		currencyFiat string, //fiat currency
		quantity string, //Quantity of crypto coins to fulfill, in currencyCrypto
		price string, //Price/Rate between crypto and fiat
		amount string, //Total amount of the order in currencyFiat
		payType uint, //Accepted payment type, chosen by trader
	)
	// SelectMerchantsToFulfillOrder - Select merchants to fulfill the specified orders.
	// The returned merchant(s) would receive the OrderToFulfill object through notification channel.
	// When there's only one merchant returned in the result, it might be exhausted matching result
	// or the first automatic processing merchant selected. No matter which situation, just send OrderToFulfill
	// message to the selected merchant. [no different process logic needed by caller]
	SelectMerchantsToFulfillOrder(order *OrderToFulfill) *[]models.Merchant
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
		merchants []models.Merchant, // a list of merchants to pick-up the order
	)
	// NotifyFulfillment - notify trader/merchant about the fulfillment.
	// Before notification, order is set to ACCEPTED
	NotifyFulfillment(
		fulfillment *OrderFulfillment, //the fulfillment choice decided by engine
		trader interface{}, //the connection end-point of trader (accountID + distributorID can't guarantee to connect back to account)
		merchant *models.Merchant, // the picked merchant
	)
}

// NewOrderFulfillmentEngine - return a new fulfillment engine.
// Factory method of fulfillment engine to adopt to future engine extension.
func NewOrderFulfillmentEngine(_ /*config*/ interface{}) OrderFulfillmentEngine {
	//may init engine by config, now ignore it
	return new(defaultEngine)
}

//defaultEngine - hidden default OrderFulfillmentEngine
type defaultEngine struct{
	*TaskQueue
}

func (engine *defaultEngine) FulfillOrder(
	accountID string,
	distributorID int,
	orderNumber string,
	direction int,
	currencyCrypto string,
	currencyFiat string,
	quantity string,
	price string,
	amount string,
	payType uint,
) {
	task := OrderToFulfill{
		AccountID: accountID,
		DistributorID: distributorID,
		OrderNumber: orderNumber,
		Direction: direction,
		CurrencyCrypto: currencyCrypto,
		CurrencyFiat: currencyFiat,
		Quantity: quantity,
		Price: price,
		Amount: amount,
		PayType: payType,
	}
	TaskQueue.AddTask(task) //process task asynchronously 
}

func (engine *defaultEngine) ReFulfillOrder(
	orderNumber string,
) {
	//get corresponding fulfillment object, update it then re-run FulfillOrder
	var lastFulfillment *OrderFulfillment
	lastFulfillment = getFufillmentByOrderNumber(orderNumber)
	engine.FulfillOrder(
		lastFulfillment.AccountID,
		lastFulfillment.DistributorID,
		lastFulfillment.OrderNumber,
		lastFulfillment.Direction,
		lastFulfillment.CurrencyCrypto,
		lastFulfillment.CurrencyFiat,
		lastFulfillment.Quantity,
		lastFulfillment.Price,
		lastFulfillment.Amount,
		lastFulfillment.PayType,
	)
}

func (engine *defaultEngine) SendOrder(
	order *OrderToFulfill,
	merchants []models.Merchant,
) {
	//send "order to fulfill" to selected merchants

}

func (engine *defaultEngine) SelectMerchantsToFulfillOrder(order *OrderToFulfill) *[]models.Merchant {
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
	trader interface{},
	merchant *models.Merchant,
) {
	//notify fulfillment information to merchant. 
	// "trader" contains arbitrary information of trader such as account id, distributor id, etc.
	// we don't define it as a struct then. [TODO: refine it later]
}

// waitWinner - wait till winner comes.
func waitWinner(
	orderNumer string,
	winner *models.Merchant,
) {
	//per each orderNumber, there will be a timer to wait till some one response to "accept order".
	//if no one accept till timeout, then ReFulfillOrder will be called.
}

func getFufillmentByOrderNumber(orderNumber string) *OrderFulfillment{
	//get current fulfillment by order number, search from cache,
	//then persistency if not found
	return &OrderFulfillment{}
}

//wrapper methods complies to goworker func.
func selectMerchantsToFulfillOrder(queue string, args ...interface{}) error {}
func sendOrder(queue string, args ...interface{}) error {}
func notifyFulfillment(queue string, args ...interface{}) error {}

func init() {
	//register worker function
	utils.RegisterWorkerFunc("SelectMerchantsToFulfillOrder", selectMerchantsToFulfillOrder)
	utils.RegisterWorkerFunc("SendOrder", sendOrder)
	utils.RegisterWorkerFunc("NotifyFulfillment", notifyFulfillment)
}