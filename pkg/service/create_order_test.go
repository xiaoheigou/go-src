package service

import (
	"github.com/jinzhu/gorm"
	"testing"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/utils"
)

func TestCreateOrder(t *testing.T) {

	var orderRequest response.OrderRequest
	orderRequest.DistributorId = 123
	orderRequest.CurrencyCrypto = "BTUSD"
	orderRequest.Quantity = 100

	//utils.Log.Debugf("distributor (id=%d) quantity = [%d], order (%s) quantity = [%d]", orderRequest.DistributorId, assets.Quantity, orderRequest.OrderNumber, orderRequest.Quantity)
	tx := utils.DB.Begin()

	var assets models.Assets

	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&assets, "distributor_id=? and currency_crypto=?", orderRequest.DistributorId, orderRequest.CurrencyCrypto).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// 没有找到，则创建记录
			if err := tx.Create(&models.Assets{DistributorId: orderRequest.DistributorId, CurrencyCrypto: orderRequest.CurrencyCrypto}).Error; err != nil {
				utils.Log.Errorf("create distributor assets fail: %v", err)
				tx.Rollback()
			}
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&assets, "distributor_id=? and currency_crypto=?", orderRequest.DistributorId, orderRequest.CurrencyCrypto).Error; err != nil {
				utils.Log.Errorf("find distributor assets fail: %v", err)
				tx.Rollback()
			}
		}
		utils.Log.Errorf("find distributor assets fail: %v", err)
		tx.Rollback()
	}

	utils.Log.Errorf("Here")
	//给平台商锁币
	if rowAffected := tx.Model(&models.Assets{}).Where("distributor_id = ? AND currency_crypto = ? AND quantity >= ?", orderRequest.DistributorId, orderRequest.CurrencyCrypto, orderRequest.Quantity).
		Updates(map[string]interface{}{"quantity": assets.Quantity - orderRequest.Quantity, "qty_frozen": assets.QtyFrozen + orderRequest.Quantity}).RowsAffected; rowAffected == 0 {
		tx.Rollback()
		utils.Log.Errorf("NOOOOOOOOO")
		utils.Log.Errorf("tx in func PlaceOrder rollback")
		utils.Log.Errorf("the distributor frozen quantity= %f, distributorId= %s", orderRequest.Quantity, orderRequest.DistributorId)
	}

	tx.Commit()
}
