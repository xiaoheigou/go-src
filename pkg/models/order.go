package models

import (
	"github.com/shopspring/decimal"
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
	Quantity decimal.Decimal `gorm:"type:decimal(30,10)"json:"quantity"`
	//成交额
	Amount float64 `gorm:"type:decimal(20,2)" json:"amount"`
	//原始成交额
	OriginAmount float64 `gorm:"type:decimal(20,2)" json:"origin_amount"`
	//手续费
	Fee        float64 `gorm:"type:decimal(20,2)" json:"fee"`
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
	// 用来标记订单BTUSD的转移过程，候选值见下文
	BTUSDFlowStatus int32 `gorm:"type:tinyint(1)"`
	// 平台手续费收入，它可能为负数。目前仅用户提现订单涉及手续费。
	TraderBTUSDFeeIncome decimal.Decimal `gorm:"type:decimal(30,10)" json:"trader_btusd_fee_income"`
	// 币商手续费收入。目前仅用户提现订单涉及手续费。
	MerchantBTUSDFeeIncome decimal.Decimal `gorm:"type:decimal(30,10)" json:"merchant_btusd_fee_income"`
	// 金融滴滴平台手续费收入。目前仅用户提现订单涉及手续费。
	JrdidiBTUSDFeeIncome decimal.Decimal `gorm:"type:decimal(30,10)" json:"jrdidi_btusd_fee_income"`
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
	//微信或支付宝账号二维码所编码的字符串
	QrCodeTxt string `gorm:"type:varchar(255)" json:"qr_code_txt"`
	//收款人姓名
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
	Timeout        int64     `gorm:"-" json:"timeout"`
	//异步通知平台商url
	AppServerNotifyUrl string `gorm:"type:varchar(255)" json:"app_server_notify_url"`
	AppReturnPageUrl   string `gorm:"type:varchar(255)" json:"app_return_page_url"`
	// 订单的接单类型，0表示手动接单订单，1表示自动接单订单。
	AcceptType int `gorm:"type:tinyint(1);default:0" json:"accept_type"`
	// 支付宝或微信的用户支付Id。仅用于自动订单。
	UserPayId string `gorm:"column:user_pay_id" json:"user_pay_id"`
	// 这个订单（用户充值订单）的实际收款金额
	ActualAmount float64 `gorm:"type:decimal(20,2);default:0" json:"actual_amount"`
	// 当H5用户付款时，把它的名字填入，以便币商确定他收到的是谁的付款
	AppUserName string `gorm:"type:varchar(191)" json:"app_user_name"`
	// 保存H5用户付款的凭证的Url
	AppUserReceiptUrl string `gorm:"type:varchar(255)" json:"app_user_receipt_url"`
	Timestamp
}

// BTUSDFlowStatus相关值
const (
	// 下面常量都关联用户提现订单
	BTUSDFlowD1TraderQtyToTraderFrozen      = 1
	BTUSDFlowD1TraderFrozenToMerchantFrozen = 2
	BTUSDFlowD1MerchantFrozenToMerchantQty  = 3

	BTUSDFlowD1MerchantFrozenToTraderFrozen = 4
	BTUSDFlowD1MerchantFrozenToTraderQty    = 5

	BTUSDFlowD1TraderFrozenToMerchantQty = 6
	BTUSDFlowD1TraderFrozenToTraderQty   = 7

	// 下面常量关联用户充值订单（目前还没有）
)

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
	// 超时没人接单的订单状态
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
