package response

import "yuudidi.com/pkg/models"

type OrdersRet struct {
	CommonRet

	Data []models.Order
}
type OrderRequest struct {
	OrderNumber string  `json:"orderNumber"`
	Price       float32 `json:"price"`
	OriginOrder string  `json:"originOrder"`
	//成交量
	Quantity string `json:"quantity"`
	//成交额
	Amount     float64 `json:"amount"`
	PaymentRef string  `json:"paymentRef"`
	//订单状态，0/1分别表示：未支付的/已支付的
	Status models.OrderStatus `json:"status"`
	//成交类型，1：买入;2：卖出。
	Direction         int    `json:"direction"`
	DistributorId     int64  `json:"distributorId"`
	MerchantId        int64  `json:"merchantId"`
	MerchantPaymentId int64  `json:"merchantPaymentId"`
	//扣除用户佣金金额
	TraderCommissionAmount string `json:"traderCommissionAmount"`
	//扣除用户佣金币的量
	TraderCommissionQty string `json:"traderCommissionQty"`
	//用户佣金比率
	TraderCommissionPercent string `json:"traderCommissionPercent"`
	//扣除币商佣金金额
	MerchantCommissionAmount string `json:"merchantCommissionAmount"`
	//扣除币商佣金币的量
	MerchantCommissionQty string `json:"merchantCommissionQty"`
	//币商佣金比率
	MerchantCommissionPercent string `json:"merchantCommissionPercent"`
	//平台扣除的佣金币的量（= trader_commision_qty+merchant_commision_qty)
	PlatformCommissionQty string `json:"platformCommissionQty"`
	//平台商用户id
	AccountId string `json:"accountId"`
	//交易币种
	CurrencyCrypto string `json:"currencyCrypto"example:"BTUSD"`
	//交易法币
	CurrencyFiat string `json:"currencyFiat" example:"RMB"`
	//交易类型 0:微信,1:支付宝,2:银行卡
	PayType uint `json:"payType"`
	//微信或支付宝二维码地址
	QrCode string `json:"qrCode"`
	//微信或支付宝账号
	Name string `json:"name"`
	//银行账号
	BankAccount string `json:"bankAccount"`
	//所属银行
	Bank string `json:"bank"`
	//所属银行分行
	BankBranch string `json:"bankBranch"`
}
