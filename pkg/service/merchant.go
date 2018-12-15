package service

import (
	"fmt"
	"strconv"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

func GetMerchants(page, size, userStatus, userCert, startTime, stopTime, timeField, sort, search string) response.PageResponse {
	var ret response.PageResponse
	var result []models.Merchant
	db := utils.DB.Model(&models.Merchant{}).Select("merchants.*,assets.quantity as quantity").Joins("left join assets on merchants.id = assets.merchant_id")
	if search != "" {
		db = db.Where(" merchants.phone = ? OR merchants.email = ?", search, search)
	} else {
		db = db.Order(fmt.Sprintf("merchants.%s %s", timeField, sort))
		pageNum, err := strconv.ParseInt(page, 10, 64)
		pageSize, err1 := strconv.ParseInt(size, 10, 64)
		if err != nil || err1 != nil {
			utils.Log.Error(pageNum, pageSize)
		}
		db = db.Offset(pageNum * pageSize).Limit(pageSize)
		if startTime != "" && stopTime != "" {
			db = db.Where(fmt.Sprintf("merchants.%s >= ? AND merchants.%s <= ?", timeField, timeField), startTime, stopTime)
		}
		if userStatus != "" {
			db = db.Where("merchants.user_status = ?", userStatus)
		}
		if userCert != "" {
			db = db.Where("merchants.user_cert = ?", userCert)
		}
		db.Count(&ret.PageCount)
		ret.PageNum = int(pageNum + 1)
		ret.PageSize = int(pageSize)
	}
	db.Find(&result)
	ret.Status = response.StatusSucc
	ret.Data = result
	return ret
}

func GetMerchant(uid string) response.EntityResponse {
	var ret response.EntityResponse
	var merchant models.Merchant

	if err := utils.DB.First(&merchant, " id = ?", uid).Error; err != nil {
		utils.Log.Warnf("not found merchant")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundMerchant.Data()
		return ret
	}
	ret.Status = response.StatusSucc
	ret.Data = []models.Merchant{merchant}
	return ret
}

func AddMerchant(arg response.RegisterArg) response.RegisterRet {
	var ret response.RegisterRet

	// 检验参数
	phone := arg.Phone
	nationCode := arg.NationCode
	email := arg.Email
	passwordPlain := arg.Password
	if ! utils.IsValidNationCode(nationCode) {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrNationCodeInvalid.Data()
		return ret
	}
	if ! utils.IsValidPhone(nationCode, phone) {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrPhoneInvalid.Data()
		return ret
	}
	if ! utils.IsValidEmail(email) {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrEmailInvalid.Data()
		return ret
	}

	// 再次验证手机随机码
	key := "app:register:" + strconv.Itoa(nationCode) + ":" + phone // example: "app:register:86:13100000000"
	value := strconv.Itoa(arg.PhoneRandomCodeSeq) + ":" + arg.PhoneRandomCode
	if err := utils.RedisVerifyValue(key, value); err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrRandomCodeVerifyFail.Data()
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

	if ! skipEmailVerify {
		key = "app:register:" + email
		value = strconv.Itoa(arg.EmailRandomCodeSeq) + ":" + arg.EmailRandomCode
		if err := utils.RedisVerifyValue(key, value); err != nil {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrRandomCodeVerifyFail.Data()
			return ret
		}
	}

	// 检测手机号或者邮箱是否已经注册过
	var user models.Merchant
	if ! utils.DB.Where("phone = ? and nation_code = ?", phone, nationCode).First(&user).RecordNotFound() {
		// 手机号已经注册过
		utils.Log.Errorf("phone [%v] nation_code [%v] is already registered.", phone, nationCode)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrPhoneAlreadyRegister.Data()
		return ret
	}

	//// 不用校验邮箱是否已经注册过
	//if ! utils.DB.Where("email = ?", email).First(&user).RecordNotFound() {
	//	// 邮箱已经注册过
	//	utils.Log.Errorf("email [%v] is already registered.", email)
	//	ret.Status = response.StatusFail
	//	ret.ErrCode, ret.ErrMsg = err_code.AppErrEmailAlreadyRegister.Data()
	//	return ret
	//}

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
	pref.TakeOrder = 1 // 默认接单
	pref.AutoOrder = 0 // 默认不自动接单
	user = models.Merchant{
		Phone:       phone,
		Email:       email,
		NationCode:  nationCode,
		Salt:        salt,
		LastLogin:   time.Now(),
		Password:    passwordEncrypted,
		Algorithm:   algorithm,
		Preferences: pref,
	}
	if err := utils.DB.Create(&user).Error; err != nil {
		utils.Log.Errorf("AddMerchant, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.RegisterData{
		Uid: user.Id,
	})
	return ret
}

func GetMerchantAuditStatus(uid int) response.GetAuditStatusRet {
	var ret response.GetAuditStatusRet

	var merchant models.Merchant
	if err := utils.DB.Where("id = ?", uid).Find(&merchant).Error; err != nil {
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

func GetMerchantProfile(uid int) response.GetProfileRet {
	var ret response.GetProfileRet

	var merchant models.Merchant
	if err := utils.DB.Where("id = ?", uid).Find(&merchant).Error; err != nil {
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
				ret.Data = append(ret.Data, response.GetProfileData{
					NickName:       nickname,
					CurrencyCrypto: asset.CurrencyCrypto,
					Quantity:       asset.Quantity,
					QtyFrozen:      asset.QtyFrozen,
				})
			}
		}
		return ret
	}
}

func SetMerchantNickname(uid int, arg response.SetNickNameArg) response.SetNickNameRet {
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

	ret.Status = response.StatusSucc
	return ret
}

func GetMerchantWorkMode(uid int) response.GetWorkModeRet {
	var ret response.GetWorkModeRet

	var merchant models.Merchant
	if err := utils.DB.Where("id = ?", uid).Find(&merchant).Error; err != nil {
		utils.Log.Errorf("GetMerchantWorkMode, find merchant(uid=[%d]) fail. [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}
	if err := utils.DB.Model(&merchant).Related(&merchant.Preferences).Error; err != nil {
		utils.Log.Errorf("GetMerchantWorkMode, can't find preference record in db for merchant(uid=[%d]),  err [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.GetWorkModeData{
		Accept: merchant.Preferences.TakeOrder,
		Auto:   merchant.Preferences.AutoOrder,
	})
	return ret
}

func SetMerchantWorkMode(uid int, arg response.SetWorkModeArg) response.SetWorkModeRet {
	var ret response.SetWorkModeRet

	auto := arg.Auto
	accept := arg.Accept
	if ! (auto == 0 || auto == 1) {
		var ret response.SetWorkModeRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
	}
	if ! (accept == 0 || accept == 1) {
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
	if err := utils.DB.Table("preferences").Where("id = ?", merchant.PreferencesId).Updates(
		map[string]interface{}{
			"take_order": accept,
			"auto_order": auto,
		}).Error; err != nil {
		utils.Log.Errorf("SetMerchantWorkMode, update preferences for merchant(uid=[%d]) fail. [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	return ret
}

func FreezeMerchant(uid string, args response.FreezeArgs) response.EntityResponse {
	if args.Operation == 1 {
		return updateMerchantStatus(uid, args.ContactPhone, args.ExtraMessage, 3)
	} else if args.Operation == 0 {
		return updateMerchantStatus(uid, args.ContactPhone, args.ExtraMessage, 1)
	} else {
		var ret response.EntityResponse
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		return ret
	}
}

func ApproveMerchant(uid string, args response.ApproveArgs) response.EntityResponse {
	if args.Operation == 1 {
		return updateMerchantStatus(uid, args.ContactPhone, args.ExtraMessage, 1)
	} else if args.Operation == 0 {
		return updateMerchantStatus(uid, args.ContactPhone, args.ExtraMessage, 2)
	} else {
		var ret response.EntityResponse
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.RequestParamErr.Data()
		return ret
	}
}

func updateMerchantStatus(merchantId, phone, msg string, userStatus int) response.EntityResponse {
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
		id, _ := strconv.ParseInt(merchantId, 10, 64)
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
	ret.Status = response.StatusSucc
	return ret
}
