package response

import "yuudidi.com/pkg/models"

type GetMerchantAssetHistoryRet struct {
	CommonRet
	Data []models.AssetHistory `json:"data"`
}

type GetRechargeApplies struct {

}