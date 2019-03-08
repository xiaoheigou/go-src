package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SendSmsUserRegister(phone string, nationCode int, randomCode string, timeout string) error {

	currentApi := Config.GetString("sms.currentapi")
	Log.Debugf("sms.currentapi is %s", currentApi)
	if currentApi == "tencent" {
		// 使用Tencent Api
		// 获取tencent短信模板id，这是提前在短信api管理台中设置的短信模板
		var tplId int64
		var err error
		if tplId, err = strconv.ParseInt(Config.GetString("sms.tencent.tplid.register"), 10, 0); err != nil {
			Log.Errorf("Wrong configuration: sms.tencent.tplid.register, should be int.")
			return errors.New("sms.tencent.tplid.register, should be int")
		}

		return SendSmsByTencentApi(phone, nationCode, tplId, randomCode, timeout)
	} else {
		// 使用Twilio Api
		// 从配置文件中拿到短信模板
		smsBodyTpl := Config.GetString("sms.twilio.smstemplate.register")
		smsBody := strings.Replace(smsBodyTpl, "{1}", randomCode, 1)
		smsBody = strings.Replace(smsBody, "{2}", timeout, 1)

		return SendSmsByTwilioApi(phone, nationCode, smsBody)
	}
}

func SendSmsByTwilioApi(phone string, nationCode int, smsBody string) error {

	var smsServiceURL = Config.GetString("sms.twilio.smssvcendpoit")
	if smsServiceURL == "" {
		Log.Warnln("Wrong configuration: sms.twilio.smssvcendpoit is empty")
		return errors.New("wrong configuration: sms.twilio.smssvcendpoit is empty")
	}

	// 下面是请求参数
	message := map[string]interface{}{
		"nation_code": nationCode,
		"phone":       phone,
		"sms_body":    smsBody,
	}

	jsonValue, err := json.Marshal(message)
	if err != nil {
		Log.Errorln(err)
		return err
	}

	request, err := http.NewRequest("POST", smsServiceURL, bytes.NewBuffer(jsonValue))
	if err != nil {
		Log.Errorln(err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		Log.Errorln(err.Error())
		return err
	} else {
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Log.Errorln(err)
			return err
		}

		Log.Debugf("Send sms by twilio api, resp body = [%+v]", string(respBody[:]))

		// 下面是返回报文的格式：
		type respMsg struct {
			Code   int    `json:"code"`
			ErrMsg string `json:"err_msg"`
		}

		var data respMsg
		data.Code = -1 // 当报文中没有code域时，json.Unmarshal(respBody, &data)后code也会为0，为避免歧义，先设置code为不为0的值
		if err := json.Unmarshal(respBody, &data); err != nil {
			Log.Errorln(err)
			return err
		}

		// 当返回报文中code为0时，表明发送成功，否则发送失败。
		if data.Code != 0 {
			errMsg := data.ErrMsg
			Log.Errorln(errMsg)
			return errors.New(errMsg)
		}
	}

	// 发送成功
	return nil
}

// 通过腾讯api发送短信
//
// 短信内容是模板中订制的 ，模板是在短信api管理台中设置的，tplId表示模板号，smsTplArg表示模板中的参数
// smsTplArg1是想要发送的短信验证码，smsTplArg2是想要发送的过期时间
func SendSmsByTencentApi(phone string, nationCode int, tplId int64, smsTplArgs ...string) error {
	// 参考 https://cloud.tencent.com/document/product/382/5976

	var SdkAppId = Config.GetString("sms.tencent.sdkappid")
	if SdkAppId == "" {
		Log.Errorln("Wrong configuration: sms.tencent.sdkappid is empty")
		return errors.New("sms.tencent.sdkappid is empty")
	}
	var AppKey = Config.GetString("sms.tencent.appkey")
	if AppKey == "" {
		Log.Errorln("Wrong configuration: sms.tencent.appkey is empty")
		return errors.New("sms.tencent.appkey is empty")
	}

	var params []string
	for _, smsTplArg := range smsTplArgs {
		params = append(params, smsTplArg)
	}

	// 按照文档要求，短信api请求url中需要提供一个随机数
	rand.Seed(time.Now().UTC().UnixNano())
	max := 100000
	min := 100
	var Random = strconv.Itoa(rand.Intn(max-min) + min)

	// 短信api请求url
	var url = "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=" + SdkAppId + "&random=" + Random

	unixSec := strconv.FormatInt(time.Now().Unix(), 10)
	sigData := "appkey=" + AppKey + "&random=" + Random + "&time=" + unixSec + "&mobile=" + phone
	// fmt.Println(tmp)
	sig := sha256.Sum256([]byte(sigData))
	sigStr := hex.EncodeToString(sig[:])

	// 下面是请求参数
	message := map[string]interface{}{
		"sig":    sigStr, // "sig" 字段根据公式sha256(appkey=$appkey&random=$random&time=$time&mobile=$mobile)生成
		"params": params,
		"tel": map[string]string{
			"mobile":     phone,
			"nationcode": strconv.Itoa(nationCode),
		},
		"time":   unixSec,
		"tpl_id": tplId,
	}

	jsonValue, err := json.Marshal(message)
	if err != nil {
		Log.Errorln(err)
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		Log.Errorln(err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		Log.Errorln(err.Error())
		return err
	} else {
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Log.Errorln(err)
			return err
		}

		Log.Debugf("Send sms, resp body = [%+v]", string(respBody[:]))

		// 下面是返回报文的实例
		//{
		//	"result": 0,
		//	"errmsg": "OK",
		//	"ext": "",
		//	"fee": 1,
		//	"sid": "xxxxxxx"
		//}
		type respMsg struct {
			Result int    `json:"result"`
			ErrMsg string `json:"errmsg"`
			Ext    string `json:"ext"`
			Fee    int    `json:"fee"`
			Sid    string `json:"sid"`
		}

		var data respMsg
		data.Result = -1
		// 如果respBody中没有result域，下面的json.Unmarshal(respBody, &data)也不会失败。
		// 先给data.Result设一个不为0的初值，下面如果data.Result中被设置为0，则说明短信api调用成功
		if err := json.Unmarshal(respBody, &data); err != nil {
			Log.Errorln(err)
			return err
		}

		// 根据短信api文档，当返回报文中result为0时，表明发送成功，否则发送失败。
		if data.Result != 0 {
			errMsg := data.ErrMsg
			Log.Errorln(errMsg)
			return errors.New(errMsg)
		}
	}

	// 发送成功
	return nil
}
