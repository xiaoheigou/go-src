package service

import (
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/utils"
)

func SetOnlineMerchant(uid int) error {
	//设置工作状态
	data := GetMerchantWorkMode(uid)
	if data.Status == response.StatusSucc && len(data.Data) > 0 {
		if data.Data[0].InWork == 1 {
			if err := utils.SetCacheSetMember(utils.UniqueMerchantInWorkKey(), 0, uid); err != nil {
				utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
				return err
			}
		}
		if data.Data[0].AutoAccept == 1 {
			if err := utils.SetCacheSetMember(utils.UniqueMerchantAutoAcceptKey(), 0, uid); err != nil {
				utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
				return err
			}
		}
		if data.Data[0].AutoConfirm == 1 {
			if err := utils.SetCacheSetMember(utils.UniqueMerchantAutoConfirmKey(), 0, uid); err != nil {
				utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
				return err
			}
		}
	}
	//设置merchant在线
	if err := utils.SetCacheSetMember(utils.UniqueMerchantOnlineKey(), 0, uid); err != nil {
		utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
		return err
	}
	return nil
}

func UpdateMerchantWorkMode(uid, workMode int, key string) error {
	if workMode == 1 {
		if err := utils.SetCacheSetMember(key, 0, uid); err != nil {
			utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
			return err
		}
	} else if workMode == 0 {
		if err := utils.DelCacheSetMember(key, uid); err != nil {
			utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
			return err
		}
	}

	return nil
}

//merchant 离线
func DelOnlineMerchant(uid int) error {
	if err := utils.DelCacheSetMember(utils.UniqueMerchantOnlineKey(), uid); err != nil {
		utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
		return err
	}
	return nil
}
