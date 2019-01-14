package models

import "yuudidi.com/pkg/utils"

type Tickets struct {
	Id int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	//工单id
	TicketId string `gorm:"type:varchar(36)" json:"ticket_id"`
	//订单号
	OrderNumber string `gorm:"type:varchar(191);not null" json:"order_number"`
	//工单类型（event, entry, note）
	TicketType string `gorm:"type:varchar(8)" json:"ticket_type"`
	//工单详情
	Content string `gorm:"type:varchar(128)" json:"content"`
	//工单编号
	TicketNo string `gorm:"type:varchar(36)" json:"ticket_no"`
	//工单标题
	Subject string `gorm:"type:varchar(36)" json:"subject"`
	//工单事件变更操作者
	Operator string `gorm:"type:varchar(36)" json:"operator"`
	//操作者Id
	OperatorId string `gorm:"type:varchar(36)" json:"operator_id"`
	//操作者类型（坐席-S/访客-U）
	OperatorType string `gorm:"type:varchar(36)" json:"operator_type"`
	//工单创建者Id
	CreatorId string `gorm:"type:varchar(36)" json:"creator_id"`
	//附件信息
	Attachments string `gorm:"type:text" json:"attachments"`
	//申诉类型
	ApplyType string `gorm:"type:varchar(36)" json:"apply_type"`
	//国家与地区代码
	CountryCode string `gorm:"type:varchar(36)" json:"country_code"`
	//联系电话
	Phone string `gorm:"type:varchar(32)" json:"phone"`
	//申诉原因
	ApplyMsg string `gorm:"type:text" json:"apply_msg"`
	//银行卡号
	BankCard string `gorm:"type:varchar(32)" json:"bank_card"`
	//开户银行
	BankName string `gorm:"type:varchar(32)" json:"bank_name"`
	//持卡人
	CardUser string `gorm:"type:varchar(32)" json:"card_user"`
	//微信支付账号
	WxAccount string `gorm:"type:varchar(32)" json:"wx_account"`
	//支付宝支付账号
	AliAccount string `gorm:"type:varchar(32)" json:"ali_account"`
	Timestamp
}

type TicketUpdate struct {
	Id int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	//工单id
	TicketId string `gorm:"type:varchar(36)" json:"ticket_id"`
	//工单类型（event, entry, note）
	TicketType string `gorm:"type:varchar(8)" json:"ticket_type"`
	//工单变化描述
	Description string `gorm:"type:text" json:"description"`

	//操作人昵称
	Nickname string `gorm:"type:varchar(36)" json:"nick_name"`
	//留言
	//Note  string `gorm:"type:text" json:"note"`
	//Entry string `gorm:"type:text" json:"entry"`
	Timestamp
}

func init() {
	utils.DB.AutoMigrate(&Tickets{})
	utils.DB.AutoMigrate(&TicketUpdate{})
}
