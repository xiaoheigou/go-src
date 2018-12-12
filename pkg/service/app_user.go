package service

import (
	"strconv"
	"strings"
	"time"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

func GetRandomCode(nationCode, account, purpose string) response.GetRandomCodeRet {
	var ret response.GetRandomCodeRet
	ret.Status = response.StatusSucc

	randCode, err := utils.GetSecuRandomCode()
	utils.Log.Debugf("random code is [%v]", randCode)
	if err != nil {
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}

	// 保存到redis
	seq := int(utils.RandomSeq.GetCount())
	key := "app:" + purpose  // example: "app:register"
	value := string(seq) + ":" + randCode
	timeoutStr := utils.Config.GetString("sms.timeout")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		utils.Log.Errorf("Convert sms.timeout [%v] to int fail", timeoutStr)
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}
	if strings.Contains(account, "@") {
		key = key + ":" + account  // example: "app:register:xxx@yyy.com"
		// 发送邮件
		// TODO
	} else {
		key = key + ":" + nationCode + ":" + account  // example: "app:register:86:13100000000"
		// 发送短信
		err = utils.SendSms(account, nationCode, randCode, timeoutStr)
		if err != nil {
			ret.ErrCode, ret.ErrMsg = err_code.AppErrSendSMSFail.Data()
			return ret
		}
	}

	// 把随机码保存到redis中，以便以后验证用户输入
	err = utils.RedisSet(key, value, time.Duration(timeout) * time.Minute)
	if err != nil {
		ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, response.GetRandomCodeData{RandomCodeSeq: seq})
	return ret
}