package service

import (
	"net/http"
	"strconv"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"

	"github.com/360EntSecGroup-Skylar/excelize"
)

/*
  A: 创单时间
  B: 最后状态时间
  C: 订单号
  D: 平台商订单号
  E: btusd 数额
  F: rmb 数目
  G: 买卖
  H: 交易状态
  I: 平台商
  J: 承兑商
  K: 承兑商电话
  L: 付款类型
  M: trader income
  N: merchant income
  O: jrdidi income
*/

func ExportExcel(data []models.Order, w http.ResponseWriter) {
	var excel = utils.Config.GetStringMapString("excel")
	xlsx := excelize.NewFile()

	//循环写入excel表格title
	for k, v := range excel {
		xlsx.SetCellValue("Sheet1", k, v)
	}

	xlsx.SetColWidth("Sheet1", "A", "C", 20)
	//遍历数据
	for i, v := range data {
		index := strconv.FormatInt(int64(i+2), 10)
		xlsx.SetCellValue("Sheet1", "A"+index, v.CreatedAt.In(time.Local))
		xlsx.SetCellValue("Sheet1", "B"+index, v.UpdatedAt.In(time.Local))
		xlsx.SetCellValue("Sheet1", "C"+index, v.OrderNumber)
		xlsx.SetCellValue("Sheet1", "D"+index, v.OriginOrder)
		xlsx.SetCellValue("Sheet1", "E"+index, v.Quantity)
		xlsx.SetCellValue("Sheet1", "F"+index, v.Amount)
		xlsx.SetCellValue("Sheet1", "G"+index, v.Direction)
		xlsx.SetCellValue("Sheet1", "H"+index, StatusReturn(v.Status))
		xlsx.SetCellValue("Sheet1", "I"+index, StatusReasonReturn(v.StatusReason))
		//显示平台商名称，查询的数据一定要有往这个字段里面赋值
		xlsx.SetCellValue("Sheet1", "J"+index, v.DistributorName)
		xlsx.SetCellValue("Sheet1", "K"+index, v.MerchantName)
		xlsx.SetCellValue("Sheet1", "L"+index, v.MerchantPhone)
		xlsx.SetCellValue("Sheet1", "M"+index, GetBankByPayTypId(v.PayType))
		//trade income
		xlsx.SetCellValue("Sheet1", "N"+index, v.TraderBTUSDFeeIncome)
		//merchant income
		xlsx.SetCellValue("Sheet1", "O"+index, v.MerchantBTUSDFeeIncome)
		//jrdidi income
		xlsx.SetCellValue("Sheet1", "P"+index, v.JrdidiBTUSDFeeIncome)
	}

	xlsx.Write(w)
}

//根据order status返回中文订单信息
func StatusReturn(status models.OrderStatus) string {
	switch status {
	case models.NEW:
		return "新建订单"
	case models.ACCEPTED:
		return "已接单"
	case models.NOTIFYPAID:
		return "标记已付款"
	case models.CONFIRMPAID:
		return "标记已收款"
	case models.SUSPENDED:
		return "订单异常"
	case models.PAYMENTMISMATCH:
		return "应收实付不符"
	case models.TRANSFERRED:
		return "订单完成"
	case models.ACCEPTTIMEOUT:
		return "接单超时"
	default:
		return ""
	}
}

func StatusReasonReturn(reason models.StatusReason) string {
	switch reason {
	case models.SYSTEMUPDATEFAIL:
		return "系统异常"
	case models.PAIDTIMEOUT:
		return "标记已付款超时"
	case models.CONFIRMTIMEOUT:
		return "标记已收款超时"
	case models.MARKCOMPLETED:
		return "客服标记完成"
	case models.CANCEL:
		return "客服取消"
	default:
		return ""
	}
}
