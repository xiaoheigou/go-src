package service

import (
	"strconv"
	"strings"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

func GetRandomCode(arg response.SendRandomCodeArg) response.SendRandomCodeRet {
	var ret response.SendRandomCodeRet

	account := arg.Account
	nationCode := arg.NationCode
	purpose := arg.Purpose

	// 检验参数
	if strings.Contains(account, "@") {
		// 邮箱
		if ! utils.IsValidEmail(account) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrEmailInvalid.Data()
			return ret
		}
	} else {
		// 手机号
		if ! utils.IsValidNationCode(nationCode) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrNationCodeInvalid.Data()
			return ret
		}
		if ! utils.IsValidPhone(nationCode, account) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrPhoneInvalid.Data()
			return ret
		}
	}
	if len(purpose) == 0 {
		utils.Log.Infof("param purpose is missing, use register as default")
		purpose = "register"
	}

	randomCode, err := utils.GetSecuRandomCode()
	utils.Log.Debugf("random code is [%v]", randomCode)
	if err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}

	// 保存到redis
	seq := int(utils.RandomSeq.GetCount())
	key := "app:" + purpose // example: "app:register"
	value := strconv.Itoa(seq) + ":" + randomCode
	var timeout int
	if strings.Contains(account, "@") {
		key = key + ":" + account // example: "app:register:xxx@yyy.com"
		if timeout, err = strconv.Atoi(utils.Config.GetString("register.timeout.email")); err != nil {
			utils.Log.Errorf("Wrong configuration: register.timeout.email [%v], should be int. Set to default 10.", timeout)
			timeout = 10
		}
		// 发送邮件
		if err = utils.SendRandomCodeToMail(account, randomCode, strconv.Itoa(timeout)); err != nil {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrSendEmailFail.Data()
			return ret
		}
	} else {
		// redis中key名称
		key = key + ":" + strconv.Itoa(nationCode) + ":" + account // example: "app:register:86:13100000000"
		// 发送短信
		if timeout, err = strconv.Atoi(utils.Config.GetString("register.timeout.sms")); err != nil {
			utils.Log.Errorf("Wrong configuration: register.timeout.sms [%v], should be int. Set to default 10.", timeout)
			timeout = 10
		}
		if err = utils.SendSms(account, nationCode, randomCode, strconv.Itoa(timeout)); err != nil {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrSendSMSFail.Data()
			return ret
		}
	}

	// 把随机码保存到redis中，以便以后验证用户输入
	err = utils.RedisSet(key, value, time.Duration(timeout)*time.Minute)
	if err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.SendRandomCodeData{RandomCodeSeq: seq})
	return ret
}

func VerifyRandomCode(arg response.VerifyRandomCodeArg) response.VerifyRandomCodeRet {
	var ret response.VerifyRandomCodeRet

	// 检验参数
	var account string = arg.Account
	var nationCode int = arg.NationCode
	var randomCode string = arg.RandomCode
	var randomCodeSeq int = arg.RandomCodeSeq
	purpose := arg.Purpose
	isEmail := true
	if strings.Contains(account, "@") {
		// 邮箱
		isEmail = true
		if ! utils.IsValidEmail(account) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrEmailInvalid.Data()
			return ret
		}
	} else {
		// 手机号
		isEmail = false
		if ! utils.IsValidNationCode(nationCode) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrNationCodeInvalid.Data()
			return ret
		}
		if ! utils.IsValidPhone(nationCode, account) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrPhoneInvalid.Data()
			return ret
		}
	}
	if len(randomCode) == 0 {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}
	if len(purpose) == 0 {
		utils.Log.Infof("param purpose is missing, use register as default")
		purpose = "register"
	}

	key := "app:" + purpose // example: "app:register"
	value := strconv.Itoa(randomCodeSeq) + ":" + randomCode
	if isEmail {
		key = key + ":" + account // example: "app:register:xxx@yyy.com"
	} else {
		key = key + ":" + strconv.Itoa(nationCode) + ":" + account // example: "app:register:86:13100000000"
	}

	skipVerify := false

	var skipEmailVerify bool
	if isEmail {
		// 邮箱服务器有时会限制，密集发送邮件可能失败。这里检测是否配置了跳过mail验证。
		var err error
		if skipEmailVerify, err = strconv.ParseBool(utils.Config.GetString("register.skipemailverify")); err != nil {
			utils.Log.Errorf("Wrong configuration: register.skipemailverify, should be boolean. Set to default false.")
			skipEmailVerify = false
		}
	}
	if isEmail && skipEmailVerify {
		skipVerify = true
	}

	if ! skipVerify {
		if err := utils.RedisVerifyValue(key, value); err != nil {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrRandomCodeVerifyFail.Data()
			return ret
		}
	}

	ret.Status = response.StatusSucc
	return ret
}

func AppLogin(arg response.LoginArg) response.LoginRet {
	var ret response.LoginRet

	// 检验参数
	if strings.Contains(arg.Account, "@") {
		// 邮箱
		if ! utils.IsValidEmail(arg.Account) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrEmailInvalid.Data()
			return ret
		}
	} else {
		// 手机号
		if ! utils.IsValidNationCode(arg.NationCode) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrNationCodeInvalid.Data()
			return ret
		}
		if ! utils.IsValidPhone(arg.NationCode, arg.Account) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrPhoneInvalid.Data()
			return ret
		}
	}

	var passwordInDB []byte
	var saltInDB []byte
	var algorithmInDB string
	var user models.Merchant
	if strings.Contains(arg.Account, "@") {
		// 邮箱
		if err := utils.DB.First(&user, "email = ?", arg.Account).Error; err != nil {
			utils.Log.Warnf("not found user :%s", arg.Account)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.NotFoundUser.Data()
			return ret
		}
		passwordInDB = user.Password
		saltInDB = user.Salt
		algorithmInDB = user.Algorithm
	} else {
		// 手机号
		if err := utils.DB.First(&user, "nation_code = ? and phone = ?", arg.NationCode, arg.Account).Error; err != nil {
			utils.Log.Warnf("not found user :%s", arg.Account)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrUserPasswordError.Data()
			return ret
		}
		passwordInDB = user.Password
		saltInDB = user.Salt
		algorithmInDB = user.Algorithm
	}

	if len(algorithmInDB) == 0 {
		utils.Log.Warnf("algorithm field is missing for user :%s, use Argon2 as default", arg.Account)
	}
	hashFunc := functionMap[algorithmInDB]
	hash := hashFunc([]byte(arg.Password), saltInDB)
	if compare(passwordInDB, hash) != 0 {
		utils.Log.Warnf("Invalid username/password set")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrUserPasswordError.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.LoginData{
		Uid:        user.Id,
		UserStatus: user.UserStatus,
		UserCert:   user.UserCert,
		NickName:   user.Nickname,
	})
	return ret
}
