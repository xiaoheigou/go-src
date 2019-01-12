package models

import "yuudidi.com/pkg/utils"

type Tickets struct {
	Id           int64  `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	TicketId     string `gorm:"type:varchar(36)" json:"ticket_id"`
	OrderNumber  string `gorm:"type:varchar(191);not null" json:"order_number"`
	TicketNo     string `gorm:"type:varchar(36)" json:"ticket_no"`
	Subject      string `gorm:"type:varchar(36)" json:"subject"`
	Operator     string `gorm:"type:varchar(36)" json:"operator"`
	OperatorId   string `gorm:"type:varchar(36)" json:"operator_id"`
	OperatorType string `gorm:"type:varchar(36)" json:"operator_type"`
	CreatorId    string `gorm:"type:varchar(36)" json:"creator_id"`
	Timestamp
}

type TicketUpdate struct {
	Id         int64  `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	TicketId   string `gorm:"type:varchar(36)" json:"ticket_id"`
	Description string `gorm:"type:text" json:"description"`
	OperatorId string `gorm:"type:varchar(36)" json:"operator_id"`
	Nickname   string `gorm:"type:varchar(36)" json:"nick_name"`
	Timestamp
}

func init() {
	utils.DB.AutoMigrate(&Tickets{})
	utils.DB.AutoMigrate(&TicketUpdate{})
}
