package service

import (
	"math/rand"
	"strconv"
	"time"

	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

// FulfillOrderByMerchant - selected merchant to fulfill the order
func FulfillOrderByMerchant(order OrderToFulfill, merchantID int64, seq int) (*OrderFulfillment, error) {
	timeoutStr := utils.Config.GetString("fulfillment.timeout.notifypaid")
	timeout, _ := strconv.ParseInt(timeoutStr, 10, 32)
	merchant := models.Merchant{}
	if err := utils.DB.First(&merchant, " id = ?", merchantID).Error; err != nil {
		utils.Log.Errorf("Invalid merchant id to match: %v", err)
		return nil, err
	}
	payment := GetBestPaymentID(&order, merchant.Id)
	fulfillment := models.Fulfillment{
		OrderNumber:       order.OrderNumber,
		SeqID:             seq,
		MerchantID:        merchant.Id,
		MerchantPaymentID: payment.Id,
		AcceptedAt:        time.Now(),
		NotifyPaidBefore:  time.Now().Add(time.Duration(timeout) * time.Second),
		Status:            models.ACCEPTED,
	}
	utils.DB.Begin()
	if err := utils.DB.Create(&fulfillment).Error; err != nil {
		utils.DB.Rollback()
		return nil, err
	}
	fulfillmentLog := models.FulfillmentLog{
		FulfillmentID: fulfillment.ID,
		OrderNumber:   order.OrderNumber,
		SeqID:         seq,
		IsSystem:      true,
		MerchantID:    merchant.Id,
		AccountID:     order.AccountID,
		DistributorID: order.DistributorID,
		OriginStatus:  models.NEW,
		UpdatedStatus: models.ACCEPTED,
	}
	if err := utils.DB.Create(&fulfillmentLog).Error; err != nil {
		utils.DB.Rollback()
		return nil, err
	}
	//update order status
	orderToUpdate := models.Order{}
	if err := utils.DB.First(&orderToUpdate, "order_number = ?", order.OrderNumber).Error; err != nil {
		utils.DB.Rollback()
		return nil, err
	}
	orderToUpdate.Status = models.ACCEPTED
	if err := utils.DB.Update(&orderToUpdate).Error; err != nil {
		utils.DB.Rollback()
		return nil, err
	}
	//lock merchant quote & payment in_use
	payment.InUse = 1
	if err := utils.DB.Update(&payment).Error; err != nil {
		utils.DB.Rollback()
		return nil, err
	}
	utils.DB.Commit()
	return &OrderFulfillment{
		OrderToFulfill:    order,
		MerchantID:        merchant.Id,
		MerchantNickName:  merchant.Nickname,
		MerchantAvatarURI: merchant.AvatarUri,
		PayType:           payment.PayType,
		PaymentInfo:       []models.PaymentInfo{payment},
	}, nil
}

// GetBestPaymentID - get best matched payment id for order:merchant combination
func GetBestPaymentID(order *OrderToFulfill, merchantID int64) models.PaymentInfo {
	if order.Direction == 1 { //Trader Sell, no need to pick for merchant payment id
		return models.PaymentInfo{}
	}
	amount := order.Amount
	payT := order.PayType // 1 - wechat, 2 - zhifubao 4 - bank, combination also supported
	payments := []models.PaymentInfo{}
	whereClause := "uid = ? AND audit_status = 1 /**audit passed**/ AND in_use = 0 /**not in use**/ AND e_amount = ? AND pay_type in "
	types := []string{}
	if payT&1 != 0 { //wechat
		types = append(types, "1")
	}
	if payT&2 != 0 { //zfb
		types = append(types, "2")
	}
	if payT&4 != 0 { //bank
		types = append(types, "4")
	}
	payTypeStr := "("
	for i, t := range types {
		if i == 0 {
			payTypeStr += t
		} else {
			payTypeStr += "," + t
		}
	}
	payTypeStr += ")"
	whereClause = whereClause + payTypeStr
	utils.DB.Find(&payments, whereClause, merchantID, amount)
	//randomly picked one
	count := len(payments)
	if count == 0 {
		return models.PaymentInfo{}
	}
	rand.Shuffle(count, func(i, j int) {
		payments[i], payments[j] = payments[j], payments[i]
	})
	return payments[0]
}
