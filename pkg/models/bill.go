package models

import (
	"yuudidi.com/pkg/utils"
)

type ReceivedBill struct {
	// 数据库自增id
	Id int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	// 上传者币商id
	UploaderUid int64 `gorm:"column:uploader_uid;index;not null" json:"uploader_uid"`
	// 账单类型 1:微信,2:支付宝
	PayType int `gorm:"column:pay_type;type:tinyint(2);unique_index:idx_pay_type_bill_id;not null" json:"pay_type"`
	// 支付宝或微信用户支付id，前端App通过hook方式可以拿到
	UserPayId string `gorm:"column:user_pay_id" json:"user_pay_id"`
	// 支付宝或微信的账单Id
	BillId string `gorm:"column:bill_id;type:varchar(191);unique_index:idx_pay_type_bill_id;not null" json:"bill_id"`
	// 从账单备注字段中得到的jrdidi订单号
	OrderNumber string `gorm:"type:varchar(191);" json:"order_number"`
	// 账单的人民币金额
	Amount float64 `gorm:"type:decimal(20,2)" json:"amount"`
	// Hook支付宝或微信账单得到的原始数据
	BillData string `gorm:"column:bill_data" json:"bill_data"`
}

func init() {
	utils.DB.AutoMigrate(&ReceivedBill{})
}
