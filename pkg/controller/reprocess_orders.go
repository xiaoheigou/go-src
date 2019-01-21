// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 失败订单再处理
// @Tags C端相关 API
// @Description 失败订单再处理api
// @Accept  json
// @Produce  json
// @Param  appOrderId query string true "平台商订单id"
// @Param  appId query string true "平台商id"
// @Success 200 {object} response.ReprocessOrderResponse "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/order/detail [get]
func ReprocessOrder(c *gin.Context) {
	var ret response.OrdersRet

	origin_order := c.Query("appOrderId")
	distributor_id := c.Query("appId")
	if origin_order == "" || distributor_id == "" {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
		return
	}
	data, err := strconv.ParseInt(distributor_id, 10, 64)
	if err != nil {
		utils.Log.Error("distributor_id convet from string to int64 wrong")
	}

	apiKey := c.Query("appApiKey")
	//签名认证
	if utils.Config.Get("signswitch.sign") == "on" {
		sign := c.Query("appSignContent")

		method := c.Request.Method
		uri := c.Request.URL.Path
		secretKey := service.GetSecretKeyByApiKey(apiKey)
		if secretKey == "" {
			utils.Log.Errorf("can not get secretkey according to apiKey=[%s] ", apiKey)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.NoSecretKeyFindErr.Data()
			c.JSON(200, ret)
			return
		}

		if err != nil {
			utils.Log.Errorf("struct convert to string wrong,[%v]", err)
		}
		str := service.GenSignatureWith2(method, uri, origin_order, distributor_id, apiKey)
		utils.Log.Debugf("%s", str)
		sign1, _ := service.HmacSha256Base64Signer(str, secretKey)
		if sign != sign1 {
			utils.Log.Errorf("sign is not right,sign=[%v]", sign1)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.IllegalSignErr.Data()
			c.JSON(200, ret)
			return
		}
	}

	contentType := c.ContentType()
	switch contentType {
	case "text/html":
		ret = service.ReprocessOrder(origin_order, data)

		if ret.Status == response.StatusFail {
			c.JSON(200, ret)
			return
		}
	case "application/json":
		var resp response.EntityResponse
		order := models.Order{}
		if utils.DB.First(&order, "origin_order = ?", origin_order).RecordNotFound() {

			resp.Status = response.StatusFail
			resp.ErrCode, resp.ErrMsg = err_code.NoOrderFindErr.Data()
			c.JSON(200, resp)
			return
		}
		//result := response.OrderRet{
		//	OrderStatus:     order.Status,
		//	Direction:       order.Direction,
		//	AppId:           order.DistributorId,
		//	AppOrderNo:      origin_order,
		//	AppCoinName:     order.AppCoinName,
		//	AppCoinRate:     order.Price,
		//	OrderPayTypeId:  order.PayType,
		//	AppUserId:       order.AccountId,
		//	AppCoinSymbol:   order.CurrencyFiat,
		//	OrderCoinAmount: order.Quantity,
		//	PayAccountUser:  order.Name,
		//	OrderRemark:     order.Remark,
		//}
		//if order.PayType <= 2 {
		//	result.PayAccountId = order.QrCode
		//} else if order.PayType > 2 {
		//	result.PayAccountId = order.BankAccount
		//	result.PayAccountInfo = order.BankBranch
		//}
		resp.Status = response.StatusSucc
		resp.Data = []models.Order{order}
		c.JSON(200, resp)
	default:
		c.JSON(400, "bad request")
	}
}
