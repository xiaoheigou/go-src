package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/service/dbcache"
	"yuudidi.com/pkg/utils"
)

// FulfillOrderByMerchant - selected merchant to fulfill the order
func FulfillOrderByMerchant(order OrderToFulfill, merchantID int64, seq int) (*OrderFulfillment, error) {
	var merchant models.Merchant
	if err := dbcache.GetMerchantById(merchantID, &merchant); err != nil {
		utils.Log.Errorf("find merchant(uid=[%d]) fail. [%v]", merchantID, err)
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

	orderToUpdate := models.Order{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&orderToUpdate, "order_number = ?", order.OrderNumber).RecordNotFound() {
		tx.Rollback()
		return nil, fmt.Errorf("Record not found of order number: %s", order.OrderNumber)
	}

	if !(orderToUpdate.Status == models.NEW || orderToUpdate.Status == models.WAITACCEPT) {
		// 订单处于除NEW和WAITACCEPT的其它状态，可能已经被其它人提前抢单了
		tx.Rollback()
		return nil, errors.New("already accepted by others")
	}

	if err := tx.Create(&fulfillment).Error; err != nil {
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
	if err := tx.Create(&fulfillmentLog).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if order.Direction == 0 { //Trader Buy, lock merchant quantity of crypto coins
		//lock merchant quote & payment in_use
		asset := models.Assets{}
		if tx.First(&asset, "merchant_id = ? AND currency_crypto = ? ", merchantID, order.CurrencyCrypto).RecordNotFound() {
			tx.Rollback()
			return nil, fmt.Errorf("Can't find corresponding asset record of merchant_id %d, currency_crypto %s", merchantID, order.CurrencyCrypto)
		}

		// 平台用户的充值订单，币商接单后，把币商的数字币从quantity列转移到qty_frozen中
		if rowsAffected := tx.Table("assets").Where("id = ? and quantity >= ?", asset.Id, order.Quantity).
			Updates(map[string]interface{}{"quantity": asset.Quantity - order.Quantity, "qty_frozen": asset.QtyFrozen + order.Quantity}).RowsAffected; rowsAffected == 0 {
			utils.Log.Errorf("Can't freeze %f %s for merchant (id=%d), asset for merchant = [%+v]", order.Quantity, order.CurrencyCrypto, merchant.Id, asset)
			tx.Rollback()
			return nil, fmt.Errorf("can't freeze %f %s for merchant (id=%d)", order.Quantity, order.CurrencyCrypto, merchant.Id)
		}

		//if err := tx.Model(&payment).Update("in_use", 1).Error; err != nil {
		//	tx.Rollback()
		//	return nil, err
		//}
		if err := tx.Model(&orderToUpdate).Updates(models.Order{MerchantId: merchant.Id, Status: models.ACCEPTED, MerchantPaymentId: payment.Id}).Error; err != nil {
			//at this timepoint only update merchant & status, payment info would be updated only once completed
			tx.Rollback()
			return nil, err
		}
	} else {
		// 金融滴滴平台赚的佣金。
		// 以平台用户提现1000个BTUSD为例，金融滴滴平台赚的BTUSD为：
		// 1000 - 6.35 * 1000 * (1 + 0.01) / 6.5 = 13.3076923077
		// var platformCommisionQty float64 = order.Quantity - (priceBuy * order.Quantity * (1 + 0.01) / priceSell)

		if err := tx.Model(&orderToUpdate).Updates(models.Order{
			MerchantId: merchant.Id,
			Status:     models.ACCEPTED}).Error; err != nil {
			//at this timepoint only update merchant & status, payment info would be updated only once completed
			tx.Rollback()
			return nil, err
		}
	} //do nothing for Direction = 1, Trader Sell
	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("error tx in func FulfillOrderByMerchant commit, err=[%v]", err)
	}
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
	utils.Log.Debugf("func GetBestPaymentID begin, merchantID = [%v]", merchantID)
	if order.Direction == 1 { //Trader Sell, no need to pick for merchant payment id
		return models.PaymentInfo{}
	}
	amount := order.Amount
	payT := order.PayType // 1 - wechat, 2 - zhifubao 4 - bank, combination also supported
	payments := []models.PaymentInfo{}
	whereClause := "uid = ? AND audit_status = 1 /**audit passed**/ AND in_use = 0 /**not in use**/ AND (e_amount = ? OR e_amount = 0) "
	types := []string{}
	types = append(types, strconv.FormatInt(int64(payT), 10))

	//if payT&1 != 0 { //wechat
	//	types = append(types, "1")
	//}
	//if payT&2 != 0 { //zfb
	//	types = append(types, "2")
	//}
	//if payT&4 != 0 { //bank
	//	types = append(types, "4")
	//}
	//payTypeStr := bytes.Buffer{}
	//payTypeStr.WriteString("(" + strings.Join(types, ",") + ")")
	//whereClause = whereClause + payTypeStr.String()

	db := utils.DB.Model(&models.PaymentInfo{}).Order("e_amount DESC").Limit(1)

	if payT >= 4 {
		db = db.Where("pay_type >= ?", 4)
	} else if payT > 0 {
		db = db.Where("pay_type = ?", payT)
	}
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
