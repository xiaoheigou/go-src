package service

import (
	"fmt"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

func GetMerchants(page,size,userStatus,userCert,startTime,stopTime,timeField,sort,search string) response.PageResponse {
	var ret response.PageResponse
	var result []models.Merchant
	db := utils.DB.Model(&models.Merchant{}).Select("merchants.*,assets.quantity as quantity").Joins("left join assets on merchants.id = assets.merchant_id")
	if search != "" {
		db = db.Where(" merchants.phone = ? OR merchants.email = ?",search,search)
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

	if err:= utils.DB.First(&merchant," id = ?" ,uid).Error;err != nil {
		utils.Log.Warnf("not found merchant")
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.NotFoundMerchant.Data()
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
		ret.ErrCode, ret.ErrMsg = err_code.AppErrNationCodeInvalid.Data()
		return ret
	}
	if ! utils.IsValidPhone(nationCode, phone) {
		ret.ErrCode, ret.ErrMsg = err_code.AppErrPhoneInvalid.Data()
		return ret
	}
	if ! utils.IsValidEmail(email) {
		ret.ErrCode, ret.ErrMsg = err_code.AppErrEmailInvalid.Data()
		return ret
	}

	// 应该再次验证手机和邮件的随机码
	// TODO

	// 检测手机号或者邮箱是否已经注册过
	var user models.Merchant
	if ! utils.DB.Where("phone = ? and nation_code = ?", phone, nationCode).First(&user).RecordNotFound() {
		// 手机号已经注册过
		utils.Log.Errorf("phone [%v] nation_code [%v] is already registered.", phone, nationCode)
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrPhoneAlreadyRegister.Data()
		return ret
	}
	if ! utils.DB.Where("email = ?", email).First(&user).RecordNotFound() {
		// 邮箱已经注册过
		utils.Log.Errorf("email [%v] is already registered.", email)
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrEmailAlreadyRegister.Data()
		return ret
	}

	algorithm := utils.Config.GetString("algorithm")
	if len(algorithm) == 0 {
		utils.Log.Errorf("Wrong configuration: algorithm [%v], it's empty. Set to default Argon2.", algorithm)
		algorithm = "Argon2"
	}

	salt, err := generateRandomBytes(64)
	if err != nil {
		utils.Log.Errorf("AddMerchant, err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}
	hashFunc := functionMap[algorithm]
	passwordEncrypted := hashFunc([]byte(passwordPlain), salt)

	// 表Merchant和Preferences是"一对一"的表
	var pref models.Preferences
	pref.TakeOrder  = 1    // 默认接单
	pref.AutoOrder = 0    // 默认不自动接单
	user = models.Merchant{
		Phone:phone,
		Email:email,
		NationCode:nationCode,
		Salt:salt,
		Password:passwordEncrypted,
		Algorithm:algorithm,
		Preferences:pref,
	}
	if err := utils.DB.Create(&user).Error; err != nil {
		utils.Log.Errorf("AddMerchant, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.RegisterData{
		Uid: user.Id,
	})
	return ret
}

func SetMerchantNickname(uid int, arg response.SetNickNameArg) response.SetNickNameRet {
	var ret response.SetNickNameRet

	// 检验参数
	nickname := arg.NickName
	if len(nickname) > 20 {
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrNicknameTooLong.Data()
		return ret
	}

	// 修改nickname字段为指定值
	if err := utils.DB.Table("merchants").Where("id = ?", uid).Update("nickname", nickname).Error; err != nil {
		utils.Log.Errorf("SetMerchantNickname, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
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
		ret.ErrCode,ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}
	if err := utils.DB.Model(&merchant).Related(&merchant.Preferences).Error; err != nil {
		utils.Log.Errorf("GetMerchantWorkMode, can't find preference record in db for merchant(uid=[%d]),  err [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
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
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
	}
	if ! (accept == 0 || accept == 1) {
		var ret response.SetWorkModeRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
	}

	var merchant models.Merchant
	if err := utils.DB.Where("id = ?", uid).Find(&merchant).Error; err != nil {
		utils.Log.Errorf("SetMerchantWorkMode, find merchant(uid=[%d]) fail. [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}
	if err := utils.DB.Table("preferences").Where("id = ?", merchant.PreferencesId).Updates(
		map[string]interface{}{
			"take_order": accept,
			"auto_order": auto,
		}).Error; err != nil {
		utils.Log.Errorf("SetMerchantWorkMode, update preferences for merchant(uid=[%d]) fail. [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	return ret
}