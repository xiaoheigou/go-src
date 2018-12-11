// +build swagger

package controller
func ReprocessOrder(c *gin.Context){
	var req response.ReprocessOrderRequest
	c.ShouldBind(&req)
	var ret response.ReprocessOrderResponse
	ret.Status=response.StatusSucc
	ret.ErrCode=123
	ret.ErrMsg="reprecess success"
	ret.Data=[]response.ReprocessOrderEntity{
		{
			Url:          "www.otc.com",
			OrderSuccess: "Notify Order Created",
			TotalCount:   "12",
			OrderNo:      "12332",
			OrderType:    "2",
		},
	}
	c.JSON(200,ret)
}
