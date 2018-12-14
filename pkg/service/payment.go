package service

import (
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

func AddPaymentInfo(uid int, payType int, name string, amount float64, account, bank, bankBranch string) response.AddPaymentRet {
	var ret response.AddPaymentRet

	var merchant models.PaymentInfo
	merchant.Uid = int64(uid)
	merchant.PayType = payType
	merchant.Name = name
	merchant.EAmount = amount
	if merchant.PayType == models.PaymentTypeWeixin || merchant.PayType == models.PaymentTypeAlipay {
		merchant.EAccount = account
		merchant.BankAccount = ""
	} else {
		merchant.EAccount = ""
		merchant.BankAccount = account
	}
	merchant.Bank = bank
	merchant.BankBranch = bankBranch
	merchant.QrCodeTxt = "TODO" // TODO
	merchant.QrCode = "TODO"

	if err := utils.DB.Create(&merchant).Error; err != nil {
		utils.Log.Errorf("AddPaymentInfo, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	return ret
}
