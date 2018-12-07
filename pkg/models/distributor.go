package models

import (
	"yuudidi.com/pkg/utils"
)

type Distributor struct {
	Id        int        `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name      string     `gorm:"not null" json:"name"`
	Status    int        `gorm:"type:tinyint(2)" json:"status"`
	PageUrl   string     `gorm:"type:varchar(255)" json:"page_url"`
	ServerUrl string     `gorm:"type:varchar(255)" json:"server_url"`
	ApiKey    string     `gorm:"type:varchar(255)" json:"api_key"`
	ApiSecret string     `gorm:"type:varchar(255)" json:"api_secret"`
	Assets    []Assets   `gorm:"foreignkey:DistributorId" json:"-"`
	Timestamp
}

func init() {
	utils.DB.AutoMigrate(&Distributor{})
}
