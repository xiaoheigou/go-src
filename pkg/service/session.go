package service

import (
	"yuudidi.com/pkg/utils"
)

func GetMerchantAutoList() []string {
	key := utils.UniqueMerchantOnlineAutoKey()
	all, err := utils.GetCacheSetMembers(key)
	if err != nil {
		return []string{}
	}
	return all
}

func GetMerchantList() []string {
	key := utils.UniqueMerchantOnlineKey()
	all, err := utils.GetCacheSetMembers(key)
	if err != nil {
		return []string{}
	}
	return all
}

func AddOnlineMerdhantAuto(merchantId int64) {
	key := utils.UniqueMerchantOnlineAutoKey()
	utils.SetCacheSetMember(key, merchantId)
}

func AddOnlineMerdhant(merchantId int64) {
	key := utils.UniqueMerchantOnlineKey()
	utils.SetCacheSetMember(key, merchantId)
}

func DelOnlineMerdhantAuto(merchantId int64) {
	key := utils.UniqueMerchantOnlineAutoKey()
	utils.DelCacheSetMember(key, merchantId)
}

func DelOnlineMerdhant(merchantId int64) {
	key := utils.UniqueMerchantOnlineKey()
	utils.DelCacheSetMember(key, merchantId)
}
