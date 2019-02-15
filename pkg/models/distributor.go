package models

import (
	"github.com/shopspring/decimal"
	"yuudidi.com/pkg/utils"
)

type Distributor struct {
	Id    int64  `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name  string `gorm:"type:varchar(191);not null" json:"name"`
	Phone string `gorm:"type:varchar(20)" json:"phone"`
	//平台商状态 0: 申请 1: 正常 2: 冻结
	Status    int      `gorm:"type:tinyint(1)" json:"status"`
	Domain    string   `gorm:"type:varchar(255)" json:"domain"`
	PageUrl   string   `gorm:"type:varchar(255)" json:"page_url"`
	ServerUrl string   `gorm:"type:varchar(255)" json:"server_url"`
	CaPem     []byte   `gorm:"type:mediumblob" json:"-"`
	ApiKey    string   `gorm:"unique;type:varchar(191)" json:"api_key"`
	ApiSecret string   `gorm:"type:varchar(255)" json:"api_secret"`
	Assets    []Assets `gorm:"foreignkey:DistributorId" json:"-"`
	Quantity  string   `gorm:"-" json:"quantity"`
	QtyFrozen string   `gorm:"-" json:"qty_frozen"`
	//自自定义币种中文文名字
	AppCoinName string `gorm:"type:varchar(8)" json:"app_coin_name"`
	//默认值是:CNY
	AppCoinSymbol string `gorm:"type:varchar(16);default:'CNY'" json:"app_coin_symbol"`
	//兑换比例
	AppCoinRate float32 `gorm:"type:decimal(10,4);default:1" json:"app_coin_rate"`
	//抽取比例
	AppUserWithdrawalFeeRate decimal.Decimal `gorm:"type:decimal(10,6);default:0.023077" json:"app_user_withdrawal_fee_rate"`
	//手续费平台所占部分
	AppUserWithdrawalFeeRateTraderPart decimal.Decimal `gorm:"type:decimal(10,6);default:0.000000" json:"app_user_withdrawal_fee_rate_trader_part"`
	//手续费jrdidi所占部分
	AppUserWithdrawalFeeRateJrdidiPart decimal.Decimal `gorm:"type:decimal(10,6);default:0.013210" json:"app_user_withdrawal_fee_rate_jrdidi_part"`
	//手续费币商所占部分
	AppUserWithdrawalFeeRateMerchantPart decimal.Decimal `gorm:"type:decimal(10,6);default:0.009867" json:"app_user_withdrawal_fee_rate_merchant_part"`
	Timestamp
}

//用户信息表
type AccountInfo struct {
	Id int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	//用户id
	AccountId     string `gorm:"type:varchar(191)" json:"account_id"`
	DistributorId int64  `gorm:"type:int(11);not null" json:"distributor_id"`
	OrderNumber   string `gorm:"type:varchar(191)" json:"order_number"`
	//成交方向，以发起方（平台商用户）为准。0表示平台商用户买入，1表示平台商用户卖出。
	Direction int     `gorm:"type:tinyint(1)" json:"direction"`
	Price     float32 `gorm:"type:decimal(10,4)" json:"price"`
	//成交量
	Quantity decimal.Decimal `gorm:"type:decimal(30,10)"json:"quantity"`
	//成交额
	Amount float64 `gorm:"type:decimal(20,5)" json:"amount"`
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
}

func init() {
	utils.DB.AutoMigrate(&Distributor{})
	utils.DB.AutoMigrate(&AccountInfo{})
}
