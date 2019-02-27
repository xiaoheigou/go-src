package service

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/utils"
)

// 得到待签名字符串
func buildWithdrawParams(arg response.DistributorWithdrawArgs, distributorId int64, apiKey string, username string) string {

	params := make(map[string]string)
	params["appId"] = strconv.FormatInt(distributorId, 10)
	params["apiKey"] = apiKey
	params["inputCharset"] = "UTF-8"
	params["apiVersion"] = "1.1"
	params["appSignType"] = "HMAC-SHA256"
	params["appUserId"] = username
	params["appOrderId"] = arg.AppOrderId
	params["orderAmount"] = arg.OrderAmount
	params["orderCoinSymbol"] = "CNY"
	params["orderPayTypeId"] = arg.OrderPayTypeId
	params["orderRemark"] = "from web console"
	params["payAccountId"] = arg.PayAccountId
	params["payQRUrl"] = ""
	params["payAccountUser"] = arg.PayAccountUser
	params["payAccountInfo"] = arg.PayAccountInfo
	params["appServerNotifyUrl"] = arg.AppServerNotifyUrl
	params["appReturnPageUrl"] = arg.AppReturnPageUrl

	urlParams := url.Values{}
	for k := range params {
		urlParams.Add(k, params[k])
	}
	return urlParams.Encode() // Encode()会按key排序
}

func buildWithdrawBodyParams(arg response.DistributorWithdrawArgs, distributorId int64, apiKey string, username string) string {
	params := make(map[string]string)
	params["appUserId"] = username
	params["appOrderId"] = arg.AppOrderId
	params["orderAmount"] = arg.OrderAmount
	params["orderCoinSymbol"] = "CNY"
	params["orderPayTypeId"] = arg.OrderPayTypeId
	params["orderRemark"] = "from web console"
	params["payAccountId"] = arg.PayAccountId
	params["payQRUrl"] = ""
	params["payAccountUser"] = arg.PayAccountUser
	params["payAccountInfo"] = arg.PayAccountInfo
	params["appServerNotifyUrl"] = arg.AppServerNotifyUrl
	params["appReturnPageUrl"] = arg.AppReturnPageUrl

	// 让下单接口返回json格式的数据
	params["responseFormat"] = "json"

	urlParams := url.Values{}
	for k := range params {
		urlParams.Add(k, params[k])
	}
	return urlParams.Encode() // Encode()会按key排序
}

func fireWithdraw(arg response.DistributorWithdrawArgs, distributorId string, username string) error {

	//
	var distributor models.Distributor

	if err := utils.DB.First(&distributor, "id = ?", distributorId).Error; err != nil {
		utils.Log.Errorf("func GetDistributorByIdAndAPIKey err: %v", err)
		return errors.New("db access error")
	}

	paramsWithUrlEncoded := buildWithdrawParams(arg, distributor.Id, distributor.ApiKey, username)
	appSignContent, _ := HmacSha256Base64Signer(paramsWithUrlEncoded, distributor.ApiSecret)

	//apiUrl := "https://jrdidi.com/order/withdraw/create?appId=10001&apiKey=c6aec828fe514980&inputCharset=UTF-8&apiVersion=1.1&appSignType=HMAC-SHA256&appSignContent=" + content
	reqUrl := "http://13.250.12.109:8084/order/withdraw/create"
	apiUrl := fmt.Sprintf("%s?appId=%s&apiKey=%s&inputCharset=UTF-8&apiVersion=1.1&appSignType=HMAC-SHA256&appSignContent=%s",
		reqUrl, distributorId, distributor.ApiKey, appSignContent)

	u, _ := url.ParseRequestURI(apiUrl)
	urlStr := u.String()

	withdrawBodyParams := buildWithdrawBodyParams(arg, distributor.Id, distributor.ApiKey, username)

	// TODO
	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(withdrawBodyParams)) // URL-encoded payload
	r.Header.Add("Content-Length", strconv.Itoa(len(withdrawBodyParams)))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(r)

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := fmt.Sprintf("%s", body)
	utils.Log.Debugf("with draw result is: %v", bodyStr)

	type respMsg struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	}

	var data respMsg
	// 如果respBody中没有result域，下面的json.Unmarshal(respBody, &data)也不会失败。
	// 先给data.Result设一个不为0的初值，下面如果data.Result中被设置为0，则说明短信api调用成功
	if err := json.Unmarshal(body, &data); err != nil {
		//
		return err
	}

	if data.Code != "000000" {
		return errors.New(data.Msg)
	}
	return nil
}

// 平台商用户登录管理后台进行提现操作的Api
func DistributorWithdraw(arg response.DistributorWithdrawArgs, distributorId string, username string) response.EntityResponse {
	var ret response.EntityResponse

	// https://jrdidi.com/order/withdraw/create?appId=[由JRDiDi平台分配的appId]&apiKey=[由JRDiDi平台分配的apiKey]&inputCharset=UTF-8&apiVersion=1.1&appSignType=HMAC-SHA256&appSignContent=[签名内容]

	if err := fireWithdraw(arg, distributorId, username); err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateWithdrawOrderErr.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	return ret
}
