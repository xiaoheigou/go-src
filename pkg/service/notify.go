package service

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/utils"
	"yuudidi.com/pkg/utils/timewheel"
)

var (
	notifywheel1 *timewheel.TimeWheel // 第一次回调没收到回复，4分钟后再通知一次
	notifyWheel2 *timewheel.TimeWheel // 第二次回调没收到回复，10分钟后再通知一次
	notifywheel3 *timewheel.TimeWheel // 第三次回调没收到回复，10分钟后再通知一次
	notifywheel4 *timewheel.TimeWheel // 第四次回调没收到回复，1小时后再通知一次
	notifywheel5 *timewheel.TimeWheel // 第四次回调没收到回复，2小时后再通知一次
	notifywheel6 *timewheel.TimeWheel // 第四次回调没收到回复，6小时后再通知一次
	notifywheel7 *timewheel.TimeWheel // 第四次回调没收到回复，15小时后再通知一次

)

type notifyItem struct {
	orderID           string
	orderStatus       models.OrderStatus
	orderStatusReason models.StatusReason
}

func (item notifyItem) MarshalBinary() (data []byte, err error) {
	return json.Marshal(item)

}

func (item notifyItem) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, item)
}

func CreateNotify(notify models.Notify) response.NotifyRet {
	//var notifyInsert models.Notify
	var ret response.NotifyRet
	if err := utils.DB.Create(&notify).Error; err != nil {
		utils.Log.Errorf("before sending notify message ,insert into notify wrong,err:[%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateNotifyErr.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = []models.Notify{notify}
	return ret

}

//根据notifyId获取notify
func GetNotifyByNotifyId(notifyId string) response.NotifyRet {
	var notify models.Notify
	var ret response.NotifyRet
	if err := utils.DB.First(&notify, "where jrdd_notify_id=?", notifyId).Error; err != nil {
		utils.Log.Errorf("find notify wrong,err:[%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.FindNotifyErr.Data()
		return ret
	}
	ret.Status = response.StatusSucc
	ret.Data = []models.Notify{notify}
	return ret
}

//获取没回调成功的notify，用来手动推送通知
func GetNotifyListBySendStatus(page, size string) response.PageResponse {
	var ret response.PageResponse
	var notifyList []models.Notify
	pageNum, err := strconv.ParseInt(page, 10, 64)
	pageSize, err1 := strconv.ParseInt(size, 10, 64)
	if err != nil || err1 != nil {
		utils.Log.Error(pageNum, pageSize)
	}
	db := utils.DB.Model(&models.Notify{}).Where("send_status < 2")

	db.Count(&ret.PageCount)
	db = db.Offset((pageNum - 1) * pageSize).Limit(pageSize)

	ret.PageNum = int(pageNum + 1)
	ret.PageSize = int(pageSize)
	db.Find(&notifyList)
	ret.Data = notifyList
	ret.Status = response.StatusSucc

	return ret

}

//根据order获取notify
func GetNotifyByOrder(order models.Order) models.Notify {
	var notify models.Notify
	orderNumber := order.OrderNumber
	status := order.Status
	statusReason := order.StatusReason
	if err := utils.DB.First(&notify, "jrdd_order_id=? and order_status=? and status_reason=?", orderNumber, status, statusReason).Error; err != nil {
		utils.Log.Errorf("get notify by order_number wrong,err=[%v]", err)
		return notify
	}

	return notify

}

//GetNotify - 根据notifyItem获取notify,除了orderID，还要同时考虑status+status_reason
func GetNotify(ni notifyItem) models.Notify {
	var notify models.Notify
	if err := utils.DB.First(&notify, "jrdd_order_id=? and order_status=? and status_reason=?", ni.orderID, ni.orderStatus, ni.orderStatusReason).Error; err != nil {
		utils.Log.Errorf("get notify by order_number wrong,err=[%v]", err)
		return notify
	}
	return notify
}

//手动推送消息
func ManualPushMessage(orderNumber string) response.CommonRet {
	var order models.Order
	var ret response.CommonRet
	if err := utils.DB.First(&order, "order_number=?", orderNumber).Error; err != nil {
		utils.Log.Errorf("get order by order_number wrong,err=[%v]", err)
		ret.Status = response.StatusFail
		return ret
	}
	notify := GetNotifyByOrder(order)
	utils.Log.Debugf("before sending mesage by hand,the notify =[%v] ", notify)
	if notify.JrddNotifyId == "" {
		utils.Log.Error("before sending mesage by hand,the notify is null")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.FindNotifyErr.Data()
		return ret
	}
	resp, err := PostNotifyToServer(order, notify)
	//发送回调消息成功并收到SUCCESS，更新notify表的send_status=2,attemps加1
	if err == nil && resp != nil && resp.Status == SUCCESS {
		utils.Log.Debugf("send message to distributor success,notify is: [%v]", notify)
		notify.Attempts += 1
		notify.SendStatus = 2
		if err := utils.DB.Model(&notify).Updates(notify).Error; err != nil {
			utils.Log.Errorf("after sending message by hand ,update the notify model wrong,notifyId=[%s],err=[%v]", notify.JrddNotifyId, err)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.UpdateNotifyErr.Data()
			return ret
		}
	}
	//发送回调消息没有收到SUCCESS，更新notify表，attemps加1
	notify.Attempts += 1
	if err := utils.DB.Model(&notify).Update(notify).Error; err != nil {
		utils.Log.Errorf("after sending message by hand ,update the notify model wrong,notifyId=[%s],err=[%v]", notify.JrddNotifyId, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.UpdateNotifyErr.Data()
		return ret
	}
	ret.Status = response.StatusSucc
	return ret

}

func Order2Notify(order models.Order) models.Notify {
	var notify models.Notify

	time := time.Now().Unix()
	notify = models.Notify{
		JrddNotifyId:       GenerateOrderNumber(),
		JrddNotifyTime:     time,
		JrddOrderId:        order.OrderNumber,
		AppOrderId:         order.OriginOrder,
		OrderAmount:        order.Amount,
		OrderCoinSymbol:    order.CurrencyFiat,
		OrderStatus:        int(order.Status),
		StatusReason:       int(order.StatusReason),
		OrderRemark:        order.Remark,
		OrderPayTypeId:     order.PayType,
		PayAccountId:       order.BankAccount,
		PayQRUrl:           order.QrCode,
		PayAccountUser:     order.Name,
		PayAccountInfo:     order.BankBranch,
		Synced:             0,
		Attempts:           0,
		SendStatus:         0,
		AppServerNotifyUrl: order.AppServerNotifyUrl,
		AppReturnPageUrl:   order.AppReturnPageUrl,
		OrderType:          order.Direction,
	}

	return notify

}

//发送回调前，先把order转换为notify，并保存在数据库，然后异步调用发送消息方法
func AsynchronousNotifyNew(order models.Order) response.NotifyRet {
	var ret response.NotifyRet
	//order转换为notify
	notify := Order2Notify(order)
	utils.Log.Debugf("before inserting into db, order convert to notify result is :[%v]", notify)
	//发送回调消息前，把消息存入数据库
	ret = CreateNotify(notify)
	//回调消息插入数据库后，异步发送通知，同时同步返回结果
	if ret.Status == response.StatusSucc {
		NotifyDistributorServerNew(order, ret.Data[0])
	}

	return ret

}

//异步回调消息
func NotifyDistributorServerNew(order models.Order, notify models.Notify) {
	serverUrl := order.AppServerNotifyUrl
	if serverUrl == "" || notify.JrddNotifyId == "" {
		utils.Log.Errorf("serverUrl or notify is null,before send to message to distributor server")
	} else {

		go func() {
			resp, err := PostNotifyToServer(order, notify)
			if err == nil && resp != nil && resp.Status == SUCCESS {
				utils.Log.Debugf("send message to distributor success,serverUrl is: [%s]", serverUrl)
				notify.Attempts += 1
				notify.SendStatus = 2
				if err := utils.DB.Model(&models.Notify{}).Updates(notify).Error; err != nil {
					utils.Log.Errorf("send message to distributor success,but update notify wrong ,notify=[%v],err=[%v]", notify, err)
				}

			} else {
				utils.Log.Errorf("send message to distributor fail,serverUrl is: [%s],err is:[%v]", serverUrl, err)
				notify.Attempts += 1
				if err := utils.DB.Model(&models.Notify{}).Update(notify).Error; err != nil {
					utils.Log.Errorf("send message to distributor fail,and update notify wrong ,notify=[%v],err=[%v]", notify, err)
				}
				notifywheel1.Add(notifyItem{
					orderID:           order.OrderNumber,
					orderStatus:       order.Status,
					orderStatusReason: order.StatusReason})
			}
		}()

	}
}

func Notify2ServerNotifyRequest(notify models.Notify) response.ServerNotifyRequest {

	var notifyRequest response.ServerNotifyRequest
	notifyRequest = response.ServerNotifyRequest{
		JrddNotifyId:    notify.JrddNotifyId,
		JrddNotifyTime:  notify.JrddNotifyTime,
		JrddOrderId:     notify.JrddOrderId,
		AppOrderId:      notify.AppOrderId,
		OrderAmount:     notify.OrderAmount,
		OrderCoinSymbol: notify.OrderCoinSymbol,
		OrderStatus:     int(notify.OrderStatus),
		StatusReason:    int(notify.StatusReason),
		OrderRemark:     notify.OrderRemark,
		OrderPayTypeId:  notify.OrderPayTypeId,
		PayAccountId:    notify.PayAccountId,
		PayAccountUser:  notify.PayAccountUser,
		PayAccountInfo:  notify.PayAccountInfo,
	}
	return notifyRequest
}

// 比apiVersion=1.0时的回调消息体多了 orderType 和 payQRUrl 两个参数
func Notify2ServerNotifyRequestV1dot1(notify models.Notify) response.ServerNotifyRequestV1dot1 {
	var notifyRequest response.ServerNotifyRequestV1dot1
	notifyRequest = response.ServerNotifyRequestV1dot1{
		JrddNotifyId:    notify.JrddNotifyId,
		JrddNotifyTime:  notify.JrddNotifyTime,
		JrddOrderId:     notify.JrddOrderId,
		AppOrderId:      notify.AppOrderId,
		OrderType:       notify.OrderType,
		OrderAmount:     notify.OrderAmount,
		OrderCoinSymbol: notify.OrderCoinSymbol,
		OrderStatus:     int(notify.OrderStatus),
		StatusReason:    int(notify.StatusReason),
		OrderRemark:     notify.OrderRemark,
		OrderPayTypeId:  notify.OrderPayTypeId,
		PayAccountId:    notify.PayAccountId,
		PayQRUrl:        notify.PayQRUrl,
		PayAccountUser:  notify.PayAccountUser,
		PayAccountInfo:  notify.PayAccountInfo,
	}
	return notifyRequest
}

//post回调消息给平台商
func PostNotifyToServer(order models.Order, notify models.Notify) (resp *http.Response, err error) {

	resp = &http.Response{}
	ul, _ := url.Parse(order.AppServerNotifyUrl)
	distributorId := strconv.FormatInt(order.DistributorId, 10)

	scheme := ul.Scheme
	// utils.Log.Debugf("appServerNotifyUrl's scheme is :[%v]", scheme)

	//兼容http及https两种格式
	var client *http.Client
	//client := &http.Client{}

	client = &http.Client{
		Timeout: 3 * time.Second,
	}
	if scheme == "https" {
		//证书认证
		pool := x509.NewCertPool()
		//根据配置文件读取证书
		caCrt := DownloadPem(distributorId)
		// utils.Log.Debugf("capem is: %v", caCrt)

		pool.AppendCertsFromPEM(caCrt)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool,
				InsecureSkipVerify: true},
		}

		client = &http.Client{Transport: tr, Timeout: 3 * time.Second}

	}

	var request *http.Request

	if AppIdVersion(distributorId) != "v1.0" {
		// apiVersion=1.1
		var notifyRequest response.ServerNotifyRequestV1dot1
		notifyRequest = Notify2ServerNotifyRequestV1dot1(notify)

		// 构建回调url
		serverUrl, err := BuildServerUrlForApi1dot1(order, notify)
		if serverUrl == "" {
			utils.Log.Errorf("buildServerUrl wrong,err=[%v]", err)
			return nil, err
		}
		str := ServerNotifyRequestV1dot12Urlencoded(notifyRequest, distributorId, "")
		utils.Log.Debugf("send to distributor server request with new signing method is :[%v] ", str)

		request, err = http.NewRequest(http.MethodPost, serverUrl, strings.NewReader(str))
		request.Header.Set(CONTENT_TYPE, "application/x-www-form-urlencoded")
	} else {
		// apiVersion=1.0
		var notifyRequest response.ServerNotifyRequest
		notifyRequest = Notify2ServerNotifyRequest(notify)

		jsonData, err := json.Marshal(notifyRequest)
		if err != nil {
			utils.Log.Errorf("order convert to json wrong,[%v]", err)
		}
		var binBody = bytes.NewReader(jsonData)

		// 构建回调url
		serverUrl, err := BuildServerUrl(order, notify)
		if serverUrl == "" {
			utils.Log.Errorf("buildServerUrl wrong,err=[%v]", err)
			return nil, err
		}
		request, err = http.NewRequest(http.MethodPost, serverUrl, binBody)
		if err != nil {
			utils.Log.Errorf("http.NewRequest wrong, err:%v", err)
			resp.Status = FAIL
			return resp, err
		}
		//orderStatus := order.Status
		Headers(request)
		utils.Log.Debugf("send to distributor server request is [%v] ", request)

	}

	resp, err = client.Do(request)
	if err != nil || resp == nil {
		utils.Log.Errorf("there is something wrong when visit distributor server,%v", err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if len(body) > 20 {
		utils.Log.Debugf("send to distributor server, response body is [%v], length is %d, only first 20 chars is printed", string(body[:20]), len(body))
	} else {
		utils.Log.Debugf("send to distributor server, response body is [%v]", string(body))
	}
	bodyStr := fmt.Sprintf("%s", body)
	if err == nil && bodyStr == SUCCESS {
		resp.Status = SUCCESS
		return resp, nil
	}

	resp.Status = FAIL
	resp.StatusCode = 200
	return resp, nil

}

func ServerNotifyRequestV1dot12Urlencoded(notifyRequest response.ServerNotifyRequestV1dot1) string {
	params := make(map[string]string)

	params["jrddNotifyId"] = notifyRequest.JrddNotifyId
	params["jrddNotifyTime"] = strconv.FormatInt(notifyRequest.JrddNotifyTime, 10)
	params["jrddOrderId"] = notifyRequest.JrddOrderId
	params["appOrderId"] = notifyRequest.AppOrderId
	params["orderType"] = strconv.Itoa(notifyRequest.OrderType)
	params["orderAmount"] = fmt.Sprintf("%.2f", notifyRequest.OrderAmount)
	params["orderCoinSymbol"] = notifyRequest.OrderCoinSymbol
	params["orderStatus"] = strconv.Itoa(notifyRequest.OrderStatus)
	params["statusReason"] = strconv.Itoa(notifyRequest.StatusReason)
	params["orderRemark"] = notifyRequest.OrderRemark
	params["orderPayTypeId"] = strconv.Itoa(int(notifyRequest.OrderPayTypeId))
	params["payAccountId"] = notifyRequest.PayAccountId
	params["payQRUrl"] = notifyRequest.PayQRUrl
	params["payAccountUser"] = notifyRequest.PayAccountUser
	params["payAccountInfo"] = notifyRequest.PayAccountInfo
	str := BuildOrderParams(params)
	utils.Log.Debugf("the ServerNotifyRequestV1dot1 convert to Urlencoded result is [%+v]", str)
	return str
}

//判断appid的版本
func AppIdVersion(distributorId string) string {
	version := "v1.1"
	oldAppId := utils.Config.GetStringMapString("appidversion")
	for appid, vs := range oldAppId {
		if appid == distributorId {
			version = vs
			break

		}
	}
	return version
}

func BuildServerUrlForApi1dot1(order models.Order, notify models.Notify) (string, error) {
	var serverUrl string
	var notifyRequest response.ServerNotifyRequestV1dot1
	notifyRequest = Notify2ServerNotifyRequestV1dot1(notify)

	utils.Log.Debugf("send to distributor server origin requestbody is notifyRequestStr=[%v]", notifyRequest)
	notifyRequestStr, _ := Struct2JsonString(notifyRequest)
	utils.Log.Debugf("send to distributor server requestbody is notifyRequestStr=[%v]", notifyRequestStr)
	distributorId := strconv.FormatInt(order.DistributorId, 10)

	var distributor models.Distributor
	if err := utils.DB.First(&distributor, "distributors.id = ?", order.DistributorId).Error; err != nil {
		utils.Log.Errorf("func BuildServerUrlForApi1dot1, not found distributor err:%v", err)
		return "", err
	}

	apiKey := distributor.ApiKey
	secretKey := distributor.ApiSecret
	params := make(map[string]string)
	params["apiKey"] = apiKey
	params["appId"] = distributorId
	params["jrddInputCharset"] = "UTF-8"
	params["jrddSignType"] = "HMAC-SHA256"
	urlStr := BuildOrderParams(params)

	params["jrddNotifyId"] = notifyRequest.JrddNotifyId
	params["jrddNotifyTime"] = strconv.FormatInt(notifyRequest.JrddNotifyTime, 10)
	params["jrddOrderId"] = notifyRequest.JrddOrderId
	params["appOrderId"] = notifyRequest.AppOrderId
	params["orderType"] = strconv.Itoa(notifyRequest.OrderType)
	params["orderAmount"] = fmt.Sprintf("%.2f", notifyRequest.OrderAmount)
	params["orderCoinSymbol"] = notifyRequest.OrderCoinSymbol
	params["orderStatus"] = strconv.Itoa(notifyRequest.OrderStatus)
	params["statusReason"] = strconv.Itoa(notifyRequest.StatusReason)
	params["orderRemark"] = notifyRequest.OrderRemark
	params["orderPayTypeId"] = strconv.Itoa(int(notifyRequest.OrderPayTypeId))
	params["payAccountId"] = notifyRequest.PayAccountId
	params["payQRUrl"] = notifyRequest.PayQRUrl
	params["payAccountUser"] = notifyRequest.PayAccountUser
	params["payAccountInfo"] = notifyRequest.PayAccountInfo
	str := BuildOrderParams(params)
	utils.Log.Debugf("func BuildServerUrlForApi1dot1, the str before signing with new method is [%v]", str)

	jrddSignContent, _ := HmacSha256Base64Signer(str, secretKey)
	utils.Log.Debugf("the signing result with new method is ,jrddSignContent = [%v]", jrddSignContent)
	serverUrl += order.AppServerNotifyUrl + "?" + urlStr + "&jrddSignContent=" + jrddSignContent
	// utils.Log.Debugf("send to distributor server url is serverUrl=[%v]", serverUrl)
	return serverUrl, nil

}

//构建回调url
func BuildServerUrl(order models.Order, notify models.Notify) (string, error) {
	var serverUrl string
	var notifyRequest response.ServerNotifyRequest

	//notifyRequest = Order2ServerNotifyReq(order)
	notifyRequest = Notify2ServerNotifyRequest(notify)

	utils.Log.Debugf("send to distributor server origin requestbody is notifyRequestStr=[%v]", notifyRequest)
	notifyRequestStr, _ := Struct2JsonString(notifyRequest)
	utils.Log.Debugf("send to distributor server requestbody is notifyRequestStr=[%v]", notifyRequestStr)
	distributorId := strconv.FormatInt(order.DistributorId, 10)

	var distributor models.Distributor
	if err := utils.DB.First(&distributor, "distributors.id = ?", order.DistributorId).Error; err != nil {
		utils.Log.Errorf("func BuildServerUrl, not found distributor err:%v", err)
		return "", err
	}

	//签名
	apiKey := distributor.ApiKey
	secretKey := distributor.ApiSecret
	originUrl := order.AppServerNotifyUrl
	ul, _ := url.Parse(originUrl)
	path := ul.Path
	params := make(map[string]string)
	params["apiKey"] = apiKey
	params["appId"] = distributorId
	params["jrddInputCharset"] = "UTF-8"
	params["jrddSignType"] = "HMAC-SHA256"
	str := BuildOrderParams(params)

	//str := "apiKey=" + apiKey + "&appId=" + distributorId + "&jrddInputCharset=UTF-8&jrddSignType=HMAC-SHA256"
	urlStr := path + "?" + str

	notifyRequestSignStr := GenSignatureWith3(http.MethodPost, urlStr, notifyRequestStr)
	utils.Log.Debugf("the str to sign when sending message to distributor server is :[%v] ", notifyRequestSignStr)

	jrddSignContent, _ := HmacSha256Base64Signer(notifyRequestSignStr, secretKey)
	// utils.Log.Debugf("jrddSignContent is [%v]", jrddSignContent)
	serverUrl += order.AppServerNotifyUrl + "?" + str + "&jrddSignContent=" + jrddSignContent
	// utils.Log.Debugf("send to distributor server url is serverUrl=[%v]", serverUrl)
	return serverUrl, nil

}

//时间轮  处理重试回调
func SendNotifyWheel1(data interface{}) {
	var order models.Order
	ni, ok := data.(notifyItem)
	if !ok {
		utils.Log.Warnf("Notify wheel throws unknown notify item")
		return
	}
	if err := utils.DB.First(&order, "order_number=?", ni.orderID).Error; err != nil {
		utils.Log.Errorf("get order by order_number wrong,err=[%v]", err)
		return
	}
	notify := GetNotify(ni)
	utils.Log.Debugf("get notify by order in  notifywheel1 is ,notify =[%v]", notify)
	if notify.SendStatus == 2 {
		utils.Log.Debugf("message has been sent successfully before notifywheel1 ")
		return
	}
	resp, err := PostNotifyToServer(order, notify)
	utils.Log.Debugf("notifywheel1 begin to run -----------------------")
	if err == nil && resp.Status == SUCCESS {
		utils.Log.Debugf("send message to distributor success,time 1,serverUrl is: [%s]", notify.AppServerNotifyUrl)
		notify.Attempts += 1
		notify.SendStatus = 2
		if err := utils.DB.Model(&models.Notify{}).Updates(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor success,time 1,but update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
	} else {
		utils.Log.Errorf("send message to distributor fail,time 1,serverUrl is: [%s],err is:[%v]", notify.AppServerNotifyUrl, err)
		notify.Attempts += 1
		if err := utils.DB.Model(&models.Notify{}).Update(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor fail,time 1,and update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
		notifyWheel2.Add(notifyItem{
			orderID:           order.OrderNumber,
			orderStatus:       order.Status,
			orderStatusReason: order.StatusReason})
	}

}

func SendNotifyWheel2(data interface{}) {

	var order models.Order
	ni, ok := data.(notifyItem)
	if !ok {
		utils.Log.Warnf("Notify wheel throws unknown notify item")
		return
	}
	if err := utils.DB.First(&order, "order_number=? and status=? and status_reason=?", ni.orderID, ni.orderStatus, ni.orderStatusReason).Error; err != nil {
		utils.Log.Errorf("get order by order_number wrong,err=[%v]", err)
		return
	}

	notify := GetNotify(ni)
	utils.Log.Debugf("get notify by order in  notifywheel2 is ,notify =[%v]", notify)

	if notify.SendStatus == 2 {
		utils.Log.Debugf("message has been sent successfully before notifywheel2 ")
		return
	}
	resp, err := PostNotifyToServer(order, notify)
	utils.Log.Debugf("notifywheel2 begin to run -----------------------")
	if err == nil && resp.Status == SUCCESS {
		utils.Log.Debugf("send message to distributor success,time 2,serverUrl is: [%s]", notify.AppServerNotifyUrl)
		notify.Attempts += 1
		notify.SendStatus = 2
		if err := utils.DB.Model(&notify).Updates(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor success,time 2,but update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
	} else {
		utils.Log.Errorf("send message to distributor fail,time 2,serverUrl is: [%s],err is:[%v]", notify.AppServerNotifyUrl, err)
		notify.Attempts += 1
		if err := utils.DB.Model(&notify).Update(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor fail,time 2,and update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
		notifywheel3.Add(notifyItem{
			orderID:           order.OrderNumber,
			orderStatus:       order.Status,
			orderStatusReason: order.StatusReason})
	}

}

func SendNotifyWheel3(data interface{}) {

	var order models.Order
	ni, ok := data.(notifyItem)
	if !ok {
		utils.Log.Warnf("Notify wheel throws unknown notify item")
		return
	}
	if err := utils.DB.First(&order, "order_number=? and status=? and status_reason=?", ni.orderID, ni.orderStatus, ni.orderStatusReason).Error; err != nil {
		utils.Log.Errorf("get order by order_number wrong,err=[%v]", err)
		return
	}

	notify := GetNotify(ni)
	utils.Log.Debugf("get notify by order in  notifywheel3 is ,notify =[%v]", notify)

	if notify.SendStatus == 2 {
		utils.Log.Debugf("message has been sent successfully before notifywheel3 ")
		return
	}
	resp, err := PostNotifyToServer(order, notify)
	utils.Log.Debugf("notifywheel3 begin to run -----------------------")
	if err == nil && resp.Status == SUCCESS {
		utils.Log.Debugf("send message to distributor success,time 3,serverUrl is: [%s]", notify.AppServerNotifyUrl)
		notify.Attempts += 1
		notify.SendStatus = 2
		if err := utils.DB.Model(&notify).Updates(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor success,time 3,but update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
	} else {
		utils.Log.Errorf("send message to distributor fail,time 3,serverUrl is: [%s],err is:[%v]", notify.AppServerNotifyUrl, err)
		notify.Attempts += 1
		if err := utils.DB.Model(&notify).Update(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor fail,time 3,and update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
		notifywheel4.Add(notifyItem{
			orderID:           order.OrderNumber,
			orderStatus:       order.Status,
			orderStatusReason: order.StatusReason})
	}

}
func SendNotifyWheel4(data interface{}) {

	var order models.Order
	ni, ok := data.(notifyItem)
	if !ok {
		utils.Log.Warnf("Notify wheel throws unknown notify item")
		return
	}
	if err := utils.DB.First(&order, "order_number=? and status=? and status_reason=?", ni.orderID, ni.orderStatus, ni.orderStatusReason).Error; err != nil {
		utils.Log.Errorf("get order by order_number wrong,err=[%v]", err)
		return
	}

	notify := GetNotify(ni)
	utils.Log.Debugf("get notify by order in  notifywheel4 is ,notify =[%v]", notify)

	if notify.SendStatus == 2 {
		utils.Log.Debugf("message has been sent successfully before notifywheel4 ")
		return
	}
	resp, err := PostNotifyToServer(order, notify)
	utils.Log.Debugf("notifywheel4 begin to run -----------------------")
	if err == nil && resp.Status == SUCCESS {
		utils.Log.Debugf("send message to distributor success,time 4,serverUrl is: [%s]", notify.AppServerNotifyUrl)
		notify.Attempts += 1
		notify.SendStatus = 2
		if err := utils.DB.Model(&notify).Updates(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor success,time 4,but update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
	} else {
		utils.Log.Errorf("send message to distributor fail,time 4,serverUrl is: [%s],err is:[%v]", notify.AppServerNotifyUrl, err)
		notify.Attempts += 1
		if err := utils.DB.Model(&notify).Update(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor fail,time 4,and update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
		notifywheel5.Add(notifyItem{
			orderID:           order.OrderNumber,
			orderStatus:       order.Status,
			orderStatusReason: order.StatusReason})
	}

}
func SendNotifyWheel5(data interface{}) {
	var order models.Order
	ni, ok := data.(notifyItem)
	if !ok {
		utils.Log.Warnf("Notify wheel throws unknown notify item")
		return
	}
	if err := utils.DB.First(&order, "order_number=? and status=? and status_reason=?", ni.orderID, ni.orderStatus, ni.orderStatusReason).Error; err != nil {
		utils.Log.Errorf("get order by order_number wrong,err=[%v]", err)
		return
	}

	notify := GetNotify(ni)
	utils.Log.Debugf("get notify by order in  notifywheel5 is ,notify =[%v]", notify)

	if notify.SendStatus == 2 {
		utils.Log.Debugf("message has been sent successfully before notifywheel5 ")
		return
	}
	resp, err := PostNotifyToServer(order, notify)
	utils.Log.Debugf("notifywheel5 begin to run -----------------------")
	if err == nil && resp.Status == SUCCESS {
		utils.Log.Debugf("send message to distributor success,time 5,serverUrl is: [%s]", notify.AppServerNotifyUrl)
		notify.Attempts += 1
		notify.SendStatus = 2
		if err := utils.DB.Model(&notify).Updates(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor success,time 5,but update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
	} else {
		utils.Log.Errorf("send message to distributor fail,time 5,serverUrl is: [%s],err is:[%v]", notify.AppServerNotifyUrl, err)
		notify.Attempts += 1
		if err := utils.DB.Model(&notify).Update(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor fail,time 5,and update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
		notifywheel6.Add(notifyItem{
			orderID:           order.OrderNumber,
			orderStatus:       order.Status,
			orderStatusReason: order.StatusReason})
	}

}
func SendNotifyWheel6(data interface{}) {

	var order models.Order
	ni, ok := data.(notifyItem)
	if !ok {
		utils.Log.Warnf("Notify wheel throws unknown notify item")
		return
	}
	if err := utils.DB.First(&order, "order_number=? and status=? and status_reason=?", ni.orderID, ni.orderStatus, ni.orderStatusReason).Error; err != nil {
		utils.Log.Errorf("get order by order_number wrong,err=[%v]", err)
		return
	}

	notify := GetNotify(ni)
	utils.Log.Debugf("get notify by order in  notifywheel6 is ,notify =[%v]", notify)

	if notify.SendStatus == 2 {
		utils.Log.Debugf("message has been sent successfully before notifywheel6 ")
		return
	}
	resp, err := PostNotifyToServer(order, notify)
	utils.Log.Debugf("notifywheel6 begin to run -----------------------")
	if err == nil && resp.Status == SUCCESS {
		utils.Log.Debugf("send message to distributor success,time 6,serverUrl is: [%s]", notify.AppServerNotifyUrl)
		notify.Attempts += 1
		notify.SendStatus = 2
		if err := utils.DB.Model(&notify).Updates(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor success,time 6,but update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
	} else {
		utils.Log.Errorf("send message to distributor fail,time 6,serverUrl is: [%s],err is:[%v]", notify.AppServerNotifyUrl, err)
		notify.Attempts += 1
		if err := utils.DB.Model(&notify).Update(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor fail,time 6,and update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
		notifywheel7.Add(notifyItem{
			orderID:           order.OrderNumber,
			orderStatus:       order.Status,
			orderStatusReason: order.StatusReason})
	}

}
func SendNotifyWheel7(data interface{}) {
	var order models.Order
	ni, ok := data.(notifyItem)
	if !ok {
		utils.Log.Warnf("Notify wheel throws unknown notify item")
		return
	}
	if err := utils.DB.First(&order, "order_number=? and status=? and status_reason=?", ni.orderID, ni.orderStatus, ni.orderStatusReason).Error; err != nil {
		utils.Log.Errorf("get order by order_number wrong,err=[%v]", err)
		return
	}

	notify := GetNotify(ni)
	utils.Log.Debugf("get notify by order in  notifywheel7 is,notify =[%v]", notify)

	if notify.SendStatus == 2 {
		utils.Log.Debugf("message has been sent successfully before notifywheel7 ")
		return
	}
	resp, err := PostNotifyToServer(order, notify)
	utils.Log.Debugf("notifywheel7 begin to run -----------------------")
	if err == nil && resp.Status == SUCCESS {
		utils.Log.Debugf("send message to distributor success,time 7,serverUrl is: [%s]", notify.AppServerNotifyUrl)
		notify.Attempts += 1
		notify.SendStatus = 2
		if err := utils.DB.Model(&notify).Updates(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor success,time 7,but update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
	} else {
		utils.Log.Errorf("send message to distributor fail,time 7,serverUrl is: [%s],err is:[%v]", notify.AppServerNotifyUrl, err)
		notify.Attempts += 1
		if err := utils.DB.Model(&notify).Update(notify).Error; err != nil {
			utils.Log.Errorf("send message to distributor fail,time 7,and update notify wrong ,notify=[%v],err=[%v]", notify, err)
		}
	}

}

/*
  build http get request params, and order
  eg:
    params := make(map[string]string)
	params["bb"] = "222"
	params["aa"] = "111"
	params["cc"] = "333"
  return string: eg: aa=111&bb=222&cc=333
*/
func BuildOrderParams(params map[string]string) string {
	urlParams := url.Values{}
	for k := range params {
		urlParams.Add(k, params[k])
	}
	return urlParams.Encode()
}

//初始化7个回调消息时间轮
func InitSendNotifyWheel() {

	key := utils.UniqueTimeWheelKey("resendnotify1")
	utils.Log.Debugf("notifywheel1 init")
	notifywheel1 = timewheel.New(1*time.Second, 240, key, SendNotifyWheel1)
	notifywheel1.Start()

	key = utils.UniqueTimeWheelKey("resendnotify2")
	utils.Log.Debugf("notifywheel2 init")
	notifyWheel2 = timewheel.New(1*time.Second, 600, key, SendNotifyWheel2)
	notifyWheel2.Start()

	key = utils.UniqueTimeWheelKey("resendnotify3")
	utils.Log.Debugf("notifywheel3 init")
	notifywheel3 = timewheel.New(1*time.Second, 600, key, SendNotifyWheel3)
	notifywheel3.Start()

	key = utils.UniqueTimeWheelKey("resendnotify4")
	utils.Log.Debugf("notifywheel4 init")
	notifywheel4 = timewheel.New(1*time.Minute, 60, key, SendNotifyWheel4)
	notifywheel4.Start()

	key = utils.UniqueTimeWheelKey("resendnotify5")
	utils.Log.Debugf("notifywheel5 init")
	notifywheel5 = timewheel.New(1*time.Minute, 120, key, SendNotifyWheel5)
	notifywheel5.Start()

	key = utils.UniqueTimeWheelKey("resendnotify6")
	utils.Log.Debugf("notifywheel6 init")
	notifywheel6 = timewheel.New(1*time.Minute, 540, key, SendNotifyWheel6)
	notifywheel6.Start()

	key = utils.UniqueTimeWheelKey("resendnotify7")
	utils.Log.Debugf("notifywheel7 init")
	notifywheel7 = timewheel.New(1*time.Minute, 900, key, SendNotifyWheel7)
	notifywheel7.Start()

}
