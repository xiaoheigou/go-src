package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

// 用户下提现订单时，冻结了平台的BTUSD
// 现在币商抢单成功了，把之前冻结的平台的BTUSD分到三家（平台自己、币商、jrdidi平台）的冻结账号（qty_frozen列）中
// 下面函数不会commit，也不会rollback，请在上层函数处理
func TransferFrozen(tx *gorm.DB, assetForTrader *models.Assets, assetForMerchant *models.Assets, assetForJrdidi *models.Assets, order *models.Order) error {
	utils.Log.Debugf("func TransferFrozen begin, order_number = %s", order.OrderNumber)

	// 用户充值订单，不收手续费，这个方法未对充值订单进行测试，不要调用它
	if order.Direction == 0 {
		utils.Log.Errorf("func TransferFrozen finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("not applicable for order with direction == 0")
	}

	var deductBTUSD float64 = order.Quantity - order.TraderBTUSDFeeIncome // 当平台自己也抽取用户的提现手续费时，TraderBTUSDFeeIncome大于0，否则TraderBTUSDFeeIncome <= 0
	// 减少平台冻结的BTUSD
	// 注：平台自己赚取的那部分用户手续费还在冻结着
	if rowsAffected := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForTrader.Id, deductBTUSD).
		Updates(map[string]interface{}{
			"qty_frozen": assetForTrader.QtyFrozen - deductBTUSD}).RowsAffected; rowsAffected == 0 {
		utils.Log.Errorf("the qty_frozen is not enough for distributor, assetForTrader = %+v", assetForTrader)
		utils.Log.Errorf("func TransferFrozen finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("the qty_frozen is not enough for distributor")
	}

	// 增加币商冻结的BTUSD
	if rowsAffected := tx.Table("assets").Where("id = ?", assetForMerchant.Id).
		Updates(map[string]interface{}{
			"qty_frozen": assetForMerchant.QtyFrozen + (order.Quantity - order.TraderBTUSDFeeIncome - order.JrdidiBTUSDFeeIncome)}).RowsAffected; rowsAffected == 0 {
		utils.Log.Errorf("func TransferFrozen finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("can not find merchant asset")
	}
	// 增加jrdidi平台冻结的BTUSD
	if rowsAffected := tx.Table("assets").Where("id = ?", assetForJrdidi.Id).
		Updates(map[string]interface{}{
			"qty_frozen": assetForJrdidi.QtyFrozen + order.JrdidiBTUSDFeeIncome}).RowsAffected; rowsAffected == 0 {
		utils.Log.Errorf("func TransferFrozen finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("can not find jrdidi asset")
	}

	utils.Log.Debugf("func TransferFrozen finished normally, order_number = %s", order.OrderNumber)
	return nil
}

type AssetHistoryOperationInfo struct {
	Operation    int
	OperatorId   int64
	OperatorName string
}

// 下面函数不会commit，也不会rollback，请在上层函数处理
func TransferNormally(tx *gorm.DB, assetForTrader *models.Assets, assetForMerchant *models.Assets, assetForJrdidi *models.Assets, order *models.Order,
	opInfo *AssetHistoryOperationInfo) error {
	utils.Log.Debugf("func TransferNormally begin, order_number = %s", order.OrderNumber)

	// 用户充值订单，不收手续费，这个方法未对充值订单进行测试，不要调用它
	if order.Direction == 0 {
		return errors.New("not applicable for order with direction == 0")
	}

	// 用户提现订单完成
	// 把各自冻结的币释放掉

	if order.TraderBTUSDFeeIncome > 0 {
		// 如果平台自己也赚取用户手续费
		// 把平台赚取的手续费（之前处于冻结状态）释放掉
		if rowsAffected := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForTrader.Id, order.TraderBTUSDFeeIncome).
			Updates(map[string]interface{}{
				"qty_frozen": assetForTrader.QtyFrozen - order.TraderBTUSDFeeIncome,
				"quantity":   assetForTrader.Quantity + order.TraderBTUSDFeeIncome}).RowsAffected; rowsAffected == 0 {
			utils.Log.Errorf("the qty_frozen is not enough for distributor, assetForTrader = %+v", assetForTrader)
			return errors.New("the qty_frozen is not enough for distributor")
		}
	} else {
		// 平台不赚手续费或者补贴用户手续费，平台没有赚取的BTUSD需要释放掉
	}
	// 把币商获得的BTUSD（包含他赚的手续费）释放掉
	if rowsAffected := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForMerchant.Id, order.Quantity-order.TraderBTUSDFeeIncome-order.JrdidiBTUSDFeeIncome).
		Updates(map[string]interface{}{
			"qty_frozen": assetForMerchant.QtyFrozen - (order.Quantity - order.TraderBTUSDFeeIncome - order.JrdidiBTUSDFeeIncome),
			"quantity":   assetForMerchant.Quantity + (order.Quantity - order.TraderBTUSDFeeIncome - order.JrdidiBTUSDFeeIncome)}).RowsAffected; rowsAffected == 0 {
		utils.Log.Errorf("the qty_frozen is not enough for merchant, assetForMerchant = %+v", assetForMerchant)
		return errors.New("the qty_frozen is not enough for merchant")
	}
	// 把jrdidi赚取的BTUSD释放掉
	if rowsAffected := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForJrdidi.Id, order.JrdidiBTUSDFeeIncome).
		Updates(map[string]interface{}{
			"qty_frozen": assetForJrdidi.QtyFrozen - order.JrdidiBTUSDFeeIncome,
			"quantity":   assetForJrdidi.Quantity + order.JrdidiBTUSDFeeIncome}).RowsAffected; rowsAffected == 0 {
		utils.Log.Errorf("the qty_frozen is not enough for jrdidi, assetForJrdidi = %+v", assetForJrdidi)
		return errors.New("the qty_frozen is not enough for jrdidi")
	}

	if opInfo == nil {

		var changesForDist float64 = order.Quantity
		if order.TraderBTUSDFeeIncome > 0 { // 没必要分条件，可以统一形式。区分条件仅仅是为了便于理解
			// 平台也赚取用户提现手续费
			changesForDist = order.Quantity - order.TraderBTUSDFeeIncome
		} else if order.TraderBTUSDFeeIncome == 0 {
			// 平台不赚取用户提现手续费
			changesForDist = order.Quantity
		} else {
			// 平台补贴用户提现手续费
			changesForDist = order.Quantity - order.TraderBTUSDFeeIncome
		}
		// Add asset history for distributor
		assetDistHistory := models.AssetHistory{
			Currency:      order.CurrencyCrypto,
			Direction:     order.Direction,
			DistributorId: order.DistributorId,
			Quantity:      -changesForDist,
			IsOrder:       1,
			OrderNumber:   order.OrderNumber,
		}
		if err := tx.Model(&models.AssetHistory{}).Create(&assetDistHistory).Error; err != nil {
			return errors.New("add asset history for distributor fail")
		}

		// Add asset history for merchant
		assetMerchantHistory := models.AssetHistory{
			Currency:    order.CurrencyCrypto,
			Direction:   order.Direction,
			MerchantId:  order.MerchantId,
			Quantity:    order.Quantity - order.TraderBTUSDFeeIncome - order.JrdidiBTUSDFeeIncome,
			IsOrder:     1,
			OrderNumber: order.OrderNumber,
		}
		if err := tx.Model(&models.AssetHistory{}).Create(&assetMerchantHistory).Error; err != nil {
			return errors.New("add asset history for merchant fail")
		}

		// Add asset history for jrdidi
		assetJididiHistory := models.AssetHistory{
			Currency:      order.CurrencyCrypto,
			Direction:     order.Direction,
			DistributorId: 1, // DistributorId为1时表示jrdidi
			Quantity:      order.JrdidiBTUSDFeeIncome,
			IsOrder:       1,
			OrderNumber:   order.OrderNumber,
		}
		if err := tx.Model(&models.AssetHistory{}).Create(&assetJididiHistory).Error; err != nil {
			return errors.New("add asset history for jrdidi fail")
		}

	} else {

		var changesForDist float64 = order.Quantity
		if order.TraderBTUSDFeeIncome > 0 { // 没必要分条件，可以统一形式。区分条件仅仅是为了便于理解
			// 平台也赚取用户提现手续费
			changesForDist = order.Quantity - order.TraderBTUSDFeeIncome
		} else if order.TraderBTUSDFeeIncome == 0 {
			// 平台不赚取用户提现手续费
			changesForDist = order.Quantity
		} else {
			// 平台补贴用户提现手续费
			changesForDist = order.Quantity - order.TraderBTUSDFeeIncome
		}
		// Add asset history for distributor
		assetDistHistory := models.AssetHistory{
			Currency:      order.CurrencyCrypto,
			Direction:     order.Direction,
			DistributorId: order.DistributorId,
			Quantity:      -changesForDist,
			IsOrder:       1,
			OrderNumber:   order.OrderNumber,
			Operation:     opInfo.Operation,
			OperatorId:    opInfo.OperatorId,
			OperatorName:  opInfo.OperatorName,
		}
		if err := tx.Model(&models.AssetHistory{}).Create(&assetDistHistory).Error; err != nil {
			return errors.New("add asset history for distributor fail")
		}

		// Add asset history for merchant
		assetMerchantHistory := models.AssetHistory{
			Currency:     order.CurrencyCrypto,
			Direction:    order.Direction,
			MerchantId:   order.MerchantId,
			Quantity:     order.Quantity - order.TraderBTUSDFeeIncome - order.JrdidiBTUSDFeeIncome,
			IsOrder:      1,
			OrderNumber:  order.OrderNumber,
			Operation:    opInfo.Operation,
			OperatorId:   opInfo.OperatorId,
			OperatorName: opInfo.OperatorName,
		}
		if err := tx.Model(&models.AssetHistory{}).Create(&assetMerchantHistory).Error; err != nil {
			return errors.New("add asset history for merchant fail")
		}

		// Add asset history for jrdidi
		assetJididiHistory := models.AssetHistory{
			Currency:      order.CurrencyCrypto,
			Direction:     order.Direction,
			DistributorId: 1, // DistributorId为1时表示jrdidi
			Quantity:      order.JrdidiBTUSDFeeIncome,
			IsOrder:       1,
			OrderNumber:   order.OrderNumber,
			Operation:     opInfo.Operation,
			OperatorId:    opInfo.OperatorId,
			OperatorName:  opInfo.OperatorName,
		}
		if err := tx.Model(&models.AssetHistory{}).Create(&assetJididiHistory).Error; err != nil {
			return errors.New("add asset history for jrdidi fail")
		}

	}

	utils.Log.Debugf("func TransferNormally finished normally, order_number = %s", order.OrderNumber)
	return nil
}

// 下面函数不会commit，也不会rollback，请在上层函数处理
func TransferAbnormally(tx *gorm.DB, assetForTrader *models.Assets, assetForMerchant *models.Assets, assetForJrdidi *models.Assets, order *models.Order) error {
	utils.Log.Debugf("func TransferAbnormally begin, order_number = %s", order.OrderNumber)

	// 用户充值订单，不收手续费，这个方法未对充值订单进行测试，不要调用它
	if order.Direction == 0 {
		utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("not applicable for order with direction == 0")
	}
	// 用户提现订单未完成，币商未真正付款
	// 把各自冻结的币退到平台的冻结账号中

	// 减少币商获得的BTUSD（包含他赚的手续费）
	if rowsAffected := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForMerchant.Id, order.Quantity-order.TraderBTUSDFeeIncome-order.JrdidiBTUSDFeeIncome).
		Updates(map[string]interface{}{
			"qty_frozen": assetForMerchant.QtyFrozen - (order.Quantity - order.TraderBTUSDFeeIncome - order.JrdidiBTUSDFeeIncome)}).RowsAffected; rowsAffected == 0 {
		utils.Log.Errorf("the qty_frozen is not enough for merchant, assetForMerchant = %+v", assetForMerchant)
		utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("the qty_frozen is not enough for merchant")
	}
	// 减少jrdidi赚取的BTUSD
	if rowsAffected := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForJrdidi.Id, order.JrdidiBTUSDFeeIncome).
		Updates(map[string]interface{}{
			"qty_frozen": assetForJrdidi.QtyFrozen - order.JrdidiBTUSDFeeIncome}).RowsAffected; rowsAffected == 0 {
		utils.Log.Errorf("the qty_frozen is not enough for jrdidi, assetForJrdidi = %+v", assetForJrdidi)
		utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("the qty_frozen is not enough for jrdidi")
	}
	// 把上面两部分减少的BTUSD全部加到平台的冻结账号中
	if rowsAffected := tx.Table("assets").Where("id = ?", assetForTrader.Id).
		Updates(map[string]interface{}{
			"qty_frozen": assetForTrader.QtyFrozen + order.Quantity - order.TraderBTUSDFeeIncome}).RowsAffected; rowsAffected == 0 {
		utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("can not find distributor asset")
	}

	utils.Log.Debugf("func TransferAbnormally finished normally, order_number = %s", order.OrderNumber)
	return nil
}
