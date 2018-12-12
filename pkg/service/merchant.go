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

	user := models.Merchant{
		Phone:phone,
		Email:email,
		NationCode:nationCode,
		Salt:salt,
		Password:passwordEncrypted,
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