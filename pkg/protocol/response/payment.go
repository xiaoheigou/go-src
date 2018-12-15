package response

import "yuudidi.com/pkg/models"

type GetPaymentsPageRet struct {
	CommonRet
	Pagination
	Data []models.PaymentInfo `json:"data"`
}

type AddPaymentRet struct {
	CommonRet
}

type DeletePaymentRet struct {
	CommonRet
}
