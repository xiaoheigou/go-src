package service

import "testing"

func TestGetBestPaymentID(t *testing.T) {
	order := OrderToFulfill{
		PayType: 5,
	}

	GetBestNormalPaymentID(&order, 1)
}
