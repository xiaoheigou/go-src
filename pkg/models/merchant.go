package models

import (
	"github.com/shopspring/decimal"
	"time"

	"yuudidi.com/pkg/utils"
)

const PaymentTypeWeixin = 1
const PaymentTypeAlipay = 2
const PaymentTypeBank = 4

const PaymentAuditNopass = 0
const PaymentAuditPass = 1

type Merchant struct {
	Id         int64  `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"uid"`
	Nickname   string `gorm:"type:varchar(20)" json:"nickname"`
	AvatarUri  string `gorm:"type:varchar(255)" json:"avatar_uri"`
	DisplayUid string `gorm:"type:varchar(20)" json:"display_uid"`
	Password   []byte `gorm:"type:varbinary(64);column:password;not null" json:"-"`
	Salt       []byte `gorm:"type:varbinary(64);column:salt;not null" json:"-"`
	Algorithm  string `gorm:"type:varchar(255);not null" json:"-"`
	Phone      string `gorm:"type:varchar(30)" json:"phone"`
	NationCode int    `gorm:"type:int" json:"nation_code"`
	Email      string `gorm:"type:varchar(50)" json:"email"`
	//user_status可以为0/1/2/3，分别表示“待审核/正常/未通过审核/冻结”
	UserStatus int `gorm:"type:tinyint(1);default:0" json:"user_status"`
	//user_cert可以为0/1，分别表示“未认证/已认证”
	UserCert int `gorm:"type:tinyint(1);default:0" json:"user_cert"`
	// role为1时表示“官方币商”
	Role          int           `gorm:"type:tinyint(1);default:0" json:"role"`
	LastLogin     time.Time     `json:"last_login"`
	Asset         []Assets      `gorm:"foreignkey:MerchantId" json:"-"`
	Quantity      string        `gorm:"-" json:"quantity"`
	QtyFrozen     string        `gorm:"-" json:"qty_frozen"`
	Payments      []PaymentInfo `gorm:"foreignkey:Uid" json:"-"`
	Preferences   Preferences   `gorm:"foreignkey:PreferencesId" json:"-"`
	PreferencesId uint64        `json:"-"`
	Timestamp
}

type Assets struct {
	Id             int64           `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	MerchantId     int64           `gorm:"column:merchant_id;not null" json:"merchant_id"`
	DistributorId  int64           `gorm:"column:distributor_id;not null" json:"distributor_id"`
	CurrencyCrypto string          `gorm:"type:varchar(20)" json:"currency_crypto"`
	Quantity       decimal.Decimal `gorm:"type:decimal(30,10);not null" json:"quantity"`
	QtyFrozen      decimal.Decimal `gorm:"type:decimal(30,10);not null" json:"qty_frozen"`
	Timestamp
}

type Preferences struct {
	Id int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	// 币商是否接单的总开关
	InWork int `gorm:"type:tinyint(2);not null" json:"in_work"`
	// 微信Hook状态(1:开启，0:关闭，-1：不做修改)
	WechatHookStatus int `gorm:"type:tinyint(2);not null" json:"wechat_hook_status"`
	// 支付宝Hook状态(1:开启，0:关闭，-1：不做修改)
	AlipayHookStatus int `gorm:"type:tinyint(2);not null" json:"alipay_hook_status"`
	// 币商是否希望收到微信收款方式的自动订单(1:开启，0:关闭，-1：不做修改)
	WechatAutoOrder int `gorm:"type:tinyint(2);not null" json:"wechat_auto_order"`
	// 币商是否希望收到支付宝收款方式的自动订单(1:开启，0:关闭，-1：不做修改)
	AlipayAutoOrder int `gorm:"type:tinyint(2);not null" json:"alipay_auto_order"`

	// 当前使用的微信自动收款账号，在PaymentInfo表中的Id值
	CurrAutoWeixinPaymentId int64 `gorm:"type:int(11);not null" json:"curr_auto_weixin_payment_id"`
	// 当前使用的支付宝自动收款账号，在PaymentInfo表中的Id值
	CurrAutoAlipayPaymentId int64  `gorm:"type:int(11);not null" json:"curr_auto_alipay_payment_id"`
	Language                string `json:"language"`
	Locale                  string `gorm:"type:varchar(12)"`
	Timestamp
}

type PaymentInfo struct {
	Id  int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Uid int64 `gorm:"column:uid;index;not null" json:"uid"`
	//支付类型 1:微信,2:支付宝,4:银行卡
	PayType int `gorm:"column:pay_type;type:tinyint(2)" json:"pay_type"`
	//微信或支付宝账号二维码（识别过后的字符串）
	QrCodeTxt string `gorm:"column:qr_code_txt" json:"qr_code_txt"`
	//微信或支付宝账号二维码（用户直接上传的图片）
	QrCodeOrigin string `gorm:"column:qr_code_origin" json:"qr_code_origin"`
	//微信或支付宝账号二维码（生成的二维码，它的size更小）
	QrCode string `gorm:"column:qr_code" json:"qr_code"`
	//微信或支付宝账号二维码对应的金额，为0时表示不固定金额
	EAmount float64 `gorm:"column:e_amount;type:decimal(20,5)" json:"e_amount"`
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
	//审核状态，表示是否通过人工审核，0：审核中，1：通过
	AuditStatus int `gorm:"column:audit_status;type:tinyint(2)" json:"audit_status"`
	//是否正在被使用，0：未被使用，1：正在被使用
	InUse int `gorm:"column:in_use;type:tinyint(2)" json:"in_use"`
	// 微信或支付宝账号Hook类型，0表示普通的收款账号，1表示Hook类型的收款账号（可以用来自动确认收款）
	PaymentAutoType int `gorm:"column:payment_auto_type;type:tinyint(2);not null" json:"payment_auto_type"`
	// 支付宝或微信用户支付id，前端App通过hook方式可以拿到
	UserPayId string `gorm:"column:user_pay_id;not null" json:"user_pay_id"`
	// 上次使用时间
	LastUseTime time.Time `json:"last_use_time"`
	Timestamp
}

//BankInfo - bank information
type BankInfo struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func init() {
	utils.DB.Set("gorm:table_options", "AUTO_INCREMENT=10001").AutoMigrate(&Merchant{}) // id从10001开始
	utils.DB.AutoMigrate(&Assets{})
	utils.DB.AutoMigrate(&Preferences{})
	utils.DB.AutoMigrate(&PaymentInfo{})
}
