package service

import (
	"fmt"
	"github.com/typa01/go-utils"
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
		db = db.Offset(pageNum * pageSize).Limit(pageSize)
		db.Count(&ret.PageCount)
		ret.PageNum = int(pageNum + 1)
		ret.PageSize = int(pageSize)
		db.Find(&data)
		ret.Data = data
		ret.Status = response.StatusSucc
	}

	return ret

}

func GetOrderByOrderNumber(orderId int64) response.OrdersRet {
	var ret response.OrdersRet
	var data models.Order
	if error := utils.DB.First(&data, "order_number=?", orderId); error != nil {
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
		pageNum, err := strconv.ParseInt(page, 10, 64)
		pageSize, err1 := strconv.ParseInt(size, 10, 64)
		if err != nil || err1 != nil {
			utils.Log.Error(pageNum, pageSize)
		}
		db = db.Offset(pageNum * pageSize).Limit(pageSize)
		if startTime != "" && stopTime != "" {
			db = db.Where(fmt.Sprintf("%s >= ? AND %s <= ?", timeField, timeField), startTime, stopTime)
		}
		if status != "" {
			db = db.Where("status = ?", status)
		}
		db.Count(&ret.PageCount)
		ret.PageNum = int(pageNum + 1)
		ret.PageSize = int(pageSize)
	}
	db.Find(&result)
	ret.Status = response.StatusSucc
	ret.Data = result
	return ret
}

//创建订单
func CreateOrder(req response.OrderRequest) response.OrdersRet {

	var ret response.OrdersRet
	order := models.Order{
		OrderNumber: GenerateOrderNumber(),
		Price:       req.Price,
		OriginOrder:req.OriginOrder,
		//成交量
		Quantity: req.Quantity,
		//成交额
		Amount:     req.Amount,
		PaymentRef: req.PaymentRef,
		//订单状态，0/1分别表示：未支付的/已支付的
		Status: 1,
		//订单类型，1为买入，2为卖出
		OrderType:         req.OrderType,
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
	if error := utils.DB.First(&order, "order_number=?", req.OrderNumber); error != nil {
		utils.Log.Error(error)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoOrderFindErr.Data()
		return ret
	}
	if req.Status!=0{
		order.Status=req.Status
	}
	order.MerchantId=req.MerchantId
	order.MerchantCommissionPercent=req.MerchantCommissionPercent
	order.MerchantCommissionQty=req.MerchantCommissionQty
	order.MerchantCommissionAmount=req.MerchantCommissionAmount
	order.MerchantPaymentId=req.MerchantPaymentId
	order.DistributorId=req.DistributorId

	if err := utils.DB.Model(&order).Updates(order).Error; err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.UpdateOrderErr.Data()
		return ret
	}
	ret.Status = response.StatusSucc
	ret.Data = append([]models.Order{}, order)
	return ret
}

//使用guid随机生成订单号方法
func GenerateOrderNumber() string {
	var guidId string
	guidId=tsgutils.GUID()
	return guidId

}