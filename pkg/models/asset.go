package models

import (
	"time"
	"yuudidi.com/pkg/utils"
)

type AssetHistory struct {
	Id int `gorm:"primary_key;AUTO_INCREMENT;type:bigint(20)" json:"Id"`
	// 承兑商ID
	MerchantId int `gorm:"type:int(11)" json:"merchant_id" example:"123"`
	//平台商ID
	DistributorId int `gorm:"type:int(11)" json:"distributor_id" example:"123"`
	// 订单编号
	OrderNumber int64 `gorm:"type:bigint(20)" json:"order_number" example:"123"`
	//是否是订单 0 不是 1 是
	IsOrder int `gorm:"type:tinyint(1)" json:"is_order" example:"0"`
	//操作0:充值申请,1:充值审核
	Operation int `gorm:"type:tinyint(1)" json:"operation" example:"0"`
	//币种
	Currency string `gorm:"type:varchar(20)" json:"currency" example:"BTUSD"`
	//数量
	Quantity float64 `gorm:"type:Decimal(15,5)" json:"quantity" example:"123"`
	//操作者id
	OperatorId int `gorm:"type:int(11)" json:"operator_id" example:"1"`
	//操作者名称
	OperatorName string    `gorm:"-" json:"operator_name" example:"test"`
	Timestamp    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"timestamp"`
}

func init() {
	utils.DB.AutoMigrate(&AssetHistory{})
}
