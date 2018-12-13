package service

import (
	"testing"

	"yuudidi.com/pkg/models"
)

func TestFulfillOrderJob(t *testing.T) {
	order := OrderToFulfill{
		AccountID:      "123",
		DistributorID:  123,
		OrderNumber:    "123",
		Direction:      1,
		CurrencyCrypto: "BTUSD",
		CurrencyFiat:   "CNY",
		Quantity:       "1.0",
		Price:          "6.2",
		Amount:         "62.0",
		PayType:        1,
	}
	engine := NewOrderFulfillmentEngine(nil)
	engine.FulfillOrder(&order)
}

func TestSendOrderJob(t *testing.T) {
	order := OrderToFulfill{
		AccountID:      "123",
		DistributorID:  123,
		OrderNumber:    "123",
		Direction:      1,
		CurrencyCrypto: "BTUSD",
		CurrencyFiat:   "CNY",
		Quantity:       "1.0",
		Price:          "6.2",
		Amount:         "62.0",
		PayType:        1,
	}
	merchants := []int64{1, 2, 3}

	engine := NewOrderFulfillmentEngine(nil)
	engine.SendOrder(&order, &merchants)
}

func TestNotifyFulfillment(t *testing.T) {
	fulfillment := OrderFulfillment{
		OrderToFulfill: OrderToFulfill{
			AccountID:      "123",
			DistributorID:  123,
			OrderNumber:    "123",
			Direction:      1,
			CurrencyCrypto: "BTUSD",
			CurrencyFiat:   "CNY",
			Quantity:       "1.0",
			Price:          "6.2",
			Amount:         "62.0",
			PayType:        1,
		},
		MerchantID:        123,
		MerchantNickName:  "yuudidi",
		MerchantAvatarURI: "yuudidi",
		PaymentInfo: models.PaymentInfo{
			PayType:     1,
			Name:        "yuudidi",
			Bank:        "yuudidi bank",
			BankAccount: "yuudidi",
			BankBranch:  "yuudidi",
		},
	}

	engine := NewOrderFulfillmentEngine(nil)
	engine.NotifyFulfillment(&fulfillment)
}
