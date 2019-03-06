package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

// 下面是调用qrcode.decodesvcendpoit返回报文对应的json格式说明
type QrcodeRespMsg struct {
	Code            int    `json:"code"`
	ErrMsg          string `json:"err_msg"`
	Amount          string `json:"amount"`             // 二维码图片中的金额
	QrCodeTxt       string `json:"qr_code_txt"`        // 解码二维码后的字符串
	NewQrCodeBase64 string `json:"new_qr_code_base64"` // 新生成的二维码，base64编码
}

func GetQrCodeInfo(src io.Reader, fileName, expectedQrCodeTxt string) (QrcodeRespMsg, error) {
	var resp []byte
	var qrcodeServiceURL = Config.GetString("qrcode.decodesvcendpoit")
	if qrcodeServiceURL == "" {
		qrcodeServiceURL = "http://localhost:8087/qrcode-tool/api/svc/upload"
		Log.Warnln("Wrong configuration: qrcode.decodesvcendpoit is empty, use [%s] as default", qrcodeServiceURL)
	}

	if expectedQrCodeTxt != "" {
		if strings.Contains(qrcodeServiceURL, "?") {
			// qrcodeServiceURL中有其它query parm
			qrcodeServiceURL = qrcodeServiceURL + "&expected_qr_code_txt=" + url.QueryEscape(expectedQrCodeTxt)
		} else {
			// qrcodeServiceURL中没有其它query parm
			qrcodeServiceURL = qrcodeServiceURL + "?expected_qr_code_txt=" + url.QueryEscape(expectedQrCodeTxt)
		}
	}

	var err error
	if resp, err = UploadFile(qrcodeServiceURL, "file", fileName, src); err != nil {
		Log.Errorf("upload file to [$s] fail: %v", qrcodeServiceURL, err)
		return QrcodeRespMsg{}, err
	}
	Log.Debugf("qrcode-tool api resp = [%+v]", string(resp[:]))

	var data QrcodeRespMsg
	data.Code = -1
	if err := json.Unmarshal(resp, &data); err != nil {
		Log.Errorf("Unmarshal data fail, error = [%+v]", err)
		return QrcodeRespMsg{}, err
	}

	if data.Code != 0 {
		// qrcodeServiceURL返回json中，code为非0时，表示失败
		return QrcodeRespMsg{}, errors.New(data.ErrMsg)
	}

	return data, nil
}

// 用户校验币商上传的维信二维码是否合法
func IsWeixinQrCode(qrCodeTxt string) bool {
	var weixinPrefix = Config.GetString("qrcode.expectprefix.weixin")
	return strings.HasPrefix(strings.ToUpper(qrCodeTxt), strings.ToUpper(weixinPrefix))
}

// 用户校验币商上传的支付宝二维码是否合法
func IsAlipayQrCode(qrCodeTxt string) bool {
	var alipayPrefix = Config.GetString("qrcode.expectprefix.alipay")
	return strings.HasPrefix(strings.ToUpper(qrCodeTxt), strings.ToUpper(alipayPrefix))
}

// 根据jrdidi订单号生成二维码备注
func GenQrCodeMark(orderNumber string) string {
	return "jrId:" + strings.TrimSpace(orderNumber)
}

// 从备注中得到jrdidi订单号并返回，如果找不到则返回空字符串
func GetOrderNumberFromQrCodeMark(mark string) string {
	var orderNumber string
	remarkWords := strings.TrimSpace(mark)
	if strings.HasPrefix(remarkWords, "jrId:") {
		orderNumber = strings.TrimPrefix(remarkWords, "jrId:")
	}

	return orderNumber
}

// 按指定的金额生成支付宝收款二维码，并把订单号设置在二维码备注中
func GenAlipayQrCodeTxt(userPayId string, amount float64, orderNumber string) string {
	// 需要把amount转换为字符串类型，这时因为：
	// 当amount为数字类型时，二维码为：
	// alipays://platformapi/startapp?appId=20000123&actionType=scan&biz_data={"a":0.14,"m":"jrId:185520141126081528","s":"money","u":"2088002015347730"}
	// 上面这种二维码用Android版本的支付宝扫码可以正常显示指定的金额，不用输入金额。但iPhone版本的支付宝却还提示要输入金额。
	// 当amount为字符串时，二维码为：
	// alipays://platformapi/startapp?appId=20000123&actionType=scan&biz_data={"a":"0.14","m":"jrId:185520141126081528","s":"money","u":"2088002015347730"}
	// 上面这种二维码，Android版本和iPhone版本的支付宝扫码都可以正常显示指定的金额，不用输入金额。
	amountStr := fmt.Sprintf("%.2f", amount)
	bizData := map[string]interface{}{
		"s": "money",
		"u": userPayId,
		"a": amountStr,
		"m": GenQrCodeMark(orderNumber),
	}
	jsonValue, err := json.Marshal(bizData)
	if err != nil {
		return "err"
	}
	return "alipays://platformapi/startapp?appId=20000123&actionType=scan&biz_data=" + string(jsonValue)
}

// 当订单金额是100的倍数，且金额小于等于2000时(由qrcode.fuzzymatch.maxamount配置)，可以应用"随机立减"匹配
func CanApplyFuzzyMatch(amount float64) bool {

	maxAmount := Config.GetString("qrcode.fuzzymatch.maxamount")
	var maxAmountFloat64 float64
	var err error
	if maxAmountFloat64, err = strconv.ParseFloat(maxAmount, 64); err != nil {
		Log.Warnf("qrcode.fuzzymatch.maxamount %s in not expected, use 2000 as default", maxAmount)
		maxAmountFloat64 = 2000
	}
	var maxAmountInt = int(maxAmountFloat64)

	// % 运算只适应于整数，先转换为整数
	var amountStr = fmt.Sprintf("%.2f", amount)
	if !strings.HasSuffix(amountStr, ".00") {
		return false
	}

	var amountInt = int(amount)

	if amountInt > maxAmountInt {
		return false
	}

	if amountInt%100 == 0 { // 是100倍数
		return true
	} else {
		return false
	}
}

// 检测支付宝的支付Id是否是合法的
// 目前通过下面3个条件判断：
// 1、个数为16个；2、全部为数字；3、以208开头
// 这是合法的支付宝的支付Id实例：2088822980739143
func IsValidAlipayUserPayId(userPayId string) bool {

	if len(userPayId) != 16 {
		Log.Warnf("Invalid alipay user pay id [%v], it must be 16 number of digits", userPayId)
		return false
	}

	for _, c := range userPayId {
		if c < '0' || c > '9' {
			Log.Warnf("Invalid alipay user pay id [%v], it can only contains 0-9", userPayId)
			return false
		}
	}

	if !strings.HasPrefix(userPayId, "208") {
		Log.Warnf("Invalid alipay user pay id [%v], it must starts with 208", userPayId)
		return false
	}

	return true
}
