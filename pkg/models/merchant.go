package models

import (
	"yuudidi.com/pkg/utils"
)

type Merchant struct {
	Id            int           `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"uid"`
	NickName      string        `gorm:"unique;not null" json:"nickname"`
	AvatarUri     string        `gorm:"type:varchar(255)" json:"avatar_uri"`
	DisplayUid    string        `gorm:"type:varchar(20)" json:"display_uid"`
	Password      string        `gorm:"type:varchar(50)"`
	Salt          string        `json:"-"`
	Algorithm     string        `json:"-"`
	Phone         string        `gorm:"type:varchar(30)"`
	Email         string        `gorm:"type:varchar(50)"`
	Asset         []Assets      `gorm:"foreignkey:MerchantId"`
	Payments      []PaymentInfo `gorm:"foreignkey:Uid"`
	Preferences   Preferences   `gorm:"foreignkey:PreferencesId"`
	PreferencesId uint          `json:"-"`
	Timestamp
}

type Assets struct {
	Id            int     `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	MerchantId    int     `gorm:"column:merchant_id" json:"merchant_id"`
	DistributorId int     `gorm:"column:distributor_id" json:"distributor_id"`
	CurrencyAsset string  `gorm:"type:varchar(20)" json:"currency_asset"`
	Quantity      float64 `gorm:"type:decimal(20,5)" json:"quantity"`
	QtyFrozen     float64 `gorm:"type:decimal(20,5)" json:"qty_frozen"`
	Timestamp
}

type Preferences struct {
	Id        int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	TakeOrder int    `gorm:"type:tinyint(2)" json:"accept"`
	AutoOrder int 	 `gorm:"type:tinyint(2)" json:"auto"`
	Language  string `json:"language"`
	Locale    string `gorm:"type:varchar(12)"`
	Timestamp
}

type PaymentInfo struct {
	Id          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Uid         uint   `gorm:"column:uid;index;not null" json:"uid"`
	PayType     uint   `gorm:"column:pay_type;type:tinyint(2)" json:"pay_type"`
	QrCode      []byte `gorm:"type:blob"`
	Name        string `gorm:"type:varchar(100)" json:"name"`
	BankAccount string `gorm:"" json:"bank_account"`
	Bank        string `json:"bank"`
	BankBranch  string `json:"bank_branch"`
	Timestamp
}

func init() {
	utils.DB.AutoMigrate(&Merchant{})
	utils.DB.AutoMigrate(&Assets{})
	utils.DB.AutoMigrate(&Preferences{})
	utils.DB.AutoMigrate(&PaymentInfo{})
}
