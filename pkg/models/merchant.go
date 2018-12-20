package models

import (
	"time"
	"yuudidi.com/pkg/utils"
)

const PaymentTypeWeixin = 0
const PaymentTypeAlipay = 1
const PaymentTypeBanck = 2

const PaymentAuditNopass = 0
const PaymentAuditPass = 1

type Merchant struct {
	Id         int64  `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"uid"`
	Nickname   string `gorm:"type:varchar(20)" json:"nickname"`
	AvatarUri  string `gorm:"type:varchar(255)" json:"avatar_uri"`
	DisplayUid string `gorm:"type:varchar(20)" json:"display_uid"`
	Password   []byte `gorm:"type:varbinary(64);column:password;not null" json:"-"`
	Salt       []byte `gorm:"type:varbinary(64);column:salt" json:"-"`
	Algorithm  string `gorm:"type:varchar(255)" json:"-"`
	Phone      string `gorm:"type:varchar(30)" json:"phone"`
	NationCode int    `gorm:"type:int" json:"nation_code"`
	Email      string `gorm:"type:varchar(50)" json:"email"`
	//user_status可以为0/1/2/3，分别表示“待审核/正常/未通过审核/冻结”
	UserStatus int `gorm:"type:tinyint(1);default:0" json:"user_status"`
	//user_cert可以为0/1，分别表示“未认证/已认证”
	UserCert      int           `gorm:"type:tinyint(1);default:0" json:"user_cert"`
	LastLogin     time.Time     `json:"last_login"`
	Asset         []Assets      `gorm:"foreignkey:MerchantId" json:"-"`
	Quantity      string        `gorm:"-" json:"quantity"`
	Payments      []PaymentInfo `gorm:"foreignkey:Uid" json:"-"`
	Preferences   Preferences   `gorm:"foreignkey:PreferencesId" json:"-"`
	PreferencesId uint          `json:"-"`
	Timestamp
}

type Assets struct {
	Id             int64   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	MerchantId     int64   `gorm:"column:merchant_id" json:"merchant_id"`
	DistributorId  int64   `gorm:"column:distributor_id" json:"distributor_id"`
	CurrencyCrypto string  `gorm:"type:varchar(20)" json:"currency_crypto"`
	Quantity       float64 `gorm:"type:decimal(20,5)" json:"quantity"`
	QtyFrozen      float64 `gorm:"type:decimal(20,5)" json:"qty_frozen"`
	Timestamp
}

type Preferences struct {
	Id        int64  `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	InWork int    `gorm:"type:tinyint(2)" json:"in_work"`
	AutoAccept int    `gorm:"type:tinyint(2)" json:"auto_accept"`
	AutoConfirm int    `gorm:"type:tinyint(2)" json:"auto_confirm"`
	Language  string `json:"language"`
	Locale    string `gorm:"type:varchar(12)"`
	Timestamp
}

type PaymentInfo struct {
	Id  int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Uid int64 `gorm:"column:uid;index;not null" json:"uid"`
	//支付类型 0:微信,1:支付宝,2:银行卡
	PayType int `gorm:"column:pay_type;type:tinyint(2)" json:"pay_type"`
	//微信或支付宝账号二维码（识别过后的字符串）
	QrCodeTxt string `gorm:"column:qr_code_txt" json:"qr_code_txt"`
	//微信或支付宝账号二维码
	QrCode string `gorm:"column:qr_code" json:"qr_code"`
	//微信或支付宝账号二维码对应的金额，为0时表示不固定金额
	EAmount float64 `gorm:"column:e_amount;" json:"e_amount"`
	//微信或支付宝账号
	EAccount string `gorm:"column:e_account" json:"e_account"`
	//收款人姓名
	Name string `gorm:"type:varchar(100)" json:"name"`
	//银行账号
	BankAccount string `gorm:"" json:"bank_account"`
	//所属银行名称
	Bank string `json:"bank"`
	//银行分行名称
	BankBranch string `json:"bank_branch"`
	//是否为默认银行卡，0：不是默认，1：默认
	AccountDefault int `gorm:"type:tinyint(2)" json:"account_default"`
	//审核状态，表示是否通过人工审核，0：未通过，1：通过
	AuditStatus int `gorm:"column:audit_status;type:tinyint(2)" json:"audit_status"`
	//是否正在被使用，0：未被使用，1：正在被使用
	InUse int `gorm:"column:in_use;type:tinyint(2)" json:"in_use"`
	Timestamp
}

func init() {
	utils.DB.AutoMigrate(&Merchant{})
	utils.DB.AutoMigrate(&Assets{})
	utils.DB.AutoMigrate(&Preferences{})
	utils.DB.AutoMigrate(&PaymentInfo{})
}
