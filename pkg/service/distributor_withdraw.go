package service

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/utils"
)

// 得到待签名字符串
func buildWithdrawParams(arg response.DistributorWithdrawArgs, distributorId int64, apiKey string) string {

	params := make(map[string]string)
	params["appId"] = strconv.FormatInt(distributorId, 10)
	params["apiKey"] = apiKey
	params["inputCharset"] = "UTF-8"
	params["apiVersion"] = "1.1"
	params["appSignType"] = "HMAC-SHA256"
	params["appUserId"] = arg.AppOrderId // TODO
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

func fireWithdraw(arg response.DistributorWithdrawArgs, distributorId string, uid string) error {

	//
	var distributor models.Distributor

	if err := utils.DB.First(&distributor, "id = ?", distributorId).Error; err != nil {
		utils.Log.Errorf("func GetDistributorByIdAndAPIKey err: %v", err)
		return errors.New("db access error")
	}

	paramsWithUrlEncoded := buildWithdrawParams(arg, distributor.Id, distributor.ApiKey)
	appSignContent, _ := HmacSha256Base64Signer(paramsWithUrlEncoded, distributor.ApiSecret)

	//apiUrl := "https://jrdidi.com/order/withdraw/create?appId=10001&apiKey=c6aec828fe514980&inputCharset=UTF-8&apiVersion=1.1&appSignType=HMAC-SHA256&appSignContent=" + content
	apiUrl := "http://13.250.12.109:8084/order/withdraw/create?appId=10001&apiKey=c6aec828fe514980&inputCharset=UTF-8&apiVersion=1.1&appSignType=HMAC-SHA256&appSignContent=" + appSignContent

	u, _ := url.ParseRequestURI(apiUrl)
	urlStr := u.String() // "https://api.com/user/"

	// TODO
	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(paramsWithUrlEncoded)) // URL-encoded payload
	r.Header.Add("Content-Length", strconv.Itoa(len(paramsWithUrlEncoded)))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(r)

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := fmt.Sprintf("%s", body)
	utils.Log.Debugf("with draw result is :[%v]", bodyStr)
	fmt.Println(resp.Status)

	return nil
}

// 平台商用户登录管理后台进行提现操作的Api
func DistributorWithdraw(arg response.DistributorWithdrawArgs, distributorId string, uid string) response.EntityResponse {
	var ret response.EntityResponse

	// https://jrdidi.com/order/withdraw/create?appId=[由JRDiDi平台分配的appId]&apiKey=[由JRDiDi平台分配的apiKey]&inputCharset=UTF-8&apiVersion=1.1&appSignType=HMAC-SHA256&appSignContent=[签名内容]

	fireWithdraw(arg, distributorId, uid)

	ret.Status = response.StatusSucc
	return ret
}
