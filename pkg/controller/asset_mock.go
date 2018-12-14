// +build swagger

package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
)

func GetMerchantAssetHistory(c *gin.Context) {
	var ret response.GetMerchantAssetHistoryRet
	ret.Status = "success"
	ret.ErrCode = 0
	ret.ErrMsg = "test"
	data := models.AssetHistory{}
	data.Id = 1
	ret.Data = []models.AssetHistory{data}
	c.JSON(200, ret)
}

func GetDistributorAssetHistory(c *gin.Context) {
	ret := `{
  "status": "success",
  "err_msg": "",
  "err_code": 0,
  "data": [
    {
      "id": 1,
      "merchant_id": 0,
      "distributor_id": 1,
      "order_number": 0,
      "is_order": 0,
      "operation": 1,
      "currency": "BTUSD",
      "quantity": 100,
      "operator_id": 1,
      "operator_name": "string",
      "create_at": "2018-12-14T17:08:43+08:00"
    }
  ],
  "page_num": 1,
  "page_size": 1,
  "page_total": 1
}`
	var result map[string]interface{}
	json.Unmarshal([]byte(ret),&result)
	c.JSON(200, result)
}

func GetRechargeApplies(c *gin.Context) {
	ret := `{
		"status": "success",
		"err_msg": "",
		"err_code": 0,
		"data": [
		{
		"id": 1,
		"merchant_id": 1,
		"phone": "",
		"email": "",
		"status": 0,
		"currency": "BTUSD",
		"quantity": 200,
		"remain_quantity": 0,
		"apply_id": 1,
		"apply_name": "string",
		"auditor_id": 2,
		"auditor_name": 0,
		"create_at": "0001-01-01T00:00:00Z"
		}
	],
		"page_num": 1,
		"page_size": 1,
		"page_total": 1
	}`
	var result map[string]interface{}
	json.Unmarshal([]byte(ret),&result)
	c.JSON(200, result)
}