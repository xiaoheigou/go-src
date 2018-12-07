package response

import "YuuPay_core-service/pkg/models"

type OrdersRet struct {
	CommonRet

	Entity struct {
		Data []models.Order
	}
}
