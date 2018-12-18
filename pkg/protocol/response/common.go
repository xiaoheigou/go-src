package response

import (
	"github.com/jinzhu/gorm"
	"math"
	"strconv"
	"yuudidi.com/pkg/utils"
)

const StatusSucc = "success"
const StatusFail = "fail"

type CommonRet struct {
	// status可以为success或者fail
	Status string `json:"status" binding:"required" example:"success"`
	// err_msg仅在失败时设置
	ErrMsg string `json:"err_msg" example:"由于xx原因，导致操作失败"`
	// err_code仅在失败时设置
	ErrCode int `json:"err_code" example:"1001"`
}

type PageResponse struct {
	EntityResponse
	Pagination
}

type EntityResponse struct {
	CommonRet
	Data interface{} `json:"data"`
}


// Pagination paging the list data
type Pagination struct {
	TotalCount int         `json:"total_count"`
	PageSize   int         `json:"page_size"`
	PageCount  int         `json:"page_count"`
	PageNum    int         `json:"page_num"`
}

// Paginate - method to execute pagination
func (p *Pagination) Paginate(pageNum string, pageSize string) *gorm.DB {
	p.PageSize, _ = strconv.Atoi(pageSize)
	p.PageCount = int(math.Ceil(float64(p.TotalCount) / float64(p.PageSize)))
	p.PageSize = int(math.Max(1, math.Min(10000, float64(p.PageSize))))

	p.PageNum, _ = strconv.Atoi(pageNum)
	p.PageNum = int(math.Max(1, float64(p.PageNum)))

	return utils.DB.Offset(p.PageSize * (p.PageNum - 1)).Limit(p.PageSize)
}
