package response

import "yuudidi.com/pkg/models"

type NotifySendMessage struct {
	JrddNotifyId    string  `json:"jrddNotifyId"`
	JrddNotifyTime  int64   `json:"jrddNotifyTime"`
	JrddOrderId     string  `json:"jrddOrderId"`
	AppOrderId      string  `json:"appOrderId"`
	OrderAmount     float64 `json:"orderAmount"`
	OrderCoinSymbol string  `json:"orderCoinSymbol"`
	OrderStatus     int     `json:"orderStatus"`
	StatusReason    int     `json:"statusReason"`
	OrderRemark     string  `json:"orderRemark"`
	OrderPayTypeId  uint    `json:"orderPayTypeId"`
	PayAccountId    string  `json:"payAccountId"`
	PayAccountUser  string  `json:"payAccountUser"`
	PayAccountInfo  string  `json:"payAccountInfo"`

	//是否发送，0：没发送，1：已经发送
	Synced uint ` json:"synced"`
	//重试次数
	Attempts uint ` json:"attempts"`
	//发送消息后是否通知成功，判断依据是返回值是否是SUCCESS，0：表示失败，1：成功
	SendStatus int `json:"sendStatus"`
	//异步通知平台商url
	AppServerNotifyUrl string ` json:"appServerNotifyUrl"`
	AppReturnPageUrl   string ` json:"appReturnPageUrl"`
}

type NotifyRet struct {
	CommonRet
	Data []models.Notify
}

//回调的消息体
type NotifyRequest struct {
	JrddNotifyId    string  `json:"jrddNotifyId"`
	JrddNotifyTime  int64   `json:"jrddNotifyTime"`
	JrddOrderId     string  `json:"jrddOrderId"`
	AppOrderId      string  `json:"appOrderId"`
	OrderAmount     float64 `json:"orderAmount"`
	OrderCoinSymbol string  `json:"orderCoinSymbol"`
	OrderStatus     int     `json:"orderStatus"`
	StatusReason    int     `json:"statusReason"`
	OrderRemark     string  `json:"orderRemark"`
	OrderPayTypeId  uint    `json:"orderPayTypeId"`
	PayAccountId    string  `json:"payAccountId"`
	PayAccountUser  string  `json:"payAccountUser"`
	PayAccountInfo  string  `json:"payAccountInfo"`
}

type NotifyListReq struct {
	OrderNumber []string `json:"orderNumber"`
}
