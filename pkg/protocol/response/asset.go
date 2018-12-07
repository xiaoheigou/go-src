package response

import "yuudidi.com/pkg/models"

type GetMerchantAssetHistoryRet struct {
	CommonRet
	Entity struct {
		Data []models.AssetHistory `json:"data"`
	}
}
