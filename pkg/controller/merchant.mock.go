// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
)

func GetAuditStatus(c *gin.Context) {
	// TODO

	var ret response.GetAuditStatusRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	ret.Entity.UserStatus = 1
	ret.Entity.ContactPhone = "13012349876"
	c.JSON(200, ret)
}

func GetProfile(c *gin.Context) {
	// TODO

	var ret response.GetProfileRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	ret.Entity.NickName = "老王"
	ret.Entity.AssetSymbol = "BTUSD"
	ret.Entity.AssetTotal = "2000"
	ret.Entity.AssetFrozen = "100"
	c.JSON(200, ret)
}

func SetNickName(c *gin.Context) {
	// TODO

	var ret response.SetNickNameRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	c.JSON(200, ret)
}

func SetWorkMode(c *gin.Context) {
	// TODO

	var ret response.SetWorkModeRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	c.JSON(200, ret)
}

func GetWorkMode(c *gin.Context) {
	// TODO

	var ret response.GetWorkModeRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	ret.Entity.Accept = 1
	ret.Entity.Auto = 1
	c.JSON(200, ret)
}

func SetIdentify(c *gin.Context) {
	// TODO

	var ret response.SetWorkModeRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	c.JSON(200, ret)
}

func GetIdentify(c *gin.Context) {
	// TODO

	var ret response.GetIdentifyRet
	ret.Status = "success"
	ret.Entity.Uid = 123
	ret.Entity.Phone = "13012341234"
	ret.Entity.Email = "xxx@xxx.com"
	ret.Entity.IdCard = "11088888888888888"
	c.JSON(200, ret)
}

func GetMerchants(c *gin.Context) {
	var ret response.MerchantRet
	ret.Status = "success"
	ret.ErrMsg = "err信息"
	ret.ErrCode = 0
	ret.Entity.Data = []models.Merchant{
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
