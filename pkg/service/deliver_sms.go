package service

import (
	"errors"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/service/dbcache"
	"yuudidi.com/pkg/utils"
)

func SendSmsOrderAccepted(merchantId int64, smsTplArg1 string) error {
	var merchant models.Merchant
	var err error
	if err = dbcache.GetMerchantById(merchantId, &merchant); err != nil {
		return err
	}

	// 短信模板id，这是提前在短信api管理台中设置的短信模板
	var tplId int64
	if tplId, err = strconv.ParseInt(utils.Config.GetString("sms.tencent.tplid.orderaccepted"), 10, 0); err != nil {
		utils.Log.Errorf("Wrong configuration: sms.tencent.tplid.orderaccepted, should be int.")
		return errors.New("sms.tencent.tplid.orderaccepted, should be int")
	}

	return utils.SendSmsByTencentApi(merchant.Phone, merchant.NationCode, tplId, smsTplArg1)
}

func SendSmsOrderPaid(merchantId int64, smsTplArg1 string) error {
	var merchant models.Merchant
	var err error
	if err = dbcache.GetMerchantById(merchantId, &merchant); err != nil {
		return err
	}

	// 短信模板id，这是提前在短信api管理台中设置的短信模板
	var tplId int64
	if tplId, err = strconv.ParseInt(utils.Config.GetString("sms.tencent.tplid.orderpaid"), 10, 0); err != nil {
		utils.Log.Errorf("Wrong configuration: sms.tencent.tplid.orderpaid, should be int.")
		return errors.New("sms.tencent.tplid.orderpaid, should be int")
	}

	return utils.SendSmsByTencentApi(merchant.Phone, merchant.NationCode, tplId, smsTplArg1)
}

func SendSmsOrderPaidTimeout(merchantId int64, smsTplArg1 string) error {
	var merchant models.Merchant
	var err error
	if err = dbcache.GetMerchantById(merchantId, &merchant); err != nil {
		return err
	}

	// 短信模板id，这是提前在短信api管理台中设置的短信模板
	var tplId int64
	if tplId, err = strconv.ParseInt(utils.Config.GetString("sms.tencent.tplid.orderpaidtimeout"), 10, 0); err != nil {
		utils.Log.Errorf("Wrong configuration: sms.tencent.tplid.orderpaidtimeout, should be int.")
		return errors.New("sms.tencent.tplid.orderpaidtimeout, should be int")
	}

	return utils.SendSmsByTencentApi(merchant.Phone, merchant.NationCode, tplId, smsTplArg1)
}
