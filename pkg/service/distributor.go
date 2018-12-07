package service

import (
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

func GetDistributors(page, size, status, startTime, stopTime, timeField, search string) []models.Distributor {
	var result []models.Distributor

	if search != "" {
		utils.DB.Where("name = ? OR id = ?",search,search).Find(&result)
	}

	return result
}

func createDistributor(distributor models.Distributor) models.Distributor {

	return models.Distributor{}
}
