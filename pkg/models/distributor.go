package models

import (
	"yuudidi.com/pkg/utils"
)

type Distributor struct {
	Id    int64    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name  string `gorm:"not null" json:"name"`
	Phone string `gorm:"type:varchar(20)" json:"phone"`
	//平台商状态 0: 申请 1: 正常 2: 冻结
	Status int      `gorm:"type:tinyint(1)" json:"status"`
	PageUrl   string   `gorm:"type:varchar(255)" json:"page_url"`
	ServerUrl string   `gorm:"type:varchar(255)" json:"server_url"`
	ApiKey    string   `gorm:"type:varchar(255)" json:"-"`
	ApiSecret string   `gorm:"type:varchar(255)" json:"-"`
	Assets    []Assets `gorm:"foreignkey:DistributorId" json:"-"`
	Quantity  string   `gorm:"-" json:"quantity"`
	Timestamp
}

func init() {
	utils.DB.AutoMigrate(&Distributor{})
}
