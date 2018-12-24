package service

import (
	"testing"
	"time"
	"yuudidi.com/pkg/utils"

	"yuudidi.com/pkg/models"
)

func TestFulfillOrderJob(t *testing.T) {
	order := OrderToFulfill{
		OrderNumber:    "111",
		AccountID:      "yuudidi",
		OriginOrder:    "1234567",
		DistributorID:  1,
		Direction:      0,
		CurrencyCrypto: "BTUSD",
		CurrencyFiat:   "CNY",
		Quantity:       1.0,
		Price:          6.5,
		Amount:         65.0,
		PayType:        1,
	}
	engine := NewOrderFulfillmentEngine(nil)
	engine.FulfillOrder(&order)
}

func TestSendOrderJob(t *testing.T) {
	order := OrderToFulfill{
		OrderNumber:    "222",
		Direction:      0,
		AccountID:      "yuudidi",
		OriginOrder:    "1234567",
		DistributorID:  1,
		CurrencyCrypto: "BTUSD",
		CurrencyFiat:   "CNY",
		Quantity:       2.0,
		Price:          6.35,
		Amount:         63.5,
		PayType:        1,
	}
	merchants := []int64{1, 2, 3}

	engine := NewOrderFulfillmentEngine(nil)
	engine.SendOrder(&order, &merchants)
}

func TestNotifyFulfillment(t *testing.T) {
	fulfillment := OrderFulfillment{
		OrderToFulfill: OrderToFulfill{
			OrderNumber:    "333",
			Direction:      1,
			AccountID:      "yuudidi",
			OriginOrder:    "1234567",
			DistributorID:  1,
			CurrencyCrypto: "BTUSD",
			CurrencyFiat:   "CNY",
			Quantity:       1.0,
			Price:          6.2,
			Amount:         62.0,
			PayType:        1,
		},
		MerchantID:        2,
		MerchantNickName:  "yuudidi",
		MerchantAvatarURI: "yuudidi",
		PaymentInfo: []models.PaymentInfo{
			{
				PayType:     4,
				Name:        "yuudidi",
				Bank:        "yuudidi bank",
				BankAccount: "yuudidi",
				BankBranch:  "yuudidi",
			},
		},
	}

	engine := NewOrderFulfillmentEngine(nil)
	engine.NotifyFulfillment(&fulfillment)
}

func TestUpdateFulfillment(t *testing.T) {
	msg := models.Msg{
		MsgType:    models.NotifyPaid,
		MerchantId: []int64{2},
		H5:         []string{"123"},
		Timeout:    600,
		Data:       []interface{}{models.Data{OrderNumber: "123", Direction: 0}},
	}
	engine := NewOrderFulfillmentEngine(nil)
	engine.UpdateFulfillment(msg)
}

func TestAcceptOrder(t *testing.T) {
	order := OrderToFulfill{
		OrderNumber:    "222",
		Direction:      0,
		AccountID:      "yuudidi",
		OriginOrder:    "1234567",
		DistributorID:  1,
		CurrencyCrypto: "BTUSD",
		CurrencyFiat:   "CNY",
		Quantity:       2.0,
		Price:          6.35,
		Amount:         63.5,
		PayType:        1,
	}
	merchantID := int64(2)
	engine := NewOrderFulfillmentEngine(nil)
	engine.AcceptOrder(order, merchantID)
}

func TestGetMerchantsQualified(t *testing.T) {
	//首先创建备用数据
	payments := []models.PaymentInfo{
		{
			Name:     "1234",
			PayType:  1,
			InUse:    0,
			EAmount:  0,
			EAccount: "13112345678",
			QrCode:   "http://13.250.12.109:8086/1_2_200_1545117100929.jpg",
		},
		{
			Name:     "1234",
			PayType:  1,
			InUse:    0,
			EAmount:  65,
			EAccount: "13112345678",
			QrCode:   "http://13.250.12.109:8086/1_2_200_1545117100929.jpg",
		},
		{
			Name:     "1234",
			PayType:  1,
			InUse:    0,
			EAmount:  650,
			EAccount: "13112345678",
			QrCode:   "http://13.250.12.109:8086/1_2_200_1545117100929.jpg",
		},
		{
			Name:     "1234",
			PayType:  1,
			InUse:    0,
			EAmount:  6500,
			EAccount: "13112345678",
			QrCode:   "http://13.250.12.109:8086/1_2_200_1545117100929.jpg",
		},
	}
	pre := models.Preferences{
		InWork:      1,
		AutoConfirm: 1,
		AutoAccept:  1,
	}
	assets := []models.Assets{
		{
			Quantity:       10000000,
			QtyFrozen:      0,
			CurrencyCrypto: "BTUSD",
		},
		{
			Quantity:       10000000,
			QtyFrozen:      0,
			CurrencyCrypto: "USDT",
		},
	}
	algorithm := utils.Config.GetString("algorithm")
	salt, pass := generatePassword("123456", algorithm)
	merchant := models.Merchant{
		Preferences: pre,
		Asset:       assets,
		Payments:    payments,
		Phone:       "13112345678",
		NationCode:  86,
		Nickname:    "test1",
		Salt:        salt,
		Password:    pass,
		Email:       "test@163.com",
		UserStatus:  0,
	}
	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Now()

	if err := utils.DB.Create(&merchant).Error; err != nil {
		utils.Log.Errorf("create merchant error,")
		t.Fail()
	}
	utils.Log.Infof("id = %d", merchant.Id)
	utils.SetCacheSetMember(utils.UniqueMerchantOnlineAutoKey(), merchant.Id)
	temp := GetMerchantsQualified(650, 650, "BTUSD", 1, true, 0, 0)
	if len(temp) <= 0 {
		t.Fail()
	}
	utils.DelCacheSetMember(utils.UniqueMerchantOnlineAutoKey(), merchant.Id)
}
