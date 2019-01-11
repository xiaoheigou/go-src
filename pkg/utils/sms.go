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
	"time"
)

func SendSmsUserRegister(phone string, nationCode int, smsTplArg1 string, smsTplArg2 string) error {

	// 短信模板id，这是提前在短信api管理台中设置的短信模板
	var tplId int64
	if tplId, err = strconv.ParseInt(Config.GetString("sms.tencent.tplid.register"), 10, 0); err != nil {
		Log.Errorf("Wrong configuration: sms.tencent.tplid.register, should be int.")
		return errors.New("sms.tencent.tplid.register, should be int")
	}

	return SendSms(phone, nationCode, tplId, smsTplArg1, smsTplArg2)
}

// 通过腾讯api发送短信
//
// 短信内容是模板中订制的 ，模板是在短信api管理台中设置的，tplId表示模板号，smsTplArg表示模板中的参数
// smsTplArg1是想要发送的短信验证码，smsTplArg2是想要发送的过期时间
func SendSms(phone string, nationCode int, tplId int64, smsTplArgs ...string) error {
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

		Log.Debugln("Send sms, resp body = [%+v]", string(respBody[:]))

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
