package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/go-redis/redis"
	"strings"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/utils"
)

//下订单
func PlaceOrder(req response.CreateOrderRequest) string {
	//var resp response.CreateOrderResult
	var ret response.OrdersRet
	var orderRequest response.OrderRequest
	var order models.Order
	var serverUrl string

	//1.todo 创建订单
	orderRequest = PlaceOrderReq2CreateOrderReq(req)
	ret = CreateOrder(orderRequest)
	//ret=UpdateOrder(orderRequest)
	if ret.Status != response.StatusSucc {
		utils.Log.Error("create order fail")
		return ""
	}
	order = ret.Data[0]
	orderNumber := order.OrderNumber //订单id

	//2.todo 创建订单成功，回调平台服务，通知创建订单成功
	serverUrl = FindServerUrl(req.ApiKey)
	if serverUrl == "" {
		utils.Log.Error("serverUrl is null")
	} else {
		utils.Log.Debugf("create order success,serverUrl is: %s", serverUrl)

		//jsonData, err1 := json.Marshal(order)
		//if err1 != nil {
		//	utils.Log.Error("order convert to json wrong,v%", err1)
		//}
		//
		//resp, err := http.Post(serverUrl, "application/json", bytes.NewBuffer(jsonData))
		//if err != nil || resp.Status != response.StatusSucc {
		//	utils.Log.Error("can not call distributor server ,v%", err)
		//	return ""
		//}

	}

	//3. todo 调用派单服务

	orderToFulfill := OrderToFulfill{
		OrderNumber:    order.OrderNumber,
		Direction:      order.Direction,
		OriginOrder:    order.OriginOrder,
		AccountID:      order.AccountId,
		DistributorID:  order.DistributorId,
		CurrencyCrypto: order.CurrencyCrypto,
		CurrencyFiat:   order.CurrencyFiat,
		Quantity:       order.Quantity,
		Price:          float32(order.Price),
		Amount:         order.Amount,
		PayType:        uint(order.PayType),
		QrCode:         order.QrCode,
		Name:           order.Name,
		BankAccount:    order.BankAccount,
		Bank:           order.Bank,
		BankBranch:     order.BankBranch,
	}
	engine := NewOrderFulfillmentEngine(nil)
	engine.FulfillOrder(&orderToFulfill)

	//4.todo 根据下单结果，重定向
	//var redirectUrl string
	//redirectUrl=utils.Config.GetString("redirectUrl.createurl")+orderNumber

	return orderNumber

}

//根据apiKey，apiSecret 到distributor表里查询serverUrl

func FindServerUrl(apiKey string) string {
	var serverUrl string
	var distributor models.Distributor
	if err := utils.DB.Model(&distributor).First(&distributor, "api_key=? ", apiKey).Error; err != nil {
		utils.Log.Error("can not find distributor by apiKey and apiSecret")
		return ""
	}
	serverUrl = distributor.ServerUrl
	return serverUrl

}

func PlaceOrderReq2CreateOrderReq(req response.CreateOrderRequest) response.OrderRequest {
	var resp response.OrderRequest

	resp.Price = req.Price
	resp.Amount = req.Amount
	resp.DistributorId = req.DistributorId
	resp.Quantity = req.TotalCount
	resp.OriginOrder = req.OrderNo
	resp.CurrencyCrypto = req.CoinType
	resp.Direction = req.OrderType
	resp.PayType = req.PayType
	resp.Name = req.Name
	resp.BankAccount = req.BankAccount
	resp.Bank = req.Bank
	resp.BankBranch = req.BankBranch
	resp.QrCode = req.QrCode

	return resp

}

/*
 struct convert json string
*/
func Struct2JsonString(structt interface{}) (jsonString string, err error) {
	data, err := json.Marshal(structt)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GenSignatureWith(mesthod string, host string, url string, str string, apikey string) string {
	return strings.ToUpper(mesthod) + host + url + str + apikey
}

func GenSignatureWith2(mesthod string, host string, url string, originOrder string, distributorId string, apikey string) string {
	return strings.ToUpper(mesthod) + host + url + originOrder + distributorId + apikey
}

//首先根据apiKey从redis里查询secretKey，若没查到，则从数据库中查询，并把apiKey，secretKey保存在redis里
func GetSecretKeyByApiKey(apiKey string) string {
	if apiKey == "" {
		utils.Log.Error("apiKey is null")
		return ""
	}
	secretKey, err := utils.RedisClient.Get(apiKey).Result()
	if err != redis.Nil {
		return secretKey

	}
	ditributor, err := GetDistributorByAPIKey(apiKey)

	if err != nil {
		utils.Log.Error("can not get secretkey according to apiKey=[%s] ", apiKey)
		return ""

	}
	secretKey = ditributor.ApiSecret
	utils.RedisSet(apiKey, secretKey, 30*time.Minute)
	return secretKey

}

func HmacSha256Base64Signer(message string, secretKey string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(message))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil

}