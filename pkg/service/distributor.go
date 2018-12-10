package service

import (
	"fmt"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

func GetDistributors(page, size, status, startTime, stopTime, sort, timeField, search string) response.PageResponse {
	var result []models.Distributor
	var ret response.PageResponse
	if search != "" {
		utils.DB.Where("name = ? OR id = ?", search, search).Find(&result)
	} else {
		db := utils.DB.Model(&models.Distributor{}).Order(fmt.Sprintf("%s %s", timeField, sort))
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
		db.Find(&result)
	}

	ret.Status = "success"
	ret.Data = append(ret.Data, result)
	return ret
}

func CreateDistributor(distributor models.Distributor) response.EntityResponse {
	var ret response.EntityResponse
	if err := utils.DB.Create(&distributor).Error;err != nil {
		ret.Status = "fail"
		ret.ErrCode,ret.ErrMsg = err_code.DistributorErr.Data()
	} else {
		ret.Status = "success"
		ret.Data = append(ret.Data,distributor)
	}

	return ret
}
