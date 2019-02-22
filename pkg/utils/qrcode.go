package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/url"
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
	bizData := map[string]interface{}{
		"s": "money",
		"u": userPayId,
		"a": amount,
		"m": GenQrCodeMark(orderNumber),
	}
	jsonValue, err := json.Marshal(bizData)
	if err != nil {
		return "err"
	}
	return "alipays://platformapi/startapp?appId=20000123&actionType=scan&biz_data=" + string(jsonValue)
}
