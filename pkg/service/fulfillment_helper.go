package service

import (
	"encoding/json"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

// 按merchants接单时间排序
func sortMerchantsByLastOrderTime(merchants []int64, direction int) []int64 {
	var redisSorted []string
	var redisSortedInt64 []int64
	var err error

	// 按redis中保存的merchants的接单时间，对merchants进行排序（接单早的排在前面）
	// 如果merchant还没有接过单，则在redis中没有记录，它也不会出现在结果集redisSorted中
	if redisSorted, err = utils.GetMerchantsSortedByLastOrderTime(direction); err != nil {
		utils.Log.Error("func sortMerchantsByLastOrderTime fail, call GetMerchantsSortedByLastOrderTime fail [%v]", err)
		return merchants
	}
	if err := utils.ConvertStringToInt(redisSorted, &redisSortedInt64); err != nil {
		utils.Log.Error("func sortMerchantsByLastOrderTime fail, call ConvertStringToInt fail [%v]", err)
		return merchants
	}

	var merchantsWithoutSuccOrder = utils.DiffSet(merchants, redisSortedInt64) // 从未接过单的merchants

	return append(merchantsWithoutSuccOrder, utils.InterSetInt64(redisSortedInt64, merchants)...)
}

func getOfficialMerchants() []int64 {
	officialMerchants := []int64{}

	// 先从redis读取
	if officialMerchantsStr, err := utils.GetCacheSetMembers(utils.RedisKeyMerchantRole1()); err != nil {
		utils.ConvertStringToInt(officialMerchantsStr, &officialMerchants)
	}

	// 读不到，则从db中读取
	if len(officialMerchants) == 0 {
		db := utils.DB.Model(&models.Merchant{}).Where("role = 1")
		if err := db.Pluck("id", &officialMerchants).Error; err != nil {
			utils.Log.Errorf("getOfficialMerchants from db failed.")
		}

		// 保存到redis中
		for _, officialMerchant := range officialMerchants {
			expireTimeInSecond := 600 // 10分钟过期，过期后重新从数据库读取
			if err := utils.SetCacheSetMember(utils.RedisKeyMerchantRole1(), expireTimeInSecond, officialMerchant); err != nil {
				utils.Log.Errorf("add official Merchant %s to redis fail, err", officialMerchant, err)
			}
		}
	}

	utils.Log.Debugf("official merchants :%v", officialMerchants)
	return officialMerchants
}

func getAutoConfirmPaidFromMessage(msg models.Msg) (merchant int64, amount float64) {
	//get merchant, amount, ts from msg.data
	if d, ok := msg.Data[0].(map[string]interface{}); ok {
		mn, ok := d["merchant_id"].(json.Number)
		if ok {
			merchant, _ = mn.Int64()
		}
		an, ok := d["amount"].(json.Number)
		if ok {
			amount, _ = an.Float64()
		}
	}
	return merchant, amount
}

func getOrderNumberAndDirectionFromMessage(msg models.Msg) (orderNumber string, direction int) {
	//get order number from msg.data.order_number
	if d, ok := msg.Data[0].(map[string]interface{}); ok {
		orderNumber = d["order_number"].(string)
		if dn, ok := d["direction"].(json.Number); ok {
			d64, _ := dn.Int64()
			direction = int(d64)
		}
	}
	return orderNumber, direction
}
