package models

import "yuudidi.com/pkg/utils"

type Order struct {
	Id          int64   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	OrderNumber string  `gorm:"type:varchar(191);unique_index;not null" json:"order_number"`
	OriginOrder string  `gorm:"type:varchar(191);unique_index:origin_distributor_order;not null" json:"origin_order"`
	Price       float32 `gorm:"type:decimal(10,4)" json:"price"`
	//成交量
	Quantity float64 `gorm:"type:decimal(20,5)"json:"quantity"`
	//成交额
	Amount     float64 `gorm:"type:decimal(20,5)" json:"amount"`
	PaymentRef string  `gorm:"type:varchar(8)" json:"payment_ref"`
	//订单状态
	Status OrderStatus `gorm:"type:tinyint(1)" json:"status"`
	//成交方向，以发起方（平台商用户）为准。0表示平台商用户买入，1表示平台商用户卖出。
	Direction         int   `gorm:"type:tinyint(1)" json:"direction"`
	DistributorId     int64 `gorm:"type:int(11);unique_index:origin_distributor_order;not null" json:"distributor_id"`
	MerchantId        int64 `gorm:"type:int(11)" json:"merchant_id"`
	MerchantPaymentId int64 `gorm:"type:int(11)" json:"merchant_payment_id"`
	//扣除用户佣金金额
	TraderCommissionAmount float64 `gorm:"type:decimal(20,5)" json:"trader_commission_amount"`
	//扣除用户佣金币的量
	TraderCommissionQty float64 `gorm:"type:decimal(20,5)" json:"trader_commission_qty"`
	//用户佣金比率
	TraderCommissionPercent float64 `gorm:"type:decimal(20,5)" json:"trader_commission_percent"`
	//扣除币商佣金金额
	MerchantCommissionAmount float64 `gorm:"type:decimal(20,5)" json:"merchant_commission_amount"`
	//扣除币商佣金币的量
	MerchantCommissionQty float64 `gorm:"type:decimal(20,5)" json:"merchant_commission_qty"`
	//币商佣金比率
	MerchantCommissionPercent float64 `gorm:"type:decimal(20,5)" json:"merchant_commission_percent"`
	//平台扣除的佣金币的量（= trader_commision_qty+merchant_commision_qty)
	PlatformCommissionQty float64 `gorm:"type:decimal(20,5)" json:"platform_commission_qty"`
	//平台商用户id
	AccountId string `gorm:"type:varchar(191)" json:"account_id"`
	//交易币种
	CurrencyCrypto string `gorm:"type:varchar(30)" json:"currency_crypto"example:"BTUSD"`
	//交易法币
	CurrencyFiat string `gorm:"type:char(3)" json:"currency_fiat" example:"RMB"`
	//交易类型
	PayType uint `gorm:"column:pay_type;type:tinyint(2)" json:"pay_type"`
	//微信或支付宝二维码地址
	QrCode string `gorm:"type:varchar(255)" json:"qr_code"`
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
	NEW        OrderStatus = 0
	WAITACCEPT OrderStatus = 1
	ACCEPTED   OrderStatus = 2
	NOTIFYPAID OrderStatus = 3
	// 确认付款
	CONFIRMPAID OrderStatus = 4
	//异常订单
	SUSPENDED OrderStatus = 5
	// 应收实付不符
	PAYMENTMISMATCH OrderStatus = 6
	// 订单完成 转账结束
	TRANSFERRED OrderStatus = 7
)

func init() {
	utils.DB.AutoMigrate(&Order{})
	utils.DB.AutoMigrate(&OrderHistory{})
}
