package service

import (
	"testing"
	"time"

	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

func TestCreateMerchant(t *testing.T) {
	merchant := models.Merchant{
		Nickname:   "yuudidi",
		AvatarUri:  "yuudidi.com",
		DisplayUid: "Y",
		Password:   []byte{'a'},
		Salt:       []byte{'b'},
		Algorithm:  "Argon2",
		Phone:      "134",
		NationCode: 86,
		Email:      "a@b.com",
		UserStatus: 1,
		UserCert:   1,
		LastLogin:  time.Now(),
		Quantity:   "10000",
	}
	if err := utils.DB.Create(&merchant).Error; err != nil {
		t.Errorf("Merchant creation failed: %v", err)
	}
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
	utils.SetCacheSetMember(utils.UniqueMerchantOnlineKey(), merchant.Id)
	utils.SetCacheSetMember(utils.UniqueMerchantAutoAcceptKey(), merchant.Id)
	utils.SetCacheSetMember(utils.UniqueMerchantAutoConfirmKey(), merchant.Id)
	utils.SetCacheSetMember(utils.UniqueMerchantInWorkKey(), merchant.Id)
	temp := GetMerchantsQualified(650, 650, "BTUSD", 1, true, 0, 0)
	if len(temp) <= 0 {
		t.Fail()
	}
	utils.DelCacheSetMember(utils.UniqueMerchantOnlineKey(), merchant.Id)
	utils.DelCacheSetMember(utils.UniqueMerchantAutoAcceptKey(), merchant.Id)
	utils.DelCacheSetMember(utils.UniqueMerchantAutoConfirmKey(), merchant.Id)
	utils.DelCacheSetMember(utils.UniqueMerchantInWorkKey(), merchant.Id)
}
