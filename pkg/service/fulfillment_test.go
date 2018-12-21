package service

import (
	"testing"

	"yuudidi.com/pkg/models"
)

func TestNotifyFulfillment(t *testing.T) {
	fulfillment := OrderFulfillment{
		OrderToFulfill: OrderToFulfill{
			OrderNumber:    "123",
			Direction:      1,
			CurrencyCrypto: "BTUSD",
			CurrencyFiat:   "CNY",
			Quantity:       1.0,
			Price:          6.2,
			Amount:         62.0,
			PayType:        1,
		},
		MerchantID:        1,
		MerchantNickName:  "yuudidi",
		MerchantAvatarURI: "yuudidi",
		PaymentInfo: []models.PaymentInfo{
			{
				PayType:     1,
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
