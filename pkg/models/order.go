package models

import (
	"time"

	"yuudidi.com/pkg/utils"
)

type Order struct {
	Id          int64   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	OrderNumber string  `gorm:"type:varchar(191);unique_index;not null" json:"order_number"`
	OriginOrder string  `gorm:"type:varchar(191);unique_index:origin_distributor_order;not null" json:"origin_order"`
	Price       float32 `gorm:"type:decimal(10,4)" json:"price"`
	//提现价格
	Price2 float32 `gorm:"type:decimal(10,4)" json:"price2"`
	//成交量
	Quantity float64 `gorm:"type:decimal(30,10)"json:"quantity"`
	//成交额
	Amount float64 `gorm:"type:decimal(20,5)" json:"amount"`
	//原始成交额
	OriginAmount float64 `gorm:"type:decimal(20,5)" json:"origin_amount"`
	//手续费
	Fee        float64 `gorm:"type:decimal(20,5)" json:"fee"`
	PaymentRef string  `gorm:"type:varchar(8)" json:"payment_ref"`
	//订单状态
	Status OrderStatus `gorm:"type:tinyint(1)" json:"status"`
	//订单异常状态原因
	StatusReason StatusReason `gorm:"type:tinyint(1);" json:"status_reason"`
	//确认收付款状态，0：没收到确认付款同步信息（没收到客户端“SUCCESS”），1：收到确认付款同步信息（收到客户端“SUCCESS”）
	Synced uint `gorm:"type:tinyint(1)" json:"synced"`
	//成交方向，以发起方（平台商用户）为准。0表示平台商用户买入，1表示平台商用户卖出。
	Direction         int    `gorm:"type:tinyint(1)" json:"direction"`
	DistributorId     int64  `gorm:"type:int(11);unique_index:origin_distributor_order;not null" json:"distributor_id"`
	DistributorName   string `gorm:"-" json:"distributor_name"`
	MerchantId        int64  `gorm:"type:int(11)" json:"merchant_id"`
	MerchantName      string `gorm:"-" json:"merchant_name"`
	MerchantPhone     string `gorm:"-" json:"merchant_phone"`
	MerchantPaymentId int64  `gorm:"type:int(11)" json:"merchant_payment_id"`
	//扣除用户佣金金额
	TraderCommissionAmount float64 `gorm:"type:decimal(20,5)" json:"trader_commission_amount"`
	//扣除用户佣金币的量
	TraderCommissionQty float64 `gorm:"type:decimal(30,10)" json:"trader_commission_qty"`
	//用户佣金比率
	TraderCommissionPercent float64 `gorm:"type:decimal(20,5)" json:"trader_commission_percent"`
	//扣除币商佣金金额
	MerchantCommissionAmount float64 `gorm:"type:decimal(20,5)" json:"merchant_commission_amount"`
	//扣除币商佣金币的量
	MerchantCommissionQty float64 `gorm:"type:decimal(30,10)" json:"merchant_commission_qty"`
	//币商佣金比率
	MerchantCommissionPercent float64 `gorm:"type:decimal(20,5)" json:"merchant_commission_percent"`
	//平台扣除的佣金币的量（= trader_commision_qty+merchant_commision_qty)
	PlatformCommissionQty float64 `gorm:"type:decimal(30,10)" json:"platform_commission_qty"`
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
	// 派单接收时间（order表中没有，返回前端时从fulfillment_events表中获得）
	AcceptedAt time.Time `gorm:"-" json:"accepted_at"`
	// 通知支付时间（order表中没有，返回前端时从fulfillment_events表中获得）
	PaidAt time.Time `gorm:"-" json:"paid_at"`
	// 确认支付时间（order表中没有，返回前端时从fulfillment_events表中获得）
	PaymentConfirmedAt time.Time `gorm:"-" json:"payment_confirmed_at"`
	// 转账时间（order表中没有，返回前端时从fulfillment_events表中获得）
	TransferredAt time.Time `gorm:"-" json:"transferred_at"`
	// 系统当前时间（order表中没有，返回前端时实时计算出来）
	SvrCurrentTime time.Time `gorm:"-" json:"svr_current_time"`
	AppCoinName    string    `gorm:"type:varchar(16)" json:"app_coin_name"`
	Remark         string    `gorm:"type:varchar(255)" json:"remark"`

	//异步通知平台商url
	AppServerNotifyUrl string `gorm:"type:varchar(255)" json:"app_server_notify_url"`
	AppReturnPageUrl   string `gorm:"type:varchar(255)" json:"app_return_page_url"`
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
	// 超时没人接单的订单状态，不要重启这样的订单
	ACCEPTTIMEOUT OrderStatus = 8
	//// 客服放币
	//RELEASE = 11
	//// 客服解冻
	//UNFREEZE OrderStatus = 10
)

type StatusReason int

const (
	//系统更新失败
	SYSTEMUPDATEFAIL StatusReason = 1
	// 付款超时
	PAIDTIMEOUT StatusReason = 2
	// 确认收款超时
	CONFIRMTIMEOUT StatusReason = 3
	// 申诉
	COMPLIANT StatusReason = 4
	// 退款进行中
	REFUNDING StatusReason = 5
	// 退款失败
	REFUNDFAIL StatusReason = 6
	// 退款成功
	REFUNDSUCCESS StatusReason = 7
	// 未真实付款
	NONPAYMENT StatusReason = 8
	// 订单有异议
	ORDERDISPUTED StatusReason = 9
	// 客服标记完成
	MARKCOMPLETED StatusReason = 19
	// 订单取消
	CANCEL StatusReason = 20
)

func init() {
	utils.DB.AutoMigrate(&Order{})
	utils.DB.AutoMigrate(&OrderHistory{})
}
