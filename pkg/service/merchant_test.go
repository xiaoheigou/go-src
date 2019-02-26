package service

import (
	"crypto/md5"
	"encoding/hex"
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

func getMd5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func createMerchantWithAsset(phone string) error {
	pre := models.Preferences{
		InWork:      1,
		AutoConfirm: 0,
		AutoAccept:  0,
	}
	assets := []models.Assets{
		{
			Quantity:       10000000,
			QtyFrozen:      0,
			CurrencyCrypto: "BTUSD",
		},
	}

	algorithm := utils.Config.GetString("algorithm")
	salt, pass := generatePassword(getMd5("123456"), algorithm)
	merchant := models.Merchant{
		Preferences: pre,
		Asset:       assets,
		Phone:       phone,
		NationCode:  86,
		Nickname:    phone,
		Salt:        salt,
		Password:    pass,
		Email:       "test@163.com",
		Algorithm:   algorithm,
		LastLogin:   time.Now(),
		UserStatus:  0,
	}
	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Now()

	if err := utils.DB.Create(&merchant).Error; err != nil {
		utils.Log.Errorf("create merchant error, %v", err)
		return err
	}
	return nil
}

func TestCreateMerchantWithAsset(t *testing.T) {
	var phones []string
	phones = append(phones, "13012340000")
	phones = append(phones, "13012340001")
	phones = append(phones, "13012340002")
	phones = append(phones, "13012340003")
	phones = append(phones, "13012340004")
	phones = append(phones, "13012340005")
	phones = append(phones, "13012340006")
	phones = append(phones, "13012340007")
	phones = append(phones, "13012340008")
	phones = append(phones, "13012340009")

	for _, phone := range phones {
		if err := createMerchantWithAsset(phone); err != nil {
			t.Fail()
		}
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
		Algorithm:   algorithm,
		UserStatus:  0,
	}
	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Now()

	if err := utils.DB.Create(&merchant).Error; err != nil {
		utils.Log.Errorf("create merchant error,")
		t.Fail()
	}
	utils.Log.Infof("id = %d", merchant.Id)
	utils.SetCacheSetMember(utils.RedisKeyMerchantOnline(), 0, merchant.Id)
	utils.SetCacheSetMember(utils.RedisKeyMerchantWechatAutoOrder(), 0, merchant.Id)
	utils.SetCacheSetMember(utils.RedisKeyMerchantAlipayAutoOrder(), 0, merchant.Id)
	utils.SetCacheSetMember(utils.RedisKeyMerchantInWork(), 0, merchant.Id)
	//temp := GetMerchantsQualified(650, 650, "BTUSD", 1, true, 0, 0)
	//if len(temp) <= 0 {
	//	t.Fail()
	//}
	utils.DelCacheSetMember(utils.RedisKeyMerchantOnline(), merchant.Id)
	utils.DelCacheSetMember(utils.RedisKeyMerchantWechatAutoOrder(), merchant.Id)
	utils.DelCacheSetMember(utils.RedisKeyMerchantAlipayAutoOrder(), merchant.Id)
	utils.DelCacheSetMember(utils.RedisKeyMerchantInWork(), merchant.Id)
}
