package service

import (
	"fmt"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/utils"
)

func GetDistributors(page, size , status,startTime,stopTime,sort, timeField, search string) (response.GetDistributorsRet,error) {
	var result []models.Distributor
	db := utils.DB.Order(fmt.Sprintf("%s %s",timeField,sort))

	switch {
	case search != "":
		db.Where("name = ? OR id = ?",search,search)
	case page != "" && size != "":
		pageNum,err := strconv.ParseInt(page,10,64)
		pageSize,err1 := strconv.ParseInt(size,10,64)
		if err != nil || err1 != nil{
			utils.Log.Error(pageNum,pageSize)
		}

		fallthrough
	case startTime != "" && stopTime != "" :
		db.Where(fmt.Sprintf("%s >= ? AND %s <= ?",timeField,timeField),startTime,stopTime)
		fallthrough
	case status != "":
		db.Where("status = ?",status)
	}
	db.Find(&result)

	return response.GetDistributorsRet{},nil
}

func createDistributor(distributor models.Distributor) models.Distributor {

	return models.Distributor{}
}
