package service

import (
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

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

func AddMerchant(phone string, email string) int {
	//  defer utils.Exit(utils.Enter("$FN(%s, %s)", phone, email))
	var ret int = -1

	user := models.Merchant{
		Phone:phone,
		Email:email,
	}
	if err := utils.DB.Create(&user).Error; err != nil {
		utils.Log.Error(err)
		// fmt.Printf("%v\n", err)
	} else {
		ret = user.Id
	}

	return ret
}