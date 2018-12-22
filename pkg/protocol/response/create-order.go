package response

type CreateOrderResult struct {
	Url          string `json:"url"`
	OrderType    string `json:"orderType"`
	TotalCount   string `json:"totalCount"`
	OrderNo      string `json:"orderNo"`
	OrderSuccess string `json:"orderSuccess"`
}

type CreateOrderRet struct {
	CommonRet
	Data [] CreateOrderResult `json:"data`
}

type CreateOrderRequest struct {
	PartnerId     PartnerId `json:"partnerId"`
	OrderNo       string    `json:"orderNo"`
	Price         float32   `json:"price"`
	Amount        float64    `json:"amount"`
	DistributorId int64    `json:"distributorId"`
	CoinType      string    `json:"coinType"`
	OrderType     int       `json:"orderType"`
	TotalCount    string    `json:"totalCount"`
	PayType       uint      `json:"payType"`
	Name          string    `json:"name"`
	BankAccount   string    `json:"bankAccount"`
	Bank          string    `json:"bank"`
	BankBranch    string    `json:"bankBranch"`
	Phone         string    `json:"phone"`
	Remark        string    `json:"remark"`
	QrCode        string    `json:"qrCode"`
	//页面回调地址
	PageUrl string `json:"pageUrl"`
	//服务端回调地址
	ServerUrl string `json:"serverUrl"`
}

type PartnerId struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}
