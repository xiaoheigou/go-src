// +build !swagger

package controller

import (
	"github.com/gin-gonic/gin"
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
// @Param search query string false "搜索值"
// @Success 200 {object} response.OrdersRet "成功（status为success）失败（status为fail）都会返回200"
// @Router /w/orders [get]
func GetOrders(c *gin.Context) {
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
	c.JSON(200, service.GetOrders(page, size, status, startTime, stopTime, sort, timeFiled, search, merchantId, distributorId, direction))
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
// @Router /w/orders/{orderNumber}/refulfill [put]
func RefulfillOrder(c *gin.Context) {
	orderNumber := c.Query("orderNumber")
	c.JSON(200, service.RefulfillOrder(orderNumber))
}
