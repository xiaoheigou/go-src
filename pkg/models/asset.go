package models

import (
	"yuudidi.com/pkg/utils"
)

type AssetHistory struct {
	ID
	// 承兑商ID
	MerchantId int64 `gorm:"type:int(11)" json:"merchant_id" example:"123"`
	//平台商ID
	DistributorId int64 `gorm:"type:int(11)" json:"distributor_id" example:"123"`
	// 订单编号
	OrderNumber string `gorm:"type:varchar(191)" json:"order_number" example:"123"`
	//成交方向，以发起方（平台商用户）为准。0表示平台商用户买入，1表示平台商用户卖出。
	Direction int `gorm:"type:tinyint(1)" json:"direction"`
	//是否是订单 0 不是 1 是
	IsOrder int `gorm:"type:tinyint(1)" json:"is_order" example:"0"`
	//操作0:充值申请,1:充值审核,2:放币,3:解冻
	Operation int `gorm:"type:tinyint(1)" json:"operation" example:"0"`
	//币种
	Currency string `gorm:"type:varchar(20);column:currency_crypto" json:"currency" example:"BTUSD"`
	//数量
	Quantity float64 `gorm:"type:Decimal(30,10)" json:"quantity" example:"123"`
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
	Status int `gorm:"type:tinyint(1);default:0" json:"status"`
	// 币种
	Currency string `gorm:"type:varchar(20)" json:"currency" example:"BTUSD"`
	// 充值数量
	Quantity float64 `gorm:"type:Decimal(20,5)" json:"quantity" example:"123"`
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
