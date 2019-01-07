package service

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/utils"
)

// FulfillOrderByMerchant - selected merchant to fulfill the order
func FulfillOrderByMerchant(order OrderToFulfill, merchantID int64, seq int) (*OrderFulfillment, error) {
	merchant := models.Merchant{}
	if utils.DB.First(&merchant, " id = ?", merchantID).RecordNotFound() {
		utils.Log.Errorf("Record not found of merchant id:", merchantID)
		return nil, fmt.Errorf("Record not found")
	}
	var payment models.PaymentInfo
	var fulfillment models.Fulfillment
	if order.Direction == 0 { //Trader Buy, select payment of merchant
		payment = GetBestPaymentID(&order, merchant.Id)
		//check payment.Id to see if valid payment
		if payment.Id == 0 {
			return nil, fmt.Errorf("No valid payment information found (pay type: %d, amount: %f)", order.PayType, order.Amount)
		}
		fulfillment = models.Fulfillment{
			OrderNumber:       order.OrderNumber,
			SeqID:             seq,
			MerchantID:        merchant.Id,
			MerchantPaymentID: payment.Id,
			AcceptedAt:        time.Now(),
			Status:            models.ACCEPTED,
		}
	} else {
		//Trader Sell, get payment info from order
		payment.PayType = int(order.PayType)
		switch order.PayType {
		case 1: //wechat
			fallthrough
		case 2: //zhifubao
			payment.EAccount = order.Name
			payment.QrCode = order.QrCode
		case 4: //bank
			payment.Bank = order.Bank
			payment.BankAccount = order.BankAccount
			payment.BankBranch = order.BankBranch
		}
		fulfillment = models.Fulfillment{
			OrderNumber: order.OrderNumber,
			SeqID:       seq,
			MerchantID:  merchant.Id,
			AcceptedAt:  time.Now(),
			Status:      models.ACCEPTED,
		}
	}
	tx := utils.DB.Begin()
	if err := utils.DB.Create(&fulfillment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	fulfillmentLog := models.FulfillmentLog{
		FulfillmentID: fulfillment.Id,
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
		tx.Rollback()
		return nil, err
	}
	//update order status
	orderToUpdate := models.Order{}
	if utils.DB.First(&orderToUpdate, "order_number = ?", order.OrderNumber).RecordNotFound() {
		tx.Rollback()
		return nil, fmt.Errorf("Record not found of order number: %s", order.OrderNumber)
	}
	if err := tx.Model(&orderToUpdate).Updates(models.Order{MerchantId: merchant.Id, Status: models.ACCEPTED}).Error; err != nil {
		//at this timepoint only update merchant & status, payment info would be updated only once completed
		tx.Rollback()
		return nil, err
	}
	if order.Direction == 0 { //Trader Buy, lock merchant quantity of crypto coins
		//lock merchant quote & payment in_use
		asset := models.Assets{}
		if utils.DB.First(&asset, "merchant_id = ? AND currency_crypto = ? ", merchantID, order.CurrencyCrypto).RecordNotFound() {
			tx.Rollback()
			return nil, fmt.Errorf("Can't find corresponding asset record of merchant_id %d, currency_crypto %s", merchantID, order.CurrencyCrypto)
		}
		if asset.Quantity < order.Quantity {
			//not enough quantity, return directly
			tx.Rollback()
			return nil, fmt.Errorf("Not enough quote for merchant %d: quantity->%f, amount->%f", merchantID, asset.Quantity, order.Amount)
		}
		if err := tx.Model(&asset).Updates(models.Assets{Quantity: asset.Quantity - order.Quantity, QtyFrozen: asset.QtyFrozen + order.Quantity}).Error; err != nil {
			utils.Log.Errorf("Can't freeze asset record: %v", err)
			tx.Rollback()
			return nil, err
		}
		if err := tx.Model(&payment).Update("in_use", 1).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} //do nothing for Direction = 1, Trader Sell
	tx.Commit()
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
	utils.Log.Debugf("func GetBestPaymentID begin, merchantID = [%v]", merchantID )
	if order.Direction == 1 { //Trader Sell, no need to pick for merchant payment id
		return models.PaymentInfo{}
	}
	amount := order.Amount
	payT := order.PayType // 1 - wechat, 2 - zhifubao 4 - bank, combination also supported
	payments := []models.PaymentInfo{}
	whereClause := "uid = ? AND audit_status = 1 /**audit passed**/ AND in_use = 0 /**not in use**/ AND (e_amount = ? OR e_amount = 0) AND pay_type in "
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
	payTypeStr := bytes.Buffer{}
	payTypeStr.WriteString("(" + strings.Join(types, ",") + ")")
	whereClause = whereClause + payTypeStr.String()

	db := utils.DB.Model(&models.PaymentInfo{}).Order("e_amount DESC").Limit(1)
	db.Where(whereClause, merchantID, amount).Find(&payments)
	//randomly picked one TODO: to support payment list in the future
	count := len(payments)
	if count == 0 {
		return models.PaymentInfo{}
	}
	rand.Shuffle(count, func(i, j int) {
		payments[i], payments[j] = payments[j], payments[i]
	})
	utils.Log.Debugf("func GetBestPaymentID finished normally.")
	return payments[0]
}
