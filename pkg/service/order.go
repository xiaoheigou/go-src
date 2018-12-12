package service

import (
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
	}else{
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
