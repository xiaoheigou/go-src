// +build swagger

package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

func GetAuditStatus(c *gin.Context) {
	// TODO

	var ret response.GetAuditStatusRet
	var data response.GetAuditStatusData
	data.UserStatus = 1
	data.ContactPhone = "13012349876"
	ret.Status = "success"
	ret.Data = []response.GetAuditStatusData{data}
	c.JSON(200, ret)
}

func GetProfile(c *gin.Context) {
	// TODO

	var ret response.GetProfileRet
	var data response.GetProfileData
	ret.Status = "success"
	data.NickName = "老王"
	data.CurrencyCrypto = "BTUSD"
	data.Quantity = 2000
	data.QtyFrozen = 100
	ret.Data = []response.GetProfileData{data}
	c.JSON(200, ret)
}

func SetNickname(c *gin.Context) {
	// TODO

	var ret response.SetNickNameRet
	ret.Status = "success"
	c.JSON(200, ret)
}

func SetWorkMode(c *gin.Context) {
	// TODO

	var ret response.SetWorkModeRet
	ret.Status = "success"
	c.JSON(200, ret)
}

func GetWorkMode(c *gin.Context) {
	// TODO

	var ret response.GetWorkModeRet
	var data response.GetWorkModeData

	ret.Status = "success"
	data.Accept = 1
	data.Auto = 1
	ret.Data = []response.GetWorkModeData{data}
	c.JSON(200, ret)
}

func SetIdentities(c *gin.Context) {
	// TODO

	var ret response.SetWorkModeRet
	ret.Status = "success"
	c.JSON(200, ret)
}

//func GetIdentities(c *gin.Context) {
//	// TODO
//
//	var ret response.GetIdentifyRet
//	ret.Status = "success"
//	ret.Uid = 123
//	ret.Phone = "13012341234"
//	ret.Email = "xxx@xxx.com"
//	ret.IdCard = "11088888888888888"
//	c.JSON(200, ret)
//}

func GetMerchants(c *gin.Context) {
	var ret response.MerchantRet
	ret.Status = "success"
	ret.ErrMsg = "err信息"
	ret.ErrCode = 0
	ret.Data = []models.Merchant{
		{
			Nickname: "1",
			Id:       1,
			Phone:    "13112345678",
		},
		{
			Nickname: "2",
			Id:       2,
			Phone:    "13112345679",
		},
	}

	c.JSON(200, ret)
}

func Recharge(c *gin.Context) {
	var ret response.RechargeRet

	ret.Status = "success"
	c.JSON(200, ret)
}

func ApproveMerchant(c *gin.Context) {
	var ret response.ApproveRet

	ret.Status = "success"
	c.JSON(200, ret)
}

func FreezeMerchant(c *gin.Context) {
	var ret response.ApproveRet
	ret.Status = "success"
	c.JSON(200, ret)
}

func UploadIdentityFile(c *gin.Context) {
	// TODO

	var ret response.UploadIdentityRet
	ret.Status = "success"
	c.JSON(200, ret)
}

func UpdateIdentities(c *gin.Context) {
	// TODO

	var ret response.SetIdentifyRet
	ret.Status = "success"
	c.JSON(200, ret)
}

func GetMerchant(c *gin.Context) {
	uid := c.Param("uid")

	c.JSON(200, service.GetMerchant(uid))
}

func RechargeConfirm(c *gin.Context) {
	session := sessions.Default(c)
	userId := utils.TransformTypeToString(session.Get("userId"))
	uid := c.Param("uid")
	assetApplyId := c.Param("applyId")

	c.JSON(200, service.RechargeConfirm(uid, assetApplyId, userId))
}

func GetBankList(c *gin.Context) {
	c.JSON(200, service.GetBankList())
}
