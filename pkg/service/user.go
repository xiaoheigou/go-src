package service

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/jinzhu/gorm"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

func Login(param response.WebLoginArgs, session sessions.Session) response.EntityResponse {
	var ret response.EntityResponse
	var user models.User

	if err := utils.DB.First(&user, "username = ?", param.Username).Error; err != nil {
		utils.Log.Warnf("not found user,username:%s", param.Username)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.LoginErr.Data()
		return ret
	}
	salt := user.Salt
	hashFunc := functionMap[user.Algorithm]
	hash := hashFunc([]byte(param.Password), salt)
	if compare(user.Password, hash) != 0 {
		utils.Log.Warnf("Invalid username/password set")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.LoginErr.Data()
		return ret
	}
	ret.Status = response.StatusSucc
	ret.Data = response.WebLoginResponse{
		Uid:      user.Id,
		Role:     user.Role,
		Username: user.Username,
	}
	session.Set("userId", user.Id)
	session.Set("userRole", user.Role)
	session.Set("username", user.Username)
	session.Options(sessions.Options{
		MaxAge: 3600,
		Path:   "/",
	})
	session.Save()
	return ret
}

func CreateUser(param response.UserArgs, tx *gorm.DB) response.EntityResponse {
	var ret response.EntityResponse

	if tx == nil {
		tx = utils.DB
	}
	algorithm := utils.Config.GetString("algorithm")

	user := models.User{
		Role:       param.Role,
		Username:   param.Username,
		Phone:      param.Phone,
		Address:    param.Address,
		Email:      param.Email,
		UserStatus: 0,
		Algorithm:  algorithm,
	}

	user.Salt, user.Password = generatePassword(param.Password, algorithm)
	if user.Salt == nil {
		utils.Log.Errorf("Unable to get random salt")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateUserErr.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	if err := tx.Create(&user).Error; err != nil {
		utils.Log.Errorf("Unable to get random salt")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateUserErr.Data()
	}
	return ret
}

func GetAgent(uid string) response.EntityResponse {
	var ret response.EntityResponse
	var agent models.User

	if err := utils.DB.First(&agent, "id = ?", uid).Error; err != nil {
		utils.Log.Warnf("not fount agent,agent user id:%s", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundUser.Data()
		return ret
	}
	ret.Status = response.StatusSucc
	ret.Data = []models.User{agent}
	return ret
}

func GetUsers(page, size, status, startTime, stopTime, sort, timeField, search, role string) response.PageResponse {
	var result []models.User
	var ret response.PageResponse
	db := utils.DB.Model(&models.User{}).Order(fmt.Sprintf("%s %s", timeField, sort)).Where("role = ?", role)
	if search != "" {
		db = db.Where("username = ? ", search)
	} else {

		if startTime != "" && stopTime != "" {
			db = db.Where(fmt.Sprintf("%s >= ? AND %s <= ?", timeField, timeField), startTime, stopTime)
		}
		if status != "" {
			db = db.Where("user_status = ?", status)
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

func UpdateUser(param response.UserArgs, uid string) response.EntityResponse {
	var ret response.EntityResponse
	var user models.User

	if err := utils.DB.First(&user, "id = ?", uid).Error; err != nil {
		utils.Log.Warnf("not fount agent,agent user id:%s", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundUser.Data()
		return ret
	}
	data := models.User{
		Phone:   param.Phone,
		Address: param.Address,
		Email:   param.Email,
	}
	if err := utils.DB.Model(&models.User{}).Where("id = ?", uid).Updates(data).Error; err != nil {
		utils.Log.Errorf("update user is failed,user id:%s", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.UpdateUserErr.Data()
		return ret
	}
	ret.Status = response.StatusSucc

	return ret
}

func UpdateUserPassword(param response.UserPasswordArgs, uid string) response.EntityResponse {
	var ret response.EntityResponse
	var user, data models.User

	if err := utils.DB.First(&user, "id = ? AND (role=1 OR role=2)", uid).Error; err != nil {
		utils.Log.Warnf("not fount agent,agent user id:%s", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundUser.Data()
		return ret
	}
	algorithm := utils.Config.GetString("algorithm")
	if !verifyPassword(param.OriginPassword, user.Algorithm, user.Password, user.Salt) {
		utils.Log.Warnf("not fount agent,agent user id:%s", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.OriginUserPasswordErr.Data()
		return ret
	}

	data.Salt, data.Password = generatePassword(param.Password, algorithm)
	data.Algorithm = algorithm
	if err := utils.DB.Model(&models.User{}).Where("id = ?", uid).Updates(data).Error; err != nil {
		utils.Log.Errorf("reset user password is failed,user id:%s", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.UpdateUserPasswordErr.Data()
		return ret
	}
	return ret
}

func ResetUserPassword(param response.UserArgs, uid string) response.EntityResponse {
	var ret response.EntityResponse
	var user, data models.User

	if err := utils.DB.First(&user, "id = ? AND (role=1 OR role=2)", uid).Error; err != nil {
		utils.Log.Warnf("not fount agent,agent user id:%s", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundUser.Data()
		return ret
	}
	algorithm := utils.Config.GetString("algorithm")
	data.Salt, data.Password = generatePassword(param.Password, algorithm)
	data.Algorithm = algorithm
	if err := utils.DB.Model(&models.User{}).Where("id = ?", uid).Updates(data).Error; err != nil {
		utils.Log.Errorf("reset user password is failed,user id:%s", uid)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.ResetUserPasswordErr.Data()
		return ret
	}
	ret.Status = response.StatusSucc
	return ret
}

func generatePassword(password, algorithm string) ([]byte, []byte) {

	salt, err := generateRandomBytes(64)
	if err != nil {
		utils.Log.Errorf("Unable to get random salt")
		return nil, nil
	}

	hashFunc := functionMap[algorithm]
	ps := hashFunc([]byte(password), salt)

	return salt, ps
}

func verifyPassword(oripassword, algorithm string, password, salt []byte) bool {

	hashFunc := functionMap[algorithm]
	hash := hashFunc([]byte(oripassword), salt)
	if compare(password, hash) != 0 {
		utils.Log.Warnf("Invalid username/password set")
		return false
	}
	return true
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
