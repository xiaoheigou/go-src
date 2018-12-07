package response

import "yuudidi.com/pkg/models"

type OrdersRet struct {
	CommonRet

	Entity struct {
		Data []models.Order
	}
}
