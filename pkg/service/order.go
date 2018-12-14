package service

import (
	"fmt"
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
