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

		if order.AcceptType == 1 {

			if order.PayType == models.PaymentTypeWeixin || order.PayType == models.PaymentTypeAlipay {
				// 对于自动接单订单，仅收款方式为微信或支付宝时，才采用自动生成的二维码
				payment = GetAutoPaymentID(&order, merchant.Id)
			} else {
				payment = GetBestPaymentID(&order, merchant.Id)
			}
		} else {
			payment = GetBestPaymentID(&order, merchant.Id)
		}

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

	orderFromDb := models.Order{}
	if tx.Set("gorm:query_option", "FOR UPDATE").First(&orderFromDb, "order_number = ?", order.OrderNumber).RecordNotFound() {
		tx.Rollback()
		return nil, fmt.Errorf("Record not found of order number: %s", order.OrderNumber)
	}

	if !(orderFromDb.Status == models.NEW || orderFromDb.Status == models.WAITACCEPT) {
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
		if asset.Quantity.GreaterThanOrEqual(orderFromDb.Quantity) { // 避免 quantity 出现负数
			if err := tx.Table("assets").Where("id = ?", asset.Id).
				Updates(map[string]interface{}{
					"quantity":   asset.Quantity.Sub(orderFromDb.Quantity),
					"qty_frozen": asset.QtyFrozen.Add(orderFromDb.Quantity)}).Error; err != nil {
				utils.Log.Errorf("update asset record for merchant fail, order_number = %s", order.OrderNumber)
				tx.Rollback()
				return nil, fmt.Errorf("update asset record for merchant fail, order_number = %s", order.OrderNumber)
			}
		} else {
			utils.Log.Errorf("Can't freeze %s %s for merchant (id=%d), asset for merchant = [%+v]", order.Quantity, order.CurrencyCrypto, merchant.Id, asset)
			tx.Rollback()
			return nil, fmt.Errorf("can't freeze %s %s for merchant (id=%d)", order.Quantity, order.CurrencyCrypto, merchant.Id)
		}

		//if err := tx.Model(&payment).Update("in_use", 1).Error; err != nil {
		//	tx.Rollback()
		//	return nil, err
		//}
		if err := tx.Model(&orderFromDb).Updates(
			models.Order{
				MerchantId:        merchant.Id,
				Status:            models.ACCEPTED,
				MerchantPaymentId: payment.Id,
				AcceptType:        order.AcceptType,
				QrCodeTxt:         payment.QrCodeTxt,
				UserPayId:         payment.UserPayId,
			}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		if err := tx.Model(&orderFromDb).Updates(models.Order{
			MerchantId: merchant.Id,
			Status:     models.ACCEPTED}).Error; err != nil {
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

func GetAutoPaymentID(order *OrderToFulfill, merchantID int64) models.PaymentInfo {
	payment := models.PaymentInfo{}

	if order.UserPayId == "" {
		utils.Log.Errorf("user_pay_id is empty, it must be set in Android app")
		utils.Log.Errorf("func GetAutoPaymentID finished abnormally.")
		return payment
	}

	// 下面从数据中获取当前币商的"支付id"，生成收款二维码时需要
	var userPayId string

	var merchant models.Merchant
	if err := dbcache.GetMerchantById(merchantID, &merchant); err != nil {
		utils.Log.Errorf("call GetMerchantById fail. [%v]", merchantID, err)
		utils.Log.Errorf("func GetAutoPaymentID finished abnormally. error %s", err)
		return payment
	}

	var pref models.Preferences
	if err := dbcache.GetPreferenceById(int64(merchant.PreferencesId), &pref); err != nil {
		utils.Log.Errorf("can't find preference record in db for merchant(uid=[%d]),  err [%v]", merchantID, err)
		utils.Log.Errorf("func GetAutoPaymentID finished abnormally. error %s", err)
		return payment
	}

	currAutoAlipayPaymentId := pref.CurrAutoAlipayPaymentId
	if err := utils.DB.Where("id = ?", currAutoAlipayPaymentId).First(&payment).Error; err != nil {
		userPayId = payment.UserPayId
	}

	if userPayId == "" {
		utils.Log.Errorf("func GetAutoPaymentID finished abnormally. can not get userPayId from db")
		return payment
	}

	// 如果从Android App传过来的user_pay_id和系统中当前配置的user_pay_id不相同，则报错
	if order.UserPayId != userPayId {
		utils.Log.Errorf("user_pay_id from Android App is %s, but current setting in db is %s, there are mismatched!", order.UserPayId, userPayId)
		utils.Log.Errorf("func GetAutoPaymentID finished abnormally")
		return payment
	}

	if order.PayType == models.PaymentTypeWeixin {
		// 对于微信，Android App要返回收款二维码，没有就报错
		if order.QrCodeTxt == "" {
			utils.Log.Errorf("qr_code_txt from Android App is empty, it must be set in Android app")
			utils.Log.Errorf("func GetAutoPaymentID finished abnormally")
			return payment
		}

		payment.QrCodeTxt = order.QrCodeTxt
		payment.UserPayId = order.UserPayId

	} else if order.PayType == models.PaymentTypeAlipay {
		// 直接在服务端生成二维码
		payment.QrCodeTxt = utils.GenAlipayQrCodeTxt(userPayId, order.Amount, order.OrderNumber)
		payment.UserPayId = order.UserPayId

	} else {
		utils.Log.Errorf("func GetAutoPaymentID finished abnormally. payType %d is not expected", order.PayType)
		return payment
	}

	return payment
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
