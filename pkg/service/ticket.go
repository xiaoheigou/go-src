package service

import (
	"fmt"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/utils"
)

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
		db.Count(&ret.TotalCount)
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
	return ret
}
