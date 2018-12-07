package response


type MerchantOrder struct {
	// 订单号码
	OrderNum int `json:"order_num" example:1`
	// 订单类型
	OrderType int `json:"order_type" example:1`
	// 订单状态
	OrderStatus int `json:"order_status" example:1`
	// 订单金额
	TotalPrice string `json:"total_price" example:"650"`
}

type GetOrderRet struct {
	CommonRet
	Entity struct{
		Data []MerchantOrder `json:"data"`
	}
}