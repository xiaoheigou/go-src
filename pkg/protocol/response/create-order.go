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
	PartnerId   string `json:"partnerId" binding:"required" example:"abcd123"`
	OrderNo     string `json:"orderNo" binding:"required" example:"2"`
	CoinType    string `json:"coinType" binding:"required" example:"2"`
	OrderType   string `json:"orderType" binding:"required" example:"12"`
	TotalCount  string `json:"totalCount" binding:"required" example:"12"`
	PayType     string `json:"payType" binding:"required" example:"1"`
	Name        string `json:"name" binding:"required" example:"hahah"`
	BankAccount string `json:"bankAccount" binding:"required" example:"test"`
	Bank        string `json:"bank" binding:"required" example:"china"`
	BankBranch  string `json:"bankBranch" binding:"required" example:"test"`
	Phone       string `json:"phone" binding:"required" example:"1380000000"`
	Remark      string `json:"remark" binding:"required" example:"test"`
}
