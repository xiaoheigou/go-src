package service

import (
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/service/dbcache"
	"yuudidi.com/pkg/utils"
)

// 设置币商的hook状态为可用，先读取数据库当前值，如果本来就可用，则什么都不做
func EnableHookStatus(merchantID int64, payType uint) {
	var merchant models.Merchant
	if err := dbcache.GetMerchantById(merchantID, &merchant); err != nil {
		utils.Log.Errorf("call GetMerchantById fail. [%v]", merchantID, err)
		utils.Log.Errorf("func EnableHookStatus finished abnormally. error %s", err)
		return
	}

	var pref models.Preferences
	if err := dbcache.GetPreferenceById(int64(merchant.PreferencesId), &pref); err != nil {
		utils.Log.Errorf("can't find preference record in db for merchant(uid=[%d]),  err [%v]", merchantID, err)
		utils.Log.Errorf("func EnableHookStatus finished abnormally. error %s", err)
		return
	}

	if payType == models.PaymentTypeWeixin {
		if pref.WechatHookStatus == 0 {
			// 修改preferences表
			if err := utils.DB.Table("preferences").Where("id = ?", merchant.PreferencesId).Update("wechat_hook_status", 1).Error; err != nil {
				utils.Log.Errorf("func EnableHookStatus, update preferences for merchant(uid=[%d]) fail. [%v]", merchantID, err)
			}
			// 使redis缓存失效
			if err := dbcache.InvalidatePreference(int64(merchant.PreferencesId)); err != nil {
				utils.Log.Errorf("func EnableHookStatus, InvalidatePreference fail, err [%v]", err)
			}
			// 如果相应的开关关掉/打开，则将merchant从相应redis key中删除/增加
			if err := UpdateMerchantWorkMode(int(merchantID), 1, utils.RedisKeyMerchantWechatHookStatus()); err != nil {
				utils.Log.Errorf("func EnableHookStatus, update preferences Redis for merchant(uid=[%d]) fail. [%v]", merchantID, err)
			}
		} else if pref.WechatHookStatus == 1 {
			// 目前已经是状态可用，什么都不用做
			utils.Log.Debugf("func EnableHookStatus do nothing, wechat_hook_status is already = 1 for merchant %d.", merchantID)
			return
		} else {
			utils.Log.Errorf("func EnableHookStatus, wechat_hook_status %d from db is not expected", pref.WechatHookStatus)
			return
		}
	} else if payType == models.PaymentTypeAlipay {
		if pref.AlipayHookStatus == 0 {
			// 修改preferences表
			if err := utils.DB.Table("preferences").Where("id = ?", merchant.PreferencesId).Update("alipay_hook_status", 1).Error; err != nil {
				utils.Log.Errorf("func EnableHookStatus, update preferences for merchant(uid=[%d]) fail. [%v]", merchantID, err)
			}
			// 使redis缓存失效
			if err := dbcache.InvalidatePreference(int64(merchant.PreferencesId)); err != nil {
				utils.Log.Errorf("func EnableHookStatus, InvalidatePreference fail, err [%v]", err)
			}
			// 如果相应的开关关掉/打开，则将merchant从相应redis key中删除/增加
			if err := UpdateMerchantWorkMode(int(merchantID), 1, utils.RedisKeyMerchantAlipayHookStatus()); err != nil {
				utils.Log.Errorf("func EnableHookStatus, update preferences Redis for merchant(uid=[%d]) fail. [%v]", merchantID, err)
			}
		} else if pref.AlipayHookStatus == 1 {
			// 目前已经是状态可用，什么都不用做
			utils.Log.Debugf("func EnableHookStatus do nothing, alipay_hook_status is already = 1 for merchant %d.", merchantID)
			return
		} else {
			utils.Log.Errorf("func EnableHookStatus, alipay_hook_status %d from db is not expected", pref.AlipayHookStatus)
			return
		}
	} else {
		utils.Log.Errorf("func EnableHookStatus, pay_type %d is not expected", payType)
		return
	}
}

// 设置币商的hook状态为不可用，先读取数据库当前值，如果本来就不可用，则什么都不做
func DisableHookStatus(merchantID int64, payType uint) {
	var merchant models.Merchant
	if err := dbcache.GetMerchantById(merchantID, &merchant); err != nil {
		utils.Log.Errorf("call GetMerchantById fail. [%v]", merchantID, err)
		utils.Log.Errorf("func DisableHookStatus finished abnormally. error %s", err)
		return
	}

	var pref models.Preferences
	if err := dbcache.GetPreferenceById(int64(merchant.PreferencesId), &pref); err != nil {
		utils.Log.Errorf("can't find preference record in db for merchant(uid=[%d]),  err [%v]", merchantID, err)
		utils.Log.Errorf("func DisableHookStatus finished abnormally. error %s", err)
		return
	}

	if payType == models.PaymentTypeWeixin {
		if pref.WechatHookStatus == 0 {
			// 目前已经是状态不可用，什么都不用做
			utils.Log.Debugf("func DisableHookStatus do nothing, wechat_hook_status is already = 0 for merchant %d.", merchantID)
			return
		} else if pref.WechatHookStatus == 1 {
			// 修改preferences表
			if err := utils.DB.Table("preferences").Where("id = ?", merchant.PreferencesId).Update("wechat_hook_status", 0).Error; err != nil {
				utils.Log.Errorf("func DisableHookStatus, update preferences for merchant(uid=[%d]) fail. [%v]", merchantID, err)
				return
			}
			// 使redis缓存失效
			if err := dbcache.InvalidatePreference(int64(merchant.PreferencesId)); err != nil {
				utils.Log.Errorf("func DisableHookStatus, InvalidatePreference fail, err [%v]", err)
				return
			}
			// 如果相应的开关关掉/打开，则将merchant从相应redis key中删除/增加
			if err := UpdateMerchantWorkMode(int(merchantID), 1, utils.RedisKeyMerchantWechatHookStatus()); err != nil {
				utils.Log.Errorf("func DisableHookStatus, update preferences Redis for merchant(uid=[%d]) fail. [%v]", merchantID, err)
				return
			}
		} else {
			utils.Log.Errorf("func DisableHookStatus, wechat_hook_status %d from db is not expected", pref.AlipayHookStatus)
			return
		}
	} else if payType == models.PaymentTypeAlipay {
		if pref.AlipayHookStatus == 0 {
			// 目前已经是状态不可用，什么都不用做
			utils.Log.Debugf("func DisableHookStatus do nothing, alipay_hook_status is already = 0 for merchant %d.", merchantID)
			return
		} else if pref.AlipayHookStatus == 1 {
			// 修改preferences表
			if err := utils.DB.Table("preferences").Where("id = ?", merchant.PreferencesId).Update("alipay_hook_status", 0).Error; err != nil {
				utils.Log.Errorf("func DisableHookStatus, update preferences for merchant(uid=[%d]) fail. [%v]", merchantID, err)
				return
			}
			// 使redis缓存失效
			if err := dbcache.InvalidatePreference(int64(merchant.PreferencesId)); err != nil {
				utils.Log.Errorf("func DisableHookStatus, InvalidatePreference fail, err [%v]", err)
				return
			}
			// 如果相应的开关关掉/打开，则将merchant从相应redis key中删除/增加
			if err := UpdateMerchantWorkMode(int(merchantID), 1, utils.RedisKeyMerchantAlipayHookStatus()); err != nil {
				utils.Log.Errorf("func DisableHookStatus, update preferences Redis for merchant(uid=[%d]) fail. [%v]", merchantID, err)
				return
			}
		} else {
			utils.Log.Errorf("func DisableHookStatus, alipay_hook_status %d from db is not expected", pref.AlipayHookStatus)
			return
		}
	} else {
		utils.Log.Errorf("func DisableHookStatus, pay_type %d is not expected", payType)
		return
	}
}
