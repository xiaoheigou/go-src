// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary c端客户下单
// @Tags C端相关 API
// @Description c端客户下单api
// @Accept  json
// @Produce  json
// @Param body body response.CreateOrderRequest true "请求体"
// @Success 200 {object} response.CreateOrderRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/create-order [post]
func CreateOrder(c *gin.Context) {
	var req response.CreateOrderRequest
	var orderNumber string

	if err := c.ShouldBind(&req); err != nil {
		utils.Log.Error("request param is error,%v", err)
	}

	//sha3签名认证

	if utils.Config.Get("signswitch.sign")=="on"{
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
		body, err := service.Struct2JsonString(req)

		if err != nil {
			utils.Log.Error("struct convert to string wrong,[%v]", err)
		}
		str := service.GenSignatureWith(method, host, uri, body, apiKey)
		sign1, _ := service.HmacSha256Base64Signer(str, secretKey)
		if sign != sign1 {
			utils.Log.Error("sign is not right,sign=[%v]", sign1)
			c.JSON(403, "you do not have the right to createOrder")
			return
		}
	}


	orderNumber = service.PlaceOrder(req)

	//c.Redirect(301, redirectUrl)
	c.JSON(200,orderNumber)

}
