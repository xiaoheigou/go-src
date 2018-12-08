package service

import (
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

func GetMerchant(uid string) {

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