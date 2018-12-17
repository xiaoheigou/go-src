package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

func verifyMerchantPw(passWord string, user models.Merchant) bool {
	var passwordInDB []byte = user.Password
	var saltInDB []byte = user.Salt
	var algorithmInDB string = user.Algorithm

	if len(algorithmInDB) == 0 {
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

	var user models.Merchant

	// 通过国家码和手机号查找记录
	if err := utils.DB.First(&user, "nation_code = ? and phone = ?", arg.NationCode, arg.Account).Error; err != nil {
		utils.Log.Warnf("not found user :%s", arg.Account)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrUserPasswordError.Data()
		return ret
	}

	if ! verifyMerchantPw(arg.Password, user) {
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

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.LoginData{
		Uid:        user.Id,
		UserStatus: user.UserStatus,
		UserCert:   user.UserCert,
		NickName:   user.Nickname,
	})
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
	if ! verifyMerchantPw(oldPw, user) {
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

func RegisterGeetest(arg response.RegisterGeetestArg) response.RegisterGeetestRet {
	var ret response.RegisterGeetestRet

	account := arg.Account
	nationCode := arg.NationCode
	purpose := arg.Purpose
	var geetestUserId string

	if purpose == "" {
		utils.Log.Infof("param purpose is missing, use register as default")
		purpose = "register"
	}

	// 检验参数
	if strings.Contains(account, "@") {
		// 邮箱
		if ! utils.IsValidEmail(account) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrEmailInvalid.Data()
			return ret
		}
		geetestUserId = utils.GetMD5Hash(purpose + account)
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
		geetestUserId = utils.GetMD5Hash(purpose + strconv.Itoa(nationCode) + account)
	}

	var url = utils.Config.GetString("register.geetestsvr")
	if url == "" {
		utils.Log.Errorln("Wrong configuration: register.geetestsvr [%v].", url)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}
	url = url + "/gt/register"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Log.Errorln("http.NewRequest fail.",)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}
	q := req.URL.Query()
	q.Add("geetest_user_id", geetestUserId) // Add a new value to the set.
	req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.

	resp, err := client.Do(req)
	if err != nil {
		utils.Log.Errorln("client.Do fail.",)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	utils.Log.Infof("respBody:[%s]", respBody)

	type respMsg struct {
		Status string    `json:"status"`
		GeetestChallenge string `json:"geetest_challenge"`
		GeetestServerStatus string `json:"geetest_server_status"`
	}

	var data respMsg
	if err := json.Unmarshal(respBody, &data); err != nil {
		utils.Log.Errorln(err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}
	utils.Log.Debugf("RegisterGeetest Status=[%s]", data.Status)
	utils.Log.Debugf("RegisterGeetest GeetestServerStatus=[%s]", data.GeetestServerStatus)
	utils.Log.Debugf("RegisterGeetest GeetestChallenge=[%s]", data.GeetestChallenge)

	if data.Status != "success" {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.RegisterGeetestData{
		GeetestServerStatus: data.GeetestServerStatus,
		GeetestChallenge:    data.GeetestChallenge,
	})
	return ret
}

func VerifyGeetest(arg response.VerifyGeetestArg) response.CommonRet {
	var ret response.CommonRet

	account := arg.Account
	nationCode := arg.NationCode
	purpose := arg.Purpose
	geetestChallenge := arg.GeetestChallenge
	geetestValidate := arg.GeetestValidate
	geetestSeccode := arg.GeetestSeccode
	var geetestUserId string

	// 检验参数
	if strings.Contains(account, "@") {
		// 邮箱
		if ! utils.IsValidEmail(account) {
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrEmailInvalid.Data()
			return ret
		}
		geetestUserId = utils.GetMD5Hash(purpose + account)
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
		geetestUserId = utils.GetMD5Hash(purpose + strconv.Itoa(nationCode) + account)
	}
	if purpose == "" {
		utils.Log.Infof("param purpose is missing, use register as default")
		purpose = "register"
	}
	if geetestChallenge == "" {
		utils.Log.Errorln("geetest_challenge is missing.")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}
	if geetestValidate == "" {
		utils.Log.Errorln("geetest_validate is missing.")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}
	if geetestSeccode == "" {
		utils.Log.Errorln("geetest_seccode is missing.")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}

	var url = utils.Config.GetString("register.geetestsvr")
	if url == "" {
		utils.Log.Errorln("Wrong configuration: register.geetestsvr [%v].", url)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}
	url = url + "/gt/validate"

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		utils.Log.Errorln("http.NewRequest fail.",)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}
	q := req.URL.Query()
	q.Add("geetest_user_id", geetestUserId) // Add a new value to the set.
	q.Add("geetest_challenge", geetestChallenge) // Add a new value to the set.
	q.Add("geetest_validate", geetestValidate) // Add a new value to the set.
	q.Add("geetest_seccode", geetestSeccode) // Add a new value to the set.
	req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.

	resp, err := client.Do(req)
	if err != nil {
		utils.Log.Errorln("client.Do fail.",)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	utils.Log.Infof("respBody:[%s]", respBody)

	type respMsg struct {
		Status string    `json:"status"`
	}

	var data respMsg
	if err := json.Unmarshal(respBody, &data); err != nil {
		utils.Log.Errorln(err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}
	utils.Log.Debugf("VerifyGeetest Status=[%s]", data.Status)
	if data.Status != "success" {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrGeetestVerifyFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	return ret
}