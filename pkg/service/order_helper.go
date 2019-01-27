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
func TransferCoinFromTraderFrozenToMerchantFrozen(tx *gorm.DB, assetForTrader *models.Assets, assetForMerchant *models.Assets, assetForJrdidi *models.Assets, order *models.Order) error {
	utils.Log.Debugf("func TransferCoinFromTraderFrozenToMerchantFrozen begin, order_number = %s", order.OrderNumber)

	// 用户充值订单，不收手续费，这个方法未对充值订单进行测试，不要调用它
	if order.Direction == 0 {
		utils.Log.Errorf("func TransferCoinFromTraderFrozenToMerchantFrozen finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("not applicable for order with direction == 0")
	}

	// TODO 增加注释
	var deductBTUSD float64 = order.Quantity - order.TraderBTUSDFeeIncome // 当平台自己也抽取用户的提现手续费时，TraderBTUSDFeeIncome大于0，否则TraderBTUSDFeeIncome <= 0
	// 减少平台冻结的BTUSD
	// 注：平台自己赚取的那部分用户手续费还在冻结着
	if utils.BtusdCompareGte(assetForTrader.QtyFrozen, deductBTUSD) { // 避免trader的qty_frozen出现负数
		if err := tx.Table("assets").Where("id = ?", assetForTrader.Id).
			Updates(map[string]interface{}{
				"qty_frozen": assetForTrader.QtyFrozen - deductBTUSD}).Error; err != nil {
			utils.Log.Errorf("can not update trader asset, err = %+v", err)
			utils.Log.Errorf("func TransferCoinFromTraderFrozenToMerchantFrozen finished abnormally, order_number = %s", order.OrderNumber)
			return errors.New("can not update trader asset")
		}
	} else {
		utils.Log.Errorf("the qty_frozen is not enough for distributor, assetForTrader = %+v", assetForTrader)
		utils.Log.Errorf("func TransferCoinFromTraderFrozenToMerchantFrozen finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("the qty_frozen is not enough for distributor")
	}

	// 增加币商冻结的BTUSD
	if err := tx.Table("assets").Where("id = ?", assetForMerchant.Id).
		Updates(map[string]interface{}{
			"qty_frozen": assetForMerchant.QtyFrozen + (order.Quantity - order.TraderBTUSDFeeIncome - order.JrdidiBTUSDFeeIncome)}).Error; err != nil {
		utils.Log.Errorf("func TransferCoinFromTraderFrozenToMerchantFrozen finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("can not update merchant asset")
	}
	// 增加jrdidi平台冻结的BTUSD
	if err := tx.Table("assets").Where("id = ?", assetForJrdidi.Id).
		Updates(map[string]interface{}{
			"qty_frozen": assetForJrdidi.QtyFrozen + order.JrdidiBTUSDFeeIncome}).Error; err != nil {
		utils.Log.Errorf("func TransferCoinFromTraderFrozenToMerchantFrozen finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("can not update jrdidi asset")
	}
	//
	if err := tx.Model(&models.Order{}).Where("id = ?", order.Id).Updates(models.Order{BTUSDFlowStatus: models.BTUSDFlowD1TraderFrozenToMerchantFrozen}).Error; err != nil {
		utils.Log.Errorf("Can't update order %s BTUSDFlowStatus to %s. error: %v", order.OrderNumber, models.BTUSDFlowD1TraderFrozenToMerchantFrozen, err)
		return errors.New("update BTUSDFlowStatus fail")
	}

	utils.Log.Debugf("func TransferCoinFromTraderFrozenToMerchantFrozen finished normally, order_number = %s", order.OrderNumber)
	return nil
}

type AssetHistoryOperationInfo struct {
	Operation    int
	OperatorId   int64
	OperatorName string
}

// 下面函数不会commit，也不会rollback，请在上层函数处理
func TransferNormally(tx *gorm.DB, assetForTrader *models.Assets, assetForMerchant *models.Assets, assetForJrdidi *models.Assets, order *models.Order, opInfo *AssetHistoryOperationInfo) error {
	utils.Log.Debugf("func TransferNormally begin, order_number = %s", order.OrderNumber)

	// 用户充值订单，不收手续费，这个方法未对充值订单进行测试，不要调用它
	if order.Direction == 0 {
		utils.Log.Errorf("func TransferNormally finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("not applicable for order with direction == 0")
	}

	// 用户提现订单完成
	// 把各自冻结的币释放掉

	// 分两种情况。
	// 情况一、币已经到币商的冻结账号中。
	// 情况二，币还在平台商的冻结账号中。
	// 通过订单的BTUSDFlowStatus字段来判断当前订单属于哪种情况。

	var alreadyAddCoinToMerchantFrozenAccount = false
	if order.BTUSDFlowStatus == models.BTUSDFlowD1TraderFrozenToMerchantFrozen {
		alreadyAddCoinToMerchantFrozenAccount = true
	}

	if alreadyAddCoinToMerchantFrozenAccount { // 情况一：币已经到币商的冻结账号中

		if utils.BtusdCompareGte(order.TraderBTUSDFeeIncome, 0) { // 如果平台自己也赚取用户手续费
			// 把平台赚取的手续费（之前处于冻结状态）释放掉
			if utils.BtusdCompareGte(assetForTrader.QtyFrozen, order.TraderBTUSDFeeIncome) { // 检查冻结的够不够，避免trader的qty_frozen出现负数
				if err := tx.Table("assets").Where("id = ?", assetForTrader.Id).
					Updates(map[string]interface{}{
						"qty_frozen": assetForTrader.QtyFrozen - order.TraderBTUSDFeeIncome,
						"quantity":   assetForTrader.Quantity + order.TraderBTUSDFeeIncome}).Error; err != nil {
					utils.Log.Errorf("update asset for trader fail, err = %+v", err)
					utils.Log.Errorf("func TransferNormally finished abnormally, order_number = %s", order.OrderNumber)
					return errors.New("update asset for trader fail")
				}
			} else {
				utils.Log.Errorf("the qty_frozen is not enough for distributor, assetForTrader = %+v", assetForTrader)
				utils.Log.Errorf("func TransferNormally finished abnormally, order_number = %s", order.OrderNumber)
				return errors.New("the qty_frozen is not enough for distributor")
			}
		} else {
			// 平台不赚手续费或者补贴用户手续费，平台没有赚取的BTUSD需要释放掉
		}

		// 把币商获得的BTUSD（包含他赚的手续费）释放掉
		// 避免merchant的qty_frozen出现负数，先检测冻的币够不够
		if utils.BtusdCompareGte(assetForMerchant.QtyFrozen, order.Quantity-order.TraderBTUSDFeeIncome-order.JrdidiBTUSDFeeIncome) {
			if err := tx.Table("assets").Where("id = ?", assetForMerchant.Id).
				Updates(map[string]interface{}{
					"qty_frozen": assetForMerchant.QtyFrozen - (order.Quantity - order.TraderBTUSDFeeIncome - order.JrdidiBTUSDFeeIncome),
					"quantity":   assetForMerchant.Quantity + (order.Quantity - order.TraderBTUSDFeeIncome - order.JrdidiBTUSDFeeIncome)}).Error; err != nil {
				utils.Log.Errorf("update asset for merchant fail, err %v assetForMerchant = %+v", err, assetForMerchant)
				utils.Log.Errorf("func TransferNormally finished abnormally, order_number = %s", order.OrderNumber)
				return errors.New("update asset for merchant fail")
			}
		} else {
			utils.Log.Errorf("the qty_frozen is not enough for merchant, assetForMerchant = %+v", assetForMerchant)
			utils.Log.Errorf("func TransferNormally finished abnormally, order_number = %s", order.OrderNumber)
			return errors.New("the qty_frozen is not enough for merchant")
		}

		// 把jrdidi赚取的BTUSD释放掉
		// 避免jrdidi的qty_frozen出现负数，先检测冻的币够不够
		if utils.BtusdCompareGte(assetForJrdidi.QtyFrozen, order.JrdidiBTUSDFeeIncome) {
			if err := tx.Table("assets").Where("id = ?", assetForJrdidi.Id).
				Updates(map[string]interface{}{
					"qty_frozen": assetForJrdidi.QtyFrozen - order.JrdidiBTUSDFeeIncome,
					"quantity":   assetForJrdidi.Quantity + order.JrdidiBTUSDFeeIncome}).Error; err != nil {
				utils.Log.Errorf("update asset for jrdidi fail, assetForJrdidi = %+v", assetForJrdidi)
				utils.Log.Errorf("func TransferNormally finished abnormally, order_number = %s", order.OrderNumber)
				return errors.New("update asset for jrdidi fail")
			}
		} else {
			utils.Log.Errorf("the qty_frozen is not enough for jrdidi, assetForJrdidi = %+v", assetForJrdidi)
			utils.Log.Errorf("func TransferNormally finished abnormally, order_number = %s", order.OrderNumber)
			return errors.New("the qty_frozen is not enough for jrdidi")
		}
		//
		if err := tx.Model(&models.Order{}).Where("id = ?", order.Id).Updates(models.Order{BTUSDFlowStatus: models.BTUSDFlowD1MerchantFrozenToMerchantQty}).Error; err != nil {
			utils.Log.Errorf("Can't update order %s BTUSDFlowStatus to %s. error: %v", order.OrderNumber, models.BTUSDFlowD1MerchantFrozenToMerchantQty, err)
			utils.Log.Errorf("func TransferNormally finished abnormally, order_number = %s", order.OrderNumber)
			return errors.New("update BTUSDFlowStatus fail")
		}
	} else { // 币没有到币商的冻结账号中（还在平台商的冻结账号中）

		// 把币从平台商的冻结账号中，分到各自的的可用账号中

		// 先计算下单时冻结了平台商多少个币
		var frozenBTUSD float64
		if utils.BtusdCompareGte(order.TraderBTUSDFeeIncome, 0) {
			// 平台商也赚用户的手续费，仅冻结了订单的BTUSD数量
			frozenBTUSD = order.Quantity
		} else {
			// 平台补贴用户的手续费，则会冻结比订单BTUSD数量更多的币
			frozenBTUSD = order.Quantity - order.TraderBTUSDFeeIncome
		}

		// 再计算平台商在本订单中有没有“最终收入”（用户手续费）
		var traderFinalBTUSDFeeIncome float64
		if utils.BtusdCompareGte(order.TraderBTUSDFeeIncome, 0) {
			// 这个订单平台商有手续费收入
			traderFinalBTUSDFeeIncome = order.TraderBTUSDFeeIncome
		} else {
			// 这个订单平台商没有手续费收入
			traderFinalBTUSDFeeIncome = 0
		}

		// 减少平台冻结的BTUSD，同时增加平台赚的手续费（如果有的话）
		if utils.BtusdCompareGte(assetForTrader.QtyFrozen, frozenBTUSD) {
			if err := tx.Table("assets").Where("id = ?", assetForTrader.Id).
				Updates(map[string]interface{}{
					"qty_frozen": assetForTrader.QtyFrozen - frozenBTUSD,
					"quantity":   assetForTrader.Quantity + traderFinalBTUSDFeeIncome}).Error; err != nil {
				utils.Log.Errorf("update asset for distributor fail, assetForTrader = %+v", assetForTrader)
				utils.Log.Errorf("func TransferNormally finished abnormally, order_number = %s", order.OrderNumber)
				return errors.New("update asset for distributor fail")
			}
		} else {
			utils.Log.Errorf("the qty_frozen is not enough for distributor, assetForTrader = %+v", assetForTrader)
			utils.Log.Errorf("func TransferNormally finished abnormally, order_number = %s", order.OrderNumber)
			return errors.New("the qty_frozen is not enough for distributor")
		}

		// 增加币商的BTUSD
		if err := tx.Table("assets").Where("id = ?", assetForMerchant.Id).
			Updates(map[string]interface{}{
				"quantity": assetForMerchant.Quantity + (order.Quantity - order.TraderBTUSDFeeIncome - order.JrdidiBTUSDFeeIncome)}).Error; err != nil {
			utils.Log.Errorf("func TransferNormally finished abnormally, order_number = %s", order.OrderNumber)
			return errors.New("can not find merchant asset")
		}
		// 增加jrdidi平台的BTUSD(赚的手续费)
		if err := tx.Table("assets").Where("id = ?", assetForJrdidi.Id).
			Updates(map[string]interface{}{
				"quantity": assetForJrdidi.Quantity + order.JrdidiBTUSDFeeIncome}).Error; err != nil {
			utils.Log.Errorf("func TransferNormally finished abnormally, order_number = %s", order.OrderNumber)
			return errors.New("can not find jrdidi asset")
		}
		//
		if err := tx.Model(&models.Order{}).Where("id = ?", order.Id).Updates(models.Order{BTUSDFlowStatus: models.BTUSDFlowD1TraderFrozenToMerchantQty}).Error; err != nil {
			utils.Log.Errorf("Can't update order %s BTUSDFlowStatus to %s. error: %v", order.OrderNumber, models.BTUSDFlowD1TraderFrozenToMerchantQty, err)
			return errors.New("update BTUSDFlowStatus fail")
		}
	}

	// 下面修改asset history
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
func TransferAbnormally(tx *gorm.DB, assetForTrader *models.Assets, assetForMerchant *models.Assets, assetForJrdidi *models.Assets,
	order *models.Order) error {
	utils.Log.Debugf("func TransferAbnormally begin, order_number = %s", order.OrderNumber)

	// 用户充值订单，不收手续费，这个方法未对充值订单进行测试，不要调用它
	if order.Direction == 0 {
		utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
		return errors.New("not applicable for order with direction == 0")
	}

	// 用户提现订单未完成，币商未真正付款
	// 把各自冻结的币退到平台的可用账号中

	// 先计算下订单时冻结了平台商多少个币
	var frozenBTUSD float64
	if utils.BtusdCompareGte(order.TraderBTUSDFeeIncome, 0) {
		// 平台商也赚用户的手续费，仅冻结了订单的BTUSD数量
		frozenBTUSD = order.Quantity
	} else {
		// 平台补贴用户的手续费，则会冻结比订单BTUSD数量更多的币
		frozenBTUSD = order.Quantity - order.TraderBTUSDFeeIncome
	}

	// 分两种情况。
	// 情况一、币已经到币商的冻结账号中。
	// 情况二，币还在平台商的冻结账号中。

	var alreadyAddCoinToMerchantFrozenAccount = false
	if order.BTUSDFlowStatus == models.BTUSDFlowD1TraderFrozenToMerchantFrozen {
		alreadyAddCoinToMerchantFrozenAccount = true
	}

	if alreadyAddCoinToMerchantFrozenAccount {
		// 情况一：币已经到币商的冻结账号中。少部分币（部分手续费）在jrdidi平台及平台商的冻结账号中
		//
		// 先计算这个订单币商qty_frozen账号中获得多少BTUSD（包含他赚的手续费）
		var btusdInMerchantFrozen float64 = order.Quantity - (order.TraderBTUSDFeeIncome + order.MerchantBTUSDFeeIncome + order.JrdidiBTUSDFeeIncome) + order.MerchantBTUSDFeeIncome
		if utils.BtusdCompareGte(assetForMerchant.QtyFrozen, btusdInMerchantFrozen) { // 避免merchant的qty_frozen出现负数
			// 减少merchant获得的BTUSD
			if err := tx.Table("assets").Where("id = ?", assetForMerchant.Id).
				Updates(map[string]interface{}{
					"qty_frozen": assetForMerchant.QtyFrozen - btusdInMerchantFrozen}).Error; err != nil {
				utils.Log.Errorf("update asset for merchant fail, error = %+v", err)
				utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
				return errors.New("update asset for merchant fail")
			}
		} else {
			utils.Log.Errorf("the qty_frozen is not enough for merchant, assetForMerchant = %+v", assetForMerchant)
			utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
			return errors.New("the qty_frozen is not enough for merchant")
		}

		// 减少jrdidi赚取的BTUSD
		if utils.BtusdCompareGte(assetForJrdidi.QtyFrozen, order.JrdidiBTUSDFeeIncome) {
			if err := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForJrdidi.Id, order.JrdidiBTUSDFeeIncome).
				Updates(map[string]interface{}{
					"qty_frozen": assetForJrdidi.QtyFrozen - order.JrdidiBTUSDFeeIncome}).Error; err != nil {
				utils.Log.Errorf("update asset for jrdidi fail, error = %+v", err)
				utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
				return errors.New("update asset for jrdidi fail")
			}
		} else {
			utils.Log.Errorf("the qty_frozen is not enough for jrdidi, assetForJrdidi = %+v", assetForJrdidi)
			utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
			return errors.New("the qty_frozen is not enough for jrdidi")
		}

		// 计算平台商在本订单中有没有“收入”（用户手续费），这个值不会是负数。
		var traderFinalBTUSDFeeIncome float64
		if utils.BtusdCompareGte(order.TraderBTUSDFeeIncome, 0) {
			// 这个订单平台商有手续费收入
			traderFinalBTUSDFeeIncome = order.TraderBTUSDFeeIncome
		} else {
			// 这个订单平台商没有手续费收入
			traderFinalBTUSDFeeIncome = 0
		}

		if utils.BtusdCompareGte(assetForTrader.QtyFrozen, traderFinalBTUSDFeeIncome) { // 避免trader的qty_frozen出现负数
			if err := tx.Table("assets").Where("id = ?", assetForTrader.Id).
				Updates(map[string]interface{}{
					"qty_frozen": assetForTrader.QtyFrozen - traderFinalBTUSDFeeIncome, // 当平台商自己也拿部分用户手续费时，这部分手续费会在平台商的冻结账号中，要把它从冻结账号中减掉
					"quantity":   assetForTrader.Quantity + frozenBTUSD}).Error;        // 平台商的可用账号应该增加创建订单时所冻结的币的总数
			err != nil {
				utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
				return errors.New("update distributor asset fail")
			}
		} else {
			utils.Log.Errorf("the qty_frozen is not enough for trader, assetForTrader = %+v", assetForTrader)
			utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
			return errors.New("update distributor asset fail")
		}
		//
		if err := tx.Model(&models.Order{}).Where("id = ?", order.Id).Updates(models.Order{BTUSDFlowStatus: models.BTUSDFlowD1MerchantFrozenToTraderQty}).Error; err != nil {
			utils.Log.Errorf("Can't update order %s BTUSDFlowStatus to %s. error: %v", order.OrderNumber, models.BTUSDFlowD1MerchantFrozenToTraderQty, err)
			return errors.New("update BTUSDFlowStatus fail")
		}
	} else { // 情况二：币没有到币商的冻结账号中（还在平台商的冻结账号中）

		// 直接把平台商的冻结账号中的钱退到平台商的账号中
		if utils.BtusdCompareGte(assetForTrader.QtyFrozen, frozenBTUSD) { // 避免trader的qty_frozen出现负数
			if err := tx.Table("assets").Where("id = ?", assetForTrader.Id).
				Updates(map[string]interface{}{
					"qty_frozen": assetForTrader.QtyFrozen - frozenBTUSD,
					"quantity":   assetForTrader.Quantity + frozenBTUSD}).Error; err != nil {
				utils.Log.Errorf("update trader asset fail, err = %s", err)
				utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
				return errors.New("update trader asset fail")
			}
		} else {
			utils.Log.Errorf("the qty_frozen is not enough for distributor, assetForTrader = %+v", assetForTrader)
			utils.Log.Errorf("func TransferAbnormally finished abnormally, order_number = %s", order.OrderNumber)
			return errors.New("the qty_frozen is not enough for distributor")
		}
		//
		if err := tx.Model(&models.Order{}).Where("id = ?", order.Id).Updates(models.Order{BTUSDFlowStatus: models.BTUSDFlowD1TraderFrozenToTraderQty}).Error; err != nil {
			utils.Log.Errorf("Can't update order %s BTUSDFlowStatus to %s. error: %v", order.OrderNumber, models.BTUSDFlowD1TraderFrozenToTraderQty, err)
			return errors.New("update BTUSDFlowStatus fail")
		}
	}

	utils.Log.Debugf("func TransferAbnormally finished normally, order_number = %s", order.OrderNumber)
	return nil
}
