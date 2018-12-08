package models

import "yuudidi.com/pkg/utils"

type Order struct {
	OrderNumber int64   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Price       float32 `gorm:"type:decimal(10,2)" json:"price"`
	//成交量
	Quantity string `gorm:"type:varchar(32)"json:"quantity"`
	//成交额
	Amount     float64 `gorm:"type:decimal(20,5)" json:"amount"`
	PaymentRef string  `gorm:"type:varchar(8)" json:"payment_ref"`
	//订单状态
	Status OrderStatus `gorm:"type:tinyint(1)" json:"status"`
	//成交方向，以发起方也就是
	Direction         int `gorm:"type:tinyint(1)" json:"direction"`
	DistributorId     int `gorm:"type:int(11)" json:"distributor_id"`
	MerchantId        int `gorm:"type:int(11)" json:"merchant_id"`
	MerchantPaymentId int `gorm:"type:int(11)" json:"merchant_payment_id"`
	//交易币种
	CurrencyCrypto string `gorm:"type:varchar(30)" json:"currency_crypto"example:"BTUSD"`
	//交易法币
	CurrencyFiat string `gorm:"type:char(3)" json:"currency_fiat" example:"RMB"`
	//交易类型 0:微信,1:支付宝,2:银行卡
	PayType uint `gorm:"column:pay_type;type:tinyint(2)" json:"pay_type"`
	//微信或支付宝二维码地址
	QrCode string `gorm:"type:varchar(255)"`
	//微信或支付宝账号
	Name string `gorm:"type:varchar(100)" json:"name"`
	//银行账号
	BankAccount string `gorm:"" json:"bank_account"`
	//所属银行
	Bank string `gorm:"" json:"bank"`
	//所属银行分行
	BankBranch string `gorm:"" json:"bank_branch"`
	Timestamp
}

type OrderHistory struct {
	Order
}

type OrderStatus int

const (
	NEW         OrderStatus = 0
	WAIT_ACCEPT OrderStatus = 1
	ACCEPTED    OrderStatus = 2
	PAID        OrderStatus = 3
	UNPAID      OrderStatus = 4
	ACCOMPLISH  OrderStatus = 5
)

func init() {
	utils.DB.AutoMigrate(&Order{})
	utils.DB.AutoMigrate(&OrderHistory{})
}