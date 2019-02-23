package models

type Msg struct {
	MsgType    MsgType       `json:"msg_type"`
	MsgId      string        `json:"msg_id"`
	ACK        MsgType       `json:"ack"`
	MerchantId []int64       `json:"merchant_id"`
	H5         []string      `json:"h5"`
	Timeout    int           `json:"timeout"`
	Data       []interface{} `json:"data"`
}

type OrderData struct {
	OrderNumber      string `json:"order_number"`
	DistributorId    int64  `json:"distributor_id"`
	AppReturnPageUrl string `json:"app_return_page_url"`
}

type MsgType string

const (
	// 下发订单需求
	SendOrder MsgType = "send_order"
	// 通知币商，用户订单的分配情况
	FulfillOrder MsgType = "fulfill_order"
	// 确认付款
	NotifyPaid MsgType = "notify_paid"
	// 确认收款
	ConfirmPaid MsgType = "confirm_paid"
	// 自动确认收款
	AutoConfirmPaid MsgType = "auto_confirm_paid"
	// 自动确认收款
	ServerConfirmPaid MsgType = "server_confirm_paid"
	// 应收实付不符
	PaymentMismatch MsgType = "payment_mismatch"
	// 订单完成 转账结束
	Transferred MsgType = "transferred"
	// 接受订单
	Accept MsgType = "accept"
	// 抢单失败
	Picked MsgType = "picked"
	// 收到请求
	ACK MsgType = "ack"
	// ping消息
	PING MsgType = "ping"
	// pong消息
	PONG MsgType = "pong"
	// 多次派单，都没有人接单，通知h5
	AcceptTimeout MsgType = "accept_timeout"
)

//Data - data field of Msg
type Data struct {
	OrderNumber string `json:"order_number"`
	Direction   int    `json:"direction"`
}
