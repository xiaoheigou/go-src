package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 上传截获的支付宝或微信账单
// @Tags 承兑商APP API
// @Description 上传截获的支付宝或微信账单
// @Accept  json
// @Produce  json
// @Param uid  path  int  true  "币商id"
// @Param body  body  response.UploadBillArg  true  "上传账单的数据格式"
// @Success 200 {object} response.CommonRet ""
// @Router /m/merchants/{uid}/bills [post]
func UploadBills(c *gin.Context) {
	var uid int64
	var err error
	if uid, err = strconv.ParseInt(c.Param("uid"), 10, 64); err != nil {
		var retFail response.CommonRet
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		retFail.Status = response.StatusFail
		retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, retFail)
		return
	}

	var json response.UploadBillArg
	if err := c.ShouldBindJSON(&json); err != nil {
		var retFail response.CommonRet
		retFail.Status = response.StatusFail
		retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, retFail)
		return
	}

	c.JSON(200, service.UploadBills(uid, json))
}
