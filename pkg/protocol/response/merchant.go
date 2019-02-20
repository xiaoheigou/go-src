package response

import "yuudidi.com/pkg/models"

type GetProfileData struct {
	// 用户昵称
	NickName       string  `json:"nickname" example:"老王"`
	CurrencyCrypto string  `json:"currency_crypto"`
	Quantity       float64 `json:"quantity"`
	QtyFrozen      float64 `json:"qty_frozen"`
}

type GetProfileRet struct {
	CommonRet
	Data []GetProfileData `json:"data"`
}

type GetAuditStatusData struct {
	// user_status可以为0/1/2/3，分别表示“待审核/正常/未通过审核/冻结”
	UserStatus int `json:"user_status" example:0`
	// 客服联系信息
	ContactPhone string `json:"contact_phone" example:"13812341234"`
	// 额外的信息
	ExtraMessage string `json:"extra_message" example:"您由于xx原因，未通过审核"`
}

type GetAuditStatusRet struct {
	CommonRet
	Data []GetAuditStatusData `json:"data"`
}

type SetNickNameArg struct {
	// 想设置的新昵称，不能超过20个字节
	NickName string `json:"nickname" binding:"required" example:"王老板"`
}

type SetNickNameRet struct {
	CommonRet
}

type SetWorkModeArg struct {
	// 是否接单(1:开启，0:关闭，-1：不做修改)
	InWork int `json:"in_work"`
	// 微信Hook状态(1:开启，0:关闭，-1：不做修改)
	WechatHookStatus int `json:"wechat_hook_status"`
	// 支付宝Hook状态(1:开启，0:关闭，-1：不做修改)
	AlipayHookStatus int `json:"alipay_hook_status"`
	// 币商是否希望收到微信收款方式的自动订单(1:开启，0:关闭，-1：不做修改)
	WechatAutoOrder int `json:"wechat_auto_order"`
	// 币商是否希望收到支付宝收款方式的自动订单(1:开启，0:关闭，-1：不做修改)
	AlipayAutoOrder int `json:"alipay_auto_order"`
}

type SetWorkModeRet struct {
	CommonRet
}

type GetWorkModeData struct {
	// 是否接单(1:开启，0:关闭)
	InWork int `json:"in_work"`
	// 微信Hook状态(1:开启，0:关闭，-1：不做修改)
	WechatHookStatus int `json:"wechat_hook_status"`
	// 支付宝Hook状态(1:开启，0:关闭，-1：不做修改)
	AlipayHookStatus int `json:"alipay_hook_status"`
	// 币商是否希望收到微信收款方式的自动订单(1:开启，0:关闭，-1：不做修改)
	WechatAutoOrder int `json:"wechat_auto_order"`
	// 币商是否希望收到支付宝收款方式的自动订单(1:开启，0:关闭，-1：不做修改)
	AlipayAutoOrder int `json:"alipay_auto_order"`
}

type GetWorkModeRet struct {
	CommonRet
	Data []GetWorkModeData `json:"data"`
}

type SetIdentifyArg struct {
	// 用户id
	Uid    int    `json:"uid" example:123`
	Phone  string `json:"phone" example:13012341234`
	Email  string `json:"email" example:"xxx@xxx.com"`
	IdCard int    `json:"idcard" example:11088888888888888`
}

type SetIdentifyRet struct {
	CommonRet
	// 用户id
	Uid int `json:"uid" example:123`
}

type UploadIdentityArg struct {
	FrontIdentityId string `json:"front-identity-id" example:123`
	BackIdentityId  string `json:"back-identity-id" example:123`
}

type UploadIdentityRet struct {
	CommonRet
}

type MerchantRet struct {
	CommonRet

	Data []models.Merchant `json:"data"`
}

type BillPayload struct {
	// 账单关联的用户支付Id
	UserPayId string `json:"user_pay_id"`
	// 账单Id
	BillId string `json:"bill_id"`
	// 账单的内容
	BillData string `json:"bill_data"`
}

// 上传账单时，接收到数据格式
type UploadBillArg struct {
	// 账单来源，1：微信，2：支付宝
	PayType int `json:"pay_type"`
	// 账单信息，支持同时上传多条
	Data []BillPayload `json:"data"`
}

type RechargeArgs struct {
	//操作人的id
	UserId int64 `json:"id" binding:"required" example:"1"`
	//充值的币种
	Currency string `json:"currency" binding:"required" example:"BTUSD"`
	//充值的数量
	Count string `json:"count" binding:"required" example:"200"`
}

type RechargeRet struct {
	CommonRet
}

type ApproveArgs struct {
	//审核操作 1：通过 0：不通过
	Operation    int    `json:"operation" binding:"required" example:"1"`
	ContactPhone string `json:"contact_phone" binding:"required" example:"13112345678"`
	ExtraMessage string `json:"extra_message" binding:"required" example:"test"`
}

type ApproveRet struct {
	CommonRet
	Data []ApproveDataResponse
}

type ApproveDataResponse struct {
	//用户id
	Uid int `json:"uid" example:"1"`
}

type FreezeArgs struct {
	//冻结操作 1：冻结 0：解冻
	Operation    int    `json:"operation" binding:"required" example:"1"`
	ContactPhone string `json:"contact_phone" binding:"required" example:"13112345678"`
	ExtraMessage string `json:"extra_message" binding:"required" example:"test"`
}

type FreezeRet struct {
	CommonRet
	Data []FreezeDataResponse `json:"data"`
}

type FreezeDataResponse struct {
	//用户ID
	Uid int `json:"uid" example:"1"`
	//用户状态
	Status int `json:"status"`
}
