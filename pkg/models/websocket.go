package models

type Msg struct {
	MsgType    MsgType       `json:"msg_type"`
	ACK        MsgType       `json:"ack"`
	MerchantId []int64       `json:"merchant_id"`
	H5         []string      `json:"h5"`
	Timeout    int           `json:"timeout"`
	Data       []interface{} `json:"data"`
}

type MsgType string

const (
	// 下发订单需求
	SendOrder MsgType = "send_order"
	// 通知币商，用户订单的分配情况
	FulfillOrder MsgType = "fulfill_order"
	// 确认收款
	NotifyPaid MsgType = "notify_paid"
	// 确认付款
	ConfirmPaid MsgType = "confirm_paid"
	// 应收实付不符
	PaymentMismatch MsgType = "payment_mismatch"
	// 订单完成 转账结束
	Transferred MsgType = "transferred"
	// 接受订单
	Accept MsgType = "accept"
	// 收到请求
	ACK MsgType = "ack"
	// 开始接单
	StartOrder MsgType = "start_order"
	// 停止接单
	StopOrder MsgType = "stop_order"
)
