package models

import (
	"time"

	"yuudidi.com/pkg/utils"
)

// Fulfillment - 订单分配历史，数据库表fulfillment_events
type Fulfillment struct {
	// ID - PK
	ID int `gorm:"primary_key;AUTO_INCREMENT;type:bigint(20)" json:"-"`
	// 订单编号
	OrderNumber int64 `gorm:"type:bigint(20);column:order_number" json:"order_number"`
	// SeqID - sequence id
	SeqID int `gorm:"type:int(2);column:seq_id" json:"seq_id"`
	// 承兑商ID
	MerchantID int `gorm:"type:int(11);column:merchant_id" json:"merchant_id"`
	// 承兑商收款方式ID （order_type = buy)
	MerchantPaymentID int `gorm:"type:int(11);column:merchant_payment_id" json:"merchant_payment_id"`
	// 派单接收时间
	AcceptedAt time.Time `gorm:"column:accepted_at" json:"accepted_at"`
	// 通知支付时间
	PaidAt time.Time `gorm:"column:paid_at" json:"paid_at"`
	// 确认支付时间
	PaymentConfirmedAt time.Time `gorm:"column:payment_confirmed_at" json:"payment_confirmed_at"`
	// 转账时间
	TransferredAt time.Time `gorm:"column:transferred_at" json:"transferred_at"`
	// Status - 订单执行状态
	Status int `gorm:"type:tinyint(1)"`
	// TimeStamp - 创建/更新/删除时间
	Timestamp
}

// TableName - Fulfillment table named as fulfillment_events
func (Fulfillment) TableName() string {
	return "fulfillment_events"
}

func init() {
	utils.DB.AutoMigrate(&Fulfillment{})
}
