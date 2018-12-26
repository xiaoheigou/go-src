package models

import (
	"time"

	"yuudidi.com/pkg/utils"
)

// Fulfillment - 订单分配历史，数据库表fulfillment_events
type Fulfillment struct {
	// ID - PK
	ID
	// 订单编号
	OrderNumber string `gorm:"type:varchar(191);column:order_number" json:"order_number"`
	// SeqID - sequence id
	SeqID int `gorm:"type:int(2);column:seq_id" json:"seq_id"`
	// 承兑商ID
	MerchantID int64 `gorm:"type:int(11);column:merchant_id" json:"merchant_id"`
	// 承兑商收款方式ID （order_type = buy)
	MerchantPaymentID int64 `gorm:"type:int(11);column:merchant_payment_id" json:"merchant_payment_id"`
	// 派单接收时间
	AcceptedAt time.Time `gorm:"column:accepted_at" json:"accepted_at"`
	// 通知支付时间
	PaidAt time.Time `gorm:"column:paid_at" json:"paid_at"`
	// 确认支付时间
	PaymentConfirmedAt time.Time `gorm:"column:payment_confirmed_at" json:"payment_confi rmed_at"`
	// 转账时间
	TransferredAt time.Time `gorm:"column:transferred_at" json:"transferred_at"`

	// Status - 订单执行状态
	Status OrderStatus `gorm:"type:tinyint(1)"`
	// fulfillmentLogs - 派单日志
	FulfillmentLogs []FulfillmentLog `gorm:"foreignkey:FulfillmentId" json:"-"`
	// TimeStamp - 创建/更新/删除时间
	Timestamp
}

//FulfillmentLog - every log records an operation upon fulfillment processing
type FulfillmentLog struct {
	ID
	FulfillmentID int64 `gorm:"type:bigint(20)" json:"-"`
	// 订单编号
	OrderNumber string `gorm:"type:varchar(191);column:order_number;index:IDX_ORDER" json:"order_number"`
	// SeqID - sequence id
	SeqID int `gorm:"type:int(2);column:seq_id" json:"seq_id"`
	// 是否系统操作 0/1 不是/是
	IsSystem      bool   `gorm:"type:tinyint(1);column:is_system" json:"is_system"`
	MerchantID    int64  `gorm:"index:IDX_MERCHANT" json:"merchant_id"`
	AccountID     string `gorm:"type:varchar(40);index:IDX_ACCOUNT" json:"account_id"`
	OriginOrder   string `gorm:"type:varchar(191);unique_index:IDX_ORDER_DISTRIBUTOR" json:"origin_order"`
	DistributorID int64  `gorm:"unique_index:IDX_ORDER_DISTRIBUTOR" json:"distributor_id"`
	//订单起始状态
	OriginStatus OrderStatus `gorm:"tinyint(1)" json:"origin_status"`
	//订单修改后状态
	UpdatedStatus OrderStatus `gorm:"tinyint(1)" json:"updated_status"`
	//额外信息
	ExtraMessage string `gorm:"type:varchar(255)" json:"extra_message"`
	Timestamp
}

// TableName - Fulfillment table named as fulfillment_events
func (FulfillmentLog) TableName() string {
	return "fulfillment_logs"
}

// TableName - Fulfillment table named as fulfillment_events
func (Fulfillment) TableName() string {
	return "fulfillment_events"
}

func init() {
	utils.DB.AutoMigrate(&Fulfillment{})
	utils.DB.AutoMigrate(&FulfillmentLog{})
}
