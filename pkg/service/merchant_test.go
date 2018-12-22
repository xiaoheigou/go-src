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
