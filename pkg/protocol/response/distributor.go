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
	Username  string `json:"username" binding:"required" example:"distributor"`
	Password  string `json:"password" binding:"required" example:"distributor password"`
	Phone     string `json:"phone" binding:"required" example:"13112345678"`
	PageUrl   string `json:"page_url" binding:"required" example:"1"`
	ServerUrl string `json:"server_url" binding:"required" example:"1"`
	ApiKey    string `json:"api_key" binding:"required" example:"13112345678"`
	ApiSecret string `json:"api_secret" binding:"required" example:"13112345678"`
}

type UpdateDistributorsRet struct {
	CommonRet
	Data []interface{}
}

type UpdateDistributorsArgs struct {
	Name  string `json:"name" binding:"required" example:"test"`
	Phone string `json:"phone" binding:"required" example:"13112345678"`
	//平台商状态 0: 申请 1: 正常 2: 冻结
	Status    int    `json:"status" binding:"required" example:"1"`
	PageUrl   string `json:"page_url" binding:"required" example:"1"`
	ServerUrl string `json:"server_url" binding:"required" example:"1"`
	ApiKey    string `json:"api_key" binding:"required" example:"13112345678"`
	ApiSecret string `json:"api_secret" binding:"required" example:"13112345678"`
}
