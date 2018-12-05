package v1

type CommonRet struct {
	// status可以为success或者fail
	Status   string `json:"status" binding:"required" example:"success"`
	// err_msg仅在失败时设置
	ErrMsg  string `json:"err_msg" example:"您输入的用户名或密码错误"`
	// err_code仅在失败时设置
	ErrCode int `json:"err_code" example:1001`
}

