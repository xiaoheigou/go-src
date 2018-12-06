package merchant

import (
	"YuuPay_core-service/pkg/models"
	"github.com/gin-gonic/gin"
)

type CommonRet struct {
	// status可以为success或者fail
	Status   string `json:"status" binding:"required" example:"success"`
	// err_msg仅在失败时设置
	ErrMsg  string `json:"err_msg" example:"由于xx原因，导致操作失败"`
	// err_code仅在失败时设置
	ErrCode int `json:"err_code" example:1001`
}


type GetAuditStatusRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      int `json:"uid" example:123`
		// user_status可以为0/1/2，分别表示“正常/待审核/冻结”
		UserStatus int `json:"user_status" example:0`
		// 客服联系信息
		ContactPhone string `json:"contact_phone" example:"13812341234"`
		// 额外的信息
		ExtraMessage string `json:"extra_message" example:"您由于xx原因，未通过审核"`
	}
}

// @Summary 获取承兑商账号审核状态
// @Tags 承兑商APP API
// @Description 获取承兑商账号审核状态
// @Accept  json
// @Produce  json
// @Param uid  query  string     true        "用户id"
// @Success 200 {object} merchant.GetAuditStatusRet ""
// @Router /merchant/auditstatus [get]
func GetAuditStatus(c *gin.Context) {
	// TODO

	var ret GetAuditStatusRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	ret.Entity.UserStatus = 1
	ret.Entity.ContactPhone = "13012349876"
	c.JSON(200, ret)
}


type GetProfileRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      int `json:"uid" example:123`
		// 用户昵称
		NickName string `json:"nickname" example:"老王"`
		// 平台币的符号
		AssetSymbol string `json:"asset_symbol" example:"BTUSD"`
		// 当前承兑商所有的平台币余额（包含被冻结的平台币）
		AssetTotal  string `json:"asset_total" example:"2000"`
		// 当前承兑商被冻结的平台币数量
		AssetFrozen  string `json:"asset_frozen" example:"100"`
	}
}

// @Summary 获取承兑商个人信息
// @Tags 承兑商APP API
// @Description 获取承兑商个人信息
// @Accept  json
// @Produce  json
// @Param uid  query  string     true        "用户id"
// @Success 200 {object} merchant.GetProfileRet ""
// @Router /merchant/profile [get]
func GetProfile(c *gin.Context) {
	// TODO

	var ret GetProfileRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	ret.Entity.NickName = "老王"
	ret.Entity.AssetSymbol = "BTUSD"
	ret.Entity.AssetTotal = "2000"
	ret.Entity.AssetFrozen = "100"
	c.JSON(200, ret)
}


type SetNickNameArg struct {
	// 用户id
	Uid      int `json:"uid" example:123`
	// 想设置的新昵称
	NickName string `json:"nickname" example:"王老板"`
}

type SetNickNameRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      int `json:"uid" example:123`
	}
}

// @Summary 设置承兑商昵称
// @Tags 承兑商APP API
// @Description 设置承兑商昵称
// @Accept  json
// @Produce  json
// @Param body  body  merchant.SetNickNameArg     true        "用户id"
// @Success 200 {object} merchant.SetNickNameRet ""
// @Router /merchant/settings/nickname [put]
func SetNickName(c *gin.Context) {
	// TODO

	var ret SetNickNameRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	c.JSON(200, ret)
}


type SetWorkModeArg struct {
	// 用户id
	Uid      int `json:"uid" example:123`
	// 是否接单(1:开启，0:关闭)
	Accept  int `json:"accept" example:1`
	// 是否自动接单(1:开启，0:关闭)
	Auto  int `json:"auto" example:1`
}

type SetWorkModeRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      int `json:"uid" example:123`
	}
}


// @Summary 承兑商设置订单推送模式和开关
// @Tags 承兑商APP API
// @Description 承兑商设置订单推送模式和开关
// @Accept  json
// @Produce  json
// @Param body  body  merchant.SetWorkModeArg     true        "用户id"
// @Success 200 {object} merchant.SetWorkModeRet ""
// @Router /merchant/settings/workmode [put]
func SetWorkMode(c *gin.Context) {
	// TODO

	var ret SetWorkModeRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	c.JSON(200, ret)
}

type GetWorkModeRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      int `json:"uid" example:123`
		// 是否接单(1:开启，0:关闭)
		Accept  int `json:"accept" example:1`
		// 是否自动接单(1:开启，0:关闭)
		Auto  int `json:"auto" example:1`
	}
}

// @Summary 获取承兑商订单推送模式和开关
// @Tags 承兑商APP API
// @Description 获取承兑商订单推送模式和开关
// @Accept  json
// @Produce  json
// @Param uid  query  string     true        "用户id"
// @Success 200 {object} merchant.GetWorkModeRet ""
// @Router /merchant/settings/workmode [get]
func GetWorkMode(c *gin.Context) {
	// TODO

	var ret GetWorkModeRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	ret.Entity.Accept = 1
	ret.Entity.Auto = 1
	c.JSON(200, ret)
}


type SetIdentifyArg struct {
	// 用户id
	Uid      int `json:"uid" example:123`
	Phone   string `json:"phone" example:13012341234`
	Email   string `json:"email" example:"xxx@xxx.com"`
	IdCard  int `json:"idcard" example:11088888888888888`
}

type SetIdentifyRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      int `json:"uid" example:123`
	}
}


// @Summary 承兑商设置自己的认证信息
// @Tags 承兑商APP API
// @Description 承兑商设置自己的认证信息
// @Accept  json
// @Produce  json
// @Param body  body  merchant.SetIdentifyArg     true        "用户id"
// @Success 200 {object} merchant.SetIdentifyRet ""
// @Router /merchant/settings/identify [put]
func SetIdentify(c *gin.Context) {
	// TODO

	var ret SetWorkModeRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	c.JSON(200, ret)
}

type GetIdentifyRet struct {
	CommonRet
	Entity struct {
		// 用户id
		Uid      int `json:"uid" example:123`
		Phone   string `json:"phone" example:13012341234`
		Email   string `json:"email" example:"xxx@xxx.com"`
		IdCard  string `json:"idcard" example:"11088888888888888"`
	}
}

// @Summary 获取承兑商订单推送模式和开关
// @Tags 承兑商APP API
// @Description 获取承兑商订单推送模式和开关
// @Accept  json
// @Produce  json
// @Param uid  query  string     true        "用户id"
// @Success 200 {object} merchant.GetIdentifyRet ""
// @Router /merchant/settings/identify [get]
func GetIdentify(c *gin.Context) {
	// TODO

	var ret GetIdentifyRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	ret.Entity.Phone = "13012341234"
	ret.Entity.Email = "xxx@xxx.com"
	ret.Entity.IdCard = "11088888888888888"
	c.JSON(200, ret)
}


type MerchantRet struct {
	CommonRet

	Entity struct {

		Data []models.Merchant `json:"data"`
	}
}

// @Summary 获取承兑商列表
// @Tags 管理后台 API
// @Description 坐席获取承兑商列表
// @Accept  json
// @Produce  json
// @Param page query int true "页数"
// @Param size query int true "每页数量"
// @Param user_status query string false "承兑商状态"
// @Param start_time query string false "筛选开始时间"
// @Param stop_time query string false "筛选截止时间"
// @Param time_field query string false "筛选字段"
// @Param search query string false "搜索值"
// @Success 200 {object} merchant.MerchantRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /merchants [get]
func GetMerchants(c *gin.Context) {

}

type RechargeArgs struct {
	Currency string `json:"currency" binding:"required" example:"BTUSD"`
	Count    string `json:"count" binding:"required" example:"200"`
}

type RechargeRet struct {
	CommonRet

	Entity struct {
		//
		Data []models.Merchant `json:"data"`
	}
}

// @Summary 充值
// @Tags 管理后台 API
// @Description 给承兑商充值
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Param body body merchant.RechargeArgs true "充值"
// @Success 200 {object} merchant.RechargeRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /merchants/{uid}/asset [put]
func Recharge(c *gin.Context) {

}

type ApproveArgs struct {
	//审核操作 1：通过 0：不通过
	Operation    int    `json:"operation" binding:"required" example:"1"`
	ContactPhone string `json:"currency" binding:"required" example:"BTUSD"`
	ExtraMessage string `json:"count" binding:"required" example:"200"`
}

type ApproveRet struct {
	CommonRet

	Entity struct {
		//
		Data []models.Merchant `json:"data"`
	}
}

// @Summary 审核
// @Tags 管理后台 API
// @Description 审核承兑商
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Param body body merchant.ApproveArgs true "充值"
// @Success 200 {object} merchant.ApproveRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /merchants/{uid}/approve [put]
func ApproveMerchant(c *gin.Context) {

}

type FreezeArgs struct {
	//冻结操作 1：冻结 0：解冻
	Operation    int    `json:"operation" binding:"required" example:"1"`
	ContactPhone string `json:"currency" binding:"required" example:"BTUSD"`
	ExtraMessage string `json:"count" binding:"required" example:"200"`
}

type FreezeRet struct {
	CommonRet
	Entity struct {


	}
}

// @Summary 冻结
// @Tags 管理后台 API
// @Description 审核承冻结或者解冻
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Param body body merchant.FreezeArgs true "冻结操作"
// @Success 200 {object} merchant.MerchantRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /merchants/{uid}/freeze [put]
func FreezeMerchant(c *gin.Context) {

}

