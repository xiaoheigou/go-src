package service

import (
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/utils"
)

func ReprocessOrder(origin_order string, distributorId int64) string {
	orderRet := GetOrderByOriginOrderAndDistributorId(origin_order, distributorId)
	if orderRet.Status==response.StatusFail||&orderRet.Data==nil{
		utils.Log.Error("can not find order according origin_order and distributorId")
		return ""
	}
	orderNumber := orderRet.Data[0].OrderNumber
	//url=utils.Config.GetString("redirecturl.reprocessurl")+orderNumber
	return orderNumber

}