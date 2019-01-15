package response

import "yuudidi.com/pkg/models"

type GetPaymentsPageRet struct {
	CommonRet
	Pagination
	Data []models.PaymentInfo `json:"data"`
}

type AddPaymentRet struct {
	CommonRet
	Data []models.PaymentInfo `json:"data"`
}

type DeletePaymentRet struct {
	CommonRet
}

type GetBankListRet struct {
	CommonRet
	Data []models.BankInfo `json:"data"`
}
