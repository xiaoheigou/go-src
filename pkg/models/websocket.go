package models

type Msg struct {
	MsgType    msgType  `json:"msg_type"`
	MerchantId []int64  `json:"merchant_id"`
	H5         []string `json:"h5"`
	Timeout    int      `json:"timeout"`
	Data       []Order  `json:"data"`
}

type msgType string

const (
	// 下发订单需求
	SendOrder msgType = "send_order"
	// 通知币商，用户订单的分配情况
	FulfillOrder msgType = "fulfill_order"
	// 确认收款
	NotifyPaid msgType = "notify_paid"
	// 确认付款　
	ConfirmPaid msgType = "confirm_paid"
	// 应收实付不符
	PaymentMismatch msgType = "payment_mismatch"
	// 订单完成 转账结束
	Transferred msgType = "transferred"
	// 接受订单
	Accept msgType = "accept"
	// 收到请求
	ACK msgType = "ack"
)
