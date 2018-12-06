package response

import "YuuPay_core-service/pkg/models"

type GetMerchantAssetHistoryRet struct {
	CommonRet
	Entity struct {
		Data []models.AssetHistory `json:"data"`
	}
}
