package service

import (
	"testing"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/utils"
)

func TestUnFreezeCoin(t *testing.T) {
	order := models.Order{
		OriginOrder:"123",
		OrderNumber:"12344123334",
		Status:models.SUSPENDED,
		MerchantId:38,
		DistributorId:1,
		CurrencyCrypto:"BTUSD",
		Quantity:100,
	}
	utils.DB.Create(&order)
	if UnFreezeCoin(order.OrderNumber,"123",1).Status == response.StatusFail {
		t.Fail()
	}

}
