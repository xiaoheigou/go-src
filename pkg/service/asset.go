package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
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
	merchant := models.Merchant{}
	if err := tx.First(&merchant, "id = ?", id).Error; err != nil {
		utils.Log.Errorf("get merchant is failed,uid:%s,params:%v", uid, params)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateMerchantRechargeErr.Data()
		tx.Rollback()
		return ret
	}
	asset.Phone = merchant.Phone
	asset.Email = merchant.Email
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
	//已经审核过，不能在审核
	if assetApply.Status == 1 {
		utils.Log.Debugf("asset apply already audited,uid:%s,assetApplyId:%s", uid, assetApplyId)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AssetApplyAlreadyAuditErr.Data()
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
	//更新充值申请为已审核状态
	if err := tx.Model(&models.AssetApply{}).Where("id = ?", assetApplyId).Updates(models.AssetApply{Status: 1, AuditorId: operatorId}).Error; err != nil {
		utils.Log.Errorf("update asset apply status is failed,uid:%s", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateMerchantRechargeErr.Data()
		tx.Rollback()
		return ret
	}
	//添加资金变动日志
	if err := tx.Model(&models.AssetHistory{}).Create(&assetHistory).Error; err != nil {
		utils.Log.Errorf("create asset history is failed,uid:%s,params:%v", uid)
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
			utils.Log.Errorf("create merchant asset is failed.err:[%v]", err)
			return err
		}
	}
	sum := asset.Quantity + quantity
	if err := tx.Model(&models.Assets{}).Where("id = ?", asset.Id).Update("quantity", sum).Error; err != nil {
		return err
	}

	return nil
}

func ReleaseCoin(orderNumber, username string, userId int64) response.EntityResponse {
	var ret response.EntityResponse
	var order models.Order

	tx := utils.DB.Begin()
	//找到订单的记录
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&order, "order_number = ? AND status = ?", orderNumber, models.SUSPENDED).RecordNotFound() {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoOrderFindErr.Data()
		tx.Rollback()
		return ret
	}

	// 找到平台商asset记录
	assetForDist := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForDist, "distributor_id = ? AND currency_crypto = ? ", order.DistributorId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundAssetErr.Data()
		tx.Rollback()
		return ret
	}

	// 找到币商asset记录
	asset := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&asset, "merchant_id = ? AND currency_crypto = ? ", order.MerchantId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundAssetErr.Data()
		tx.Rollback()
		return ret
	}

	// 找到金融滴滴平台asset记录
	assetForPlatform := models.Assets{}
	platformDistId := 1 // 金融滴滴平台的distributor_id为1
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForPlatform, "distributor_id = ? AND currency_crypto = ? ",
		platformDistId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundAssetErr.Data()
		tx.Rollback()
		return ret
	}
	//修改订单状态
	if err := tx.Model(&order).Updates(models.Order{Status: models.RELEASE}).Error; err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
		tx.Rollback()
		return ret
	}
	if order.Direction == 0 {
		//扣除币商冻结的币
		if err := tx.Table("assets").Where("id = ? and qty_frozen >= ?", asset.Id, order.Quantity).Update("qty_frozen", asset.QtyFrozen-order.Quantity).Error; err != nil {
			utils.Log.Errorf("Can't freeze asset for merchant (uid=[%v]). err: %v", asset.MerchantId, err)
			utils.Log.Errorf("func doTransfer finished abnormally.")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			return ret
		}
		//释放币种给平台商
		if err := tx.Table("assets").Where("id = ? ", assetForDist.Id).Update("quantity", assetForDist.Quantity+order.Quantity).Error; err != nil {
			utils.Log.Errorf("Can't transfer asset to distributor (distributor_id=[%v]). err: %v", assetForDist.DistributorId, err)
			utils.Log.Errorf("func doTransfer finished abnormally.")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			return ret
		}

	} else if order.Direction == 1 {
		//客户提现
		// 扣除平台商冻结的币
		if order.Quantity < order.MerchantCommissionQty {
			utils.Log.Errorf("order.Quantity < order.MerchantCommissionQty, invalid order [%s]", order.OrderNumber)
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("func doTransfer finished abnormally.")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			return ret
		}

		// 释放币商冻结的币
		if err := tx.Table("assets").Where("id = ? and qty_frozen >= ?", asset.Id, order.Quantity-order.MerchantCommissionQty).
			Updates(map[string]interface{}{"qty_frozen": asset.QtyFrozen - (order.Quantity - order.MerchantCommissionQty), "quantity": asset.Quantity + (order.Quantity - order.MerchantCommissionQty)}).Error; err != nil {
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("Can't unfrozen [%d] [%s] for merchant (uid=[%v]): %v", order.Quantity-order.MerchantCommissionQty, order.CurrencyCrypto, asset.MerchantId, err)
			utils.Log.Errorf("func doTransfer finished abnormally.")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			return ret
		}

		// 释放金融滴滴平台冻结的币
		if err := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForPlatform.Id).
			Updates(map[string]interface{}{"qty_frozen": assetForPlatform.QtyFrozen - (order.Quantity + order.PlatformCommissionQty), "quantity": assetForPlatform.Quantity + (order.Quantity + order.PlatformCommissionQty)}).Error; err != nil {
			utils.Log.Errorf("Can't unfrozen [%d] [%s] for platform (id=[%v]): %v", order.Quantity+order.PlatformCommissionQty, order.CurrencyCrypto, assetForPlatform.Id, err)
			utils.Log.Errorf("tx in func doTransfer rollback, tx=[%v]", tx)
			utils.Log.Errorf("func doTransfer finished abnormally.")
			tx.Rollback()
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			return ret
		}

	} else {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.OrderDirectionErr.Data()
		tx.Rollback()
		return ret
	}
	assetLog := models.AssetHistory{
		IsOrder:       0,
		Quantity:      order.Quantity,
		MerchantId:    order.MerchantId,
		DistributorId: order.DistributorId,
		Operation:     2,
		OperatorId:    userId,
		OperatorName:  username,
	}
	if err := tx.Create(&assetLog).Error; err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.OrderDirectionErr.Data()
		tx.Rollback()
		return ret
	}
	tx.Commit()
	ret.Status = response.StatusSucc
	return ret
}

func UnFreezeCoin(orderNumber, username string, userId int64) response.EntityResponse {
	var ret response.EntityResponse
	var order models.Order

	//获取订单
	tx := utils.DB.Begin()
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&order, "order_number = ? AND status = ?", orderNumber, models.SUSPENDED).RecordNotFound() {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoOrderFindErr.Data()
		tx.Rollback()
		return ret
	}

	// 找到平台商asset记录
	assetForDist := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForDist, "distributor_id = ? AND currency_crypto = ? ", order.DistributorId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundAssetErr.Data()
		tx.Rollback()
		return ret
	}

	// 找到币商asset记录
	asset := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&asset, "merchant_id = ? AND currency_crypto = ? ", order.MerchantId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundAssetErr.Data()
		tx.Rollback()
		return ret
	}

	// 找到金融滴滴平台asset记录
	assetForPlatform := models.Assets{}
	platformDistId := 1 // 金融滴滴平台的distributor_id为1
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForPlatform, "distributor_id = ? AND currency_crypto = ? ",
		platformDistId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundAssetErr.Data()
		tx.Rollback()
		return ret
	}
	//修改订单状态为解冻状态
	if err := tx.Model(&order).Updates(models.Order{Status: models.UNFREEZE}).Error; err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
		tx.Rollback()
		return ret
	}
	if order.Direction == 0 {

		//解除币商冻结的币
		if err := tx.Table("assets").Where("id = ? and qty_frozen >= ?", asset.Id, order.Quantity).
			Updates(map[string]interface{}{"qty_frozen": asset.QtyFrozen - order.Quantity, "quantity": asset.Quantity + order.Quantity}).Error; err != nil {
			utils.Log.Errorf("Can't freeze asset for merchant (uid=[%v]). err: %v", asset.MerchantId, err)
			utils.Log.Errorf("func UnfreezeCoin finished abnormally.")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			return ret
		}
	} else if order.Direction == 1 {
		// 释放币商冻结的币
		if err := tx.Table("assets").Where("id = ? and qty_frozen >= ?", asset.Id, order.Quantity-order.MerchantCommissionQty).
			Update("qty_frozen", asset.QtyFrozen-(order.Quantity-order.MerchantCommissionQty)).Error; err != nil {
			utils.Log.Errorf("Can't unfrozen [%d] [%s] for merchant (uid=[%v]): %v", order.Quantity-order.MerchantCommissionQty, order.CurrencyCrypto, asset.MerchantId, err)
			utils.Log.Errorf("func doTransfer finished abnormally.")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			return ret
		}

		// 释放币商冻结的币
		if err := tx.Table("assets").Where("id = ?", assetForDist.Id).
			Update("qty_frozen", assetForDist.QtyFrozen+(order.Quantity+order.PlatformCommissionQty)).Error; err != nil {
			utils.Log.Errorf("Can't unfrozen [%d] [%s] for merchant (uid=[%v]): %v", order.Quantity-order.MerchantCommissionQty, order.CurrencyCrypto, asset.MerchantId, err)
			utils.Log.Errorf("func UnFreezeCoin finished abnormally.")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			return ret
		}

		// 释放金融滴滴平台冻结的币
		if err := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForPlatform.Id).
			Update("qty_frozen", assetForPlatform.QtyFrozen-order.PlatformCommissionQty).Error; err != nil {
			utils.Log.Errorf("Can't unfrozen [%d] [%s] for platform (id=[%v]): %v", order.Quantity+order.PlatformCommissionQty, order.CurrencyCrypto, assetForPlatform.Id, err)
			utils.Log.Errorf("func UnFreezeCoin finished abnormally.")
			tx.Rollback()
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			return ret
		}
	} else {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.OrderDirectionErr.Data()
		tx.Rollback()
		return ret
	}
	//资金变动历史添加
	assetLog := models.AssetHistory{
		IsOrder:       0,
		Quantity:      order.Quantity,
		MerchantId:    order.MerchantId,
		DistributorId: order.DistributorId,
		Operation:     3,
		OperatorId:    userId,
		OperatorName:  username,
	}
	if err := tx.Create(&assetLog).Error; err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.OrderDirectionErr.Data()
		tx.Rollback()
		return ret
	}
	tx.Commit()
	ret.Status = response.StatusSucc
	return ret
}
