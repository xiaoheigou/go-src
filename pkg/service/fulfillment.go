package service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"math/rand"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/service/dbcache"
	"yuudidi.com/pkg/utils"
)

// FulfillOrderByMerchant - selected merchant to fulfill the order
func FulfillOrderByMerchant(order OrderToFulfill, merchantID int64, seq int) (*OrderFulfillment, error) {
	utils.Log.Debugf("func FulfillOrderByMerchant, order = %+v, merchantID = %d, seq = %d", order, merchantID, seq)

	var merchant models.Merchant
	if err := dbcache.GetMerchantById(merchantID, &merchant); err != nil {
		utils.Log.Errorf("find merchant(uid=[%d]) fail. [%v]", merchantID, err)
		return nil, fmt.Errorf("Record not found")
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

	var payment models.PaymentInfo
	var fulfillment models.Fulfillment
	if order.Direction == 0 { //Trader Buy, select payment of merchant
		if order.AcceptType == 1 {
			if order.PayType == models.PaymentTypeWeixin || order.PayType == models.PaymentTypeAlipay {
				// 对于自动接单订单，才采用自动生成的二维码
				payment = GetAutoPaymentID(tx, &order, merchant.Id)
			} else {
				payment = GetBestNormalPaymentID(tx, &order, merchant.Id)
			}
		} else {
			payment = GetBestNormalPaymentID(tx, &order, merchant.Id)
		}

		//check payment.Id to see if valid payment
		if payment.Id == 0 {
			tx.Rollback()
			return nil, fmt.Errorf("no valid payment information found (pay type: %d, accept_type: %d, amount: %f)",
				order.PayType, order.AcceptType, order.Amount)
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

	actualAmount := orderFromDb.Amount // 默认，实际金额就是订单金额
	if order.Direction == 0 {          //Trader Buy, lock merchant quantity of crypto coins
		//lock merchant quote & payment in_use
		asset := models.Assets{}
		if tx.Set("gorm:query_option", "FOR UPDATE").First(&asset, "merchant_id = ? AND currency_crypto = ? ", merchantID, order.CurrencyCrypto).RecordNotFound() {
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

		if order.PayType == models.PaymentTypeWeixin || order.PayType == models.PaymentTypeAlipay {
			// 固定金额二维码都有占用锁定状态，但针对“不固定金额”二维码，因被同时派到相同金额订单的概率较小，暂时不考虑占用锁定
			if payment.EAmount > 0 {
				if err := tx.Model(&payment).Update("in_use", 1).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}

		if (order.PayType == models.PaymentTypeWeixin || order.PayType == models.PaymentTypeAlipay) && payment.EAmount > 0 {
			// 出现"随机立减"二维码收款方式时，实际金额可能小于订单金额
			actualAmount = payment.EAmount
		}
		if err := tx.Model(&orderFromDb).Updates(
			models.Order{
				MerchantId:        merchant.Id,
				Status:            models.ACCEPTED,
				MerchantPaymentId: payment.Id,
				Bank:              payment.Bank,
				BankAccount:       payment.BankAccount,
				BankBranch:        payment.BankBranch,
				AcceptType:        order.AcceptType,
				QrCode:            payment.QrCode,
				QrCodeTxt:         payment.QrCodeTxt,
				Name:              payment.Name,
				UserPayId:         payment.UserPayId,
				ActualAmount:      actualAmount,
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

	// 有的二维码是币商上传的，有的是Android自动生成的，有的是服务器自动生成的，它们都在payment中，复制到order中，保证order的QrCodeTxt总有值
	order.QrCodeTxt = payment.QrCodeTxt
	order.Name = payment.Name

	order.ActualAmount = actualAmount
	return &OrderFulfillment{
		OrderToFulfill:    order,
		MerchantID:        merchant.Id,
		MerchantNickName:  merchant.Nickname,
		MerchantAvatarURI: merchant.AvatarUri,
		PayType:           payment.PayType,
		PaymentInfo:       []models.PaymentInfo{payment},
	}, nil
}

func GetAutoPaymentID(tx *gorm.DB, order *OrderToFulfill, merchantID int64) models.PaymentInfo {
	payment := models.PaymentInfo{}

	// 目前不检查前端有没有传支付ID
	//if order.UserPayId == "" {
	//	utils.Log.Errorf("user_pay_id is empty, it must be set in Android app. order = %s, pay_type = %d", order.OrderNumber, order.PayType)
	//	utils.Log.Errorf("func GetAutoPaymentID finished abnormally. user_pay_id is empty, order = %s, merchant = %d", order.OrderNumber, merchantID)
	//	return models.PaymentInfo{}
	//}

	// 下面从数据中获取当前币商的"支付id"，生成收款二维码时需要
	var userPayId string

	var merchant models.Merchant
	if err := dbcache.GetMerchantById(merchantID, &merchant); err != nil {
		utils.Log.Errorf("call GetMerchantById fail. [%v]", merchantID, err)
		utils.Log.Errorf("func GetAutoPaymentID finished abnormally. error %s", err)
		return models.PaymentInfo{}
	}

	var pref models.Preferences
	if err := dbcache.GetPreferenceById(int64(merchant.PreferencesId), &pref); err != nil {
		utils.Log.Errorf("can't find preference record in db for merchant(uid=[%d]),  err [%v]", merchantID, err)
		utils.Log.Errorf("func GetAutoPaymentID finished abnormally. error %s", err)
		return models.PaymentInfo{}
	}

	if order.PayType == models.PaymentTypeWeixin {
		if order.QrCodeFromSvr == 1 {
			// 使用币商上传的二维码
			return GetBestNormalPaymentID(tx, order, merchant.Id)
		} else { // 使用App返回的二维码
			currAutoWechatPaymentId := pref.CurrAutoWeixinPaymentId
			if err := utils.DB.Where("id = ?", currAutoWechatPaymentId).First(&payment).Error; err != nil {
				utils.Log.Errorf("can't find payment info in db for merchant(uid=[%d]),  err [%v]", merchantID, err)
				utils.Log.Errorf("func GetAutoPaymentID finished abnormally. error %s", err)
				return models.PaymentInfo{}
			}

			userPayId = payment.UserPayId

			if userPayId == "" {
				utils.Log.Errorf("func GetAutoPaymentID finished abnormally. can not get userPayId from db, order = %s, merchant = %d", order.OrderNumber, merchantID)
				return models.PaymentInfo{}
			}

			// 如果从Android App传过来的user_pay_id和系统中当前配置的user_pay_id不相同，则报错
			if order.UserPayId != userPayId {
				utils.Log.Warnf("for order %s, merchant %d, user_pay_id from Android App is %s, but current setting in db is %s, there are mismatched!", order.OrderNumber, merchantID, order.UserPayId, userPayId)
				// utils.Log.Errorf("func GetAutoPaymentID finished abnormally. user_pay_id from Android App is mismatched with db, order = %s, merchant = %d", order.OrderNumber, merchantID)
				// return models.PaymentInfo{}
			}

			// 对于微信，Android App要返回收款二维码，没有就报错
			if order.QrCodeTxt == "" {
				utils.Log.Errorf("qr_code_txt from Android App is empty, it must be set in Android app. order = %s", order.OrderNumber)
				utils.Log.Errorf("func GetAutoPaymentID finished abnormally. qr_code_txt from Android App is empty, order = %s, merchant = %d", order.OrderNumber, merchantID)
				return models.PaymentInfo{}
			}

			// 对于微信，使用Android App端传过来的二维码
			payment.QrCodeTxt = order.QrCodeTxt
			payment.UserPayId = userPayId
		}
	} else if order.PayType == models.PaymentTypeAlipay {
		// 获取当前的支付ID
		//currAutoAlipayPaymentId := pref.CurrAutoAlipayPaymentId
		//if err := utils.DB.Where("id = ?", currAutoAlipayPaymentId).First(&payment).Error; err != nil {
		//	utils.Log.Errorf("can't find payment info in db for merchant(uid=[%d]),  err [%v]", merchantID, err)
		//	utils.Log.Errorf("func GetAutoPaymentID finished abnormally. error %s", err)
		//	return models.PaymentInfo{}
		//}
		//
		//userPayId = payment.UserPayId

		var userPayIds []string
		db := utils.DB.Model(&models.PaymentInfo{}).Where("uid = ? AND pay_type = ? AND payment_auto_type = 1 AND enable = 1", merchantID, models.PaymentTypeAlipay)
		if err := db.Pluck("user_pay_id", &userPayIds).Error; err != nil {
			utils.Log.Errorf("Gets user_pay_id list fail for merchant %d. err = %s", merchantID, err)
			return models.PaymentInfo{}
		}

		if len(userPayIds) == 0 {
			utils.Log.Errorf("Gets user_pay_id list fail for merchant %d. list is empty", merchantID)
			return models.PaymentInfo{}
		}

		// 下面对userPayIds去重、同时去除非法的支付宝支付Id，保存到validUserPayIds中
		keys := make(map[string]bool)
		validUserPayIds := []string{}
		for _, entry := range userPayIds {
			if _, ok := keys[entry]; !ok {
				keys[entry] = true
				if utils.IsValidAlipayUserPayId(entry) {
					validUserPayIds = append(validUserPayIds, entry)
				}
			}
		}

		if len(validUserPayIds) == 0 {
			utils.Log.Errorf("Gets user_pay_id list fail for merchant %d. validUserPayIds is empty after filter out invalid items", merchantID)
			return models.PaymentInfo{}
		}

		// 从validUserPayIds中随机选一个
		rand.Shuffle(len(validUserPayIds), func(i, j int) {
			validUserPayIds[i], validUserPayIds[j] = validUserPayIds[j], validUserPayIds[i]
		})
		userPayId = validUserPayIds[0]

		// 对于支付宝，直接在服务端生成二维码
		payment.QrCodeTxt = utils.GenAlipayQrCodeTxt(userPayId, order.Amount, order.OrderNumber)
		payment.UserPayId = userPayId

	} else {
		utils.Log.Errorf("func GetAutoPaymentID finished abnormally. payType %d is not expected", order.PayType)
		return models.PaymentInfo{}
	}

	return payment
}

// GetBestNormalPaymentID - get best matched payment id for order:merchant combination
func GetBestNormalPaymentID(tx *gorm.DB, order *OrderToFulfill, merchantID int64) models.PaymentInfo {
	utils.Log.Debugf("func GetBestNormalPaymentID begin, order %s, merchantID = [%v]", order.OrderNumber, merchantID)
	if order.Direction == 1 { //Trader Sell, no need to pick for merchant payment id
		return models.PaymentInfo{}
	}
	amount := order.Amount
	payT := order.PayType // 1 - wechat, 2 - zhifubao 4 - bank, combination also supported
	payment := models.PaymentInfo{}

	db := tx.Set("gorm:query_option", "FOR UPDATE")

	if payT >= 4 {
		// 对于银行卡，分
		db = db.Where("pay_type = ?", payT)

		orderByStatement := fmt.Sprintf("ABS( %d - pay_type)", payT) // 优先和payT能精确匹配上的银行
		db = db.Order(orderByStatement)
	} else if payT == models.PaymentTypeWeixin {
		db = db.Where("pay_type = ?", payT)
		// 微信支付方式支持"随机立减"的二维码：比如匹配不到空闲的200二维码，就匹配199.99，199.98等金额的二维码
		db = db.Where("(e_amount > 0 AND e_amount >= ? AND e_amount <= ?) OR e_amount = 0", amount-0.09-0.00001, amount) // 0.00001用来避免人民币金额浮点误差（目前仅BTUSD使用了没有浮点误差的decimal.Decimal类型）

		db = db.Order("e_amount DESC") // 优先固定二维码（e_amount对于固定二维码就是二维码金额，对于非固定二维码为0）
	} else if payT == models.PaymentTypeAlipay {
		db = db.Where("pay_type = ?", payT)

		db = db.Where("e_amount = ? OR e_amount = 0", amount)

		db = db.Order("e_amount DESC") // 优先固定二维码（e_amount对于固定二维码就是二维码金额，对于非固定二维码为0）
	} else {
		utils.Log.Warnf("payT %d is invalid", payT)
		return models.PaymentInfo{}
	}

	db = db.Where("uid = ? AND audit_status = 1 /**audit passed**/ AND in_use = 0 /**not in use**/", merchantID)
	db = db.Where("payment_auto_type = 0") // 仅查询手动收款账号

	if db.First(&payment).RecordNotFound() {
		return models.PaymentInfo{}
	}

	if payT == models.PaymentTypeWeixin || payT == models.PaymentTypeAlipay {
		utils.Log.Debugf("func GetBestNormalPaymentID, order %s, amount %f, e_amount %f in payment (id %d)", order.OrderNumber, order.Amount, payment.EAmount, payment.Id)
	}

	utils.Log.Debugf("func GetBestNormalPaymentID finished normally.")
	return payment
}
