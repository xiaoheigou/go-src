package service

import (
	"encoding/json"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/utils"
)

var RmbPatten, _ = regexp.Compile("^\\d+\\.\\d\\d$") // 小数点后两位小数

var EngineUsedByAppSvr = NewOrderFulfillmentEngine(nil)

func parseAlipayBillData(billData string, receivedBill *models.ReceivedBill) error {
	// 支付宝的账单数据格式如下：
	// {"content":"￥0.02","assistMsg1":"二维码收款到账通知","assistMsg2":"jrId:162918537667547921","linkName":"","buttonLink":"","templateId":"WALLET-FWC@remindDefaultText"}

	type AlipyBillData struct {
		Content    string `json:"content"` // 人民币金额在这个字段中
		AssistMsg1 string `json:"assistMsg1"`
		AssistMsg2 string `json:"assistMsg2"` // 备注在这个字段中
	}

	var data AlipyBillData
	if err := json.Unmarshal([]byte(billData), &data); err != nil {
		utils.Log.Errorf("unmarshal alipay bill data fail, err %s", err)
		return err
	}

	// 分析账单中的人民币金额
	var amount float64
	content := strings.TrimSpace(data.Content)
	if strings.HasPrefix(content, "￥") {
		rmb := strings.TrimPrefix(content, "￥")
		if RmbPatten.MatchString(rmb) {
			amount, _ = strconv.ParseFloat(rmb, 64)
		} else {
			utils.Log.Errorf("can not get rmb amount from alipay bill data")
			return errors.New("can not get rmb amount from alipay bill data")
		}
	}

	// 分析账单中的备注字段，从中提取出jrdidi订单号
	var orderNumber string
	remarkWords := strings.TrimSpace(data.AssistMsg2)
	if strings.HasPrefix(remarkWords, "jrId:") {
		orderNumber = strings.TrimPrefix(content, "jrId:")
	} else {
		// 备注中没有jrId字样
		utils.Log.Infof("got a alipay bill without jrId:XXX")
	}

	receivedBill.OrderNumber = orderNumber
	receivedBill.Amount = amount

	return nil
}

func checkBillAndTryConfirmPaid(receivedBill *models.ReceivedBill) {

	order := models.Order{}
	if err := utils.DB.First(&order, "order_number = ?", receivedBill.OrderNumber).Error; err != nil {
		utils.Log.Errorf("find order %s error: %s", receivedBill.OrderNumber, err)
		return
	}

	if order.Direction == 1 {
		// 目前，所有自动确认收款的订单都是"用户充值订单"
		utils.Log.Errorf("order %s direction is 1, it is not expected for auto order", order.OrderNumber)
		return
	}

	// TODO
	if receivedBill.Amount >= order.Amount {
		// 自动确认收款
		message := models.Msg{
			MsgType: models.ConfirmPaid,
			Data: []interface{}{
				map[string]interface{}{
					"order_number": order.OrderNumber,
					"direction":    order.Direction,
				},
			},
		}
		EngineUsedByAppSvr.UpdateFulfillment(message)
	}
}

func UploadBills(uid int64, arg response.UploadBillArg) response.CommonRet {
	var ret response.CommonRet

	if arg.PayType == models.PaymentTypeWeixin {
		// TODO
	} else if arg.PayType == models.PaymentTypeAlipay {
		for _, bill := range arg.Data {

			var receivedBill models.ReceivedBill

			receivedBill.UploaderUid = uid
			receivedBill.UserPayId = bill.UserPayId
			receivedBill.BillId = bill.BillId
			receivedBill.BillData = bill.BillData

			// 从bill.BillData中分析金额和jrdidi订单号
			if err := parseAlipayBillData(bill.BillData, &receivedBill); err != nil {
				utils.Log.Errorf("parse alipay bill data fail, err = %s", err)
			}

			// 保存到数据库
			if err := utils.DB.Save(&receivedBill).Error; err != nil {
				utils.Log.Errorf("UploadBills fail, db err [%v]", err)
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
				return ret
			}

			checkBillAndTryConfirmPaid(&receivedBill)
		}
	} else {
		var retFail response.CommonRet
		utils.Log.Errorf("pay_type %d is invalid, expect 1 or 2", arg.PayType)
		retFail.Status = response.StatusFail
		retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
		return retFail
	}

	return ret
}
