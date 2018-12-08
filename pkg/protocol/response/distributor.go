package response

import "yuudidi.com/pkg/models"

type GetDistributorsRet struct {
	CommonRet

	Data []models.Distributor `json:"data"`
}

type CreateDistributorsRet struct {
	CommonRet
	Data []interface{}
}

type CreateDistributorsArgs struct {
	Name      string `json:"name" binding:"required" example:"test"`
	Phone     string `json:"phone" binding:"required" example:"13112345678"`
	Status    int    `json:"status" binding:"required" example:"1"`
	Url       string `json:"url" binding:"required" example:"13112345678"`
	ApiKey    string `json:"api_key" binding:"required" example:"13112345678"`
	ApiSecret string `json:"api_secret" binding:"required" example:"13112345678"`
}

type UpdateDistributorsRet struct {
	CommonRet
	Data []interface{}
}

type UpdateDistributorsArgs struct {
	Name      string `json:"name" binding:"required" example:"test"`
	Phone     string `json:"phone" binding:"required" example:"13112345678"`
	Status    int    `json:"status" binding:"required" example:"13112345678"`
	Url       string `json:"url" binding:"required" example:"13112345678"`
	ApiKey    string `json:"api_key" binding:"required" example:"13112345678"`
	ApiSecret string `json:"api_secret" binding:"required" example:"13112345678"`
}
