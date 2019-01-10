package service

import (
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/utils"

	"github.com/dgrijalva/jwt-go"
)

func GetRandomCode(arg response.SendRandomCodeArg) response.SendRandomCodeRet {
	var ret response.SendRandomCodeRet

	account := arg.Account
	nationCode := arg.NationCode
	purpose := arg.Purpose

	// 检验参数
	if strings.Contains(account, "@") {
		// 邮箱
		if !utils.IsValidEmail(account) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrEmailInvalid.Data()
			return ret
		}
	} else {
		// 手机号
		if !utils.IsValidNationCode(nationCode) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrNationCodeInvalid.Data()
			return ret
		}
		if !utils.IsValidPhone(nationCode, account) {
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

	// 发送短信（邮件）前，测试是否通过geetest测试，如果通过了测试，则captcha服务会把redis中对应的key会被设置为success
	// 下面检测这个key对应的value是否为success
	geetestKey := "app:" + purpose + ":captcha" // example: "app:register:captcha"
	if strings.Contains(account, "@") {
		geetestKey = geetestKey + ":" + account // key example: "app:register:captcha:xx@yy.com"
	} else {
		geetestKey = geetestKey + ":" + strconv.Itoa(nationCode) + ":" + account // key example: "app:register:captcha:86:13100000000"
	}
	geetestValue := "success"
	if err := utils.RedisVerifyValue(geetestKey, geetestValue); err != nil {
		utils.Log.Errorf("check captcha verify result fail, can not send sms/email. error [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrCaptchaVerifyFail.Data()
		return ret
	}

	key := "app:" + purpose // example: "app:register"
	seq := int(utils.RandomSeq.GetCount())
	value := strconv.Itoa(seq) + ":" + randomCode
	var timeout int
	if strings.Contains(account, "@") {
		key = key + ":" + account // key example: "app:register:xxx@yyy.com"
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
		key = key + ":" + strconv.Itoa(nationCode) + ":" + account // key example: "app:register:86:13100000000"
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
		if !utils.IsValidEmail(account) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrEmailInvalid.Data()
			return ret
		}
	} else {
		// 手机号
		isEmail = false
		if !utils.IsValidNationCode(nationCode) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrNationCodeInvalid.Data()
			return ret
		}
		if !utils.IsValidPhone(nationCode, account) {
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

	if !skipVerify {
		if err := utils.RedisVerifyValue(key, value); err != nil {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrRandomCodeVerifyFail.Data()
			return ret
		}
	}

	// 检测手机号是否已经注册过，如果已经注册过，则提示错误
	if purpose == "register" && (!isEmail) {
		var registered bool
		var err error
		if registered, err = IsMerchantPhoneRegistered(arg.NationCode, arg.Account); err != nil {
			utils.Log.Errorf("database access err = [%v]", err)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
			return ret
		}

		if registered == true {
			// 手机号已经注册过
			utils.Log.Errorf("phone [%v] nation_code [%v] is already registered.", arg.Account, nationCode)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrPhoneAlreadyRegister.Data()
			return ret
		}
	}

	ret.Status = response.StatusSucc
	return ret
}

func verifyMerchantPw(passWord string, user models.Merchant) bool {
	var passwordInDB []byte = user.Password
	var saltInDB []byte = user.Salt
	var algorithmInDB string = user.Algorithm

	if len(algorithmInDB) == 0 {
		algorithmInDB = "Argon2"
		utils.Log.Warnf("algorithm field is missing for user [%s], use Argon2 as default", user.Phone)
	}
	hashFunc := functionMap[algorithmInDB]
	hash := hashFunc([]byte(passWord), saltInDB)
	if compare(passwordInDB, hash) != 0 {
		return false
	} else {
		return true
	}

}

func AppLogin(arg response.LoginArg) response.LoginRet {
	var ret response.LoginRet

	// 检验参数
	if !utils.IsValidNationCode(arg.NationCode) {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrNationCodeInvalid.Data()
		return ret
	}
	if !utils.IsValidPhone(arg.NationCode, arg.Account) {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrPhoneInvalid.Data()
		return ret
	}

	var user models.Merchant

	// 通过国家码和手机号查找记录
	if err := utils.DB.First(&user, "nation_code = ? and phone = ?", arg.NationCode, arg.Account).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// 找不到记录
			utils.Log.Warnf("not found user :%s", arg.Account)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrUserPasswordError.Data()
			return ret
		} else {
			// 其它错误
			utils.Log.Errorf("database access err = [%v]", err)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
			return ret
		}
	}

	if !verifyMerchantPw(arg.Password, user) {
		utils.Log.Warnf("Invalid username/password set")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrUserPasswordError.Data()
		return ret
	}

	// 更新上一次登录时间
	if err := utils.DB.Table("merchants").Where("id = ?", user.Id).Update("last_login", time.Now()).Error; err != nil {
		utils.Log.Errorf("AppLogin, modify last_login faile err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	// 生成一个jwt
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	tokenExpire := time.Now().Add(time.Hour * 1).Unix() // TODO 可配置
	// Set some claims
	token.Claims = jwt.MapClaims{
		"uid": strconv.FormatInt(user.Id, 10), // 为方便校验合法性时分析token，转换为字符串
		"exp": tokenExpire,
	}
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(utils.Config.GetString("appauth.authkey")))
	if err != nil {
		utils.Log.Errorf("Can't generate jwt token [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}

	var svrConfig = getSvrConfigFromFile()

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.LoginData{
		Uid:           user.Id,
		UserStatus:    user.UserStatus,
		UserCert:      user.UserCert,
		NickName:      user.Nickname,
		Token:         tokenString,
		TokenExpire:   tokenExpire,
		SvrConfigData: svrConfig,
	})
	return ret
}

func RefreshToken(uid int) response.RefreshTokenRet {
	var ret response.RefreshTokenRet
	// 生成一个jwt
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	tokenExpire := time.Now().Add(time.Hour * 1).Unix()
	// Set some claims
	token.Claims = jwt.MapClaims{
		"uid": strconv.FormatInt(int64(uid), 10), // 为方便校验合法性时分析token，转换为字符串
		"exp": tokenExpire,
	}
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(utils.Config.GetString("appauth.authkey")))
	if err != nil {
		utils.Log.Errorf("Can't generate jwt token [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.RefreshTokenData{
		Token:       tokenString,
		TokenExpire: tokenExpire,
	})
	return ret
}

func getSvrConfigFromFile() response.SvrConfigData {
	var svrConfig response.SvrConfigData

	svrConfig.SvrCurrentTime = time.Now().UTC()
	svrConfig.LatestApkVerCode = utils.Config.GetInt("app.latestapkvercode")
	svrConfig.LatestApkVerName = utils.Config.GetString("app.latestapkvername")
	svrConfig.LatestApkUrl = utils.Config.GetString("app.latestapkurl")
	svrConfig.QrcodePrefixAlipay = utils.Config.GetString("qrcode.expectprefix.alipay")
	svrConfig.QrcodePrefixWeixin = utils.Config.GetString("qrcode.expectprefix.weixin")
	svrConfig.TimeoutAwaitAccept = utils.Config.GetInt("fulfillment.timeout.awaitaccept")
	svrConfig.TimeoutNotifyPaid = utils.Config.GetInt("fulfillment.timeout.notifypaid")
	svrConfig.TimeoutNotifyPaymentConfirmed = utils.Config.GetInt("fulfillment.timeout.notifypaymentconfirmed")
	return svrConfig
}

func GetSvrConfig(uid int64) response.SvrConfigRet {
	var ret response.SvrConfigRet

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, getSvrConfigFromFile())
	return ret
}

func ChangeMerchantPassword(uid int, arg response.ChangePasswordArg) response.ChangePasswordRet {
	var ret response.ChangePasswordRet

	oldPw := arg.OldPassword
	newPw := arg.NewPassword

	var user models.Merchant
	if err := utils.DB.First(&user, "id = ?", uid).Error; err != nil {
		utils.Log.Warnf("found merchant(id=[%v]) err %v", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrUserPasswordError.Data()
		return ret
	}

	// 验证旧密码
	if !verifyMerchantPw(oldPw, user) {
		utils.Log.Warnf("old password invalid")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrOldPasswordError.Data()
		return ret
	}

	// 设置新密码
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
	passwordEncrypted := hashFunc([]byte(newPw), salt)

	if err := utils.DB.Model(&user).Updates(map[string]interface{}{"algorithm": algorithm, "salt": salt, "password": passwordEncrypted}).Error; err != nil {
		utils.Log.Errorf("ChangeMerchantPassword fail, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	utils.Log.Info("Merchant (uid=[%v] phone=[%v]) update password successful", uid, user.Phone)
	ret.Status = response.StatusSucc
	return ret
}
