package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/bitly/go-simplejson"
	"sort"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/utils"
)

func DealTicket(body []byte) response.CommonRet {
	var ret response.CommonRet
	var tickets models.Tickets
	var ticketUpdate models.TicketUpdate
	reqBody, _ := simplejson.NewJson(body)
	ticketType, _ := reqBody.Get("type").String()
	ticketBody, _ := reqBody.Get("body").Encode()

	//创建工单表

	tickets = CreateTickets(ticketBody, ticketType)

	if err := utils.DB.Create(&tickets).Error; err != nil {
		utils.Log.Errorf("create tickets error,err:[%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateTicketsErr.Data()
		return ret
	}

	//创建工单描述表
	ticketUpdate = CreateTicketsUpdate(ticketBody, ticketType)
	if err := utils.DB.Create(&ticketUpdate).Error; err != nil {
		utils.Log.Errorf("create ticketUpdate error,err:[%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateTicketUpdateErr.Data()
		return ret
	}

	return ret

}

func CreateTickets(body []byte, ticketType string) models.Tickets {
	var tickets models.Tickets

	js2, _ := simplejson.NewJson([]byte(body))
	//工单id
	ticket_id, _ := js2.Get("ticket_id").String()
	//工单号
	ticket_no, _ := js2.Get("ticket_no").String()
	//工单主题
	subject, _ := js2.Get("creator_id").String()
	//操作人id
	operator_id, _ := js2.Get("operator_id").String()
	//操作人
	operator, _ := js2.Get("operator").String()
	//操作人类型
	operator_type, _ := js2.Get("operator_type").String()
	//工单创建者
	creator_id, _ := js2.Get("creator_id").String()
	//附件
	attachments, _ := js2.Get("detail_info").Get("attachments").Encode()
	content, _ := js2.Get("detail_info").Get("content").String()
	attachment := string(attachments)
	//订单号
	order_number, _ := js2.Get("detail_info").Get("form_value").Get("order_number").String()
	//申诉类型
	apply_type, _ := js2.Get("detail_info").Get("form_value").Get("apply_type").String()
	//国家与地区代码
	country_code, _ := js2.Get("detail_info").Get("form_value").Get("country_code").String()
	//联系电话
	phone, _ := js2.Get("detail_info").Get("form_value").Get("phone").String()
	//申诉原因
	apply_msg, _ := js2.Get("detail_info").Get("form_value").Get("apply_msg").String()
	//银行卡号
	bank_card, _ := js2.Get("detail_info").Get("form_value").Get("bank_card").String()
	//开户银行
	bank_name, _ := js2.Get("detail_info").Get("form_value").Get("bank_name").String()
	//持卡人
	card_user, _ := js2.Get("detail_info").Get("form_value").Get("card_user").String()
	//微信支付账号
	wx_account, _ := js2.Get("detail_info").Get("form_value").Get("wx_account").String()
	//支付宝支付账号
	ali_account, _ := js2.Get("detail_info").Get("form_value").Get("ali_account").String()

	tickets = models.Tickets{
		TicketId:     ticket_id,
		OrderNumber:  order_number,
		TicketType:   ticketType,
		Content:      content,
		TicketNo:     ticket_no,
		Subject:      subject,
		Operator:     operator,
		OperatorId:   operator_id,
		OperatorType: operator_type,
		CreatorId:    creator_id,
		Attachments:  attachment,

		ApplyType: apply_type,
		//国家与地区代码
		CountryCode: country_code,
		//联系电话
		Phone: phone,
		//申诉原因
		ApplyMsg: apply_msg,
		//银行卡号
		BankCard: bank_card,
		//开户银行
		BankName: bank_name,
		//持卡人
		CardUser: card_user,
		//微信支付账号
		WxAccount: wx_account,
		//支付宝支付账号
		AliAccount: ali_account,
	}

	return tickets

}

func CreateTicketsUpdate(body []byte, ticketType string) models.TicketUpdate {
	var ticketUpdate models.TicketUpdate
	var description string
	var note string
	var entry string
	js2, _ := simplejson.NewJson([]byte(body))
	//工单号
	ticket_id, _ := js2.Get("ticket_id").String()
	//操作人昵称
	nickname, _ := js2.Get("operator").String()
	if ticketType == "event" {
		//主题变更
		topic, _ := js2.Get("events").Get("topic").StringArray()
		if topic != nil {
			topicChange := nickname + "把主题从:" + topic[0] + " 变化为:" + topic[1] + ";\n"

			description += topicChange
		}
		//工单来源变更
		source, _ := js2.Get("events").Get("source").StringArray()
		if source != nil {
			sourceChange := nickname + "把工单来源从:" + source[0] + " 变化为:" + source[1] + ";\n"
			description += sourceChange
		}
		//优先级变更
		priority, _ := js2.Get("events").Get("priority").StringArray()
		if priority != nil {
			priorityChange := nickname + "把优先级从:" + priority[0] + " 变化为:" + priority[1] + ";\n"
			description += priorityChange
		}
		////工单数据字段发生变化
		//fields, _ := js2.Get("events").Get("fields").StringArray()
		//工单状态被更改
		status, _ := js2.Get("events").Get("status").StringArray()
		if status != nil {
			statusChange := nickname + "把单状态从:" + status[0] + " 变化为:" + status[1] + ";\n"
			description += statusChange
		}
		//工单被转移到部门名称
		transferred, _ := js2.Get("events").Get("transferred").StringArray()
		if transferred != nil {
			transferredChange := nickname + "把单被从部门:" + transferred[0] + " 转移到部门:" + transferred[1] + ";\n"
			description += transferredChange
		}
		//工单被分配到团队名称或者坐席名称
		assigned, _ := js2.Get("events").Get("assigned").StringArray()
		if assigned != nil {
			assignedChange := nickname + "把单被从:" + assigned[0] + " 分配到:" + assigned[1] + ";\n"
			description += assignedChange
		}
		//SLA计划从A变更到B
		sla, _ := js2.Get("events").Get("sla").StringArray()
		if sla != nil {
			slaChange := nickname + "把SLA计划从:" + sla[0] + " 变更到:" + sla[1] + ";\n"
			description += slaChange

		}
		//工单所有人变更从A到B
		user, _ := js2.Get("events").Get("user").StringArray()
		if user != nil {
			userChange := nickname + "把工单所有人从:" + user[0] + " 变更为:" + user[1] + ";\n"
			description += userChange
		}
		//工单过期时间从A变更到B
		duedate, _ := js2.Get("events").Get("duedate").StringArray()
		if duedate != nil {
			duedateChange := nickname + "把工单过期时间从:" + duedate[0] + "变更到:" + duedate[1] + "\n"
			description += duedateChange

		}

	} else if ticketType == "note" {
		note = js2.Get("events").Get("note").MustString()
		description += nickname + "留言:" + note

	} else {
		entry = js2.Get("events").Get("entry").MustString()
		description += nickname + "回复:" + entry

	}

	ticketUpdate = models.TicketUpdate{
		//工单id
		TicketId:   ticket_id,
		TicketType: ticketType,
		//工单变化描述
		Description: description,
		//操作人昵称
		Nickname: nickname,
		//Note:     note,
		//Entry:    entry,
	}

	return ticketUpdate

}

//对字符串做排序
func SortString(token string, timestamp string, nonce string, message string) string {

	var str string
	strList := []string{token, timestamp, nonce, message}
	sort.Strings(strList)
	for i := 0; i < len(strList); i++ {
		str += strList[i]
	}
	//str:=strList[0]+strList[1]+strList[2]+strList[3]

	return str
}

//sha1 加密
func Sha1(data string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return fmt.Sprintf("%x", sha1.Sum(nil))
}

func GetTicket(orderNumber string) response.EntityResponse {
	var ret response.EntityResponse
	var ticket models.Tickets

	if err := utils.DB.First(&ticket, "order_number = ?", orderNumber).Error; err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.NotFoundTicketErr.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = []models.Tickets{ticket}
	return ret
}

func GetTicketUpdates(page, size, startTime, stopTime, sort, timeField, search, ticketId string) response.PageResponse {
	var ret response.PageResponse
	var result []models.TicketUpdate
	db := utils.DB.Order(fmt.Sprintf("%s %s", timeField, sort))
	if search != "" {
		db = db.Where("nickname like ?", search+"%")
	} else {
		db = db.Where("ticket_id = ?", ticketId)
		if startTime != "" && stopTime != "" {
			db = db.Where(fmt.Sprintf("%s >= ? AND %s <= ?", timeField, timeField), startTime, stopTime)
		}
		db.Model(&models.TicketUpdate{}).Count(&ret.TotalCount)
		pageNum, err := strconv.ParseInt(page, 10, 64)
		pageSize, err1 := strconv.ParseInt(size, 10, 64)
		if err != nil || err1 != nil {
			utils.Log.Error(pageNum, pageSize)
		}
		db = db.Offset((pageNum - 1) * pageSize).Limit(pageSize)
		ret.PageNum = int(pageNum)
		ret.PageSize = int(pageSize)
	}
	db.Find(&result)
	ret.PageCount = len(result)

	ret.Data = result
	ret.Status = response.StatusSucc
	return ret
}
