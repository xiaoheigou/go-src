package service

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/service/dbcache"
	"yuudidi.com/pkg/utils"
)

func GetMerchants(page, size, userStatus, userCert, startTime, stopTime, timeField, sort, search string) response.PageResponse {
	var ret response.PageResponse
	var result []models.Merchant
	db := utils.DB.Model(&models.Merchant{}).Select("merchants.*,assets.quantity as quantity, assets.qty_frozen as qty_frozen").Joins("left join assets on merchants.id = assets.merchant_id")
	if search != "" {
		db = db.Where(" merchants.phone like ? OR merchants.email like ? OR merchants.nickname like ?", search+"%", search+"%", search+"%")
	} else {
		db = db.Order(fmt.Sprintf("merchants.%s %s", timeField, sort))
		if startTime != "" && stopTime != "" {
			db = db.Where(fmt.Sprintf("merchants.%s >= ? AND merchants.%s <= ?", timeField, timeField), startTime, stopTime)
		}
		if userStatus != "" {
			db = db.Where("merchants.user_status = ?", userStatus)
		}
		if userCert != "" {
			db = db.Where("merchants.user_cert = ?", userCert)
		}
		db.Count(&ret.TotalCount)
		pageNum, err := strconv.ParseInt(page, 10, 64)
		pageSize, err1 := strconv.ParseInt(size, 10, 64)
		if err != nil || err1 != nil {
			utils.Log.Error(pageNum, pageSize)
		}
		db = db.Offset((pageNum - 1) * pageSize).Limit(pageSize)
		ret.PageNum = int(pageNum)
		ret.PageSize = int(pageSize)
	}
	db.Find(&result)
	ret.PageCount = len(result)
	ret.Status = response.StatusSucc
	ret.Data = result
	return ret
}

func GetMerchant(uid string) response.EntityResponse {
	var ret response.EntityResponse
	var merchants []models.Merchant

	if err := utils.DB.Select("merchants.*,assets.quantity as quantity").Joins("left join assets on merchants.id = assets.merchant_id").Find(&merchants, " merchants.id = ?", uid).Error; err != nil {
		utils.Log.Warnf("not found merchant")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundMerchant.Data()
		return ret
	}
	ret.Status = response.StatusSucc
	ret.Data = merchants
	return ret
}

func IsMerchantPhoneRegistered(nationCode int, phone string) (bool, error) {
	var user models.Merchant
	if err := utils.DB.First(&user, "nation_code = ? and phone = ?", nationCode, phone).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// 数据库找不到手机号相关记录，则没有注册过
			return false, nil
		} else {
			// 其它错误
			utils.Log.Errorf("database access err = [%v]", err)
			return false, err
		}
	}

	// 找到了记录，说明已经注册
	return true, nil
}

func AddMerchant(arg response.RegisterArg) response.RegisterRet {
	var ret response.RegisterRet

	// 检验参数
	phone := arg.Phone
	nationCode := arg.NationCode
	email := arg.Email
	passwordPlain := arg.Password
	if !utils.IsValidNationCode(nationCode) {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrNationCodeInvalid.Data()
		return ret
	}
	if !utils.IsValidPhone(nationCode, phone) {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrPhoneInvalid.Data()
		return ret
	}
	if !utils.IsValidEmail(email) {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrEmailInvalid.Data()
		return ret
	}

	// 再次验证手机随机码
	key := "app:register:" + strconv.Itoa(nationCode) + ":" + phone // example: "app:register:86:13100000000"
	value := strconv.Itoa(arg.PhoneRandomCodeSeq) + ":" + arg.PhoneRandomCode
	if err := utils.RedisVerifyValue(key, value); err != nil {
		utils.Log.Errorf("registering %s, re-verify sms random code fail. %s", phone, err.Error())
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrRegisterRandomCodeReVerifyFail.Data()
		return ret
	}

	// 再次验证邮箱随机码
	var skipEmailVerify bool
	// 邮箱服务器有时会限制，密集发送邮件可能失败。这里检测是否配置了跳过mail验证。
	var err error
	if skipEmailVerify, err = strconv.ParseBool(utils.Config.GetString("register.skipemailverify")); err != nil {
		utils.Log.Errorf("Wrong configuration: register.skipemailverify, should be boolean. Set to default false.")
		skipEmailVerify = false
	}

	if !skipEmailVerify {
		key = "app:register:" + email
		value = strconv.Itoa(arg.EmailRandomCodeSeq) + ":" + arg.EmailRandomCode
		if err := utils.RedisVerifyValue(key, value); err != nil {
			utils.Log.Errorf("registering %s, re-verify email random code fail. %s", email, err.Error())
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrRegisterRandomCodeReVerifyFail.Data()
			return ret
		}
	}

	// 检测手机号是否已经注册过
	var registered bool
	if registered, err = IsMerchantPhoneRegistered(nationCode, phone); err != nil {
		utils.Log.Errorf("database access err = [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	if registered == true {
		// 手机号已经注册过
		utils.Log.Errorf("phone [%v] nation_code [%v] is already registered.", phone, nationCode)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrPhoneAlreadyRegister.Data()
		return ret
	}

	// 不校验邮箱是否已经注册过

	algorithm := utils.Config.GetString("algorithm")
	if len(algorithm) == 0 {
		utils.Log.Errorf("Wrong configuration: algorithm [%v], it's empty. Set to default Argon2.", algorithm)
		algorithm = "Argon2"
	}

	salt, err := generateRandomBytes(64)
	if err != nil {
		utils.Log.Errorf("AddMerchant, err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}
	hashFunc := functionMap[algorithm]
	passwordEncrypted := hashFunc([]byte(passwordPlain), salt)

	// 表Merchant和Preferences是"一对一"的表
	var pref models.Preferences
	pref.InWork = 0
	pref.WechatHookStatus = 0
	pref.AlipayHookStatus = 0
	pref.WechatAutoOrder = 0
	pref.AlipayAutoOrder = 0

	tx := utils.DB.Begin()
	if err := tx.Model(&models.Preferences{}).Create(&pref).Error; err != nil {
		utils.Log.Errorf("create preference for user [%v] fail. err = [%v]", phone, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		tx.Rollback()
		return ret
	}
	var user models.Merchant
	user = models.Merchant{
		Phone:      phone,
		Email:      email,
		NationCode: nationCode,
		Salt:       salt,
		LastLogin:  time.Now(),
		Password:   passwordEncrypted,
		Algorithm:  algorithm,
		// Preferences: pref,
		PreferencesId: uint64(pref.Id),
	}
	if err := tx.Model(&models.Merchant{}).Create(&user).Error; err != nil {
		utils.Log.Errorf("AddMerchant, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		tx.Rollback()
		return ret
	}
	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("AddMerchant commit fail, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	// 注册后，删除以前redis中记录的登录失败次数
	// 用户在注册前，可能尝试登录（用户忘记自己是否已经注册），这样系统中会记录下登录失败次数。注册成功后，需要删除这些记录，否则会影响新注册账号的登录
	utils.ClearAppLoginFailTimes(arg.NationCode, arg.Phone)

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.RegisterData{
		Uid: user.Id,
	})
	return ret
}

func GetMerchantAuditStatus(uid int64) response.GetAuditStatusRet {
	var ret response.GetAuditStatusRet

	var merchant models.Merchant
	if err := dbcache.GetMerchantById(uid, &merchant); err != nil {
		utils.Log.Errorf("GetMerchantAuditStatus, find merchant(uid=[%d]) fail. [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	// user_status可以为0/1/2/3，分别表示“待审核/正常/未通过审核/冻结”
	userStatus := merchant.UserStatus
	if userStatus == 0 || userStatus == 1 {
		// 当状态为0/1（“待审核/正常”）时，不查AuditMessage对应的表了
		ret.Status = response.StatusSucc
		ret.Data = append(ret.Data, response.GetAuditStatusData{
			UserStatus:   userStatus,
			ContactPhone: "",
			ExtraMessage: "",
		})
		return ret
	} else {
		// 当状态为2/3（“未通过审核/冻结”）时，查AuditMessage对应的表
		var audit models.AuditMessage
		if err := utils.DB.Where("merchant_id = ?", uid).Find(&audit).Error; err != nil {
			utils.Log.Errorf("GetMerchantAuditStatus, find AuditMessage for merchant(uid=[%d]) fail. [%v]", uid, err)
			// 查不到联系信息等也没必要报错给前端，返回空即可
		}
		ret.Status = response.StatusSucc
		ret.Data = append(ret.Data, response.GetAuditStatusData{
			UserStatus:   userStatus,
			ContactPhone: audit.ContactPhone,
			ExtraMessage: audit.ExtraMessage,
		})
		return ret
	}
}

func GetMerchantProfile(uid int64) response.GetProfileRet {
	var ret response.GetProfileRet

	var merchant models.Merchant
	if err := dbcache.GetMerchantById(uid, &merchant); err != nil {
		utils.Log.Errorf("GetMerchantProfile, find merchant(uid=[%d]) fail. [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	nickname := merchant.Nickname

	var assets []models.Assets
	if err := utils.DB.Where(&models.Assets{MerchantId: int64(uid)}).Find(&assets).Error; err != nil {
		utils.Log.Errorf("GetMerchantProfile, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	} else {
		if len(assets) == 0 {
			utils.Log.Errorf("GetMerchantProfile, can't find assets for merchant(uid=[%d]).", uid)
			// 查不到没必要报错给前端，返回空即可
			ret.Status = response.StatusSucc
			ret.Data = append(ret.Data, response.GetProfileData{
				NickName:       nickname,
				CurrencyCrypto: "BTUSD",
				Quantity:       0,
				QtyFrozen:      0,
			})
		} else {
			ret.Status = response.StatusSucc
			for _, asset := range assets {
				quantityFloat, _ := asset.Quantity.Float64()
				qtyFrozenFloat, _ := asset.QtyFrozen.Float64()
				ret.Data = append(ret.Data, response.GetProfileData{
					NickName:       nickname,
					CurrencyCrypto: asset.CurrencyCrypto,
					Quantity:       quantityFloat,
					QtyFrozen:      qtyFrozenFloat,
				})
			}
		}
		return ret
	}
}

func SetMerchantNickname(uid int64, arg response.SetNickNameArg) response.SetNickNameRet {
	var ret response.SetNickNameRet

	// 检验参数
	nickname := arg.NickName
	if len(nickname) > 20 {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrNicknameTooLong.Data()
		return ret
	}

	// 修改nickname字段为指定值
	if err := utils.DB.Table("merchants").Where("id = ?", uid).Update("nickname", nickname).Error; err != nil {
		utils.Log.Errorf("SetMerchantNickname, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}
	// 使缓存失效
	if err := dbcache.InvalidateMerchant(uid); err != nil {
		utils.Log.Warnf("SetMerchantNickname, db err [%v]", err)
	}

	ret.Status = response.StatusSucc
	return ret
}

func GetMerchantWorkMode(uid int) response.GetWorkModeRet {
	var ret response.GetWorkModeRet

	var merchant models.Merchant
	if err := dbcache.GetMerchantById(int64(uid), &merchant); err != nil {
		utils.Log.Errorf("GetMerchantWorkMode, find merchant(uid=[%d]) fail. [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	var pref models.Preferences
	if err := dbcache.GetPreferenceById(int64(merchant.PreferencesId), &pref); err != nil {
		utils.Log.Errorf("GetMerchantWorkMode, can't find preference record in db for merchant(uid=[%d]),  err [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.GetWorkModeData{
		InWork:           pref.InWork,
		WechatHookStatus: pref.WechatHookStatus,
		AlipayHookStatus: pref.AlipayHookStatus,
		WechatAutoOrder:  pref.WechatAutoOrder,
		AlipayAutoOrder:  pref.AlipayAutoOrder,
	})
	return ret
}

func SetMerchantWorkMode(uid int, arg response.SetWorkModeArg) response.SetWorkModeRet {
	var ret response.SetWorkModeRet

	inWork := arg.InWork
	wechatHookStatus := arg.WechatHookStatus
	alipayHookStatus := arg.AlipayHookStatus
	wechatAutoOrder := arg.WechatAutoOrder
	alipayAutoOrder := arg.AlipayAutoOrder
	if !(inWork == 0 || inWork == 1 || inWork == -1) { // -1 表示不进行修改，下同
		var ret response.SetWorkModeRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
	}
	if !(wechatHookStatus == 0 || wechatHookStatus == 1 || wechatHookStatus == -1) {
		var ret response.SetWorkModeRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
	}
	if !(alipayHookStatus == 0 || alipayHookStatus == 1 || alipayHookStatus == -1) {
		var ret response.SetWorkModeRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
	}
	if !(wechatAutoOrder == 0 || wechatAutoOrder == 1 || wechatAutoOrder == -1) {
		var ret response.SetWorkModeRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
	}
	if !(alipayAutoOrder == 0 || alipayAutoOrder == 1 || alipayAutoOrder == -1) {
		var ret response.SetWorkModeRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
	}
	if inWork == -1 && wechatHookStatus == -1 && alipayHookStatus == -1 && wechatAutoOrder == -1 && alipayAutoOrder == -1 {
		var ret response.SetWorkModeRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
	}

	var merchant models.Merchant
	if err := utils.DB.Where("id = ?", uid).Find(&merchant).Error; err != nil {
		utils.Log.Errorf("SetMerchantWorkMode, find merchant(uid=[%d]) fail. [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}
	changeParam := make(map[string]interface{})
	if inWork == 0 || inWork == 1 {
		changeParam["in_work"] = inWork
	}
	if wechatHookStatus == 0 || wechatHookStatus == 1 {
		changeParam["wechat_hook_status"] = wechatHookStatus
	}
	if alipayHookStatus == 0 || alipayHookStatus == 1 {
		changeParam["alipay_hook_status"] = alipayHookStatus
	}
	if wechatAutoOrder == 0 || wechatAutoOrder == 1 {
		changeParam["wechat_auto_order"] = wechatAutoOrder
	}
	if alipayAutoOrder == 0 || alipayAutoOrder == 1 {
		changeParam["alipay_auto_order"] = alipayAutoOrder
	}
	// 修改preferences表
	if err := utils.DB.Table("preferences").Where("id = ?", merchant.PreferencesId).Updates(changeParam).Error; err != nil {
		utils.Log.Errorf("SetMerchantWorkMode, update preferences for merchant(uid=[%d]) fail. [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}
	// 使缓存失效
	if err := dbcache.InvalidatePreference(int64(merchant.PreferencesId)); err != nil {
		utils.Log.Errorf("InvalidatePreference fail, err [%v]", err)
	}

	// 如果相应的开关关掉/打开，则将merchant从相应redis key中删除/增加
	if err := UpdateMerchantWorkMode(uid, inWork, utils.RedisKeyMerchantInWork()); err != nil {
		utils.Log.Errorf("SetMerchantWorkMode, update preferences Redis for merchant(uid=[%d]) fail. [%v]", uid, err)
	}
	if err := UpdateMerchantWorkMode(uid, wechatHookStatus, utils.RedisKeyMerchantWechatHookStatus()); err != nil {
		utils.Log.Errorf("SetMerchantWorkMode, update preferences Redis for merchant(uid=[%d]) fail. [%v]", uid, err)
	}
	if err := UpdateMerchantWorkMode(uid, alipayHookStatus, utils.RedisKeyMerchantAlipayHookStatus()); err != nil {
		utils.Log.Errorf("SetMerchantWorkMode, update preferences Redis for merchant(uid=[%d]) fail. [%v]", uid, err)
	}
	if err := UpdateMerchantWorkMode(uid, wechatAutoOrder, utils.RedisKeyMerchantWechatAutoOrder()); err != nil {
		utils.Log.Errorf("SetMerchantWorkMode, update preferences Redis for merchant(uid=[%d]) fail. [%v]", uid, err)
	}
	if err := UpdateMerchantWorkMode(uid, alipayAutoOrder, utils.RedisKeyMerchantAlipayAutoOrder()); err != nil {
		utils.Log.Errorf("SetMerchantWorkMode, update preferences Redis for merchant(uid=[%d]) fail. [%v]", uid, err)
	}
	ret.Status = response.StatusSucc
	return ret
}

func FreezeMerchant(uid string, args response.FreezeArgs) response.EntityResponse {
	if args.Operation == 1 {
		return UpdateMerchantStatus(uid, args.ContactPhone, args.ExtraMessage, 3)
	} else if args.Operation == 0 {
		return UpdateMerchantStatus(uid, args.ContactPhone, args.ExtraMessage, 1)
	} else {
		var ret response.EntityResponse
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		return ret
	}
}

func ApproveMerchant(uid string, args response.ApproveArgs) response.EntityResponse {
	if args.Operation == 1 {
		return UpdateMerchantStatus(uid, args.ContactPhone, args.ExtraMessage, 1)
	} else if args.Operation == 0 {
		return UpdateMerchantStatus(uid, args.ContactPhone, args.ExtraMessage, 2)
	} else {
		var ret response.EntityResponse
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		return ret
	}
}

func UpdateMerchantStatus(merchantId, phone, msg string, userStatus int) response.EntityResponse {
	var ret response.EntityResponse
	var merchant models.Merchant
	ret.Status = response.StatusSucc
	if err := utils.DB.First(&merchant, "id = ?", merchantId).Error; err != nil {
		utils.Log.Warnf("not found merchant")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundMerchant.Data()
		return ret
	}
	if merchant.UserStatus == userStatus {
		ret.Data = []models.Merchant{merchant}
		return ret
	}
	tx := utils.DB.Begin()
	id, _ := strconv.ParseInt(merchantId, 10, 64)
	switch userStatus {
	case 1:
		if err := tx.Delete(&models.AuditMessage{}, "merchant_id = ?", merchantId).Error; err != nil {
			utils.Log.Errorf("delete audit message is failed,uid:%s,userStatus:%v", merchantId, userStatus)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.UpdateMerchantStatusErr.Data()
			tx.Rollback()
			return ret
		}
		if err := tx.Model(&models.Merchant{}).Where("id = ?", merchantId).Update("user_status", userStatus).Error; err != nil {
			utils.Log.Errorf("update merchant status is failed,uid:%s,userStatus:%v", merchantId, userStatus)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.UpdateMerchantStatusErr.Data()
			tx.Rollback()
			return ret
		}
	case 2:
		fallthrough
	case 3:
		audit := models.AuditMessage{
			MerchantId:   id,
			ContactPhone: phone,
			ExtraMessage: msg,
			OperatorId:   1,
		}
		if err := tx.Model(&models.AuditMessage{}).Create(&audit).Error; err != nil {
			utils.Log.Errorf("create audit message status is failed,uid:%s,userStatus:%v", merchantId, userStatus)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.UpdateMerchantStatusErr.Data()
			tx.Rollback()
			return ret
		}
		if err := tx.Model(&models.Merchant{}).Where("id = ?", merchantId).Update("user_status", userStatus).Error; err != nil {
			utils.Log.Errorf("update merchant status is failed,uid:%s,userStatus:%v", merchantId, userStatus)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.UpdateMerchantStatusErr.Data()
			tx.Rollback()
			return ret
		}
	}
	tx.Commit()
	// 修改了Merchant表，使对应缓存失效
	if err := dbcache.InvalidateMerchant(id); err != nil {
		utils.Log.Infof("InvalidateMerchant fail, err [%v]", err)
	}
	ret.Status = response.StatusSucc
	return ret
}

//GetMerchantsQualified - return mock data
func GetMerchantsQualified(orderNumber string, amount float64, quantity decimal.Decimal, currencyCrypto string, payType uint, fix bool, autoOrder bool, limit, direction uint8) []int64 {
	var merchantIds []int64
	var assetMerchantIds []int64
	var paymentMerchantIds []int64
	result := []int64{}
	var tempIds []string
	if autoOrder {
		if payType == models.PaymentTypeWeixin {
			if err := utils.GetCacheSetInterMembers(&tempIds,
				utils.RedisKeyMerchantOnline(),
				utils.RedisKeyMerchantInWork(),
				utils.RedisKeyMerchantWechatAutoOrder(),
				utils.RedisKeyMerchantWechatHookStatus(),
			); err != nil {
				utils.Log.Errorf("get intersection merchants fail, [%v]", tempIds)
				return result
			}

			utils.Log.Debugf("for order %s, merchants (in_work = 1 and wechat_auto_order = 1 and wechat_hook_status = 1) %s", orderNumber, tempIds)
		} else if payType == models.PaymentTypeAlipay {
			if err := utils.GetCacheSetInterMembers(&tempIds,
				utils.RedisKeyMerchantOnline(),
				utils.RedisKeyMerchantInWork(),
				utils.RedisKeyMerchantAlipayAutoOrder(),
				utils.RedisKeyMerchantAlipayHookStatus(),
			); err != nil {
				utils.Log.Errorf("get intersection merchants fail, [%v]", tempIds)
				return result
			}

			utils.Log.Debugf("for order %s, merchants (in_work = 1 and alipay_auto_order = 1 and alipay_hook_status = 1) %s", orderNumber, tempIds)
		} else {
			utils.Log.Errorf("payType %d is not expected, auto order only applicable for wechat and alipay", payType)
			return result
		}
	} else {
		//获取承兑商在线列表
		if err := utils.GetCacheSetInterMembers(&tempIds,
			utils.RedisKeyMerchantOnline(),
			utils.RedisKeyMerchantInWork()); err != nil {
			utils.Log.Errorf("get intersection merchants fail, [%v]", tempIds)
			return result
		}
	}

	if err := utils.ConvertStringToInt(tempIds, &merchantIds); err != nil {
		utils.Log.Errorf("convert string to int is failed,merchantIds = %v,err= %v ", merchantIds, err)
		return result
	}

	if len(merchantIds) == 0 {
		utils.Log.Debugf("func GetMerchantsQualified finished, order %s, merchants: %v", orderNumber, merchantIds)
		return result
	}

	//充值单过滤资产，提现单不需要过滤
	if direction == 0 {
		//查询资产符合情况的币商列表
		db := utils.DB.Model(&models.Assets{}).Where("currency_crypto = ? AND quantity >= ? AND merchant_id > 0", currencyCrypto, quantity)
		if err := db.Pluck("merchant_id", &assetMerchantIds).Error; err != nil {
			utils.Log.Errorf("Gets merchants that have enough btusd failed. err = %s", err)
			return result
		}
		utils.Log.Debugf("for order %s, merchants have enough assets: %v", orderNumber, assetMerchantIds)
	}

	//通过支付方式过滤
	db := utils.DB.Model(&models.PaymentInfo{}).Where("in_use = ? AND audit_status = ?", 0, 1)

	switch {
	case payType == 1:
		fallthrough
	case payType == 2:
		//fix - 是否只查询具有固定支付金额对应（支付宝，微信）二维码的币商
		//true - 只查询固定支付金额二维码（支付宝，微信）如果支付方式为银行卡，匹配相同银行
		//false - 查询所有支付方式（即只要支付方式满足即可）,如果支付方式为银行卡，匹配所有的银行
		if fix {
			db = db.Where("e_amount = ?", amount)
		} else {
			//0表示非固定金额
			db = db.Where("e_amount = ?", 0)
		}
		db = db.Where("pay_type = ?", payType)
	case payType >= 4:
		//两轮查询，如果是固定金额，匹配银行相同的，如是不是，随意找一个
		if fix {
			db = db.Where("pay_type = ?", payType)
		} else {
			db = db.Where("pay_type >= ?", 4)
		}
	default:
		return result
	}

	if autoOrder {
		db = db.Where("payment_auto_type = 1") // 仅查询自动收款账号
	} else {
		db = db.Where("payment_auto_type = 0") // 仅查询手动收款账号
	}

	if err := db.Pluck("uid", &paymentMerchantIds).Error; err != nil {
		utils.Log.Errorf("Gets a list of payment conformance is failed.")
		return result
	}

	utils.Log.Debugf("for order %s, merchants can match payment: %v", orderNumber, paymentMerchantIds)

	if direction == 0 {
		//
		merchantIds = utils.IntersectList(merchantIds, assetMerchantIds, paymentMerchantIds)
	} else if direction == 1 {
		//
		merchantIds = utils.IntersectList(merchantIds, paymentMerchantIds)
	} else {
		return result
	}

	//限制返回条数 0 代表全部返回
	utils.Log.Debugf("func GetMerchantsQualified finished, order %s, merchants: %v", orderNumber, merchantIds)
	if limit == 0 {
		return merchantIds
	} else if limit > 0 {
		return merchantIds[0:limit]
	}
	return result
}
