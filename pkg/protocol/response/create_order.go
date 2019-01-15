package response

type CreateOrderResult struct {
	OrderNumber string `json:"orderNumber"`
}

type CreateOrderRet struct {
	CommonRet
	Data []CreateOrderResult `json:"data`
}

type CreateOrderRequest struct {
	ApiKey        string  `json:"apiKey"`
	OrderNo       string  `json:"orderNo"`
	Price         float32 `json:"price,string"`
	Amount        float64 `json:"amount,string"`
	DistributorId int64   `json:"distributorId,string"`
	CoinType      string  `json:"coinType"`
	OrderType     int     `json:"orderType,string"`
	TotalCount    float64 `json:"totalCount,string"`
	PayType       uint    `json:"payType,string"`
	Name          string  `json:"name"`
	BankAccount   string  `json:"bankAccount"`
	Bank          string  `json:"bank"`
	BankBranch    string  `json:"bankBranch"`
	Phone         string  `json:"phone"`
	Remark        string  `json:"remark"`
	QrCode        string  `json:"qrCode"`
	//页面回调地址
	PageUrl string `json:"pageUrl"`
	//服务端回调地址
	ServerUrl    string  `json:"serverUrl"`
	CurrencyFiat string  `json:"currencyFiat"`
	AccountId    string  `json:"accountId"`
	OriginAmount float64 `json:"originAmount,string"`
	Fee          float64 `json:"fee"`
	Price2       float32 `json:"price2,string"`
	AppCoinName  string  `json:"appCoinName"`
}

type PartnerId struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

type BuyOrderRequest struct {
	//平台商的公钥
	AppApiKey string `json:"appApiKey"`
	//平台商的ID
	AppId int64 `json:"appId,string"`
	//平台商当前用用户的ID
	AppUserId string `json:"appUserId"`
	//签名算法名称
	AppSignType string `json:"appSignType"`
	//订单完成后⻚页面面跳转地址
	AppReturnPageUrl string `json:"appReturnPageUrl"`
	//订单完成后异步调用用传参通知给平台商的API地址
	AppServerAPI string `json:"appServerAPI"`
	//平台商生生成的订单ID
	AppOrderNo string `json:"appOrderNo"`
	//币种名称
	AppCoinName string `json:"appCoinName"`
	//币种符号
	AppCoinSymbol string `json:"appCoinSymbol"`
	//币种和BTUSD之间的汇率值
	AppCoinRate float32 `json:"appCoinRate,string"`
	//本次订单中下单的币的数量
	OrderCoinAmount float64 `json:"orderCoinAmount,string"`
	//本次订单中使用用的收款方方式
	OrderPayTypeId uint `json:"orderPayTypeId,string"`
	//订单备注
	OrderRemark string `json:"orderRemark"`
}

type SellOrderRequest struct {
	//平台商的公钥
	AppApiKey string `json:"appApiKey"`
	//平台商的ID
	AppId int64 `json:"appId,string"`
	//平台商当前用用户的ID
	AppUserId string `json:"appUserId"`
	//签名算法名称
	AppSignType string `json:"appSignType"`
	//订单完成后⻚页面面跳转地址
	AppReturnPageUrl string `json:"appReturnPageUrl"`
	//订单完成后异步调用用传参通知给平台商的API地址
	AppServerAPI string `json:"appServerAPI"`
	//平台商生生成的订单ID
	AppOrderNo string `json:"appOrderNo"`
	//币种名称
	AppCoinName string `json:"appCoinName"`
	//币种符号
	AppCoinSymbol string `json:"appCoinSymbol"`
	//币种和BTUSD之间的汇率值
	AppCoinRate float32 `json:"appCoinRate,string"`
	//本次订单中下单的币的数量
	OrderCoinAmount float64 `json:"orderCoinAmount,string"`
	//本次订单中使用用的收款方方式
	OrderPayTypeId uint `json:"orderPayTypeId,string"`
	//收款账户
	PayAccountId string `json:"payAccountId"`
	//收款账户姓名
	PayAccountUser string `json:"payAccountUser"`
	//收款账户信息
	PayAccountInfo string `json:"payAccountInfo"`
	//订单备注
	OrderRemark string `json:"orderRemark"`
}
