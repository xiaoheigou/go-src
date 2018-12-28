package service

import (
	"strconv"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/utils"
)

func SetOnlineMerchant(uid int) error {
	//设置工作状态
	if err := SetMerchantAutoOrAccept(uid); err != nil {
		utils.Log.Errorf("set merchant work mode is failed,merchantId:%d", uid)
		return err

	}
	//设置merchant在线
	if err := utils.SetCacheSetMember(utils.UniqueMerchantOnlineKey(), uid); err != nil {
		utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
		return err
	}
	return nil
}

func SetMerchantInWork(uid int) error {
	online, err := utils.GetCacheSetMembers(utils.UniqueMerchantOnlineKey())
	if err != nil {
		utils.Log.Errorf("get online merchant is failed")
		return err
	}
	temp := strconv.FormatInt(int64(uid), 10)
	isOnline := false
	for _, v := range online {
		if v == temp {
			isOnline = true
			break
		}
	}
	if isOnline {
		if err := SetMerchantAutoOrAccept(uid); err != nil {
			utils.Log.Errorf("set merchant word mode is failed,merchantId:%d", uid)
			return err
		}
	}
	return nil
}

func SetMerchantAutoOrAccept(uid int) error {
	data := GetMerchantWorkMode(uid)
	if data.Status == response.StatusSucc && len(data.Data) > 0 {
		temp := data.Data[0]
		if temp.InWork == 1 && temp.AutoAccept == 1 && temp.AutoConfirm == 1 {
			if err := utils.SetCacheSetMember(utils.UniqueMerchantOnlineAutoKey(), uid); err != nil {
				utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
				return err
			}
		} else if temp.InWork == 1 {
			if err := utils.SetCacheSetMember(utils.UniqueMerchantOnlineAcceptKey(), uid); err != nil {
				utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
				return err
			}
		}
	}
	return nil
}

//merchant 不自动接单
func DelAuto(uid int) error {
	if err := utils.DelCacheSetMember(utils.UniqueMerchantOnlineAutoKey(), uid); err != nil {
		utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
		return err
	}
	return nil
}

//merchant 不工作
func DelInWork(uid int) error {
	if err := DelAuto(uid); err != nil {
		utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
		return err
	}
	if err := utils.DelCacheSetMember(utils.UniqueMerchantOnlineAcceptKey(), uid); err != nil {
		utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
		return err
	}
	return nil
}

//merchant 离线
func DelOnlineMerchant(uid int) error {

	if err := DelInWork(uid); err != nil {
		utils.Log.Errorf("set merchant online is failed,merchantId:%d", uid)
		return err
	}

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

func convertStringToInt(ids []string) ([]int64, error) {
	var result []int64
	for _, id := range ids {
		if temp, err := strconv.ParseInt(id, 10, 64); err != nil {
			return nil, err
		} else {
			result = append(result, temp)
		}
	}
	return result, nil
}
