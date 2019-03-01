package models

import "yuudidi.com/pkg/utils"

type Notify struct {
	Id              int64   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	JrddNotifyId    string  `gorm:"type:varchar(255)" json:"jrdd_notify_id"`
	JrddNotifyTime  int64   `gorm:"type:int(11)" json:"jrdd_notify_time"`
	JrddOrderId     string  `gorm:"type:varchar(191)" json:"jrdd_order_id"`
	AppOrderId      string  `gorm:"type:varchar(191)" json:"app_order_id"`
	OrderAmount     float64 `gorm:"type:decimal(20,5)" json:"order_amount"`
	OrderCoinSymbol string  `gorm:"type:varchar(12)" json:"order_coin_symbol"`
	OrderStatus     int     `gorm:"type:tinyint(1)" json:"order_status"`
	StatusReason    int     `gorm:"type:tinyint(1)" json:"status_reason"`
	OrderRemark     string  `gorm:"type:varchar(191)" json:"order_remark"`
	OrderPayTypeId  uint    `gorm:"type:tinyint(2)" json:"order_pay_type_id"`
	PayAccountId    string  `gorm:"type:varchar(191)" json:"pay_account_id"`
	PayQRUrl        string  `gorm:"type:varchar(255)" json:"pay_qr_url"`
	PayAccountUser  string  `gorm:"type:varchar(191)" json:"pay_account_user"`
	PayAccountInfo  string  `gorm:"type:varchar(191)" json:"pay_account_info"`

	//是否发送，0：没发送，1：已经发送
	Synced uint `gorm:"type:tinyint(1)" json:"synced"`
	//重试次数
	Attempts uint `gorm:"type:tinyint(1)" json:"attempts"`
	//发送消息后是否通知成功，判断依据是返回值是否是SUCCESS，0：表示失败，1：重试8次之后仍未成功，2：成功
	SendStatus int `gorm:"type:tinyint(1)" json:"send_status"`
	//异步通知平台商url
	AppServerNotifyUrl string `gorm:"type:varchar(255)" json:"app_server_notify_url"`
	AppReturnPageUrl   string `gorm:"type:varchar(255)" json:"app_return_page_url"`
	//订单类型 0：买单  1：提现
	OrderType int `gorm:"type:tinyint(1)" json:"order_type"`
	Timestamp
}

func init() {
	utils.DB.AutoMigrate(&Notify{})
}
