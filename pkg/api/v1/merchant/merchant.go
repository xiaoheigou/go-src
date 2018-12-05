package merchant

import (
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
		// user_status可以为0/1/2，分别表示“正常/待审核/冻结”
		UserStatus int `json:"user_status" example:0`
		// 客服联系信息
		ContactPhone string `json:"contact_phone" example:0`
		// 额外的信息
		ExtraMessage string `json:"extra_message" example:"您由于xx原因，未通过审核"`
	}
}

// @Summary 获取承兑商账号审核状态
// @Tags 承兑商APP API
// @Description 获取承兑商账号审核状态
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Success 200 {object} merchant.GetAuditStatusRet ""
// @Router /merchant/auditstatus [get]
func GetAuditStatus(c *gin.Context) {
	// TODO

	var ret GetAuditStatusRet
	ret.Status = "success"
	ret.Entity.UserStatus = 1
	ret.Entity.ContactPhone = "13012349876"
	c.JSON(200, ret)
}


type GetProfileRet struct {
	CommonRet
	Entity struct {
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
// @Param uid  path  string     true        "用户id"
// @Success 200 {object} merchant.GetProfileRet ""
// @Router /merchant/profile [get]
func GetProfile(c *gin.Context) {
	// TODO

	var ret GetProfileRet
	ret.Status = "success"
	ret.Entity.NickName = "老王"
	ret.Entity.AssetSymbol = "BTUSD"
	ret.Entity.AssetTotal = "2000"
	ret.Entity.AssetFrozen = "100"
	c.JSON(200, ret)
}
