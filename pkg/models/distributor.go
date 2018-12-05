package models

import "time"

type Distributor struct {
	Id         int           `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Uid        int64         `gorm:"AUTO_INCREMENT;unique" json:"uid"`
	NickName   string        `gorm:"unique;not null" json:"nick_name"`
	ApiKey     string        `gorm:"type:varchar(255)" json:"avatar_uri"`
	ApiSecret  Assets        `gorm:"foreignkey:AssetId"`
	Payments   []PaymentInfo `gorm:"foreignkey:UserPaymentId"`
	Identify   Preferences   `gorm:"foreignkey:IdentifyId"`
	IdentifyId uint          `json:"-"`
	CreatedAt  time.Time     `json:"-"`
	UpdatedAt  time.Time     `json:"-"`
	DeletedAt  *time.Time    `json:"-"`
}
