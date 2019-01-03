// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 失败订单再处理
// @Tags C端相关 API
// @Description 失败订单再处理api
// @Accept  json
// @Produce  json
// @Param  origin_order query string true "平台商订单id"
// @Param  distributor_id query string true "平台商id"
// @Success 200 {object} response.ReprocessOrderResponse "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/order/reprocess [get]
func ReprocessOrder(c *gin.Context) {

	origin_order := c.Query("origin_order")
	distributor_id := c.Query("distributor_id")
	data, err := strconv.ParseInt(distributor_id, 10, 64)
	if err != nil {
		utils.Log.Error("distributor_id convet from string to int64 wrong")
	}

	//签名认证

	apiKey := c.Query("apiKey")
	sign := c.Query("sign")

	method := c.Request.Method
	host := c.Request.Host
	uri := c.Request.URL.Path
	secretKey := service.GetSecretKeyByApiKey(apiKey)
	if secretKey == "" {
		utils.Log.Error("can not get secretkey according to apiKey=[%s] ", apiKey)
		return
	}

	if err != nil {
		utils.Log.Error("struct convert to string wrong,[%v]", err)
	}
	str := service.GenSignatureWith2(method, host, uri, origin_order, distributor_id, apiKey)
	sign1, _ := service.HmacSha256Base64Signer(str, secretKey)
	if sign != sign1 {
		utils.Log.Error("sign is not right,sign=[%v]", sign1)
		c.JSON(403, "you do not have the right to createOrder")
		return
	}


	orderNumber := service.ReprocessOrder(origin_order, data)

	//c.Redirect(301, url)
	c.JSON(200,orderNumber)

}

