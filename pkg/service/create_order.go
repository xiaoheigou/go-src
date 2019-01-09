package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/utils"
)

const (
	CONTENT_TYPE          = "Content-Type"
	ACCEPT                = "Accept"
	FAIL                  = "Fail"
	SUCCESS               = "Success"
	APPLICATION_JSON      = "application/json"
	APPLICATION_JSON_UTF8 = "application/json; charset=UTF-8"
)

//下订单
func PlaceOrder(req response.CreateOrderRequest) response.CreateOrderRet {
	var orderRet response.OrdersRet
	var orderRequest response.OrderRequest
	var order models.Order
	var serverUrl string
	var ret response.CreateOrderRet
	var createOrderResult response.CreateOrderResult

	//1.todo 创建订单
	orderRequest = PlaceOrderReq2CreateOrderReq(req)

	tx := utils.DB.Begin()
	assets, err := GetCoinQuantity(strconv.FormatInt(orderRequest.DistributorId, 10), orderRequest.CurrencyCrypto)
	if err != nil {
		utils.Log.Debugf("get the coin number of distributor wrong,to create, distributorId= %s", orderRequest.DistributorId)
		if err = tx.Create(&models.Assets{DistributorId: orderRequest.DistributorId, CurrencyCrypto: orderRequest.CurrencyCrypto}).Error; err != nil {
			utils.Log.Errorf("create distributor assets is error:%v", err)
			tx.Rollback()
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.QuantityNotEnoughErr.Data()
			return ret
		}
	}
	//创建订单前，判断平台商是否有足够的币用于交易
	if orderRequest.Direction == 1 {
		//check := CheckCoinQuantity(orderRequest.DistributorId, orderRequest.CurrencyCrypto, orderRequest.Quantity)
		if assets.Quantity < orderRequest.Quantity {
			utils.Log.Errorf("the distributor do not have enough coin so sell, distributorId= %s", orderRequest.DistributorId)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.QuantityNotEnoughErr.Data()
			return ret
		}
		//if check == false {
		//	return ""
		//}
		//给平台商锁币
		if err := tx.Model(&models.Assets{}).Where("distributor_id = ? AND currency_crypto = ? ", orderRequest.DistributorId,orderRequest.CurrencyCrypto).Updates(map[string]interface{}{"quantity": assets.Quantity - order.Quantity, "qty_frozen": assets.QtyFrozen + order.Quantity}).Error; err != nil {
			utils.Log.Errorf("the distributor lock quantity= %f, distributorId= %s", orderRequest.Quantity, orderRequest.DistributorId)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.QuantityNotEnoughErr.Data()
			return ret
		}
		//if check == false {
		//	utils.Log.Errorf("there is something wrong when checking the distributor's coin number")
		//	ret.Status = response.StatusFail
		//	ret.ErrCode, ret.ErrMsg = err_code.QuantityNotEnoughErr.Data()
		//	return ret
		//}
	}
	//创建订单
	orderRet = CreateOrder(orderRequest)
	if orderRet.Status != response.StatusSucc {
		utils.Log.Error("create order fail")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateOrderErr.Data()
		return ret
	}
	order = orderRet.Data[0]
	orderNumber := order.OrderNumber //订单id

	//2.todo 创建订单成功，回调平台服务，通知创建订单成功

	serverUrl = GetServerUrlByApiKey(req.ApiKey)
	if serverUrl == "" {
		utils.Log.Errorf("serverUrl is null")
	} else {

		resp, _ := NotifyDistributorServer(serverUrl, order)
		if resp != nil && resp.Status == SUCCESS {
			utils.Log.Debugf("create order success,serverUrl is: [%s],order is :[%v]", serverUrl, order)
		} else {
			utils.Log.Errorf("send message to distributor fail,serverUrl is: %s", serverUrl)
		}
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

	ret.Status = response.StatusSucc
	createOrderResult.OrderNumber = orderNumber
	ret.Data = []response.CreateOrderResult{createOrderResult}
	tx.Commit()
	return ret

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
	resp.CurrencyFiat = req.CurrencyFiat
	resp.AccountId = req.AccountId

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

func GenSignatureWith(mesthod string, url string, str string, apikey string) string {
	return strings.ToUpper(mesthod) + url + str + apikey
}

func GenSignatureWith2(mesthod string, url string, originOrder string, distributorId string, apikey string) string {
	return strings.ToUpper(mesthod) + url + originOrder + distributorId + apikey
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
	utils.Log.Debugf("message:%s", message)
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(message))
	if err != nil {
		return "", err
	}
	h := fmt.Sprintf("%x", mac.Sum(nil))
	utils.Log.Debugf("h is %s", h)

	//return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
	return h, nil

}

func GetServerUrlByApiKey(apikey string) string {
	if apikey == "" {
		utils.Log.Errorf("apiKey is null,can not find serverUrl according to apiKey")
		return ""
	}
	ditributor, err := GetDistributorByAPIKey(apikey)
	if err != nil {
		utils.Log.Errorf("can not get serverUrl, apiKey is %s", apikey)
		return ""
	}
	return ditributor.ServerUrl
}

//send message to distributor server
func NotifyDistributorServer(serverUrl string, order models.Order) (resp *http.Response, err error) {

	//证书认证
	pool := x509.NewCertPool()
	//根据配置文件读取证书
	//caCrt, err := ioutil.ReadFile(utils.Config.GetString("certificate.path"))
	distributorId := strconv.FormatInt(order.DistributorId, 10)
	caCrt := DownloadPem(distributorId)
	utils.Log.Debugf("capem is: %v", caCrt)

	pool.AppendCertsFromPEM(caCrt)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}

	client := &http.Client{Transport: tr}

	jsonData, err := json.Marshal(order)
	if err != nil {
		utils.Log.Errorf("order convert to json wrong,[%v]", err)
	}
	var binBody = bytes.NewReader(jsonData)
	request, err := http.NewRequest(http.MethodPost, serverUrl, binBody)
	if err != nil {
		utils.Log.Errorf("http.NewRequest wrong, err:%v", err)
		resp.Status = FAIL
		return resp, err
	}
	orderStatus := order.Status
	Headers(request)

	if orderStatus == 1 {
		resp, err = client.Do(request)
		if err != nil || resp == nil {
			utils.Log.Errorf("there is something wrong when visit distributor server,%v", err)
			return nil, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil && string(body) == SUCCESS {
			resp.Status = SUCCESS
			return resp, nil
		}

	} else if orderStatus == 4 {

		resp, err = client.Do(request)
		if err != nil || resp == nil {
			utils.Log.Errorf("there is something wrong when visit distributor server,%v", err)
			return nil, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil && string(body) == SUCCESS && UpdateOrderSyncd(order).Status == response.StatusSucc {
			resp.Status = SUCCESS
			return resp, nil
		}

		timer1 := time.NewTimer(10 * time.Minute)
		utils.Log.Debugf("wait for 10 minutes and send message for the second time.......")
		<-timer1.C
		resp, err = client.Do(request)
		if err != nil || resp == nil {
			utils.Log.Errorf("there is something wrong when visit distributor server,%v", err)
			return nil, err
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil && string(body) == SUCCESS && UpdateOrderSyncd(order).Status == response.StatusSucc {
			resp.Status = SUCCESS
			return resp, nil
		}

		timer2 := time.NewTimer(30 * time.Minute)
		utils.Log.Debugf("wait for 30 minutes and send message for the third time........")
		<-timer2.C
		resp, err = client.Do(request)
		if err != nil || resp == nil {
			utils.Log.Errorf("there is something wrong when visit distributor server,%v", err)
			return nil, err
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil && string(body) == SUCCESS && UpdateOrderSyncd(order).Status == response.StatusSucc {
			resp.Status = SUCCESS
			return resp, nil
		}

		timer3 := time.NewTimer(2 * time.Hour)
		utils.Log.Debugf("wait for 2 hours and send message for the fourth time........")
		<-timer3.C
		resp, err = client.Do(request)
		if err != nil || resp == nil {
			utils.Log.Errorf("there is something wrong when visit distributor server,%v", err)
			return nil, err
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil && string(body) == SUCCESS && UpdateOrderSyncd(order).Status == response.StatusSucc {
			resp.Status = SUCCESS
			return resp, nil
		}

	}
	resp.Status = FAIL
	resp.StatusCode = 200
	return resp, nil

}

func Headers(request *http.Request) {
	request.Header.Add(ACCEPT, APPLICATION_JSON)
	request.Header.Add(CONTENT_TYPE, APPLICATION_JSON_UTF8)
}

//下单前检验平台商币的数量
func CheckCoinQuantity(distributorId int64, currencyCrypto string, quantity float64) bool {
	assets, err := GetCoinQuantity(strconv.FormatInt(distributorId, 10), currencyCrypto)
	if err != nil {
		utils.Log.Errorf("get the coin number of distributor wrong, distributorId= %s", distributorId)
		return false
	}
	if assets.Quantity < quantity {
		utils.Log.Errorf("the distributor do not have enough coin so sell, distributorId= %s", distributorId)
		return false
	}
	return true

}
