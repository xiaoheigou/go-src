// +build swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
)

func GetMerchantAssetHistory(c *gin.Context) {
	var ret response.GetMerchantAssetHistoryRet
	ret.Status = "success"
	ret.ErrCode = 0
	ret.ErrMsg = "test"
	ret.Entity.Data = []models.AssetHistory{
		{
			Id:         1,
			Msg:        "123充值了 500",
			MerchantId: 1,
			Currency: "BTUSD",
			IsOrder: 0,
			Operation: 0,
			OperatorId: "1",
			DistributorId: 0,
			Timestamp: "2018-12-12 12:12:12",
		},
	}
	c.JSON(200, ret)
}
