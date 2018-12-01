package models

import (
	"otc-project/pkg/utils"
	"time"
)

type User struct {
	Id         int           `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Uid        int64         `gorm:"AUTO_INCREMENT;unique" json:"uid"`
	NickName   string        `gorm:"unique;not null" json:"nick_name"`
	AvatarUri  string        `gorm:"type:varchar(255)" json:"avatar_uri"`
	Role       string        `gorm:"not null;type:ENUM('Admin', 'Merchant', 'User')" json:"-"`
	Asset      Assets        `gorm:"foreignkey:AssetId"`
	Payments   []PaymentInfo `gorm:"foreignkey:UserPaymentId"`
	Identify   Identities    `gorm:"foreignkey:IdentifyId"`
	AssetId    uint          `json:"-"`
	IdentifyId uint          `json:"-"`
	CreatedAt  time.Time     `json:"-"`
	UpdatedAt  time.Time     `json:"-"`
	DeletedAt  *time.Time    `json:"-"`
}

type Assets struct {
	Id        int        `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Quantity  float64    `json:"quantity"`
	QtyFrozen float32    `json:"frozen_quantity"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type Identities struct {
	Id        int        `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name      string     `json:"name"`
	IdCard    string     `gorm:"type:varchar(12)"`
	CreatedAt time.Time  `gorm:""json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type PaymentInfo struct {
	Id            int        `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserPaymentId uint       `gorm:"column:uid;index;not null"`
	Account       string     `json:"account"`
	CreatedAt     time.Time  `json:"-"`
	UpdatedAt     time.Time  `json:"-"`
	DeletedAt     *time.Time `json:"-"`
}

func init() {
	utils.DB.AutoMigrate(&User{})
	utils.DB.AutoMigrate(&Assets{})
	utils.DB.AutoMigrate(&Identities{})
	utils.DB.AutoMigrate(&PaymentInfo{})
}
