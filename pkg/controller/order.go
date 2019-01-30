// +build !swagger

package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

// @Summary 获取订单列表
// @Tags 管理后台 API
// @Description 坐席获取订单列表
// @Accept  json
// @Produce  json
// @Param page query int true "页数"
// @Param size query int true "每页数量"
// @Param status query string false "订单状态"
// @Param distributorId query string false "平台商"
// @Param merchantId query string false "承兑商"
// @Param start_time query string false "筛选开始时间"
// @Param stop_time query string false "筛选截止时间"
// @Param time_field query string false "筛选字段"
// @Param origin_order query string false "商户订单号"
// @Param direction query string false "订单类型，0用户充值，1用户提现"
// @Param search query string false "搜索值"
// @Success 200 {object} response.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/orders [get]
func GetOrders(c *gin.Context) {
	session := sessions.Default(c)
	distributor := session.Get("distributor")
	role := session.Get("userRole")

	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")
	status := c.Query("status")
	startTime := c.Query("start_time")
	stopTime := c.Query("stop_time")
	sort := c.DefaultQuery("sort", "desc")
	merchantId := c.Query("merchantId")
	distributorId := c.Query("distributorId")
	timeFiled := c.DefaultQuery("time_field", "created_at")
	direction := c.Query("direction")
	//search only match distributorId and name
	search := c.Query("search")
	originOrder := strings.TrimSpace(c.Query("origin_order"))

	distributorIdTemp := distributor.(int64)
	if distributorIdTemp > 0 && role == 2 {
		c.JSON(200, service.GetOrdersByDistributor(page, size, status, startTime, stopTime, sort, timeFiled, distributorIdTemp, search, originOrder, direction))
	} else if role == 1 && distributorIdTemp == 0 {
		c.JSON(200, service.GetOrders(page, size, status, startTime, stopTime, sort, timeFiled, search, merchantId, distributorId, originOrder, direction))
	} else {
		c.JSON(400,"bad request")
	}
}

// @Summary 获取订单详情
// @Tags 管理后台 API
// @Description 坐席获取订单详情
// @Accept  json
// @Produce  json
// @Param orderNumber path string true "订单号"
// @Success 200 {object} response.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/orders/details/{orderNumber} [get]
func GetOrder(c *gin.Context) {
	orderNumber := c.Param("orderNumber")
	c.JSON(200, service.GetOrderByOrderNumber(orderNumber))
}

// @Summary 获取订单
// @Tags C端相关 API
// @Description 根据订单id查询订单
// @Accept  json
// @Produce  json
// @Param orderNumber path int true "订单id"
// @Success 200 {object} response.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/order/query/{orderNumber} [get]
func GetOrderByOrderNumber(c *gin.Context) {

	id := c.Param("orderNumber")
	var ret response.OrdersRet

	if id == "" {
		utils.Log.Error("orderNumber is null")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoOrderNumberErr.Data()
		c.JSON(200, ret)
		return
	}

	//签名认证
	if utils.Config.Get("signswitch.sign") == "on" {
		method := c.Request.Method
		uri := c.Request.URL.Path
		apiKey := c.Query("apiKey")
		sign := c.Query("sign")
		secretKey := service.GetSecretKeyByApiKey(apiKey)
		if secretKey == "" {
			utils.Log.Errorf("can not get secretkey according to apiKey=[%s] ", apiKey)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.NoSecretKeyFindErr.Data()
			c.JSON(200, ret)
			return
		}
		str := service.GenSignatureWith(method, uri, id, apiKey)
		sign1, _ := service.HmacSha256Base64Signer(str, secretKey)
		if sign != sign1 {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.IllegalSignErr.Data()
			c.JSON(200, ret)
			return

		}

	}

	ret = service.GetOrderByOrderNumber(id)

	c.JSON(200, ret)
}

// @Summary 获取订单列表
// @Tags C端相关 API
// @Description 根据accountId及distributorId获取订单列表
// @Accept  json
// @Produce  json
// @Param page query int false "页数"
// @Param size query int false "每页数量"
// @Param status query string false "订单状态"
// @Param distributor_id query string true "平台商id"
//@Param  account_id query string true "账户id"
// @Param start_time query string false "筛选开始时间"
// @Param stop_time query string false "筛选截止时间"
// @Success 200 {object} response.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/order/list [get]
func GetOrderList(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	size := c.DefaultQuery("size", "10")
	accountId := c.Query("accountId")
	distributorId := c.Query("distributorId")

	var ret response.PageResponse
	ret = service.GetOrderList(page, size, accountId, distributorId)

	c.JSON(200, ret)

}

// @Summary 更新订单
// @Tags C端相关 API
// @Description 更新订单
// @Accept  json
// @Produce  json
// @Param body body response.OrderRequest true "输入参数"
// @Success 200 {object} response.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/order/update [put]
func UpdateOrder(c *gin.Context) {
	var req response.OrderRequest
	var ret response.OrdersRet
	if err := c.ShouldBind(&req); err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
	} else {
		ret = service.UpdateOrder(req)
		ret.Status = response.StatusSucc
		c.JSON(200, ret)
	}
}

// @Summary 创建订单
// @Tags C端相关 API
// @Description 创建订单
// @Accept  json
// @Produce  json
// @Param body body response.OrderRequest true "输入参数"
// @Success 200 {object} response.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/order/add [post]
func AddOrder(c *gin.Context) {
	var req response.OrderRequest
	var ret response.OrdersRet
	if err := c.ShouldBind(&req); err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		c.JSON(200, ret)
	} else {
		ret = service.CreateOrder(req)
		ret.Status = response.StatusSucc
		c.JSON(200, ret)
	}
}

// @Summary 获取订单状态
// @Tags 管理后台 API
// @Description 坐席获取订单列表
// @Accept  json
// @Produce  json
// @Success 200 {object} response.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/orders/status [get]
func GetOrderStatus(c *gin.Context) {
	c.JSON(200, service.GetOrderStatus())
}

// @Summary 重新派单
// @Tags 管理后台 API
// @Description 坐席针对异常订单重新派单
// @Accept  json
// @Produce  json
// @Param orderNumber path int true "订单id"
// @Success 200 {object} response.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/orders/refulfill/{orderNumber} [put]
func RefulfillOrder(c *gin.Context) {
	orderNumber := c.Query("orderNumber")
	c.JSON(200, service.RefulfillOrder(orderNumber))
}

// @Summary 客服放币
// @Tags 管理后台 API
// @Description 客服根据订单放币
// @Accept  json
// @Produce  json
// @Param orderNumber path int true "订单号"
// @Success 200 {object} response.EntityResponse "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/orders/release/{orderNumber} [put]
func ReleaseCoin(c *gin.Context) {
	session := sessions.Default(c)
	userId := utils.TransformTypeToString(session.Get("userId"))
	userName := utils.TransformTypeToString(session.Get("username"))
	orderNumber := c.Param("orderNumber")
	id, _ := strconv.ParseInt(userId, 10, 64)

	c.JSON(200, service.ReleaseCoin(orderNumber, userName, id))
}

// @Summary 客服解冻
// @Tags 管理后台 API
// @Description 客服根据订单解冻
// @Accept  json
// @Produce  json
// @Param orderNumber path int true "订单号"
// @Success 200 {object} response.RechargeRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/orders/unfreeze/{orderNumber} [put]
func UnFreezeCoin(c *gin.Context) {
	session := sessions.Default(c)
	userId := utils.TransformTypeToString(session.Get("userId"))
	userName := utils.TransformTypeToString(session.Get("username"))
	orderNumber := c.Param("orderNumber")
	id, _ := strconv.ParseInt(userId, 10, 64)

	c.JSON(200, service.UnFreezeCoin(orderNumber, userName, id))
}

// @Summary 申诉订单
// @Tags C端相关 API
// @Description 客服根据订单解冻
// @Accept  json
// @Produce  json
// @Param orderNumber path int true "订单号"
// @Success 200 {object} response.RechargeRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /c/orders/compliant/{orderNumber} [post]
func Compliant(c *gin.Context) {
	var ret response.EntityResponse
	orderNumber := c.Param("orderNumber")
	body, _ := ioutil.ReadAll(c.Request.Body)

	c.Header("order", string(body))

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

	ret.Status = response.StatusSucc
	engine := service.NewOrderFulfillmentEngine(nil)
	engine.DeleteWheel(orderNumber)
	if err := service.ModifyOrderAsCompliant(orderNumber); err == nil {
		utils.Log.Debugf("update order status is suspended and status reason is compliant success,orderNumber:%s", orderNumber)
	} else if err.Error() == "order status is final,not allow to update" || err.Error() == "not found order number" {
		utils.Log.Errorf("order status is already final,stop ticker")
	} else {
		// 当更新订单申诉状态失败时，定时修改订单，知道修改完成
		go func() {
			ticker := time.NewTicker(time.Duration(30) * time.Second)
			for {
				select {
				case <-ticker.C:
					if err := service.ModifyOrderAsCompliant(orderNumber); err == nil {
						utils.Log.Debugf("update order status is suspended and status reason is compliant,orderNumber:%s", orderNumber)
						ticker.Stop()
						return
					} else if err.Error() == "order status is final,not allow to update" || err.Error() == "not found order number" {
						utils.Log.Errorf("order status or not found order,stop ticker")
						ticker.Stop()
						return
					}
					utils.Log.Errorf("update order status is suspended and status reason is compliant,orderNumber:%s", orderNumber)
				}
			}

		}()
	}

	c.JSON(200, ret)
}
