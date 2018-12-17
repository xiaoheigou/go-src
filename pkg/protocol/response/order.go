package response

import "yuudidi.com/pkg/models"

type OrdersRet struct {
	CommonRet

	Data []models.Order
}
type OrderRequest struct {
	OrderNumber string   `json:"order_number"`
	Price       float32 `json:"price"`
	//成交量
	Quantity string `json:"quantity"`
	//成交额
	Amount     float64 `json:"amount"`
	PaymentRef string  `json:"payment_ref"`
	//订单状态，0/1分别表示：未支付的/已支付的
	Status models.OrderStatus `json:"status"`
	//成交方向，以发起方（平台商用户）为准。0表示平台商用户买入，1表示平台商用户卖出。
	Direction         int   `json:"direction"`
	DistributorId     int64 `json:"distributor_id"`
	MerchantId        int64 `json:"merchant_id"`
	MerchantPaymentId int64 `json:"merchant_payment_id"`
	//扣除用户佣金金额
	TraderCommissionAmount string `json:"trader_commission_amount"`
	//扣除用户佣金币的量
	TraderCommissionQty string `json:"trader_commission_qty"`
	//用户佣金比率
	TraderCommissionPercent string `json:"trader_commission_percent"`
	//扣除币商佣金金额
	MerchantCommissionAmount string `json:"merchant_commission_amount"`
	//扣除币商佣金币的量
	MerchantCommissionQty string `json:"merchant_commission_qty"`
	//币商佣金比率
	MerchantCommissionPercent string `json:"merchant_commission_percent"`
	//平台扣除的佣金币的量（= trader_commision_qty+merchant_commision_qty)
	PlatformCommissionQty string `json:"platform_commission_qty"`
	//平台商用户id
	AccountId string `json:"account_id"`
	//交易币种
	CurrencyCrypto string `json:"currency_crypto"example:"BTUSD"`
	//交易法币
	CurrencyFiat string `json:"currency_fiat" example:"RMB"`
	//交易类型 0:微信,1:支付宝,2:银行卡
	PayType uint `json:"pay_type"`
	//微信或支付宝二维码地址
	QrCode string `json:"qr_code"`
	//微信或支付宝账号
	Name string `json:"name"`
	//银行账号
	BankAccount string `json:"bank_account"`
	//所属银行
	Bank string `json:"bank"`
	//所属银行分行
	BankBranch string `json:"bank_branch"`
}
