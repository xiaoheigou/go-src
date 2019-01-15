package service

import (
	"yuudidi.com/pkg/protocol/response"
)

func ReprocessOrder(origin_order string, distributorId int64) response.OrdersRet {
	//var ret response.CreateOrderRet
	//var result response.CreateOrderResult
	//var ret response.OrdersRet
	orderRet := GetOrderByOriginOrderAndDistributorId(origin_order, distributorId)
	//
	//if orderRet.Status == response.StatusFail || orderRet.Data == nil {
	//	utils.Log.Error("can not find order according origin_order and distributorId")
	//	ret.Status = response.StatusFail
	//	ret.ErrCode, ret.ErrMsg = err_code.NoOrderFindErr.Data()
	//	return ret
	//}
	////orderNumber := orderRet.Data[0].OrderNumber
	//ret.Status = response.StatusSucc
	////result.OrderNumber = orderNumber
	//ret.Data = []models.Order{orderRet.Data[0]}

	//重定向



	return orderRet

}
