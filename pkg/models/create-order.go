package models

type CreateOrderResult struct {
	Url          string `json:"url"`
	OrderType    string `json:"orderType"`
	TotalCount   string `json:"totalCount"`
	OrderNo      string `json:"orderNo"`
	OrderSuccess string `json:"orderSuccess"`
}
