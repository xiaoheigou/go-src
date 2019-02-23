package service

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/utils"
)

var RmbPatten, _ = regexp.Compile("^\\d+\\.\\d\\d$") // 小数点后两位小数

var EngineUsedByAppSvr = NewOrderFulfillmentEngine(nil)

func parseWechatBillData(billData string, receivedBill *models.ReceivedBill) error {
	// 微信的账单数据格式如下：
	/*
		{
		  "showtype": "1",
		  "appid": "",
		  "contentattr": "0",
		  "title": "微信支付收款0.01元(朋友到店)",
		  "url": "https://payapp.weixin.qq.com/payf2f/jumpf2fbill?timestamp=1550819646&openid=AlWTyWBNY3d0rEqoO1pQoM8MBIyVGtau2MCGlDuQf28=",
		  "lowurl": "",
		  "thumburl": "",
		  "ext_pay_info": {
		    "pay_type": "wx_f2f",
		    "pay_outtradeno": "9ZZPNdZE_OPJXNbBXIQYWnOvwG0cPv6OPkPgX06lgJgXws46-RYUjrz3GVPYx4eLAlDCQMXF2Rfj7vMUIUkBzA",
		    "pay_fee": "1",
		    "pay_feetype": "1"
		  },
		  "extinfo": "",
		  "sourcedisplayname": "",
		  "action": "",
		  "template_id": "ey45ZWkUmYUBk_fMgxBLvyaFqVop1rmoWLFd62OXGiU",
		  "appattach": {
		    "fileext": "",
		    "totallen": "0",
		    "aeskey": "",
		    "cdnthumbaeskey": "",
		    "attachid": "",
		    "cdnthumburl": ""
		  },
		  "type": "5",
		  "mmreader": {
		    "category": {
		      "topnew": {
		        "digest": "收款金额￥0.01 收款方备注jrId:184067407654951928 汇总今日第7笔收款，共计￥0.07 备注收款成功，已存入零钱。点击可查看详情",
		        "width": "0",
		        "cover": "",
		        "height": "0"
		      },
		      "count": "1",
		      "item": {
		        "weapp_path": "pages/index/index.html",
		        "del_flag": "0",
		        "fileid": "0",
		        "tweetid": "",
		        "player": "",
		        "music_source": "0",
		        "weapp_state": "0",
		        "pic_urls": "",
		        "show_complaint_button": "0",
		        "play_url": "",
		        "url": "https://payapp.weixin.qq.com/payf2f/jumpf2fbill?timestamp=1550819646&openid=AlWTyWBNY3d0rEqoO1pQoM8MBIyVGtau2MCGlDuQf28=",
		        "shorturl": "",
		        "longurl": "",
		        "digest": "收款金额￥0.01 收款方备注jrId:184067407654951928 汇总今日第7笔收款，共计￥0.07 备注收款成功，已存入零钱。点击可查看详情",
		        "styles": {
		          "style": [
		            {
		              "color": "#000000",
		              "range": "{4,5}",
		              "font": "s"
		            },
		            {
		              "color": "#000000",
		              "range": "{15,23}",
		              "font": "s"
		            },
		            {
		              "color": "#000000",
		              "range": "{41,15}",
		              "font": "s"
		            },
		            {
		              "color": "#000000",
		              "range": "{59,18}",
		              "font": "s"
		            }
		          ],
		          "topColor": ""
		        },
		        "comment_topic_id": "0",
		        "cover": "",
		        "vid": "",
		        "contentattr": "0",
		        "recommendation": "",
		        "pic_num": "0",
		        "cover_235_1": "",
		        "itemshowtype": "4",
		        "weapp_username": "gh_fac0ad4c321d@app",
		        "cover_1_1": "",
		        "pub_time": "1550819646",
		        "appmsg_like_type": "0",
		        "template_op_type": "1",
		        "sources": {
		          "source": {
		            "name": "微信支付"
		          }
		        },
		        "play_length": "0",
		        "title": "收款到账通知",
		        "native_url": "",
		        "weapp_version": "144"
		      },
		      "name": "微信支付",
		      "type": "0"
		    },
		    "template_header": {
		      "pub_time": "1550819646",
		      "hide_time": "1",
		      "pay_style": "1",
		      "show_icon_and_display_name": "0",
		      "title_color": "",
		      "hide_title_and_time": "1",
		      "title": "收款到账通知",
		      "icon_url": "",
		      "hide_icon_and_display_name_line": "1",
		      "display_name": "",
		      "shortcut_icon_url": "",
		      "first_data": "",
		      "header_jump_url": "",
		      "ignore_hide_title_and_time": "1",
		      "first_color": ""
		    },
		    "forbid_forward": "0",
		    "publisher": {
		      "nickname": "微信支付",
		      "username": "wxzhifu"
		    },
		    "template_detail": {
		      "text_content": {
		        "cover": "",
		        "color": "",
		        "text": ""
		      },
		      "line_content": {
		        "topline": {
		          "value": {
		            "color": "#000000",
		            "word": "￥0.01",
		            "small_text_count": "1"
		          },
		          "key": {
		            "hide_dash_line": "1",
		            "color": "#888888",
		            "word": "收款金额"
		          }
		        },
		        "lines": {
		          "line": [
		            {
		              "value": {
		                "color": "#000000",
		                "word": "jrId:184067407654951928"
		              },
		              "key": {
		                "color": "#888888",
		                "word": "收款方备注"
		              }
		            },
		            {
		              "value": {
		                "color": "#000000",
		                "word": "今日第7笔收款，共计￥0.07"
		              },
		              "key": {
		                "color": "#888888",
		                "word": "汇总"
		              }
		            },
		            {
		              "value": {
		                "color": "#000000",
		                "word": "收款成功，已存入零钱。点击可查看详情"
		              },
		              "key": {
		                "color": "#888888",
		                "word": "备注"
		              }
		            }
		          ]
		        }
		      },
		      "template_show_type": "1",
		      "opitems": {
		        "opitem": {
		          "hint_word": "",
		          "op_type": "1",
		          "display_line_number": "0",
		          "weapp_version": "144",
		          "is_rich_text": "0",
		          "word": "收款小账本",
		          "weapp_path": "pages/index/index.html",
		          "url": "",
		          "color": "#000000",
		          "weapp_username": "gh_fac0ad4c321d@app",
		          "icon": "",
		          "weapp_state": "0"
		        },
		        "show_type": "1"
		      }
		    }
		  },
		  "soundtype": "0",
		  "des": "收款金额￥0.01 收款方备注jrId:184067407654951928 汇总今日第7笔收款，共计￥0.07 备注收款成功，已存入零钱。点击可查看详情",
		  "sdkver": "0",
		  "content": "",
		  "sourceusername": ""
		}
	*/

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
					Topline struct {
						Value struct {
							Word string `json:"word"`
						} `json:"value"`
						Key struct {
							Word string `json:"word"`
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
	// 先从"收款方备注"中查找jrdidi订单号
	var orderNumber string
	for _, line := range data.Mmreader.TemplateDetail.LineContent.Lines.Line {
		if line.Key.Word == "收款方备注" {
			orderNumber = utils.GetOrderNumberFromQrCodeMark(line.Value.Word)
			break
		}
	}
	if orderNumber == "" {
		// 备注中找不到订单号
		utils.Log.Infof("can not get order_number for wechat bill %s", receivedBill.BillId)
	}

	receivedBill.OrderNumber = orderNumber
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

	for _, bill := range arg.Data {
		utils.Log.Debugf("upload bill, uploader_uid = %s, pay_type = %d, bill_id = %s", uid, arg.PayType, bill.BillId)

		if bill.BillData == "" {
			var retFail response.CommonRet
			utils.Log.Errorf("bill_data is empty for bill %s", bill.BillId)
			retFail.Status = response.StatusFail
			retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
			return retFail
		}

		var receivedBill models.ReceivedBill

		receivedBill.UploaderUid = uid
		receivedBill.PayType = arg.PayType
		receivedBill.UserPayId = bill.UserPayId
		receivedBill.BillId = bill.BillId
		receivedBill.BillData = bill.BillData

		if arg.PayType == models.PaymentTypeWeixin {
			// 从bill.BillData中分析金额和jrdidi订单号
			if err := parseWechatBillData(bill.BillData, &receivedBill); err != nil {
				utils.Log.Errorf("parse wechat bill data fail, err = %s", err)
			}
		} else if arg.PayType == models.PaymentTypeAlipay {
			// 从bill.BillData中分析金额和jrdidi订单号
			if err := parseAlipayBillData(bill.BillData, &receivedBill); err != nil {
				utils.Log.Errorf("parse alipay bill data fail, err = %s", err)
			}
		} else {
			var retFail response.CommonRet
			utils.Log.Errorf("pay_type %d is invalid, expect 1 or 2", arg.PayType)
			retFail.Status = response.StatusFail
			retFail.ErrCode, retFail.ErrMsg = err_code.AppErrArgInvalid.Data()
			return retFail
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

		// 更新币商自动账号的last_use_time字段
		if err := utils.DB.Model(&models.PaymentInfo{}).Where("uid = ? and user_pay_id = ? and payment_auto_type = 1", uid, bill.UserPayId).
			Update("last_use_time", time.Now()).Error; err != nil {
			utils.Log.Warnf("update last_use_time fail. uid = %s, user_pay_id = %s, err = %s", uid, bill.UserPayId, err)
		}

		if receivedBill.OrderNumber != "" {
			checkBillAndTryConfirmPaid(&receivedBill)
		}
	}

	return ret
}
