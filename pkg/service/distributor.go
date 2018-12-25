package service

import (
	"fmt"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

func GetDistributors(page, size, status, startTime, stopTime, sort, timeField, search string) response.PageResponse {
	var result []models.Distributor
	var ret response.PageResponse
	db := utils.DB.Model(&models.Distributor{}).Order(fmt.Sprintf("distributors.%s %s", timeField, sort)).Select("distributors.*,assets.quantity as quantity").Joins("left join assets on distributors.id = assets.distributor_id")
	if search != "" {
		db = db.Where("name = ? OR id = ?", search, search)
	} else {
		pageNum, err := strconv.ParseInt(page, 10, 64)
		pageSize, err1 := strconv.ParseInt(size, 10, 64)
		if err != nil || err1 != nil {
			utils.Log.Error(pageNum, pageSize)
		}
		db = db.Offset((pageNum - 1) * pageSize).Limit(pageSize)
		if startTime != "" && stopTime != "" {
			db = db.Where(fmt.Sprintf("distributors.%s >= ? AND distributors.%s <= ?", timeField, timeField), startTime, stopTime)
		}
		if status != "" {
			db = db.Where("status = ?", status)
		}
		db.Count(&ret.TotalCount)
		ret.PageNum = int(pageNum)
		ret.PageSize = int(pageSize)
	}
	db.Find(&result)
	ret.PageCount = len(result)
	ret.Status = response.StatusSucc
	ret.Data = result
	return ret
}

func CreateDistributor(param response.CreateDistributorsArgs) response.EntityResponse {
	var ret response.EntityResponse
	distributor := models.Distributor{
		Name:      param.Name,
		Phone:     param.Phone,
		ServerUrl: param.ServerUrl,
		PageUrl:   param.PageUrl,
		ApiKey:    param.ApiKey,
		ApiSecret: param.ApiSecret,
	}
	tx := utils.DB.Begin()

	res := CreateUser(response.UserArgs{
		Role:     2,
		Phone:    param.Phone,
		Username: param.Username,
		Password: param.Password,
	}, tx)
	if res.Status == response.StatusFail {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateDistributorErr.Data()
		tx.Rollback()
		return ret
	}
	if err := tx.Create(&distributor).Error; err != nil {
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.CreateDistributorErr.Data()
		tx.Rollback()
		return ret
	}
	ret.Status = response.StatusSucc
	ret.Data = append([]models.Distributor{}, distributor)

	tx.Commit()
	return ret
}

func UpdateDistributor(param response.UpdateDistributorsArgs, uid string) response.EntityResponse {
	var ret response.EntityResponse
	var distributor models.Distributor
	if err := utils.DB.Model(&distributor).Where("distributors.id = ?", uid).Find(&distributor).Error; err != nil {
		utils.Log.Errorf("update distributor find distributor is failed,uid:%s,%v", uid, err)
	} else {
		changeParam := make(map[string]interface{})
		if param.Name != distributor.Name {
			changeParam["name"] = param.Name
		}
		if param.Phone != distributor.Phone {
			changeParam["phone"] = param.Phone
		}
		if param.Status != distributor.Status {
			changeParam["status"] = param.Status
		}
		if param.ServerUrl != distributor.ServerUrl {
			changeParam["server_url"] = param.ServerUrl
		}
		if param.PageUrl != distributor.PageUrl {
			changeParam["page_url"] = param.PageUrl
		}
		if param.ApiKey != distributor.ApiKey {
			changeParam["api_key"] = param.ApiKey
		}
		if param.ApiSecret != distributor.ApiSecret {
			changeParam["api_secret"] = param.ApiSecret
		}
		utils.DB.Model(&distributor).Updates(changeParam)
	}
	ret.Status = response.StatusSucc
	ret.Data = append([]models.Distributor{}, distributor)
	return ret
}

func GetDistributor(uid string) response.EntityResponse {
	var distributors []models.Distributor
	ret := response.EntityResponse{}
	ret.Status = response.StatusSucc
	db := utils.DB.Where("distributors.id = ?", uid).Select("distributors.*,assets.quantity as quantity").Joins("left join assets on distributors.id = assets.distributor_id")

	if err := db.Find(&distributors).Error; err != nil {
		utils.Log.Debugf("err:%v", err)
	}
	ret.Data = distributors
	return ret
}

func GetDistributorByAPIKey(apiKey string) (models.Distributor, error) {
	var distributor models.Distributor

	if err := utils.DB.First(&distributor, "api_key = ?", apiKey).Error; err != nil {
		utils.Log.Debugf("err:%v", err)
		return models.Distributor{},err
	}
	return distributor,nil
}
