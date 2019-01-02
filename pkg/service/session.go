package service

import (
	"strconv"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/utils"
)

func SetOnlineMerchant(uid int) error {
	//设置工作状态
	data := GetMerchantWorkMode(uid)
	if data.Status == response.StatusSucc && len(data.Data) > 0 {
		if data.Data[0].InWork == 1 {
			if err := utils.SetCacheSetMember(utils.UniqueMerchantInWorkKey(), uid); err != nil {
				utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
				return err
			}
		}
		if data.Data[0].AutoAccept == 1 {
			if err := utils.SetCacheSetMember(utils.UniqueMerchantAutoAcceptKey(), uid); err != nil {
				utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
				return err
			}
		}
		if data.Data[0].AutoConfirm == 1 {
			if err := utils.SetCacheSetMember(utils.UniqueMerchantAutoConfirmKey(), uid); err != nil {
				utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
				return err
			}
		}
	}
	//设置merchant在线
	if err := utils.SetCacheSetMember(utils.UniqueMerchantOnlineKey(), uid); err != nil {
		utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
		return err
	}
	return nil
}

func UpdateMerchantWorkMode(uid, workMode int, key string) error {
	if workMode == 1 {
		if err := utils.SetCacheSetMember(key, uid); err != nil {
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

func mergeList(l1, l2, l3 []int64) []int64 {
	var result []int64
	tempMap := make(map[int64]int)

	for _, v := range l1 {
		tempMap[v] = 1
	}

	for _, v := range l2 {
		if tempMap[v] == 1 {
			tempMap[v] = 2
		}
	}

	for _, v := range l3 {
		if tempMap[v] == 2 {
			tempMap[v] = 3
			result = append(result, v)
		}
	}

	return result
}

func convertStringToInt(ids []string,results *[]int64) error {
	var result []int64
	for _, id := range ids {
		if temp, err := strconv.ParseInt(id, 10, 64); err != nil {
			return  err
		} else {
			result = append(result, temp)
		}
	}
	*results = result
	return  nil
}
