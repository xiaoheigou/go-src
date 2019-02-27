package response

import "yuudidi.com/pkg/models"

type GetDistributorsRet struct {
	CommonRet

	Data []models.Distributor `json:"data"`
}

type CreateDistributorsRet struct {
	CommonRet
	Data []interface{}
}

type CreateDistributorsArgs struct {
	Name      string `json:"name" binding:"required" example:"test"`
	Phone     string `json:"phone" example:"13112345678"`
	Domain    string `json:"domain" binding:"required" example:"baidu.com"`
	PageUrl   string `json:"page_url" example:"1"`
	ServerUrl string `json:"server_url" example:"1"`
	ApiKey    string `json:"api_key" binding:"required" example:"13112345678"`
	ApiSecret string `json:"api_secret" binding:"required" example:"13112345678"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	//AppUserWithdrawalFeeRate float64 `json:"app_user_withdrawal_fee_rate"`
}

type UpdateDistributorsRet struct {
	CommonRet
	Data []interface{}
}

type UpdateDistributorsArgs struct {
	Name   string `json:"name" binding:"required" example:"test"`
	Phone  string `json:"phone" example:"13112345678"`
	Domain string `json:"domain" binding:"required" example:"baidu.com"`
}

type DistributorWithdrawArgs struct {
	AppOrderId         string `json:"app_order_id" binding:"required"`      // 由商家系统内部生成的订单ID
	OrderAmount        string `json:"order_amount" binding:"required"`      // 本次订单中下单的金额
	OrderPayTypeId     string `json:"order_pay_type_id" binding:"required"` // 收款方式，详见《JRDiDi平台支付方式对应ID列表》
	PayAccountId       string `json:"pay_account_id" binding:"required"`    // 银行卡卡号
	PayAccountUser     string `json:"pay_account_user" binding:"required"`  // 收款人的真实姓名
	PayAccountInfo     string `json:"pay_account_info" binding:"required"`  // 分行或支行名称
	AppServerNotifyUrl string `json:"app_server_notify_url"`                // 异步通知的接口
	AppReturnPageUrl   string `json:"app_return_page_url"`
}
