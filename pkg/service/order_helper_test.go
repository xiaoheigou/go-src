package service

import (
	"testing"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

func getTestOrder() models.Order {

	var orderId int64 = 3 // 测试时修改它，保证它不存在，否则有错误 Duplicate entry '' for key 'PRIMARY'
	var orderNumber string = "xxx"
	var testDistributorId int64 = 100 // 测试时修改它
	var testMerchantId int64 = 10003  // 测试时修改它

	var traderBTUSDFeeIncome = 2.4
	var merchantBTUSDFeeIncome = 1.0
	var jrdidiBTUSDFeeIncome = 1.3

	var btusdFlowStatus int32 = models.BTUSDFlowD1TraderFrozenToMerchantFrozen //BTUSDFlowD1TraderFrozenToMerchantFrozen

	order := models.Order{
		BTUSDFlowStatus:        btusdFlowStatus,
		Id:                     orderId,
		OrderNumber:            orderNumber,
		OriginOrder:            "",
		Price:                  0,
		Price2:                 0,
		Quantity:               100,
		Amount:                 0,
		OriginAmount:           0,
		Fee:                    0,
		PaymentRef:             "",
		Status:                 0,
		StatusReason:           0,
		Synced:                 0,
		Direction:              1,
		DistributorId:          testDistributorId,
		DistributorName:        "",
		MerchantId:             testMerchantId,
		MerchantName:           "",
		MerchantPhone:          "",
		MerchantPaymentId:      0,
		TraderBTUSDFeeIncome:   traderBTUSDFeeIncome,
		MerchantBTUSDFeeIncome: merchantBTUSDFeeIncome,
		JrdidiBTUSDFeeIncome:   jrdidiBTUSDFeeIncome,
		AccountId:              "",
		CurrencyCrypto:         "BTUSD",
		CurrencyFiat:           "",
		PayType:                0,
		QrCode:                 "",
		Name:                   "",
		BankAccount:            "",
		Bank:                   "",
		BankBranch:             "",
		AcceptedAt:             time.Time{},
		PaidAt:                 time.Time{},
		PaymentConfirmedAt:     time.Time{},
		TransferredAt:          time.Time{},
		SvrCurrentTime:         time.Time{},
		AppCoinName:            "",
		Remark:                 "",
		Timeout:                0,
		AppServerNotifyUrl:     "",
		AppReturnPageUrl:       "",
		Timestamp:              models.Timestamp{},
	}

	//if db := utils.DB.Create(&order); db.Error != nil {
	//	utils.Log.Error("create order fail")
	//}
	return order
}

func TestTransferCoinFromTraderFrozenToMerchantFrozen(t *testing.T) {
	tx := utils.DB.Begin()

	order := getTestOrder()

	// 找到平台asset记录
	assetForTrader := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForTrader, "distributor_id = ? AND currency_crypto = ? ", order.DistributorId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		t.Fail()
	}

	// 找到币商asset记录
	assetForMerchant := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForMerchant, "merchant_id = ? AND currency_crypto = ? ", order.MerchantId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		t.Fail()
	}

	// 找到金融滴滴平台asset记录
	assetForPlatform := models.Assets{}
	platformDistId := 1 // 金融滴滴平台的distributor_id为1
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForPlatform, "distributor_id = ? AND currency_crypto = ? ",
		platformDistId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		t.Fail()
	}

	if err := TransferCoinFromTraderFrozenToMerchantFrozen(tx, &assetForTrader, &assetForMerchant, &assetForPlatform, &order); err != nil {
		utils.Log.Errorf("TransferCoinFromTraderFrozenToMerchantFrozen err ", err)
	}

	tx.Commit()

}

func TestTransferAbnormally(t *testing.T) {
	tx := utils.DB.Begin()

	order := getTestOrder()

	// 找到平台asset记录
	assetForTrader := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForTrader, "distributor_id = ? AND currency_crypto = ? ", order.DistributorId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		t.Fail()
	}

	// 找到币商asset记录
	assetForMerchant := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForMerchant, "merchant_id = ? AND currency_crypto = ? ", order.MerchantId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		t.Fail()
	}

	// 找到金融滴滴平台asset记录
	assetForPlatform := models.Assets{}
	platformDistId := 1 // 金融滴滴平台的distributor_id为1
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForPlatform, "distributor_id = ? AND currency_crypto = ? ",
		platformDistId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		t.Fail()
	}

	if err := TransferAbnormally(tx, &assetForTrader, &assetForMerchant, &assetForPlatform, &order); err != nil {
		tx.Rollback()
		utils.Log.Errorf("TransferCoinFromTraderFrozenToMerchantFrozen err ", err)
	}

	tx.Commit()

}

func TestTransferNormally(t *testing.T) {
	tx := utils.DB.Begin()

	order := getTestOrder()

	// 找到平台asset记录
	assetForTrader := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForTrader, "distributor_id = ? AND currency_crypto = ? ", order.DistributorId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		t.Fail()
	}

	// 找到币商asset记录
	assetForMerchant := models.Assets{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForMerchant, "merchant_id = ? AND currency_crypto = ? ", order.MerchantId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		t.Fail()
	}

	// 找到金融滴滴平台asset记录
	assetForPlatform := models.Assets{}
	platformDistId := 1 // 金融滴滴平台的distributor_id为1
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&assetForPlatform, "distributor_id = ? AND currency_crypto = ? ",
		platformDistId, order.CurrencyCrypto).RecordNotFound() {
		tx.Rollback()
		t.Fail()

	}

	if err := TransferNormally(tx, &assetForTrader, &assetForMerchant, &assetForPlatform, &order, nil); err != nil {
		tx.Rollback()
		utils.Log.Errorf("TransferCoinFromTraderFrozenToMerchantFrozen err ", err)
	}

	tx.Commit()

}
