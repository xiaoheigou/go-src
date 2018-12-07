package response

type HandleComplaintsArgs struct {
	Note string `json:"note" binding:"required" example:"处理意见"`
}

type HandleComplaintsRet struct {
	CommonRet
	Entity struct{
		Id       string `json:"id"`
		IssuedBy int    `json:"issued_by"`
		//平台用户id
		Account int `json:"account" example:"1"`
		//平台商id
		DistributorId int    `json:"distributor_id" example:"1"`
		MerchantId    int    `json:"merchant_id" example:"1"`
		OrderNumber   int    `json:"order_number" example:"1"`
		OrderType     int    `json:"order_type" example:"1"`
		PayType       int    `json:"pay_type" example:"1"`
		Status        string `json:"status" example:"1"`
		CreatedAt     string `json:"created_at" example:"2016-09-21T08:50:08"`
		EndTime       string `json:"end_time" example:"2016-09-21T08:50:08"`
		Content       string `json:"content" example:"申诉内容详情"`
		Note          string `json:"note" example:"处理意见"`
	}
}

type GetComplaintsRet struct {
	CommonRet
	Entity struct {
		Id       string `json:"id"`
		IssuedBy int    `json:"issued_by"`
		//平台用户id
		Account int `json:"account" example:"1"`
		//平台商id
		DistributorId int    `json:"distributor_id" example:"1"`
		MerchantId    int    `json:"merchant_id" example:"1"`
		OrderNumber   int    `json:"order_number" example:"1"`
		OrderType     int    `json:"order_type" example:"1"`
		PayType       int    `json:"pay_type" example:"1"`
		Status        string `json:"status" example:"1"`
		CreatedAt     string `json:"created_at" example:"2016-09-21T08:50:08"`
		EndTime       string `json:"end_time" example:"2016-09-21T08:50:08"`
		Content       string `json:"content" example:"申诉内容详情"`
		Note          string `json:"note" example:"处理意见"`
	}
}
