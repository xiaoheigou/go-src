package service

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math"
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

func parseWechatBillData(billData string, receivedBill *models.ReceivedBill) error {

	type WechatBillLine struct {
		Value struct {
			Color string `json:"color"`
			Word  string `json:"word"`
		} `json:"value"`
		Key struct {
			Color string `json:"color"`
			Word  string `json:"word"`
		} `json:"key"`
	}
	type WechatBillData struct {
		TemplateId string `json:"template_id"`
		Mmreader   struct {
			TemplateDetail struct {
				LineContent struct {
					Topline struct{
						Value struct {
							Word  string `json:"word"`
						} `json:"value"`
						Key struct {
							Word  string `json:"word"`
						} `json:"key"`
					} `json:"topline"`
					Lines struct {
						Line []WechatBillLine `json:"line"`
					} `json:"lines"`
				} `json:"line_content"`
			} `json:"template_detail"`
		} `json:"mmreader"`
	}

	var data WechatBillData
	if err := json.Unmarshal([]byte(billData), &data); err != nil {
		utils.Log.Errorf("unmarshal wechat bill data fail, err %s", err)
		return err
	}

	// 分析账单中的人民币金额
	var amount float64
	if data.Mmreader.TemplateDetail.LineContent.Topline.Key.Word == "收款金额" {
		value := data.Mmreader.TemplateDetail.LineContent.Topline.Value.Word // "￥0.01"
		if strings.HasPrefix(value, "￥") {
			rmb := strings.TrimPrefix(value, "￥")
			if RmbPatten.MatchString(rmb) {
				amount, _ = strconv.ParseFloat(rmb, 64)
			} else {
				msg := fmt.Sprintf("can not get rmb amount from wechat bill %s", receivedBill.BillId)
				utils.Log.Errorf("%s", msg)
				return errors.New(msg)
			}
		} else {
			msg := fmt.Sprintf("can not get rmb amount from wechat bill %s", receivedBill.BillId)
			utils.Log.Errorf("%s", msg)
			return errors.New(msg)
		}
	}

	// 分析账单中的备注字段，从中提取出jrdidi订单号
	// TODO

	receivedBill.OrderNumber = ""
	receivedBill.Amount = amount

	return nil
}

// 从支付宝账单数据中分析金额和jrdidi订单号
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
			msg := fmt.Sprintf("can not get rmb amount from alipay bill %s, the field is %s", receivedBill.BillId, data.Content)
			utils.Log.Errorf("%s", msg)
			return errors.New(msg)
		}
	} else {
		msg := fmt.Sprintf("can not get rmb amount from alipay bill %s, the field is %s", receivedBill.BillId, data.Content)
		utils.Log.Errorf("%s", msg)
		return errors.New(msg)
	}

	// 分析账单中的备注字段，从中提取出jrdidi订单号
	var orderNumber = utils.GetOrderNumberFromQrCodeMark(data.AssistMsg2)
	if orderNumber == "" {
		// 备注中找不到订单号
		utils.Log.Infof("can not get order_number for alipay bill %s, the mark in bill is %s", receivedBill.BillId, data.AssistMsg2)
	}

	receivedBill.OrderNumber = orderNumber
	receivedBill.Amount = amount

	return nil
}

func rmbCompareEq(v1, v2 float64) bool {
	epsilon := 0.01
	return math.Abs(v1-v2) <= epsilon
}

func rmbCompareGte(v1, v2 float64) bool {
	if rmbCompareEq(v1, v2) {
		return true
	}
	return v1 > v2
}

func checkBillAndTryConfirmPaid(receivedBill *models.ReceivedBill) {

	if receivedBill.OrderNumber == "" {
		utils.Log.Infof("bill %s don't contains jrdidi order number, skip it", receivedBill.BillId)
		return
	}

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

	if rmbCompareGte(receivedBill.Amount, order.Amount) {
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
	} else {
		utils.Log.Warnf("amount in received bill (%f) less than amount in order (%s)", receivedBill.Amount, order.Amount)
	}
}

func UploadBills(uid int64, arg response.UploadBillArg) response.CommonRet {
	var ret response.CommonRet

	if arg.PayType == models.PaymentTypeWeixin {
		for _, bill := range arg.Data {

			if bill.BillData == "" {
				var retFail response.CommonRet
				utils.Log.Errorf("bill_data is empty for bill %s", bill.BillId)
				retFail.Status = response.StatusFail
				retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
				return retFail
			}

			var receivedBill models.ReceivedBill

			receivedBill.UploaderUid = uid
			receivedBill.PayType = models.PaymentTypeWeixin
			receivedBill.UserPayId = bill.UserPayId
			receivedBill.BillId = bill.BillId
			receivedBill.BillData = bill.BillData

			// 从bill.BillData中分析金额和jrdidi订单号
			if err := parseWechatBillData(bill.BillData, &receivedBill); err != nil {
				utils.Log.Errorf("parse wechat bill data fail, err = %s", err)
			}

			// TODO
		}
	} else if arg.PayType == models.PaymentTypeAlipay {
		for _, bill := range arg.Data {

			if bill.BillData == "" {
				var retFail response.CommonRet
				utils.Log.Errorf("bill_data is empty for bill %s", bill.BillId)
				retFail.Status = response.StatusFail
				retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
				return retFail
			}

			var receivedBill models.ReceivedBill

			receivedBill.UploaderUid = uid
			receivedBill.PayType = models.PaymentTypeAlipay
			receivedBill.UserPayId = bill.UserPayId
			receivedBill.BillId = bill.BillId
			receivedBill.BillData = bill.BillData

			// 从bill.BillData中分析金额和jrdidi订单号
			if err := parseAlipayBillData(bill.BillData, &receivedBill); err != nil {
				utils.Log.Errorf("parse alipay bill data fail, err = %s", err)
			}

			// 保存到数据库
			if err := utils.DB.Save(&receivedBill).Error; err != nil {
				// 如果账单之前上传过，并成功保存到数据库，这时再保存会报错误：Duplicate entry 'xxx' for key 'idx_pay_type_bill_id'
				if strings.Contains(err.Error(), "Duplicate entry") {
					// 忽略重复数据
					utils.Log.Infof("bill %s is already uploaded before", bill.BillId)
				} else {
					utils.Log.Errorf("UploadBills fail, db err [%v]", err)
					ret.Status = response.StatusFail
					ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
					return ret
				}
			}

			if receivedBill.OrderNumber != "" {
				checkBillAndTryConfirmPaid(&receivedBill)
			}
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
