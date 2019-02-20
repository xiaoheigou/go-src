package service

import (
	"encoding/json"
	"regexp"
	"strings"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/utils"
)

var RmbPatten, _ = regexp.Compile("\\d+\\.\\d+")

func parseAlipayBillData(billData string, receivedBill *models.ReceivedBill) error {
	// 支付宝的账单数据格式如下：
	// {"content":"￥0.02","assistMsg1":"二维码收款到账通知","assistMsg2":"jrId:162918537667547921","linkName":"","buttonLink":"","templateId":"WALLET-FWC@remindDefaultText"}

	type AlipyBillData struct {
		Content    string `json:"content"` // 金额在这个字段中
		AssistMsg1 string `json:"assistMsg1"`
		AssistMsg2 string `json:"assistMsg2"` // 备注在这个字段中
	}

	var data AlipyBillData
	if err := json.Unmarshal([]byte(billData), &data); err != nil {
		utils.Log.Errorln("unmarshal alipay bill data fail, err %s", err)
		return err
	}

	content := strings.TrimSpace(data.Content)
	if strings.HasPrefix(content, "￥") {

	}

	receivedBill.OrderNumber = ""
	receivedBill.Amount = 1

	return nil
}

func UploadBills(uid int64, arg response.UploadBillArg) response.CommonRet {
	var ret response.CommonRet

	if arg.PayType == models.PaymentTypeWeixin {

	} else if arg.PayType == models.PaymentTypeAlipay {
		for _, bill := range arg.Data {

			var receivedBill models.ReceivedBill

			receivedBill.UploaderUid = uid
			receivedBill.UserPayId = bill.UserPayId
			receivedBill.BillId = bill.BillId
			receivedBill.BillData = bill.BillData

			// 从bill.BillData中分析金额和jrdidi订单号
			parseAlipayBillData(bill.BillData, &receivedBill)

			// 保存到数据库
			if err := utils.DB.Save(&receivedBill).Error; err != nil {
				utils.Log.Errorf("UploadBills fail, db err [%v]", err)
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
				return ret
			}
		}
	} else {

	}

	return ret
}
