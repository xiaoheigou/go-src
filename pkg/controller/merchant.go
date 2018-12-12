// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 获取承兑商账号审核状态
// @Tags 承兑商APP API
// @Description 获取承兑商账号审核状态
// @Accept  json
// @Produce  json
// @Param uid  path  int  true  "用户id"
// @Success 200 {object} response.GetAuditStatusRet ""
// @Router /m/merchants/{uid}/audit-status [get]
func GetAuditStatus(c *gin.Context) {
	// TODO

	var ret response.GetAuditStatusRet
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.GetAuditStatusData{
		UserStatus:   0,
		ContactPhone: "13000000000",
		ExtraMessage: "由于xx原因，您没有通过审核。",
	})
	c.JSON(200, ret)
}

// @Summary 获取承兑商个人信息
// @Tags 承兑商APP API
// @Description 获取承兑商个人信息
// @Accept  json
// @Produce  json
// @Param uid  path  int  true  "用户id"
// @Success 200 {object} response.GetProfileRet ""
// @Router /m/merchants/{uid}/profile [get]
func GetProfile(c *gin.Context) {
	// TODO

	var ret response.GetProfileRet
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.GetProfileData{
		NickName:       "老王",
		CurrencyCrypto: "BTUSD",
		Quantity:       10000,
		QtyFrozen:      200,
	})
	c.JSON(200, ret)
}

// @Summary 设置承兑商昵称
// @Tags 承兑商APP API
// @Description 设置承兑商昵称
// @Accept  json
// @Produce  json
// @Param uid  path  int  true  "用户id"
// @Param body  body  response.SetNickNameArg     true        "新参数"
// @Success 200 {object} response.SetNickNameRet ""
// @Router /m/merchants/{uid}/settings/nickname [put]
func SetNickName(c *gin.Context) {
	var json response.SetNickNameArg
	if err := c.ShouldBindJSON(&json); err != nil {
		var retFail response.SetNickNameRet
		retFail.Status = response.StatusFail
		retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, retFail)
		return
	}

	var uid int
	var err error
	if uid, err = strconv.Atoi(c.Param("uid")); err != nil {
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		var ret response.SetNickNameRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return
	}

	c.JSON(200, service.SetMerchantNickName(uid, json))
	return

	//var ret response.SetNickNameRet
	//ret.Status = response.StatusSucc
	//c.JSON(200, ret)
}

// @Summary 承兑商设置订单推送模式和开关
// @Tags 承兑商APP API
// @Description 承兑商设置订单推送模式和开关
// @Accept  json
// @Produce  json
// @Param uid  path  int  true  "用户id"
// @Param body  body  response.SetWorkModeArg     true        "新参数"
// @Success 200 {object} response.SetWorkModeRet ""
// @Router /m/merchants/{uid}/settings/work-mode [put]
func SetWorkMode(c *gin.Context) {
	// TODO

	var ret response.SetWorkModeRet
	ret.Status = response.StatusSucc
	c.JSON(200, ret)
}

// @Summary 获取承兑商订单推送模式和开关
// @Tags 承兑商APP API
// @Description 获取承兑商订单推送模式和开关
// @Accept  json
// @Produce  json
// @Param uid  path  string  true  "用户id"
// @Success 200 {object} response.GetWorkModeRet ""
// @Router /m/merchants/{uid}/settings/work-mode [get]
func GetWorkMode(c *gin.Context) {
	// TODO

	var ret response.GetWorkModeRet
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.GetWorkModeData{
		Accept: 1,
		Auto:   1,
	})
	c.JSON(200, ret)
}

// @Summary 承兑商设置自己的认证信息
// @Tags 承兑商APP API
// @Description 承兑商设置自己的认证信息
// @Accept  json
// @Produce  json
// @Param uid  path  int  true  "用户id"
// @Param body  body  response.SetIdentifyArg     true        "新参数"
// @Success 200 {object} response.SetIdentifyRet ""
// @Router /m/merchants/{uid}/settings/identities [post]
func SetIdentities(c *gin.Context) {
	// TODO

	var ret response.SetIdentifyRet
	ret.Status = "success"
	c.JSON(200, ret)
}

// @Summary 承兑商未通过认证时更新认证信息
// @Tags 承兑商APP API
// @Description 承兑商未通过认证时更新认证信息
// @Accept  json
// @Produce  json
// @Param uid  path  int  true  "用户id"
// @Param body  body  response.SetIdentifyArg  true  "新参数"
// @Success 200 {object} response.SetIdentifyRet ""
// @Router /m/merchants/{uid}/settings/identities [put]
func UpdateIdentities(c *gin.Context) {
	// TODO

	var ret response.SetIdentifyRet
	ret.Status = "success"
	c.JSON(200, ret)
}

// @Summary 承兑商上传认证图片（身份证图片）
// @Tags 承兑商APP API
// @Description 承兑商上传认证图片（身份证图片）
// @Accept  jpeg
// @Produce  json
// @Param uid  path  string     true        "用户id"
// @Param type  query  string     true        "0/1分别表示图片正面/反面"
// @Success 200 {object} response.UploadIdentityRet ""
// @Router /m/merchants/{uid}/settings/identity/upload [post]
func UploadIdentityFile(c *gin.Context) {
	// TODO

	var ret response.UploadIdentityRet
	ret.Status = "success"
	c.JSON(200, ret)
}

// @Summary 承兑商提交申诉
// @Tags 承兑商APP API
// @Description 承兑商提交申诉
// @Accept  json
// @Produce  json
// @Param order-id  path  string     true        "订单id"
// @Param body body response.OrderComplainArg true "参数"
// @Success 200 {object} response.OrderComplaintRet ""
// @Router /m/orders/{order-id}/complaint [post]
func OrderComplaint(c *gin.Context) {
	// TODO

	var ret response.OrderComplaintRet
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
// @Param user_cert query string false "承兑商认证状态"
// @Param start_time query string false "筛选开始时间"
// @Param stop_time query string false "筛选截止时间"
// @Param time_field query string false "筛选字段 created_at updated_at"
// @Param sort query string false "排序方式 desc asc"
// @Param search query string false "搜索值"
// @Success 200 {object} response.MerchantRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/merchants [get]
func GetMerchants(c *gin.Context) {
	page := c.DefaultQuery("page","0")
	size := c.DefaultQuery("size","10")
	userStatus := c.Query("user_status")
	userCert := c.Query("user_cert")
	startTime := c.Query("start_time")
	stopTime := c.Query("stop_time")
	timeField := c.DefaultQuery("time_field","created_at")
	sort := c.DefaultQuery("sort","desc")
	search := c.Query("search")
	c.JSON(200, service.GetMerchants(page,size,userStatus,userCert,startTime,stopTime,timeField,sort,search))
}

// @Summary 审核
// @Tags 管理后台 API
// @Description 审核承兑商
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Param body body response.ApproveArgs true "审核参数"
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
	c.JSON(200, ret)
}

// @Summary 冻结
// @Tags 管理后台 API
// @Description 审核承冻结或者解冻
// @Accept  json
// @Produce  json
// @Param uid path int true "用户id"
// @Param body body response.FreezeArgs true "冻结操作"
// @Success 200 {object} response.FreezeRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/merchants/{uid}/freeze [put]
func FreezeMerchant(c *gin.Context) {
	var args response.FreezeArgs
	err := c.ShouldBind(&args)
	var ret response.FreezeRet
	ret.Status = "fail"
	ret.ErrCode = 0
	ret.ErrMsg = "test1"
	if err != nil {
		c.JSON(200, ret)
	}
	ret.Status = "success"
	ret.Data = make([]response.FreezeDataResponse,1)
	ret.Data[0].Uid = 1
	ret.Data[0].Status = 1
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
	uid := c.Param("uid")

	c.JSON(200, service.GetMerchant(uid))
}