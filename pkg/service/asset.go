package service

import (
	"fmt"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/utils"
)

func GetAssetHistories(page, size, startTime, stopTime, sort, timeField, search, uid string, isMerchant bool) response.PageResponse {
	var result []models.AssetHistory
	var ret response.PageResponse
	db := utils.DB.Model(&models.AssetHistory{}).Order(fmt.Sprintf("asset_histories.%s %s", timeField, sort)).
		Select("asset_histories.*,users.username as operator_name").
		Joins("left join users on asset_histories.operator_id = users.id")
	if isMerchant {
		db = db.Where("asset_histories.merchant_id = ?", uid)
	} else {
		db = db.Where("asset_histories.distributor_id = ?", uid)
	}
	if search != "" {
		db = db.Where("asset_histories.order_number ?", search)
	} else {
		pageNum, err := strconv.ParseInt(page, 10, 64)
		pageSize, err1 := strconv.ParseInt(size, 10, 64)
		if err != nil || err1 != nil {
			utils.Log.Error(pageNum, pageSize)
		}
		db = db.Offset(pageNum * pageSize).Limit(pageSize)
		if startTime != "" && stopTime != "" {
			db = db.Where(fmt.Sprintf("asset_histories.%s >= ? AND asset_histories.%s <= ?", timeField, timeField), startTime, stopTime)
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

func GetAssetApplies(page, size, status, startTime, stopTime, sort, timeField, search string) response.PageResponse {
	var result []models.AssetApply
	var ret response.PageResponse
	db := utils.DB.Model(&models.AssetApply{}).Order(fmt.Sprintf("asset_applies.%s %s", timeField, sort)).
		Select("asset_applies.*,assets.quantity as remain_quantity,users.username as apply_name").
		Joins("left join assets on assets.merchant_id = asset_applies.merchant_id left join users on asset_applies.apply_id = users.id")
	if search != "" {
		db = db.Where("phone = ? OR email = ?", search, search)
	} else {
		pageNum, err := strconv.ParseInt(page, 10, 64)
		pageSize, err1 := strconv.ParseInt(size, 10, 64)
		if err != nil || err1 != nil {
			utils.Log.Error(pageNum, pageSize)
		}
		db = db.Offset(pageNum * pageSize).Limit(pageSize)
		if startTime != "" && stopTime != "" {
			db = db.Where(fmt.Sprintf("asset_applies.%s >= ? AND asset_applies.%s <= ?", timeField, timeField), startTime, stopTime)
		}
		if status != "" {
			db = db.Where("asset_applies.status = ?", status)
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
