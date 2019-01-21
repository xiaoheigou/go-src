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
	// 用户充值订单，不收手续费，这个方法未对充值订单进行测试，不要调用它
	if order.Direction == 0 {
		return errors.New("not applicable for order with direction == 0")
	}

	var deductBTUSD float64 = order.Quantity - order.TraderBTUSDFeeIncome // 当平台自己也抽取用户的提现手续费时，TraderBTUSDFeeIncome大于0，否则TraderBTUSDFeeIncome <= 0
	// 减少平台冻结的BTUSD
	// 注：平台自己赚取的那部分用户手续费还在冻结着
	if rowsAffected := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForTrader.Id, deductBTUSD).
		Updates(map[string]interface{}{
			"qty_frozen": assetForTrader.QtyFrozen - deductBTUSD}).RowsAffected; rowsAffected == 0 {
		utils.Log.Errorf("the qty_frozen is not enough for distributor, assetForTrader = %+v", assetForTrader)
		return errors.New("the qty_frozen is not enough for distributor")
	}

	// 增加币商冻结的BTUSD
	if rowsAffected := tx.Table("assets").Where("id = ?", assetForMerchant.Id).
		Updates(map[string]interface{}{
			"qty_frozen": assetForMerchant.QtyFrozen + (order.Quantity - order.TraderBTUSDFeeIncome - order.JrdidiBTUSDFeeIncome)}).RowsAffected; rowsAffected == 0 {
		return errors.New("can not find merchant asset")
	}
	// 增加jrdidi平台冻结的BTUSD
	if rowsAffected := tx.Table("assets").Where("id = ?", assetForJrdidi.Id).
		Updates(map[string]interface{}{
			"qty_frozen": assetForJrdidi.QtyFrozen + order.JrdidiBTUSDFeeIncome}).RowsAffected; rowsAffected == 0 {
		return errors.New("can not find jrdidi asset")
	}

	return nil
}

// 下面函数不会commit，也不会rollback，请在上层函数处理
func TransferNormally(tx *gorm.DB, assetForTrader *models.Assets, assetForMerchant *models.Assets, assetForJrdidi *models.Assets, order *models.Order) error {
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

	return nil
}

// 下面函数不会commit，也不会rollback，请在上层函数处理
func TransferAbnormally(tx *gorm.DB, assetForTrader *models.Assets, assetForMerchant *models.Assets, assetForJrdidi *models.Assets, order *models.Order) error {
	// 用户充值订单，不收手续费，这个方法未对充值订单进行测试，不要调用它
	if order.Direction == 0 {
		return errors.New("not applicable for order with direction == 0")
	}
	// 用户提现订单未完成，币商未真正付款
	// 把各自冻结的币退到平台的冻结账号中

	// 减少币商获得的BTUSD（包含他赚的手续费）
	if rowsAffected := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForMerchant.Id, order.Quantity-order.TraderBTUSDFeeIncome-order.JrdidiBTUSDFeeIncome).
		Updates(map[string]interface{}{
			"qty_frozen": assetForMerchant.QtyFrozen - (order.Quantity - order.TraderBTUSDFeeIncome - order.JrdidiBTUSDFeeIncome)}).RowsAffected; rowsAffected == 0 {
		utils.Log.Errorf("the qty_frozen is not enough for merchant, assetForMerchant = %+v", assetForMerchant)
		return errors.New("the qty_frozen is not enough for merchant")
	}
	// 减少jrdidi赚取的BTUSD
	if rowsAffected := tx.Table("assets").Where("id = ? and qty_frozen >= ?", assetForJrdidi.Id, order.JrdidiBTUSDFeeIncome).
		Updates(map[string]interface{}{
			"qty_frozen": assetForJrdidi.QtyFrozen - order.JrdidiBTUSDFeeIncome}).RowsAffected; rowsAffected == 0 {
		utils.Log.Errorf("the qty_frozen is not enough for jrdidi, assetForJrdidi = %+v", assetForJrdidi)
		return errors.New("the qty_frozen is not enough for jrdidi")
	}
	// 把上面两部分减少的BTUSD全部加到平台的冻结账号中
	if rowsAffected := tx.Table("assets").Where("id = ?", assetForTrader.Id).
		Updates(map[string]interface{}{
			"qty_frozen": assetForTrader.QtyFrozen + order.Quantity - order.TraderBTUSDFeeIncome}).RowsAffected; rowsAffected == 0 {
		return errors.New("can not find distributor asset")
	}

	return nil
}
