package response

import "yuudidi.com/pkg/models"

type GetOrderRet struct {
	CommonRet
	Data []models.Order `json:"data"`
	PageNum int `json:"page_num" example:100`
	PageSize int `json:"page_size" example:10`
	PageCount int `json:"page_count" example:5`
}


type GetOrderDetailRet struct {
	CommonRet
	Data []models.Order `json:"data"`
}