package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
	"strings"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err-code"
	"yuudidi.com/pkg/utils"
)

func AddPaymentInfo(c *gin.Context) response.CommonRet {
	return addOrUpdatePaymentInfo(c, false)
}

func UpdatePaymentInfo(c *gin.Context) response.CommonRet {
	return addOrUpdatePaymentInfo(c, true)
}

func addOrUpdatePaymentInfo(c *gin.Context, isUpdate bool) response.CommonRet {
	var uid int
	var err error
	if uid, err = strconv.Atoi(c.Param("uid")); err != nil {
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		var ret response.CommonRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}

	var paymentId int
	if isUpdate {
		// id仅在更新信息时需要
		if paymentId, err = strconv.Atoi(c.Param("id")); err != nil {
			utils.Log.Errorf("id [%v] is invalid, expect a integer", c.Param("uid"))
			var ret response.CommonRet
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
	}

	var payType int = 1
	if payType, err = strconv.Atoi(c.Query("pay_type")); err != nil {
		utils.Log.Errorf("pay_type [%v] is invalid, expect a integer", c.Param("pay_type"))
		var ret response.CommonRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}
	if ! (payType == models.PaymentTypeWeixin || payType == models.PaymentTypeAlipay || payType == models.PaymentTypeBanck) {
		utils.Log.Errorf("pay_type [%v] is invalid", payType)
		var ret response.CommonRet
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}
	name := c.Query("name")
	amount := c.Query("amount")
	account := c.Query("account")
	var amountFloat float64
	bank := c.Query("bank")
	bankBranch := c.Query("bank_branch")
	accountDefault := c.Query("account_default")
	var accountDefaultInt int64

	var imgFilename string
	var qrCodeTxt = ""
	var qrCode = ""
	if payType == models.PaymentTypeWeixin || payType == models.PaymentTypeAlipay {
		// 检测方式为Weixin或者Alipay时的参数
		if amountFloat, err = strconv.ParseFloat(amount, 32); err != nil {
			utils.Log.Errorf("amount [%v] is invalid", amount)
			var ret response.CommonRet
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}

		file, err := c.FormFile("file")
		if err != nil {
			utils.Log.Errorf("get form err: [%v]", err)
			var ret response.CommonRet
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}

		var imgPath = utils.Config.GetString("qrcode.imgpath")
		if imgPath == "" {
			utils.Log.Errorf("missing configuration qrcode.imgpath")
			var ret response.CommonRet
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
			return ret
		}

		imgFilename = fmt.Sprintf("%s_%s_%s_%s", strconv.Itoa(uid), strconv.Itoa(payType), amount, file.Filename)
		remoteSvr := utils.Config.GetString("qrcode.remotesvr")
		if remoteSvr == "" {
			utils.Log.Errorf("missing configuration qrcode.remotesvr")
			var ret response.CommonRet
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
			return ret
		}
		// 下面把上传的图片（收款二维码）保存到本地文件中
		if err := c.SaveUploadedFile(file, imgPath+"/"+imgFilename); err != nil {
			utils.Log.Errorf("save upload file err: [%v]", err)
			var ret response.CommonRet
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
			return ret
		}
		qrCodeTxt = "TODO" // TODO
		qrCode = remoteSvr + "/" + imgFilename
	} else {
		// 检测方式为银行时的参数
		if strings.TrimSpace(name) == "" {
			utils.Log.Errorf("Add bank payment info, but missing name")
			var ret response.CommonRet
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
		if strings.TrimSpace(account) == "" {
			utils.Log.Errorf("Add bank payment info, but missing account")
			var ret response.CommonRet
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
		if strings.TrimSpace(bank) == "" {
			utils.Log.Errorf("Add bank payment info, but missing bank")
			var ret response.CommonRet
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
		if strings.TrimSpace(bankBranch) == "" {
			utils.Log.Errorf("Add bank payment info, but missing bank_branch")
			var ret response.CommonRet
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}

		if accountDefaultInt, err = strconv.ParseInt(accountDefault, 10, 0); err != nil {
			utils.Log.Errorf("account_default [%v] is invalid", accountDefault)
			var ret response.CommonRet
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
	}

	if isUpdate {
		return updatePaymentInfoToDB(uid, paymentId, payType, name, amountFloat, qrCodeTxt, qrCode, account, bank, bankBranch, int(accountDefaultInt))
	} else {
		return addPaymentInfoToDB(uid, payType, name, amountFloat, qrCodeTxt, qrCode, account, bank, bankBranch, int(accountDefaultInt))
	}
}

func addPaymentInfoToDB(uid int, payType int, name string, amount float64, qrCodeTxt, qrCode, account, bank, bankBranch string, accountDefault int) response.CommonRet {
	var ret response.CommonRet

	var paymentInfo models.PaymentInfo
	paymentInfo.Uid = int64(uid)
	paymentInfo.PayType = payType
	paymentInfo.Name = name
	paymentInfo.EAmount = amount
	if payType == models.PaymentTypeWeixin || payType == models.PaymentTypeAlipay {
		paymentInfo.EAccount = account
		paymentInfo.BankAccount = ""
	} else {
		paymentInfo.EAccount = ""
		paymentInfo.BankAccount = account
	}
	paymentInfo.Bank = bank
	paymentInfo.BankBranch = bankBranch
	paymentInfo.QrCodeTxt = qrCodeTxt
	paymentInfo.QrCode = qrCode
	paymentInfo.AuditStatus = models.PaymentAuditNopass // 新增加的收款信息，为未审核的状态
	paymentInfo.AccountDefault = accountDefault

	if err := utils.DB.Create(&paymentInfo).Error; err != nil {
		utils.Log.Errorf("AddPaymentInfo, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	return ret
}

func updatePaymentInfoToDB(uid int, id int, payType int, name string, amount float64, qrCodeTxt, qrCode, account, bank, bankBranch string, accountDefault int) response.CommonRet {
	var ret response.CommonRet

	var paymentInfo models.PaymentInfo
	if err := utils.DB.First(&paymentInfo, "uid = ? and id = ?", uid, id).Error; err != nil {
		utils.Log.Warnf("not found payment[id=%v] for user[%v], err", id, uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}

	paymentInfo.PayType = payType
	paymentInfo.Name = name
	paymentInfo.EAmount = amount
	if payType == models.PaymentTypeWeixin || payType == models.PaymentTypeAlipay {
		paymentInfo.EAccount = account
		paymentInfo.BankAccount = ""
	} else {
		paymentInfo.EAccount = ""
		paymentInfo.BankAccount = account
	}
	paymentInfo.Bank = bank
	paymentInfo.BankBranch = bankBranch
	paymentInfo.QrCodeTxt = qrCodeTxt
	paymentInfo.QrCode = qrCode
	paymentInfo.AuditStatus = models.PaymentAuditNopass // 更新信息后，要重置审核状态
	paymentInfo.AccountDefault = accountDefault

	if err := utils.DB.Save(&paymentInfo).Error; err != nil {
		utils.Log.Errorf("updatePaymentInfoToDB fail, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	utils.Log.Infof("update payment[id=%v] for user[%v] successful", id, uid)
	ret.Status = response.StatusSucc
	return ret
}

func GetPaymentInfo(uid int, c *gin.Context) response.GetPaymentsPageRet {
	var ret response.GetPaymentsPageRet

	var err error

	db := utils.DB.Model(&models.PaymentInfo{Uid: int64(uid)})

	pageNumStr := c.Query("page_num")
	pageSizeStr := c.Query("page_size")
	var pageNum int
	if pageNumStr == "" {
		utils.Log.Warnf("page_num is missing. Set to default 1.")
		pageNum = 1
	} else if pageNum, err = strconv.Atoi(pageNumStr); err != nil {
		utils.Log.Errorf("page_num [%v] is invalid, should be int.", pageNumStr)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}

	var pageSize int
	if pageSizeStr == "" {
		utils.Log.Warnf("page_size is missing. Set to default 10.")
		pageSize = 10
	} else if pageSize, err = strconv.Atoi(pageSizeStr); err != nil {
		utils.Log.Errorf("page_size [%v] is invalid, should be int.", pageSizeStr)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}

	db = db.Offset(pageSize * (pageNum - 1)).Limit(pageSize)

	var payType int
	payTypeStr := c.Query("pay_type")
	if payTypeStr == "" {
		utils.Log.Warnf("pay_type is missing. query all payments for merchant(uid=[%d])", uid)
	} else {
		if payType, err = strconv.Atoi(payTypeStr); err != nil {
			utils.Log.Errorf("pay_type [%v] is invalid, expect a integer", c.Param("pay_type"))
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
		if payType == -1 {
			// -1表示查询所有的
			utils.Log.Debugf("GetPaymentInfo, query all payments for merchant(uid=[%d])", uid)
		} else {
			if ! (payType == models.PaymentTypeWeixin || payType == models.PaymentTypeAlipay || payType == models.PaymentTypeBanck) {
				utils.Log.Warnln("pay_type [%v] is invalid", payType)
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
				return ret
			}
			// 增加一个查询条件
			db = db.Where("pay_type = ? ", payType)
		}
	}

	// 设置page参数，返回给前端
	db.Count(&ret.TotalCount)
	ret.PageCount = int(math.Ceil(float64(ret.TotalCount) / float64(pageSize)))
	ret.PageNum = pageNum  // 前端提供的查询参数或默认参数，返回给前端
	ret.PageSize = pageSize  // 前端提供的查询参数或默认参数，返回给前端

	var payments []models.PaymentInfo
	if err := db.Find(&payments).Error; err != nil {
		utils.Log.Errorf("GetPaymentInfo, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	} else {
		if len(payments) == 0 {
			utils.Log.Errorf("GetPaymentInfo, can't find assets for merchant(uid=[%d]).", uid)
			// 查不到没必要报错给前端，返回空即可
			ret.Status = response.StatusSucc
			ret.Data = make([]models.PaymentInfo, 0, 1)
			return ret
		} else {
			ret.Status = response.StatusSucc
			for _, payment := range payments {
				ret.Data = append(ret.Data, models.PaymentInfo{
					Id:             payment.Id,
					Uid:            payment.Uid,
					PayType:        payment.PayType,
					QrCodeTxt:      payment.QrCodeTxt,
					QrCode:         payment.QrCode,
					EAmount:        payment.EAmount,
					EAccount:       payment.EAccount,
					Name:           payment.Name,
					BankAccount:    payment.BankAccount,
					Bank:           payment.Bank,
					BankBranch:     payment.BankBranch,
					AccountDefault: payment.AccountDefault,
					AuditStatus:    payment.AuditStatus,
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
