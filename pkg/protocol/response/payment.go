package response

import "yuudidi.com/pkg/models"

type GetPaymentsRet struct {
	CommonRet
	Data []models.PaymentInfo `json:"data"`
}

type AddPaymentRet struct {
	CommonRet
}

type DeletePaymentRet struct {
	CommonRet
}
