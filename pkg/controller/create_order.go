// +build !swagger

package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
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
	var ret response.CreateOrderRet

	body, _ := ioutil.ReadAll(c.Request.Body)
	utils.Log.Debugf("%s", body)
	err := json.Unmarshal(body, &req)
	if err != nil {
		utils.Log.Error("err,%v", err)

	}

	//sha3签名认证

	if utils.Config.Get("signswitch.sign") == "on" {
		apiKey := c.Query("apiKey")
		sign := c.Query("sign")

		method := c.Request.Method
		uri := c.Request.URL.Path
		secretKey := service.GetSecretKeyByApiKey(apiKey)
		if secretKey == "" {
			utils.Log.Error("can not get secretkey according to apiKey=[%s] ", apiKey)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.NoSecretKeyFindErr.Data()
			c.JSON(200, ret)
			return

		}
		utils.Log.Debugf("body is --------:%s", string(body))
		utils.Log.Debugf("method is :%s ,url is:%s,apikey is :%s", method, uri, apiKey)

		str := service.GenSignatureWith(method, uri, string(body), apiKey)
		utils.Log.Debugf("str is +++++++++:%s", str)

		sign1, _ := service.HmacSha256Base64Signer(str, secretKey)
		if sign != sign1 {
			utils.Log.Error("sign is not right,sign=[%v]", sign1)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.IllegalSignErr.Data()
			c.JSON(200, ret)
			return
		}
	}

	ret = service.PlaceOrder(req)

	c.JSON(200, ret)

}
