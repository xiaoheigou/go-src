package service

import (
	"testing"

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
		MerchantID:        1,
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
		MerchantId: []int64{1},
		H5:         []string{"123"},
		Timeout:    600,
		Data:       []interface{}{models.Data{OrderNumber: "123", Direction: 0}},
	}
	engine := NewOrderFulfillmentEngine(nil)
	engine.UpdateFulfillment(msg)
}
