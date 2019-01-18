// +build !swagger

package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
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
// @Param body body response.BuyOrderRequest true "请求体"
// @Success 200 {object} response.CreateOrderRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/create-order/buy [post]
func BuyOrder(c *gin.Context) {
	var req response.BuyOrderRequest
	var ret response.CreateOrderRet
	var createOrderReq response.CreateOrderRequest

	body, _ := ioutil.ReadAll(c.Request.Body)
	utils.Log.Debugf("%s", body)
	err := json.Unmarshal(body, &req)
	if err != nil {
		utils.Log.Error("err,%v", err)

	}
	c.Header("order", string(body))
	createOrderReq = service.BuyOrderReq2CreateOrderReq(req)

	//sha3签名认证

	if utils.Config.Get("signswitch.sign") == "on" {
		apiKey := c.Query("appApiKey")
		sign := c.Query("appSignContent")

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

	ret = service.PlaceOrder(createOrderReq, c)
	//if ret.Status == response.StatusFail {
	c.JSON(200, ret)
	//}

}

// @Summary c端客户下单
// @Tags C端相关 API
// @Description c端客户下单api(提现)
// @Accept  json
// @Produce  json
// @Param body body response.SellOrderRequest true "请求体"
// @Success 200 {object} response.CreateOrderRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/create-order/sell [post]
func SellOrder(c *gin.Context) {
	var req response.SellOrderRequest
	var ret response.CreateOrderRet
	var createOrderReq response.CreateOrderRequest

	body, _ := ioutil.ReadAll(c.Request.Body)
	utils.Log.Debugf("%s", body)
	err := json.Unmarshal(body, &req)
	if err != nil {
		utils.Log.Error("err,%v", err)

	}

	createOrderReq = service.SellOrderReq2CreateOrderReq(req)

	//sha3签名认证

	if utils.Config.Get("signswitch.sign") == "on" {
		apiKey := c.Query("appApiKey")
		sign := c.Query("appSignContent")

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

	ret = service.PlaceOrder(createOrderReq, c)

	//if ret.Status == response.StatusFail {
	c.JSON(200, ret)
	//}
}

// @Summary C端客户下单签名
// @Tags C端相关 API
// @Description C端客户下单签名
// @Accept  json
// @Produce  json
// @Param appId query string true "平台商id"
// @Param appKey query string true "平台商appKey"
// @Param body body response.SignatureRequest true "请求体"
// @Success 200 {object} response.SignatureRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/signature [post]
func SignFor(c *gin.Context) {
	var ret response.SignatureRet

	var json response.SignatureRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		utils.Log.Error("func SignFor, invalid arg")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
		return
	}

	appId := c.Query("appId") // 平台商id
	appKey := c.Query("appKey")

	utils.Log.Debugf("query param appId = %s, appKey = %s", appId, appKey)
	if strings.TrimSpace(appId) == "" {
		utils.Log.Error("func SignFor, appId is empty")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
		return
	}
	if strings.TrimSpace(appKey) == "" {
		utils.Log.Error("func SignFor, appKey is empty")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
		return
	}

	var secretKey string
	var err error
	if secretKey, err = service.GetApiSecretByIdAndAPIKey(appId, appKey); err != nil {
		utils.Log.Error("can not get secretkey for apiKey=[%s] (distributor = %s)", appKey, appId)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoSecretKeyFindErr.Data()
		c.JSON(200, ret)
		return
	}
	if secretKey == "" {
		utils.Log.Error("secretKey is empty for apiKey=[%s] (distributor = %s)", appKey, appId)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoSecretKeyFindErr.Data()
		c.JSON(200, ret)
		return
	}

	signData, err := base64.StdEncoding.DecodeString(json.SignDataBase64)
	if err != nil {
		utils.Log.Errorf("signDataBase64 is not encoded with base64. error:", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
		return
	}
	utils.Log.Debugf("func SignFor, the request data = [%+v]", json)

	hasher := hmac.New(sha256.New, []byte(secretKey))
	if _, err := hasher.Write(signData); err != nil {
		utils.Log.Error("")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoSecretKeyFindErr.Data()
		c.JSON(200, ret)
		return
	}
	sign := fmt.Sprintf("%x", hasher.Sum(nil))

	utils.Log.Debugf("func SignFor, sign = %s", sign)

	ret.Data = append(ret.Data, response.SignatureRetData{
		AppSignContent: sign,
	})

	c.JSON(200, ret)
}
