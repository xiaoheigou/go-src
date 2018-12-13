package models

import (
	"yuudidi.com/pkg/utils"
)

type AssetHistory struct {
	Id int64 `gorm:"primary_key;AUTO_INCREMENT;type:bigint(20)" json:"Id"`
	// 承兑商ID
	MerchantId int64 `gorm:"type:int(11)" json:"merchant_id" example:"123"`
	//平台商ID
	DistributorId int64 `gorm:"type:int(11)" json:"distributor_id" example:"123"`
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
	OperatorId int64 `gorm:"type:int(11)" json:"operator_id" example:"1"`
	//操作者名称
	OperatorName string `gorm:"-" json:"operator_name" example:"test"`
	Timestamp
}

type AssetApply struct {
	ID
	// 承兑商ID
	MerchantId int64 `gorm:"type:int(11)" json:"merchant_id" example:"123"`
	// 承兑商手机号
	Phone string `gorm:"varchar(30)" json:"phone"`
	// 承兑商邮箱
	Email string `gorm:"varchar(50)" json:"email"`
	// 充值申请状态 0/1 未审核/已审核
	Status int64 `gorm:"type:tinyint(1);default:0" json:"status"`
	// 币种
	Currency string `gorm:"type:varchar(20)" json:"currency" example:"BTUSD"`
	// 充值数量
	Quantity float64 `gorm:"type:Decimal(15,5)" json:"quantity" example:"123"`
	// 剩余数量
	RemainQuantity float64 `gorm:"-" json:"remain_quantity" example:"123"`
	// 申请人ID
	ApplyId int64 `gorm:"type:int(11)" json:"apply_id"`
	// 申请人username
	ApplyName string `gorm:"-" json:"apply_name"`
	// 审核人ID
	AuditorId int64 `gorm:"type:int(11)" json:"auditor_id"`
	// 审核人姓名
	AuditorName int `gorm:"-" json:"auditor_name"`
	Timestamp
}

func init() {
	utils.DB.AutoMigrate(&AssetHistory{})
	utils.DB.AutoMigrate(&AssetApply{})
}
