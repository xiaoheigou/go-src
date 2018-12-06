package models

type AssetHistory struct {
	Id           int    `json:"Id"`
	MerchantId   int    `json:"merchant_id"`
	Msg          string `json:"msg"`
	Timestamp    string `json:"timestamp"`
	OperatorId   int    `json:"operator_id"`
	OperatorName string `json:"operator_name"`
}
