package utils

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

// 下面是调用qrcode.decodesvcendpoit返回报文对应的json格式说明
type QrcodeRespMsg struct {
	Code            int    `json:"code"`
	ErrMsg          string `json:"err_msg"`
	Amount          string `json:"amount"` // 二维码图片中的金额
	QrCodeTxt       string `json:"qr_code_txt"` // 解码二维码后的字符串
	NewQrCodeBase64 string `json:"new_qr_code_base64"` // 新生成的二维码，base64编码
}

func GetQrcodeInfo(src io.Reader) (QrcodeRespMsg, error) {
	var resp []byte
	var qrcodeServiceURL = Config.GetString("qrcode.decodesvcendpoit")
	if qrcodeServiceURL == "" {
		qrcodeServiceURL = "http://localhost:8087/qrcode-tool/api/svc/upload"
		Log.Warnln("Wrong configuration: qrcode.decodesvcendpoit is empty, use [%s] as default", qrcodeServiceURL)
	}

	if resp, err = UploadFile(qrcodeServiceURL, "file", "123.jpg", src); err != nil {
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

func IsWeixinQrCode(qrCodeTxt string) bool {
	var weixinPrefix = Config.GetString("qrcode.expectprefix.weixin")
	return strings.HasPrefix(qrCodeTxt, weixinPrefix)
}

func IsAlipayQrCode(qrCodeTxt string) bool {
	var alipayPrefix = Config.GetString("qrcode.expectprefix.alipay")
	return strings.HasPrefix(qrCodeTxt, alipayPrefix)
}