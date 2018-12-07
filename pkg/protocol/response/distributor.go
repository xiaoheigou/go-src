package response

import "YuuPay_core-service/pkg/models"

type GetDistributorsRet struct {
	CommonRet

	Entity struct {
		Data []models.Distributor `json:"data"`
	}
}

type CreateDistributorsRet struct {
	CommonRet
	Entity struct {
	}
}

type CreateDistributorsArgs struct {
	Name      string `json:"name" binding:"required" example:"test"`
	Phone     string `json:"phone" binding:"required" example:"13112345678"`
	Status    int    `json:"status" binding:"required" example:"13112345678"`
	Url       string `json:"url" binding:"required" example:"13112345678"`
	ApiKey    string `json:"api_key" binding:"required" example:"13112345678"`
	ApiSecret string `json:"api_secret" binding:"required" example:"13112345678"`
}

type UpdateDistributorsRet struct {
	CommonRet
	Entity struct {
	}
}

type UpdateDistributorsArgs struct {
	Name      string `json:"name" binding:"required" example:"test"`
	Phone     string `json:"phone" binding:"required" example:"13112345678"`
	Status    int    `json:"status" binding:"required" example:"13112345678"`
	Url       string `json:"url" binding:"required" example:"13112345678"`
	ApiKey    string `json:"api_key" binding:"required" example:"13112345678"`
	ApiSecret string `json:"api_secret" binding:"required" example:"13112345678"`
}
