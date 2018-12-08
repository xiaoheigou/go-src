package response

import "yuudidi.com/pkg/models"

type OrdersRet struct {
	CommonRet

	Data []models.Order
}
