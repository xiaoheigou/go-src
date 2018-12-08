// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
)

// @Summary 获取承兑商账号审核状态
// @Tags 承兑商APP API
// @Description 获取承兑商账号审核状态
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Success 200 {object} response.GetAuditStatusRet ""
// @Router /m/merchants/{uid}/audit-status [get]
func GetAuditStatus(c *gin.Context) {
	// TODO

	var ret response.GetAuditStatusRet
	ret.Status = "success"
	ret.Uid = 123
	ret.UserStatus = 1
	ret.ContactPhone = "13012349876"
	c.JSON(200, ret)
}

// @Summary 获取承兑商个人信息
// @Tags 承兑商APP API
// @Description 获取承兑商个人信息
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Success 200 {object} response.GetProfileRet ""
// @Router /m/merchants/{uid}/profile [get]
func GetProfile(c *gin.Context) {
	// TODO

	var ret response.GetProfileRet
	ret.Status = "success"
	ret.Uid = 123
	ret.NickName = "老王"
	ret.AssetSymbol = "BTUSD"
	ret.AssetTotal = "2000"
	ret.AssetFrozen = "100"
	c.JSON(200, ret)
}

// @Summary 设置承兑商昵称
// @Tags 承兑商APP API
// @Description 设置承兑商昵称
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Param body  body  response.SetNickNameArg     true        "新参数"
// @Success 200 {object} response.SetNickNameRet ""
// @Router /m/merchants/{uid}/settings/nickname [put]
func SetNickName(c *gin.Context) {
	// TODO

	var ret response.SetNickNameRet
	ret.Status = "success"
	ret.Uid = 123
	c.JSON(200, ret)
}

// @Summary 承兑商设置订单推送模式和开关
// @Tags 承兑商APP API
// @Description 承兑商设置订单推送模式和开关
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Param body  body  response.SetWorkModeArg     true        "新参数"
// @Success 200 {object} response.SetWorkModeRet ""
// @Router /m/merchants/{uid}/settings/work-mode [put]
func SetWorkMode(c *gin.Context) {
	// TODO

	var ret response.SetWorkModeRet
	ret.Status = "success"
	ret.Uid = 123
	c.JSON(200, ret)
}

// @Summary 获取承兑商订单推送模式和开关
// @Tags 承兑商APP API
// @Description 获取承兑商订单推送模式和开关
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Success 200 {object} response.GetWorkModeRet ""
// @Router /m/merchants/{uid}/settings/work-mode [get]
func GetWorkMode(c *gin.Context) {
	// TODO

	var ret response.GetWorkModeRet
	ret.Status = "success"
	ret.Uid = 123
	ret.Accept = 1
	ret.Auto = 1
	c.JSON(200, ret)
}

// @Summary 承兑商设置自己的认证信息
// @Tags 承兑商APP API
// @Description 承兑商设置自己的认证信息
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Param body  body  response.SetIdentifyArg     true        "新参数"
// @Success 200 {object} response.SetIdentifyRet ""
// @Router /m/merchants/{uid}/settings/identities [put]
func SetIdentities(c *gin.Context) {
	// TODO

	var ret response.SetWorkModeRet
	ret.Status = "success"
	ret.Uid = 123
	c.JSON(200, ret)
}

// @Summary 获取承兑商的认证信息
// @Tags 承兑商APP API
// @Description 获取承兑商的认证信息
// @Accept  json
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Success 200 {object} response.GetIdentifyRet ""
// @Router /m/merchants/{uid}/settings/identities [get]
func GetIdentities(c *gin.Context) {
	// TODO

	var ret response.GetIdentifyRet
	ret.Status = "success"
	ret.Uid = 123
	ret.Phone = "13012341234"
	ret.Email = "xxx@xxx.com"
	ret.IdCard = "11088888888888888"
	c.JSON(200, ret)
}

// @Summary 承兑商提交申诉
// @Tags 承兑商APP API
// @Description 承兑商提交申诉
// @Accept  json
// @Produce  json
// @Param order-id  path  string     true        "订单id"
// @Param body body response.OrderComplainArg true "参数"
// @Success 200 {object} response.OrderComplainRet ""
// @Router /m/orders/{order-id}/complain [post]
func OrderComplain(c *gin.Context) {
	// TODO

	var ret response.OrderComplainRet
	ret.Status = "success"
	c.JSON(200, ret)
}

// @Summary 承兑商获取它提交的申诉列表
// @Tags 承兑商APP API
// @Description 承兑商获取它提交的申诉列表
// @Accept  json
// @Produce  json
// @Param uid  query  string     true        "承兑商用户id"
// @Success 200 {object} response.GetComplainsRet ""
// @Router /m/merchant/complains [post]
func GetComplains(c *gin.Context) {

	
	var ret response.OrderComplainRet
	ret.Status = "success"
	c.JSON(200, ret)
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
// @Success 200 {object} response.MerchantRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/merchants [get]
func GetMerchants(c *gin.Context) {
	var ret response.MerchantRet
	ret.Status = "success"
	ret.ErrMsg = "err信息"
	ret.ErrCode = 0
	ret.Data = []models.Merchant{
		{
			NickName: "1",
			Id:       1,
			Phone:    "13112345678",
		},
		{
			NickName: "2",
			Id:       2,
			Phone:    "13112345679",
		},
	}

	c.JSON(200, ret)
}

// @Summary 充值
// @Tags 管理后台 API
// @Description 给承兑商充值
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Param body body response.RechargeArgs true "充值"
// @Success 200 {object} response.RechargeRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/merchants/{uid}/asset [put]
func Recharge(c *gin.Context) {
	var args response.RechargeArgs
	err := c.ShouldBind(&args)
	var ret response.RechargeRet
	ret.Status = "fail"
	ret.ErrCode = 0
	ret.ErrMsg = "test1"
	if err != nil {
		c.JSON(200, ret)
	}
	ret.Status = "success"
	ret.Entity.Balance = args.Count
	c.JSON(200, ret)
}

// @Summary 审核
// @Tags 管理后台 API
// @Description 审核承兑商
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Param body body response.ApproveArgs true "充值"
// @Success 200 {object} response.ApproveRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/merchants/{uid}/approve [put]
func ApproveMerchant(c *gin.Context) {
	var args response.ApproveArgs
	err := c.ShouldBind(&args)
	var ret response.ApproveRet
	ret.Status = "fail"
	ret.ErrCode = 0
	ret.ErrMsg = "test1"
	if err != nil {
		c.JSON(200, ret)
	}
	ret.Status = "success"
	ret.Entity.Uid = 1
	ret.Entity.Status = 1
	c.JSON(200, ret)
}

// @Summary 冻结
// @Tags 管理后台 API
// @Description 审核承冻结或者解冻
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Param body body response.FreezeArgs true "冻结操作"
// @Success 200 {object} response.MerchantRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/merchants/{uid}/freeze [put]
func FreezeMerchant(c *gin.Context) {
	var args response.ApproveArgs
	err := c.ShouldBind(&args)
	var ret response.ApproveRet
	ret.Status = "fail"
	ret.ErrCode = 0
	ret.ErrMsg = "test1"
	if err != nil {
		c.JSON(200, ret)
	}
	ret.Status = "success"
	ret.Entity.Uid = 1
	ret.Entity.Status = 1
	c.JSON(200, ret)
}

// @Summary 获取承兑商
// @Tags 管理后台 API
// @Description 审核承冻结或者解冻
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Success 200 {object} response.MerchantRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/merchants/{uid} [get]
func GetMerchant(c *gin.Context) {
	var ret response.MerchantRet

	ret.Status = "success"
	ret.Data = []models.Merchant{{
		Id:1,
		NickName:"test",
		Phone:"13112345678",
	}}

	c.JSON(200, ret)
}