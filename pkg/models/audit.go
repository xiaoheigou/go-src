package models

import "yuudidi.com/pkg/utils"

type AuditMessage struct {
	ID
	MerchantId    int64 `gorm:"index:IDX_MERCHANT" json:"merchant_id"`
	DistributorId int64 `gorm:"index:IDX_DISTRIBUTOR" json:"distributor_id"`
	OperatorId    int64 `gorm:"type:int(11)" json:"operator_id"`
	//联系电话
	ContactPhone string `gorm:"type:varchar(20)" json:"contact_phone"`
	//原因
	ExtraMessage string `gorm:"type:varchar(255)" json:"extra_message"`
	Timestamp
}

func (AuditMessage) TableName() string {
	return "audit_message"
}

func init() {
	utils.DB.AutoMigrate(&AuditMessage{})
}
