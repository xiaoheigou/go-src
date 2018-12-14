package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

func AddPaymentInfo(c *gin.Context) response.AddPaymentRet {
	var uid int
	var err error
	if uid, err = strconv.Atoi(c.Param("uid")); err != nil {
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		var ret response.AddPaymentRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return ret
	}

	var payType int  = 1
	if payType, err = strconv.Atoi(c.Query("pay_type")); err != nil {
		utils.Log.Errorf("pay_type [%v] is invalid, expect a integer", c.Param("pay_type"))
		var ret response.AddPaymentRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return ret
	}
	if ! (payType == models.PaymentTypeWeixin || payType == models.PaymentTypeAlipay || payType == models.PaymentTypeBanck) {
		utils.Log.Errorf("pay_type [%v] is invalid", payType)
		var ret response.AddPaymentRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return ret
	}
	name := c.Query("name")
	amount := c.Query("amount")
	account := c.Query("account")
	bank := c.Query("bank")
	bankBranch := c.Query("bank_branch")
	var amountFloat float64
	if amountFloat, err = strconv.ParseFloat(amount, 32); err != nil {
		utils.Log.Errorf("amount [%v] is invalid", amount)
		var ret response.AddPaymentRet
		ret.Status = response.StatusFail
		ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		c.JSON(200, ret)
		return ret
	}

	var imgFilename string
	var qrCodeTxt = ""
	var qrCode = ""
	if payType == models.PaymentTypeWeixin || payType == models.PaymentTypeAlipay {
		file, err := c.FormFile("file")
		if err != nil {
			utils.Log.Errorf("get form err: [%v]", err)
			var ret response.AddPaymentRet
			ret.Status = response.StatusFail
			ret.ErrCode,ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			c.JSON(200, ret)
			return ret
		}

		var imgPath = utils.Config.GetString("qrcode.imgpath")
		if imgPath == "" {
			utils.Log.Errorf("missing configuration qrcode.imgpath")
			var ret response.AddPaymentRet
			ret.Status = response.StatusFail
			ret.ErrCode,ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
			c.JSON(200, ret)
			return ret
		}

		imgFilename = fmt.Sprintf("%s_%s_%s_%s", strconv.Itoa(uid), strconv.Itoa(payType), amount, file.Filename)
		remoteSvr := utils.Config.GetString("qrcode.remotesvr")
		if remoteSvr == "" {
			utils.Log.Errorf("missing configuration qrcode.remotesvr")
			var ret response.AddPaymentRet
			ret.Status = response.StatusFail
			ret.ErrCode,ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
			c.JSON(200, ret)
			return ret
		}
		// 下面把上传的图片（收款二维码）保存到本地文件中
		if err := c.SaveUploadedFile(file, imgPath + "/" + imgFilename); err != nil {
			utils.Log.Errorf("save upload file err: [%v]", err)
			var ret response.AddPaymentRet
			ret.Status = response.StatusFail
			ret.ErrCode,ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
			c.JSON(200, ret)
			return ret
		}
		qrCodeTxt = "TODO" // TODO
		qrCode = remoteSvr + "/" + imgFilename
	}

	return AddPaymentInfoToDB(uid, payType, name, amountFloat, qrCodeTxt, qrCode, account, bank, bankBranch)
}

func AddPaymentInfoToDB(uid int, payType int, name string, amount float64, qrCodeTxt, qrCode, account, bank, bankBranch string) response.AddPaymentRet {
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
	merchant.QrCodeTxt = qrCodeTxt
	merchant.QrCode = qrCode

	if err := utils.DB.Create(&merchant).Error; err != nil {
		utils.Log.Errorf("AddPaymentInfo, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	return ret
}


func GetPaymentInfo(uid int) response.GetPaymentsRet {
	var ret response.GetPaymentsRet

	var payments []models.PaymentInfo
	if err := utils.DB.Where(&models.PaymentInfo{Uid: int64(uid)}).Find(&payments).Error; err != nil {
		utils.Log.Errorf("GetPaymentInfo, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	} else {
		if len(payments) == 0 {
			utils.Log.Errorf("GetPaymentInfo, can't find assets for merchant(uid=[%d]).", uid)
			// 查不到没必要报错给前端，返回空即可
			ret.Status = response.StatusSucc
			ret.Data = append(ret.Data, models.PaymentInfo{
				Id:          0,
				Uid:         0,
				PayType:     0,
				QrCodeTxt:   "",
				QrCode:      "",
				EAmount:     0,
				EAccount:    "",
				Name:        "",
				BankAccount: "",
				Bank:        "",
				BankBranch:  "",
			})
			return ret
		} else {
			ret.Status = response.StatusSucc
			for _, payment := range payments {
				ret.Data = append(ret.Data, models.PaymentInfo{
					Id:          payment.Id,
					Uid:         payment.Uid,
					PayType:     payment.PayType,
					QrCodeTxt:   payment.QrCodeTxt,
					QrCode:      payment.QrCode,
					EAmount:     payment.EAmount,
					EAccount:    payment.EAccount,
					Name:        payment.Name,
					BankAccount: payment.BankAccount,
					Bank:        payment.Bank,
					BankBranch:  payment.BankBranch,
				})
			}
		}
		return ret
	}
}


func DeletePaymentInfo(uid int, paymentId int) response.DeletePaymentRet {
	var ret response.DeletePaymentRet

	var payment models.PaymentInfo
	if err := utils.DB.Table("payment_infos").Where("uid = ? and id = ?", uid, paymentId).Delete(&payment).Error; err != nil {
		utils.Log.Errorf("DeletePaymentInfo, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	return ret
}
