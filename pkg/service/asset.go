package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
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

	dec, err := decimal.NewFromString(params.Count)
	if err != nil {
		utils.Log.Errorf("invalid param: count, err %s", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		return ret
	}

	asset := models.AssetApply{
		MerchantId: id,
		Currency:   params.Currency,
		ApplyId:    params.UserId,
		Quantity:   dec,
	}
	assetHistory := models.AssetHistory{
		Currency:   params.Currency,
		MerchantId: id,
		Quantity:   dec,
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

func recharge(merchantId, currency string, quantity decimal.Decimal, tx *gorm.DB) error {
	var asset models.Assets

	if err := tx.First(&asset, "merchant_id = ? and currency_crypto = ?", merchantId, currency).Error; err != nil {
		merchantIdInt, _ := strconv.ParseInt(merchantId, 10, 64)
		asset.MerchantId = merchantIdInt
		asset.Quantity = decimal.Zero
		asset.CurrencyCrypto = currency
		if err := tx.Model(&models.Assets{}).Create(&asset).Error; err != nil {
			utils.Log.Errorf("create merchant asset is failed.err:[%v]", err)
			return err
		}
	}
	sum := asset.Quantity.Add(quantity)
	if err := tx.Model(&models.Assets{}).Where("id = ?", asset.Id).Update("quantity", sum).Error; err != nil {
		return err
	}

	return nil
}

// 和订单原始预期一致
func ReleaseCoin(orderNumber, username string, userId int64) response.EntityResponse {
	var ret response.EntityResponse
	var order models.Order

	tx := utils.DB.Begin()
	//找到订单的记录
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&order, "order_number = ?", orderNumber).RecordNotFound() {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoOrderFindErr.Data()
		tx.Rollback()
		return ret
	}

	// 判断订单当前状态，决定是否允许放币操作
	if order.Direction == 0 {
		if order.Status == models.SUSPENDED {
			if order.StatusReason == 19 || order.StatusReason == 20 {
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
				tx.Rollback()
				utils.Log.Debugf("func ReleaseCoin finished abnormally, order_number = %s", orderNumber)
				return ret
			} else {
				// pass
			}
		} else { // 状态不是SUSPENDED
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			utils.Log.Debugf("func ReleaseCoin finished abnormally, order_number = %s", orderNumber)
			return ret
		}
	} else if order.Direction == 1 {
		if order.Status == models.SUSPENDED {
			if order.StatusReason == 19 || order.StatusReason == 20 {
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
				tx.Rollback()
				utils.Log.Debugf("func ReleaseCoin finished abnormally, order_number = %s", orderNumber)
				return ret
			} else {
				// pass
			}
		} else if order.Status == models.CONFIRMPAID { // status 4
			// pass
		} else {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			utils.Log.Debugf("func ReleaseCoin finished abnormally, order_number = %s", orderNumber)
			return ret
		}
	} else {
		utils.Log.Errorf("direction [%d] is invalid for order %s", order.Direction, order.OrderNumber)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
		tx.Rollback()
		utils.Log.Debugf("func ReleaseCoin finished abnormally, order_number = %s", orderNumber)
		return ret
	}

	// 找到平台商asset记录
	assetForDist := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForDist, "distributor_id = ? AND currency_crypto = ? ", order.DistributorId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundAssetErr.Data()
		return ret
	}

	// 找到币商asset记录
	asset := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&asset, "merchant_id = ? AND currency_crypto = ? ", order.MerchantId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundAssetErr.Data()
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
		return ret
	}

	if order.Direction == 0 {
		//扣除币商冻结的币
		if asset.QtyFrozen.GreaterThanOrEqual(order.Quantity) { // 避免merchant的qty_frozen列扣成负数
			if err := tx.Table("assets").Where("id = ?", asset.Id).Update("qty_frozen", asset.QtyFrozen.Sub(order.Quantity)).Error; err != nil {
				utils.Log.Errorf("update asset for merchant fail. err %s", err)
				utils.Log.Errorf("func ReleaseCoin finished abnormally.")
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
				tx.Rollback()
				return ret
			}
		} else {
			utils.Log.Errorf("Can't deduct %f %s for merchant (uid=[%v]), the qty_frozen is not enough (%f). asset for merchant = [%+v], order_number = %s",
				order.Quantity, order.CurrencyCrypto, asset.MerchantId, asset.QtyFrozen, asset, order.OrderNumber)
			utils.Log.Errorf("func ReleaseCoin finished abnormally.")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			return ret
		}
		//释放币种给平台商
		if err := tx.Table("assets").Where("id = ? ", assetForDist.Id).Update("quantity", assetForDist.Quantity.Add(order.Quantity)).Error; err != nil {
			utils.Log.Errorf("Can't transfer asset to distributor (distributor_id=[%v]). err: %v", assetForDist.DistributorId, err)
			utils.Log.Errorf("func ReleaseCoin finished abnormally.")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			return ret
		}

		assetLog := models.AssetHistory{
			IsOrder:       1,
			OrderNumber:   order.OrderNumber,
			Quantity:      order.Quantity,
			DistributorId: order.DistributorId,
			Operation:     2, // 放币
			OperatorId:    userId,
			OperatorName:  username,
		}
		if err := tx.Create(&assetLog).Error; err != nil {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			return ret
		}

		assetMerchantLog := models.AssetHistory{
			IsOrder:      1,
			OrderNumber:  order.OrderNumber,
			Quantity:     order.Quantity.Neg(),
			MerchantId:   order.MerchantId,
			Operation:    2, // 放币
			OperatorId:   userId,
			OperatorName: username,
		}
		if err := tx.Create(&assetMerchantLog).Error; err != nil {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			return ret
		}

	} else if order.Direction == 1 {
		//客户提现
		if err := TransferNormally(tx, &assetForDist, &asset, &assetForPlatform, &order, &AssetHistoryOperationInfo{
			Operation:    2,
			OperatorId:   userId,
			OperatorName: username,
		}); err != nil {
			utils.Log.Errorf("func TransferNormally fail, err: %s", err)
			utils.Log.Errorf("func ReleaseCoin finished abnormally.")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.ReleaseCoinErr.Data()
			tx.Rollback()
			return ret
		}
	} else {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.OrderDirectionErr.Data()
		tx.Rollback()
		return ret
	}

	//修改订单状态
	if err := tx.Model(&order).Where("order_number = ?", orderNumber).Updates(
		models.Order{
			Status:       models.SUSPENDED,
			StatusReason: models.MARKCOMPLETED}).Error; err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.UpdateOrderErr.Data()
		tx.Rollback()
		return ret
	}

	tx.Commit()

	utils.Log.Infof("call AsynchronousNotifyDistributor for %s, order status is 5.19 (MARKCOMPLETED)", order.OrderNumber)
	AsynchronousNotifyDistributor(order)

	// 删除这个订单在时间轮中的记录，避免这个订单被时间轮再次处理
	engine := NewOrderFulfillmentEngine(nil)
	engine.DeleteWheel(orderNumber)

	ret.Status = response.StatusSucc
	return ret
}

// 和订单原始预期不一致
func UnFreezeCoin(orderNumber, username string, userId int64) response.EntityResponse {
	utils.Log.Debugf("func UnFreezeCoin begin, order_number = %s", orderNumber)
	var ret response.EntityResponse
	var order models.Order

	//获取订单
	tx := utils.DB.Begin()
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&order, "order_number = ?", orderNumber).RecordNotFound() {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NoOrderFindErr.Data()
		tx.Rollback()
		utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
		return ret
	}

	// 判断订单当前状态，决定是否允许解冻操作
	if order.Direction == 0 {
		// 用户充值订单，目前仅当status为5，且status_reason不为19/20时，才可以解冻。
		if order.Status == models.SUSPENDED {
			if order.StatusReason == 19 || order.StatusReason == 20 {
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.UnFreezeCoinErr.Data()
				tx.Rollback()
				utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
				return ret
			} else {
				// pass
			}
		} else { // 状态不是SUSPENDED
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.UnFreezeCoinErr.Data()
			tx.Rollback()
			utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
			return ret
		}
	} else if order.Direction == 1 {
		if order.Status == models.SUSPENDED {
			if order.StatusReason == 19 || order.StatusReason == 20 {
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.UnFreezeCoinErr.Data()
				tx.Rollback()
				utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
				return ret
			} else {
				// pass
			}
		} else if order.Status == models.CONFIRMPAID { // status 4
			// pass
		} else if order.Status == models.ACCEPTTIMEOUT { // status 8
			// pass
		} else {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.UnFreezeCoinErr.Data()
			tx.Rollback()
			utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
			return ret
		}
	} else {
		utils.Log.Errorf("direction [%d] is invalid for order %s", order.Direction, order.OrderNumber)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.UnFreezeCoinErr.Data()
		tx.Rollback()
		utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
		return ret
	}

	// 找到平台商asset记录
	assetForDist := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForDist, "distributor_id = ? AND currency_crypto = ? ", order.DistributorId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundAssetErr.Data()
		utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
		return ret
	}

	// 找到币商asset记录
	asset := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&asset, "merchant_id = ? AND currency_crypto = ? ", order.MerchantId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundAssetErr.Data()
		utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
		return ret
	}

	// 找到金融滴滴平台asset记录
	assetForPlatform := models.Assets{}
	platformDistId := 1 // 金融滴滴平台的distributor_id为1
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForPlatform, "distributor_id = ? AND currency_crypto = ? ",
		platformDistId, order.CurrencyCrypto).RecordNotFound() {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundAssetErr.Data()
		tx.Rollback()
		utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
		return ret
	}

	if order.Direction == 0 {

		//解除币商冻结的币
		if asset.QtyFrozen.GreaterThanOrEqual(order.Quantity) { // 避免qty_frozen扣为负数
			if err := tx.Table("assets").Where("id = ?", asset.Id).
				Updates(map[string]interface{}{
					"qty_frozen": asset.QtyFrozen.Sub(order.Quantity),
					"quantity":   asset.Quantity.Add(order.Quantity)}).Error; err != nil {
				utils.Log.Errorf("update asset for merchant fail")
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.UnFreezeCoinErr.Data()
				tx.Rollback()
				utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
				return ret
			}
		} else {
			utils.Log.Errorf("Can't unfreeze asset for merchant (uid=[%v]). Donot have enough frozen btusd, asset for merchant = [%+v], order_number = %s",
				asset.MerchantId, asset, order.OrderNumber)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.UnFreezeCoinErr.Data()
			tx.Rollback()
			utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
			return ret
		}
	} else if order.Direction == 1 {
		// 平台用户提现订单，币商抢了单，却未付款的情况
		if err := TransferAbnormally(tx, &assetForDist, &asset, &assetForPlatform, &order); err != nil {
			utils.Log.Errorf("func TransferAbnormally err %v", err)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.UnFreezeCoinErr.Data()
			tx.Rollback()
			utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
			return ret
		}
	} else {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.OrderDirectionErr.Data()
		tx.Rollback()
		utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
		return ret
	}

	//修改订单原因状态为订单已取消状态
	if err := tx.Model(&order).Where("order_number = ?", orderNumber).Updates(
		models.Order{
			Status:       models.SUSPENDED,
			StatusReason: models.CANCEL}).Error; err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.UpdateOrderErr.Data()
		tx.Rollback()
		utils.Log.Debugf("func UnFreezeCoin finished abnormally, order_number = %s", orderNumber)
		return ret
	}

	tx.Commit()

	utils.Log.Infof("call AsynchronousNotifyDistributor for %s, order status is 5.20 (CANCEL)", order.OrderNumber)
	AsynchronousNotifyDistributor(order)

	// 删除这个订单在时间轮中的记录，避免这个订单被时间轮再次处理
	engine := NewOrderFulfillmentEngine(nil)
	engine.DeleteWheel(orderNumber)

	ret.Status = response.StatusSucc
	utils.Log.Debugf("func UnFreezeCoin finished normally, order_number = %s", orderNumber)
	return ret
}
