package response

type CommonRet struct {
	// status可以为success或者fail
	Status string `json:"status" binding:"required" example:"success"`
	// err_msg仅在失败时设置
	ErrMsg string `json:"err_msg" example:"由于xx原因，导致操作失败"`
	// err_code仅在失败时设置
	ErrCode int `json:"err_code" example:"1001"`
}

type PageResponse struct {
	CommonRet
	Data      []interface{} `json:"data"`
	PageNum   int           `json:"page_num"`
	PageSize  int           `json:"page_size"`
	PageCount int           `json:"page_total"`
}

type EntityResponse struct {
	CommonRet
	Data []interface{}
}
