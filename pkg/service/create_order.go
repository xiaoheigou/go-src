package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
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
	FAIL                  = "FAIL"
	SUCCESS               = "SUCCESS"
	APPLICATION_JSON      = "application/json"
	APPLICATION_JSON_UTF8 = "application/json; charset=UTF-8"
)

//下订单
func PlaceOrder(req response.CreateOrderRequest) response.CreateOrderRet {
	//var orderRequest response.OrderRequest
	var order models.Order
	//var serverUrl string
	var ret response.CreateOrderRet

	//1. 创建订单
	orderRequest, err := PlaceOrderReq2CreateOrderReq(req)
	if err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoDistributorFindErr.Data()
		return ret
	}
	utils.Log.Debugf("orderRequest = [%+v]", orderRequest)

	//distributorId := strconv.FormatInt(orderRequest.DistributorId, 10)
	//currencyCrypto := orderRequest.CurrencyCrypto

	tx := utils.DB.Begin()

	//创建订单
	order = models.Order{
		OrderNumber: GenerateOrderNumByFastId(),
		Price:       orderRequest.Price,
		OriginOrder: orderRequest.OriginOrder,
		//成交量
		Quantity: orderRequest.Quantity,
		//成交额
		Amount:     orderRequest.Amount,
		PaymentRef: orderRequest.PaymentRef,
		//订单状态，0/1分别表示：未支付的/已支付的
		Status: 1,
		//订单类型，1为买入，2为卖出
		Direction:         orderRequest.Direction,
		DistributorId:     orderRequest.DistributorId,
		MerchantId:        orderRequest.MerchantId,
		MerchantPaymentId: orderRequest.MerchantPaymentId,
		//扣除用户佣金金额
		TraderCommissionAmount: orderRequest.TraderCommissionAmount,
		//扣除用户佣金币的量
		TraderCommissionQty: orderRequest.TraderCommissionQty,
		//用户佣金比率
		TraderCommissionPercent: orderRequest.TraderCommissionPercent,
		//扣除币商佣金金额
		MerchantCommissionAmount: orderRequest.MerchantCommissionAmount,
		//扣除币商佣金币的量
		MerchantCommissionQty: orderRequest.MerchantCommissionQty,
		//币商佣金比率
		MerchantCommissionPercent: orderRequest.MerchantCommissionPercent,
		//平台扣除的佣金币的量（= trader_commision_qty+merchant_commision_qty)
		PlatformCommissionQty: orderRequest.PlatformCommissionQty,
		//平台商用户id
		AccountId: orderRequest.AccountId,
		//交易币种
		CurrencyCrypto: orderRequest.CurrencyCrypto,
		//交易法币
		CurrencyFiat: orderRequest.CurrencyFiat,
		//交易类型 0:微信,1:支付宝,2:银行卡
		PayType: orderRequest.PayType,
		//微信或支付宝二维码地址
		QrCode: orderRequest.QrCode,
		//微信或支付宝账号
		Name: orderRequest.Name,
		//银行账号
		BankAccount: orderRequest.BankAccount,
		//所属银行
		Bank: orderRequest.Bank,
		//所属银行分行
		BankBranch:         orderRequest.BankBranch,
		Fee:                orderRequest.Fee,
		OriginAmount:       orderRequest.OriginAmount,
		Price2:             orderRequest.Price2,
		AppCoinName:        orderRequest.AppCoinName,
		Remark:             orderRequest.Remark,
		AppReturnPageUrl:   orderRequest.AppReturnPageUrl,
		AppServerNotifyUrl: orderRequest.AppServerNotifyUrl,
	}
	if db := tx.Create(&order); db.Error != nil {
		tx.Rollback()
		utils.Log.Errorf("tx in func PlaceOrder rollback")
		utils.Log.Error("create order fail")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateOrderErr.Data()
		return ret
	}
	orderNumber := order.OrderNumber //订单id

	//异步保存用户信息
	AsynchronousSaveAccountInfo(order)

	//utils.Log.Debugf("get the coin number of distributor wrong,to create, distributorId= %s", orderRequest.DistributorId)
	var assets models.Assets

	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&assets, "distributor_id=? and currency_crypto=?", orderRequest.DistributorId, orderRequest.CurrencyCrypto).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// 没有找到，则为平台商创建asset记录
			if err := tx.Create(&models.Assets{DistributorId: orderRequest.DistributorId, CurrencyCrypto: orderRequest.CurrencyCrypto}).Error; err != nil {
				utils.Log.Errorf("create distributor assets fail: %v", err)
				tx.Rollback()
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.DatabaseErr.Data()
				return ret
			}
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&assets, "distributor_id=? and currency_crypto=?", orderRequest.DistributorId, orderRequest.CurrencyCrypto).Error; err != nil {
				utils.Log.Errorf("find distributor assets fail: %v", err)
				tx.Rollback()
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.DatabaseErr.Data()
				return ret
			}
		} else {
			// 其它错误
			utils.Log.Errorf("find distributor assets fail: %v", err)
			tx.Rollback()
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.DatabaseErr.Data()
			return ret
		}
	}

	//判断平台商是否有足够的币用于交易，并冻结相应的币
	if orderRequest.Direction == 1 {
		utils.Log.Debugf("distributor (id=%d) quantity = [%d], order (%s) quantity = [%d]", orderRequest.DistributorId, assets.Quantity, orderRequest.OrderNumber, orderRequest.Quantity)
		//给平台商锁币
		if rowsAffected := tx.Model(&models.Assets{}).Where("distributor_id = ? AND currency_crypto = ? AND quantity >= ?", orderRequest.DistributorId, orderRequest.CurrencyCrypto, orderRequest.Quantity).
			Updates(map[string]interface{}{"quantity": assets.Quantity - orderRequest.Quantity, "qty_frozen": assets.QtyFrozen + orderRequest.Quantity}).RowsAffected; rowsAffected == 0 {
			tx.Rollback()
			utils.Log.Errorf("tx in func PlaceOrder rollback")
			utils.Log.Errorf("the distributor (distributor_id=%s) only has %f %s, but want to freeze %f. Operation fail. assert for distributor = [%+v]",
				orderRequest.DistributorId, assets.Quantity, orderRequest.CurrencyCrypto, orderRequest.Quantity, assets)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.QuantityNotEnoughErr.Data()
			return ret
		}

	}

	tx.Commit()
	utils.Log.Debugf("tx in func PlaceOrder commit")

	//2. 创建订单成功，重定向
	//createurl := utils.Config.Get("redirecturl.createurl")
	//url := fmt.Sprintf("%v", createurl)
	//orderStr, _ := Struct2JsonString(order)
	//c.Request.Header.Add("order", orderStr)
	//c.Redirect(301, url)
	//serverUrl = GetServerUrlByApiKey(req.ApiKey)
	//serverUrl = order.AppServerNotifyUrl

	//3.异步通知平台商
	AsynchronousNotify( order)
	//4. 调用派单服务

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
	createOrderResult := response.CreateOrderResult{
		OrderNumber: orderNumber,
		//RedirectUrl:    url,
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
	ret.Data = []response.CreateOrderResult{createOrderResult}
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

func BuyOrderReq2CreateOrderReq(buyOrderReq response.BuyOrderRequest) response.CreateOrderRequest {
	var req response.CreateOrderRequest
	req = response.CreateOrderRequest{
		ApiKey:        buyOrderReq.AppApiKey,
		OrderNo:       buyOrderReq.AppOrderNo,
		Price:         buyOrderReq.AppCoinRate,
		Amount:        buyOrderReq.OrderCoinAmount,
		DistributorId: buyOrderReq.AppId,
		AppCoinName:   buyOrderReq.AppCoinName,
		OrderType:     0,
		TotalCount:    0,
		PayType:       buyOrderReq.OrderPayTypeId,

		Remark: buyOrderReq.OrderRemark,
		//页面回调地址
		PageUrl: buyOrderReq.AppReturnPageUrl,
		//服务端回调地址
		ServerUrl:    buyOrderReq.AppServerNotifyUrl,
		CurrencyFiat: buyOrderReq.AppCoinSymbol,
		AccountId:    buyOrderReq.AppUserId,
	}

	return req

}

func SellOrderReq2CreateOrderReq(sellOrderReq response.SellOrderRequest) response.CreateOrderRequest {

	var req response.CreateOrderRequest
	req = response.CreateOrderRequest{
		ApiKey:        sellOrderReq.AppApiKey,
		OrderNo:       sellOrderReq.AppOrderNo,
		Price:         sellOrderReq.AppCoinRate,
		Amount:        sellOrderReq.OrderCoinAmount,
		DistributorId: sellOrderReq.AppId,
		AppCoinName:   sellOrderReq.AppCoinName,
		OrderType:     1,
		TotalCount:    0,
		PayType:       sellOrderReq.OrderPayTypeId,
		Name:          sellOrderReq.PayAccountUser,
		BankAccount:   sellOrderReq.PayAccountId,
		Bank:          "",
		BankBranch:    sellOrderReq.PayAccountInfo,
		Phone:         "",
		Remark:        sellOrderReq.OrderRemark,
		QrCode:        sellOrderReq.PayQRUrl,
		//页面回调地址
		PageUrl: sellOrderReq.AppReturnPageUrl,
		//服务端回调地址
		ServerUrl:    sellOrderReq.AppServerNotifyUrl,
		CurrencyFiat: sellOrderReq.AppCoinSymbol,
		AccountId:    sellOrderReq.AppUserId,
	}

	return req

}

func PlaceOrderReq2CreateOrderReq(req response.CreateOrderRequest) (response.OrderRequest, error) {
	var resp response.OrderRequest
	var fee float64
	var originAmount float64
	var amount float64
	var quantity float64

	buyPrice := utils.Config.GetFloat64("price.buyprice")
	sellPrice := utils.Config.GetFloat64("price.sellprice")
	originAmount = req.Amount
	if req.OrderType == 0 {
		amount = req.Amount
		quantity = originAmount / buyPrice
	} else {
		var distributor models.Distributor
		var appCoinSymbol string
		if req.CurrencyFiat != "" {
			appCoinSymbol = req.CurrencyFiat
		} else {
			appCoinSymbol = "CNY"
		}
		if err := utils.DB.First(&distributor, "distributors.id = ? and distributors.app_coin_symbol=?", req.DistributorId, appCoinSymbol).Error; err != nil {
			utils.Log.Errorf("func AsynchronousNotifyDistributor, not found distributor err:%v", err)
			return response.OrderRequest{},err
		}
		appCoinRate := distributor.AppCoinRate
		appUserWithdrawalFeeRate := distributor.AppUserWithdrawalFeeRate
		appCNY := originAmount * float64(appCoinRate)
		quantity = appCNY / buyPrice
		amount = originAmount * float64((100-appUserWithdrawalFeeRate)/100)

		//amount = originAmount * sellPrice / buyPrice
		//quantity = originAmount / buyPrice
		fee = originAmount - amount

	}

	resp.Price = float32(buyPrice)
	resp.Amount = amount
	resp.DistributorId = req.DistributorId
	resp.Quantity = quantity
	resp.OriginOrder = req.OrderNo
	resp.CurrencyCrypto = "BTUSD"
	resp.Direction = req.OrderType
	resp.PayType = req.PayType
	resp.Name = req.Name
	resp.BankAccount = req.BankAccount
	resp.Bank = req.Bank
	resp.BankBranch = req.BankBranch
	resp.QrCode = req.QrCode
	resp.CurrencyFiat = req.CurrencyFiat
	resp.AccountId = req.AccountId
	resp.OriginAmount = originAmount
	resp.Fee = fee
	resp.Price2 = float32(sellPrice)
	resp.AppCoinName = req.AppCoinName
	resp.Remark = req.Remark
	resp.AppReturnPageUrl = req.PageUrl
	resp.AppServerNotifyUrl = req.ServerUrl

	return resp,nil

}

//保存用户信息
func SaveAccountIdInfo(order models.Order) {
	var accountInfo models.AccountInfo
	accountInfo = models.AccountInfo{
		AccountId:     order.AccountId,
		DistributorId: order.DistributorId,
		OrderNumber:   order.OrderNumber,
		//成交方向，以发起方（平台商用户）为准。0表示平台商用户买入，1表示平台商用户卖出。
		Direction: order.Direction,
		Price:     order.Price,
		//成交量
		Quantity: order.Quantity,
		//成交额
		Amount: order.Amount,
		//交易币种
		CurrencyCrypto: order.CurrencyCrypto,
		//交易法币
		CurrencyFiat: order.CurrencyFiat,
		//交易类型
		PayType: order.PayType,
		//微信或支付宝二维码地址
		QrCode: order.QrCode,
		//微信或支付宝账号
		Name: order.Name,
		//银行账号
		BankAccount: order.BankAccount,
		//所属银行
		Bank: order.Bank,
		//所属银行分行
		BankBranch: order.BankBranch,
	}
	if err := utils.DB.Create(&accountInfo).Error; err != nil {
		utils.Log.Errorf("save accountInfo wrong  after creating order,err:[%v]", err)
	}

}

//异步执行保存用户信息
func AsynchronousSaveAccountInfo(order models.Order) {
	go func() {
		SaveAccountIdInfo(order)
	}()
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
	apiKeyStr := "apiKey:" + apiKey
	if apiKey == "" {
		utils.Log.Error("apiKey is null")
		return ""
	}
	secretKey, err := utils.RedisClient.Get(apiKeyStr).Result()
	if err == nil && secretKey != "" {
		return secretKey

	}
	ditributor, err := GetDistributorByAPIKey(apiKey)

	if err != nil {
		utils.Log.Error("can not get secretkey according to apiKey=[%s] ", apiKey)
		return ""

	}
	secretKey = ditributor.ApiSecret
	utils.RedisSet(apiKeyStr, secretKey, 30*time.Minute)
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

func AsynchronousNotifyDistributor(order models.Order) {
	//var distributor models.Distributor
	//if err := utils.DB.First(&distributor, "distributors.id = ?", order.DistributorId).Error; err != nil {
	//	utils.Log.Errorf("func AsynchronousNotifyDistributor, not found distributor err:%v", err)
	//	return
	//}
	//AsynchronousNotify(distributor.ServerUrl, order)
	AsynchronousNotify( order)
}

func AsynchronousNotify( order models.Order) {
	serverUrl:=order.AppServerNotifyUrl
	if serverUrl == "" {
		utils.Log.Errorf("serverUrl is null")
	} else {

		go func() {
			resp, err := NotifyDistributorServer(order)
			if err == nil && resp != nil && resp.Status == SUCCESS {
				utils.Log.Debugf("send message to distributor success,serverUrl is: [%s]", serverUrl)
			} else {
				utils.Log.Errorf("send message to distributor fail,serverUrl is: [%s],err is:[%v]", serverUrl, err)
			}
		}()

	}
}

//send message to distributor server
func NotifyDistributorServer(order models.Order) (resp *http.Response, err error) {
	var serverUrl string
	var notifyRequest response.ServerNotifyRequest
	notifyRequest = Order2ServerNotifyReq(order)
	notifyRequestStr, _ := Struct2JsonString(notifyRequest)
	distributorId := strconv.FormatInt(order.DistributorId, 10)

	var distributor models.Distributor
	if err := utils.DB.First(&distributor, "distributors.id = ?", order.DistributorId).Error; err != nil {
		utils.Log.Errorf("func AsynchronousNotifyDistributor, not found distributor err:%v", err)
		resp.Status = response.StatusFail
		return resp, err
	}
	apiKey := distributor.ApiKey
	secretKey := distributor.ApiSecret
	jrddSignContent, _ := HmacSha256Base64Signer(notifyRequestStr, secretKey)
	serverUrl += order.AppServerNotifyUrl + "?" +"appId="+ distributorId +"&apiKey="+apiKey + "&jrddSignType=HMAC-SHA256" + "&" +"jrddSignContent="+ jrddSignContent + "&" + "jrddInputCharset=UTF-8"

	//证书认证
	pool := x509.NewCertPool()
	//根据配置文件读取证书
	//caCrt, err := ioutil.ReadFile(utils.Config.GetString("certificate.path"))

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

	} else if orderStatus == 7 {

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

func Order2ServerNotifyReq(order models.Order) response.ServerNotifyRequest {
	var req response.ServerNotifyRequest
	time := time.Now().Unix()
	req = response.ServerNotifyRequest{
		JrddNotifyId:    GenerateOrderNumber(),
		JrddNotifyTime:  time,
		JrddOrderId:     order.OrderNumber,
		AppOrderId:      order.OriginOrder,
		OrderAmount:     order.Amount,
		OrderCoinSymbol: order.CurrencyCrypto,
		OrderStatus:     int(order.Status),
		StatusReason:    int(order.StatusReason),
		OrderRemark:     order.Remark,
		OrderPayTypeId:  order.PayType,
		PayAccountId:    order.BankAccount,
		PayAccountUser:  order.Name,
		PayAccountInfo:  order.BankBranch,
	}
	return req
}

func Headers(request *http.Request) {
	request.Header.Add(ACCEPT, APPLICATION_JSON)
	request.Header.Add(CONTENT_TYPE, APPLICATION_JSON_UTF8)
}

func Redirect301Handler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://taadis.com", http.StatusMovedPermanently)
}
