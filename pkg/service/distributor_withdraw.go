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

// 准备提现订单的待签名字符串
func buildWithdrawSignatureParams(arg response.DistributorWithdrawArgs, distributorId int64, apiKey string, username string) string {

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

	// 让下单接口返回json格式的数据
	params["responseFormat"] = "json"

	urlParams := url.Values{}
	for k := range params {
		urlParams.Add(k, params[k])
	}
	return urlParams.Encode() // Encode()会按key排序
}

// 准备提现订单的body
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

// 提交提现请求
func fireWithdraw(arg response.DistributorWithdrawArgs, distributorId string, username string) error {

	var distributor models.Distributor
	if err := utils.DB.First(&distributor, "id = ?", distributorId).Error; err != nil {
		utils.Log.Errorf("func GetDistributorByIdAndAPIKey err: %v", err)
		return errors.New("db access error")
	}

	paramsWithUrlEncoded := buildWithdrawSignatureParams(arg, distributor.Id, distributor.ApiKey, username)
	appSignContent, _ := HmacSha256Base64Signer(paramsWithUrlEncoded, distributor.ApiSecret)

	reqUrl := utils.Config.GetString("jrdidiurl.api")
	apiUrl := fmt.Sprintf("/order/withdraw/create%s?appId=%s&apiKey=%s&inputCharset=UTF-8&apiVersion=1.1&appSignType=HMAC-SHA256&appSignContent=%s",
		reqUrl, distributorId, distributor.ApiKey, appSignContent)

	u, _ := url.ParseRequestURI(apiUrl)
	urlStr := u.String()

	// 准备请求body
	withdrawBodyParams := buildWithdrawBodyParams(arg, distributor.Id, distributor.ApiKey, username)

	client := http.Client{}
	request, _ := http.NewRequest("POST", urlStr, strings.NewReader(withdrawBodyParams)) // URL-encoded payload
	request.Header.Add("Content-Length", strconv.Itoa(len(withdrawBodyParams)))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(request)
	if err != nil {
		utils.Log.Errorln(err.Error())
		return err
	} else {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := fmt.Sprintf("%s", body)
		utils.Log.Debugf("func fireWithdraw, withdraw result is: %v", bodyStr)

		// 下面是下单api返回的消息格式
		type respMsg struct {
			Code string `json:"code"`
			Msg  string `json:"msg"`
		}

		var data respMsg
		if err := json.Unmarshal(body, &data); err != nil {
			utils.Log.Errorf("func fireWithdraw, unmarshal fail %s", err)
			return err
		}

		if data.Code != "000000" { // 下提现订单的接口，只有 "000000" 是正常的
			utils.Log.Errorf("func fireWithdraw, request to %s fail, err : %s", apiUrl, data.Msg)
			return errors.New(data.Msg)
		}
		return nil
	}
}

// 平台商用户登录管理后台进行提现操作的Api
func DistributorWithdraw(arg response.DistributorWithdrawArgs, distributorId string, username string) response.EntityResponse {
	var ret response.EntityResponse

	if err := fireWithdraw(arg, distributorId, username); err != nil {
		utils.Log.Errorf("func DistributorWithdraw, call fireWithdraw fail. err %s", err)
		ret.Status = response.StatusFail
		ret.ErrCode = err_code.CreateWithdrawOrderErr.ErrCode
		ret.ErrMsg = err.Error()
		return ret
	}

	ret.Status = response.StatusSucc
	return ret
}
