package response

import (
	"github.com/shopspring/decimal"
	"yuudidi.com/pkg/models"
)

type OrdersRet struct {
	CommonRet

	Data []models.Order `json:"data"`
}
type OrderRequest struct {
	OrderNumber string  `json:"orderNumber"`
	Price       float32 `json:"price"`
	OriginOrder string  `json:"originOrder"`
	//成交量
	Quantity decimal.Decimal `json:"quantity"`
	//成交额
	Amount     float64 `json:"amount"`
	PaymentRef string  `json:"paymentRef"`
	//订单状态，0/1分别表示：未支付的/已支付的
	Status models.OrderStatus `json:"status"`
	//成交类型，1：买入;2：卖出。
	Direction         int   `json:"direction"`
	DistributorId     int64 `json:"distributorId"`
	MerchantId        int64 `json:"merchantId"`
	MerchantPaymentId int64 `json:"merchantPaymentId"`
	// 平台手续费收入，它可能为负数。目前仅用户提现订单涉及手续费。
	TraderBTUSDFeeIncome decimal.Decimal `json:"trader_btusd_fee_income"`
	// 币商手续费收入。目前仅用户提现订单涉及手续费。
	MerchantBTUSDFeeIncome decimal.Decimal `json:"merchant_btusd_fee_income"`
	// 金融滴滴平台手续费收入。目前仅用户提现订单涉及手续费。
	JrdidiBTUSDFeeIncome decimal.Decimal `json:"jrdidi_btusd_fee_income"`
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

	OriginAmount float64 `json:"originAmount,string"`
	Fee          float64 `json:"fee"`
	Price2       float32 `json:"price2,string"`
	AppCoinName  string  `json:"appCoinName"`
	Remark       string  `json:"remark"`
	//异步通知平台商url
	AppServerNotifyUrl string `json:"appServerNotifyUrl"`
	AppReturnPageUrl   string `json:"appReturnPageUrl"`
}
