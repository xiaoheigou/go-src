package service

import (
	"fmt"
	"github.com/typa01/go-utils"
	"math"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

func GetOrderList(page, size, accountId string, distributorId string) response.PageResponse {
	var ret response.PageResponse
	var data []models.Order
	if accountId == "" || distributorId == "" {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoAccountIdOrDistributorIdErr.Data()
	} else {
		db := utils.DB.Model(&models.Order{}).Where("account_id=? and distributor_id=?", accountId, distributorId)
		pageNum, err := strconv.ParseInt(page, 10, 64)
		pageSize, err1 := strconv.ParseInt(size, 10, 64)
		if err != nil || err1 != nil {
			utils.Log.Error(pageNum, pageSize)
		}
		db = db.Offset((pageNum - 1) * pageSize).Limit(pageSize)
		db.Count(&ret.PageCount)
		ret.PageNum = int(pageNum + 1)
		ret.PageSize = int(pageSize)
		db.Find(&data)
		ret.Data = data
		ret.Status = response.StatusSucc
	}

	return ret

}

func GetOrderByOrderNumber(orderId string) response.OrdersRet {
	var ret response.OrdersRet
	var data models.Order
	if error := utils.DB.First(&data, "order_number=?", orderId).Error; error != nil {
		utils.Log.Error(error)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoOrderFindErr.Data()
		return ret
	}
	ret.Data = []models.Order{data}
	ret.Status = response.StatusSucc
	return ret

}

func GetOrderByMerchantIdAndOrderNumber(merchantId int64, orderNumber string) response.OrdersRet {
	var ret response.OrdersRet
	var data models.Order
	if error := utils.DB.First(&data, "order_number=? and merchant_id = ?", orderNumber, merchantId).Error; error != nil {
		utils.Log.Error(error)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoOrderFindErr.Data()
		return ret
	}
	ret.Data = []models.Order{data}
	ret.Status = response.StatusSucc
	return ret
}

func GetOrders(page, size, status, startTime, stopTime, sort, timeField, search string) response.PageResponse {
	var result []models.Order
	var ret response.PageResponse
	db := utils.DB.Model(&models.Order{}).Order(fmt.Sprintf("%s %s", timeField, sort))
	if search != "" {
		db = db.Where("merchant_id = ? OR distributor_id = ?", search, search)
	} else {

		if startTime != "" && stopTime != "" {
			db = db.Where(fmt.Sprintf("%s >= ? AND %s <= ?", timeField, timeField), startTime, stopTime)
		}
		if status != "" {
			db = db.Where("status = ?", status)
		}
		db.Count(&ret.TotalCount)
		pageNum, err := strconv.ParseInt(page, 10, 64)
		pageSize, err1 := strconv.ParseInt(size, 10, 64)
		if err != nil || err1 != nil {
			utils.Log.Error(pageNum, pageSize)
		}
		db = db.Offset((pageNum - 1) * pageSize).Limit(pageSize)
		ret.PageNum = int(pageNum)
		ret.PageSize = int(pageSize)
	}
	db.Find(&result)
	ret.PageCount = len(result)

	var merchants []models.Merchant
	var distributors []models.Distributor
	var merchantIds, distributorIds []int64
	for _, order := range result {
		merchantIds = append(merchantIds, order.MerchantId)
		distributorIds = append(distributorIds, order.DistributorId)
	}
	//查询符合订单的币商和平台商
	if err := utils.DB.Find(&merchants, "id in (?)", merchantIds).Error; err != nil {
		utils.Log.Errorf("get merchant name is failed,merchantIds is %v", merchantIds)
	}
	if err := utils.DB.Find(&distributors, "id in (?)", distributorIds).Error; err != nil {
		utils.Log.Errorf("get distributor name is failed,distributorIds is %v", distributorIds)
	}

	//遍历获取平台商和币商的名字
	for i, order := range result {
		for _, merchant := range merchants {
			if order.MerchantId == merchant.Id {
				order.MerchantName = merchant.Nickname

				break
			}
		}

		for _, distributor := range distributors {
			if order.DistributorId == distributor.Id {
				order.DistributorName = distributor.Name
				break
			}
		}
		result[i] = order
	}

	ret.Status = response.StatusSucc
	ret.Data = result
	return ret
}

//平台管理员按照创建时间（start-end),订单状态，平台商标识，承兑商标识组合搜索条件查询订单列表；
func GetOrdersByAdmin(page int, size int, status int, startTime string, stopTime string, sort string, timeField string, distributorId int64, merchantId int64, orderNumber string) response.PageResponse {
	var order models.Order
	var orderList []models.Order
	var ret response.PageResponse
	if orderNumber != "" {
		resp := GetOrderByOrderNumber(orderNumber)
		ret.EntityResponse.CommonRet = resp.CommonRet
		ret.EntityResponse.Data = resp.Data
		return ret
	}
	db := utils.DB.Model(&order).Order(fmt.Sprintf("%s %s", timeField, sort))
	db = db.Offset((page - 1) * size).Limit(size)
	if startTime != "" && stopTime != "" {
		db = db.Where(fmt.Sprintf("%s >= ? AND %s <= ?", timeField, timeField), startTime, stopTime)
	}
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	if distributorId != 0 {
		db = db.Where("distributor_id=?", distributorId)
	}
	if merchantId != 0 {
		db = db.Where("merchant_id=?", merchantId)
	}
	db.Count(&ret.TotalCount)
	db.Find(&orderList)
	ret.PageNum = page
	ret.PageSize = size
	ret.PageCount = len(orderList)
	ret.Status = response.StatusSucc
	ret.Data = orderList
	return ret

}

//平台商管理界面：（默认指定平台商distributor-id相关订单）， 按照订单号查询；按照创建时间，订单状态组合搜索条件查询订单列表

func GetOrdersByDistributor(page int, size int, status int, startTime string, stopTime string, sort string, timeField string, distributorId int64, orderNumber string) response.PageResponse {
	var order models.Order
	var orderList []models.Order
	var ret response.PageResponse
	if distributorId == 0 {
		utils.Log.Error("distributorId is null when getOrdersByDistributor")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		return ret
	}
	if orderNumber != "" {
		resp := GetOrderByOrderNumber(orderNumber)
		ret.EntityResponse.CommonRet = resp.CommonRet
		ret.EntityResponse.Data = resp.Data
		return ret
	}
	db := utils.DB.Model(&order).Order(fmt.Sprintf("%s %s", timeField, sort))
	db = db.Offset((page - 1) * size).Limit(size)
	db = db.Where("distributor_id=?", distributorId)
	if startTime != "" && stopTime != "" {
		db = db.Where(fmt.Sprintf("%s >= ? AND %s <= ?", timeField, timeField), startTime, stopTime)
	}
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	db.Count(&ret.TotalCount)
	db.Find(&orderList)
	ret.PageNum = page
	ret.PageSize = size
	ret.PageCount = len(orderList)
	ret.Status = response.StatusSucc
	ret.Data = orderList
	return ret

}

//根据origin_order及distributorId查询订单详情
func GetOrderByOriginOrderAndDistributorId(origin_order string, distributorId int64) response.OrdersRet {
	var ret response.OrdersRet
	var order models.Order
	if origin_order == "" || distributorId == 0 {
		utils.Log.Error("origin_order or distributorId is null")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		return ret
	}
	if err := utils.DB.First(&order, "origin_order=? and distributor_id=?", origin_order, distributorId).Error; err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoOrderFindErr.Data()
		return ret
	}
	ret.Status = response.StatusSucc
	ret.Data = []models.Order{order}
	return ret

}

//承兑商查询订单方法,direction(0:买入，1：卖出，-1：买入和卖出)，in_progress（0：订单完成，1：订单正在进行中，-1：订单完成和正在进行中）
func GetOrdersByMerchant(page int, size int, direction int, in_progress int, merchantId int64) response.PageResponse {
	var ret response.PageResponse
	var orderList []models.Order
	if merchantId == 0 {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		return ret
	}
	// 注：下面left join查询假定：在fulfillment_events表中，列merchant_id和列order_number的组合仅能找到一个记录。
	db := utils.DB.Model(&models.Order{}).Select("orders.*, fulfillment_events.accepted_at as accepted_at, " +
		"fulfillment_events.paid_at as paid_at, fulfillment_events.payment_confirmed_at as payment_confirmed_at, " +
		"fulfillment_events.transferred_at as transferred_at").
		Joins("left join fulfillment_events on orders.merchant_id = fulfillment_events.merchant_id and orders.order_number = fulfillment_events.order_number")

	db = db.Where("orders.merchant_id = ?", merchantId)

	if direction == 0 {
		db = db.Where("orders.direction = ?", direction)
	} else if direction == 1 {
		db = db.Where("orders.direction = ?", direction)
	}
	if in_progress == 0 {
		db = db.Where("orders.status = 7")
	} else if in_progress == 1 {
		db = db.Where("orders.status > 1 && orders.status < 7")
	}
	db.Count(&ret.TotalCount)

	db = db.Order("orders.updated_at desc") // 最近更新的订单放在前面
	db = db.Offset((page - 1) * size).Limit(size)

	db.Find(&orderList)
	ret.PageNum = page
	ret.PageSize = size
	ret.PageCount = int(math.Ceil(float64(ret.TotalCount) / float64(size)))
	ret.Status = response.StatusSucc
	ret.Data = orderList
	return ret

}

//创建订单
func CreateOrder(req response.OrderRequest) response.OrdersRet {

	var ret response.OrdersRet
	order := models.Order{
		OrderNumber: GenerateOrderNumber(),
		Price:       req.Price,
		OriginOrder: req.OriginOrder,
		//成交量
		Quantity: req.Quantity,
		//成交额
		Amount:     req.Amount,
		PaymentRef: req.PaymentRef,
		//订单状态，0/1分别表示：未支付的/已支付的
		Status: 1,
		//订单类型，1为买入，2为卖出
		Direction:         req.Direction,
		DistributorId:     req.DistributorId,
		MerchantId:        req.MerchantId,
		MerchantPaymentId: req.MerchantPaymentId,
		//扣除用户佣金金额
		TraderCommissionAmount: req.TraderCommissionAmount,
		//扣除用户佣金币的量
		TraderCommissionQty: req.TraderCommissionQty,
		//用户佣金比率
		TraderCommissionPercent: req.TraderCommissionPercent,
		//扣除币商佣金金额
		MerchantCommissionAmount: req.MerchantCommissionAmount,
		//扣除币商佣金币的量
		MerchantCommissionQty: req.MerchantCommissionQty,
		//币商佣金比率
		MerchantCommissionPercent: req.MerchantCommissionPercent,
		//平台扣除的佣金币的量（= trader_commision_qty+merchant_commision_qty)
		PlatformCommissionQty: req.PlatformCommissionQty,
		//平台商用户id
		AccountId: req.AccountId,
		//交易币种
		CurrencyCrypto: req.CurrencyCrypto,
		//交易法币
		CurrencyFiat: req.CurrencyFiat,
		//交易类型 0:微信,1:支付宝,2:银行卡
		PayType: req.PayType,
		//微信或支付宝二维码地址
		QrCode: req.QrCode,
		//微信或支付宝账号
		Name: req.Name,
		//银行账号
		BankAccount: req.BankAccount,
		//所属银行
		Bank: req.Bank,
		//所属银行分行
		BankBranch: req.BankBranch,
	}
	if db := utils.DB.Create(&order); db.Error != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateOrderErr.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = append([]models.Order{}, order)
	return ret
}

//修改订单
func UpdateOrder(req response.OrderRequest) response.OrdersRet {
	var ret response.OrdersRet
	var order models.Order
	if error := utils.DB.First(&order, "order_number=?", req.OrderNumber).Error; error != nil {
		utils.Log.Error(error)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoOrderFindErr.Data()
		return ret
	}
	if req.Status != 0 {
		order.Status = req.Status
	}
	if req.MerchantId != 0 {
		order.MerchantId = req.MerchantId
	}
	if req.Price != 0 {
		order.Price = req.Price
	}
	if req.MerchantCommissionPercent != 0 {
		order.MerchantCommissionPercent = req.MerchantCommissionPercent
	}
	if req.MerchantCommissionQty != 0 {
		order.MerchantCommissionQty = req.MerchantCommissionQty
	}
	if req.MerchantCommissionAmount != 0 {
		order.MerchantCommissionAmount = req.MerchantCommissionAmount
	}
	if req.MerchantPaymentId != 0 {
		order.MerchantPaymentId = req.MerchantPaymentId
	}

	if err := utils.DB.Model(&order).Updates(order).Error; err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.UpdateOrderErr.Data()
		return ret
	}
	ret.Status = response.StatusSucc
	ret.Data = append([]models.Order{}, order)
	return ret
}

func GetOrderStatus() response.EntityResponse {
	var ret response.EntityResponse
	ret.Status = response.StatusSucc
	data := make(map[string]models.OrderStatus)
	data["new"] = models.NEW
	data["wait_accept"] = models.WAITACCEPT
	data["accepted"] = models.ACCEPTED
	data["notify_paid"] = models.NOTIFYPAID
	data["confirm_paid"] = models.CONFIRMPAID
	data["suspended"] = models.SUSPENDED
	data["payment_mismatch"] = models.PAYMENTMISMATCH
	data["transferred"] = models.TRANSFERRED
	ret.Data = data
	return ret
}

//使用guid随机生成订单号方法
func GenerateOrderNumber() string {
	var guidId string
	guidId = tsgutils.GUID()
	return guidId

}
