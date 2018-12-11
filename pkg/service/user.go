package service

import (
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

func Login(param response.WebLoginArgs) response.EntityResponse {
	var ret response.EntityResponse
	var user models.User

	if err := utils.DB.First(&user, "username = ?", param.Username).Error; err != nil {
		utils.Log.Warnf("not found user,username:%s", param.Username)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundUser.Data()
		return ret
	}
	salt := user.Salt
	hashFunc := functionMap[user.Algorithm]
	hash := hashFunc([]byte(param.Password), salt)
	if compare(user.Password,hash) != 0 {
		utils.Log.Warnf("Invalid username/password set")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.UserPasswordError.Data()
		return ret
	}
	ret.Status = response.StatusSucc
	ret.Data = response.WebLoginResponse{
		Uid:      user.Id,
		Role:     user.Role,
		Username: user.Username,
	}
	return ret
}

func CreateUser(param response.UserArgs) response.EntityResponse {
	var ret response.EntityResponse

	algorithm := utils.Config.GetString("algorithm")

	user := models.User{
		Role:       1,
		Username:   param.Username,
		Phone:      param.Phone,
		Address:    param.Address,
		Email:      param.Email,
		UserStatus: 0,
		Algorithm:  algorithm,
	}

	salt, err := generateRandomBytes(64)
	if err != nil {
		utils.Log.Errorf("Unable to get random salt")
		panic(err)
	}
	user.Salt = salt
	hashFunc := functionMap[user.Algorithm]
	user.Password = hashFunc([]byte(param.Password), salt)

	ret.Status = response.StatusSucc
	if err := utils.DB.Create(&user).Error; err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateUserErr.Data()
	}
	return ret
}

func compare(a, b []byte) int {
	if len(a) != len(b) {
		return len(a) - len(b)
	}
	for idx := range a {
		if a[idx] == b[idx] {
			continue
		}
		return int(a[idx]) - int(b[idx])
	}
	return 0
}
