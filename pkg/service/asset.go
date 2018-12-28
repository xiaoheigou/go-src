package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
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
		if startTime != "" && stopTime != "" {
			db = db.Where(fmt.Sprintf("asset_histories.%s >= ? AND asset_histories.%s <= ?", timeField, timeField), startTime, stopTime)
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

		if startTime != "" && stopTime != "" {
			db = db.Where(fmt.Sprintf("asset_applies.%s >= ? AND asset_applies.%s <= ?", timeField, timeField), startTime, stopTime)
		}
		if status != "" {
			db = db.Where("asset_applies.status = ?", status)
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
	ret.Status = response.StatusSucc
	ret.Data = result
	return ret
}

//充值申请
func RechargeApply(uid string, params response.RechargeArgs) response.EntityResponse {
	var ret response.EntityResponse
	id, _ := strconv.ParseInt(uid, 10, 64)
	asset := models.AssetApply{
		MerchantId: id,
		Currency:   params.Currency,
		ApplyId:    params.UserId,
		Quantity:   params.Count,
	}
	assetHistory := models.AssetHistory{
		Currency:   params.Currency,
		MerchantId: id,
		Quantity:   params.Count,
		OperatorId: params.UserId,
		IsOrder:    0,
		Operation:  0,
	}
	tx := utils.DB.Begin()
	if err := tx.Model(&models.AssetApply{}).Create(&asset).Error; err != nil {
		utils.Log.Errorf("create asset apply is failed,uid:%s,params:%v", uid, params)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateMerchantRechargeErr.Data()
		tx.Rollback()
		return ret
	}
	if err := tx.Model(&models.AssetHistory{}).Create(&assetHistory).Error; err != nil {
		utils.Log.Errorf("create asset history is failed,uid:%s,params:%v,err:%v", uid, params)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateMerchantRechargeErr.Data()
		tx.Rollback()
		return ret
	}
	tx.Commit()
	ret.Status = response.StatusSucc
	return ret
}

//充值确认
func RechargeConfirm(uid, assetApplyId, userId string) response.EntityResponse {
	var ret response.EntityResponse
	id, _ := strconv.ParseInt(uid, 10, 64)
	//assetId,_ := strconv.ParseInt(assetApplyId,10,64)
	operatorId, _ := strconv.ParseInt(userId, 10, 64)
	var assetApply models.AssetApply
	if err := utils.DB.First(&assetApply, "id = ?", assetApplyId).Error; err != nil {
		utils.Log.Errorf("not fount asset apply,uid:%s,assetApplyId:%s", uid, assetApplyId)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundAssetApplyErr.Data()
		return ret
	}
	assetHistory := models.AssetHistory{
		Currency:   assetApply.Currency,
		MerchantId: id,
		Quantity:   assetApply.Quantity,
		OperatorId: operatorId,
		IsOrder:    0,
		Operation:  1,
	}
	tx := utils.DB.Begin()
	//添加资金变动日志
	if err := tx.Model(&models.AssetHistory{}).Create(&assetHistory).Error; err != nil {
		utils.Log.Errorf("create asset history is failed,uid:%s,params:%v", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateMerchantRechargeErr.Data()
		tx.Rollback()
		return ret
	}
	if err := tx.Model(&models.AssetApply{}).Update("status", 1).Error; err != nil {
		utils.Log.Errorf("create asset apply is failed,uid:%s,params:%v", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateMerchantRechargeErr.Data()
		tx.Rollback()
		return ret
	}
	//更新充值申请为已审核状态
	if err := tx.Model(&models.AssetApply{}).Update("status", 1).Error; err != nil {
		utils.Log.Errorf("update asset apply status is failed,uid:%s", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateMerchantRechargeErr.Data()
		tx.Rollback()
		return ret
	}
	//添加用户的资产
	if err := recharge(uid, assetApply.Currency, assetApply.Quantity, tx); err != nil {
		utils.Log.Errorf("update asset is failed,uid:%s", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateMerchantRechargeErr.Data()
		tx.Rollback()
		return ret
	}
	tx.Commit()
	ret.Status = response.StatusSucc
	return ret
}

func recharge(merchantId, currency string, quantity float64, tx *gorm.DB) error {
	var asset models.Assets

	if err := tx.First(&asset, "merchant_id = ? and currency_crypto = ?", merchantId, currency).Error; err != nil {
		merchantIdInt, _ := strconv.ParseInt(merchantId, 10, 64)
		asset.MerchantId = merchantIdInt
		asset.Quantity = 0
		asset.CurrencyCrypto = currency
		if err := tx.Model(&models.Assets{}).Create(&asset).Error; err != nil {
			return err
		}
	}
	sum := asset.Quantity + quantity
	if err := tx.Model(&models.Assets{}).Update("quantity", sum).Error; err != nil {

		return err
	}

	return nil
}
