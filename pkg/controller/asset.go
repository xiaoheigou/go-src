package controller

import (
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"github.com/gin-gonic/gin"
)

// @Summary 获取资金变动历史
// @Tags 管理后台 API
// @Description 查看资金变动历史
// @Accept  json
// @Produce  json
// @Param page query int true "页数"
// @Param size query int true "每页数量"
// @Param search query string false "搜索值"
// @Success 200 {object} response.GetMerchantAssetHistoryRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /merchants/{uid}/assets/history [get]
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
		},
	}
	c.JSON(200,ret)
}
