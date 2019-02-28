package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/service/dbcache"
	"yuudidi.com/pkg/utils"
)

func AddPaymentInfo(c *gin.Context) response.AddPaymentRet {
	return addOrUpdatePaymentInfo(c, false)
}

func UpdatePaymentInfo(c *gin.Context) response.AddPaymentRet {
	return addOrUpdatePaymentInfo(c, true)
}

// 下面函数主要做参数校验等工作，主要业务逻辑在updatePaymentInfoToDB或者addPaymentInfoToDB中
func addOrUpdatePaymentInfo(c *gin.Context, isUpdate bool) response.AddPaymentRet {
	var ret response.AddPaymentRet

	var uid int64
	var err error
	if uid, err = strconv.ParseInt(c.Param("uid"), 10, 64); err != nil {
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}

	var paymentId int64
	if isUpdate {
		// id仅在更新信息时需要
		if paymentId, err = strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
			utils.Log.Errorf("id [%v] is invalid, expect a integer", c.Param("uid"))
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
	}

	var payType int = 1
	if payType, err = strconv.Atoi(c.Query("pay_type")); err != nil {
		utils.Log.Errorf("pay_type [%v] is invalid, expect a integer", c.Param("pay_type"))
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}

	var paymentAutoType int = 0 // 默认为手动账号
	if strings.TrimSpace(c.Query("payment_auto_type")) != "" {
		if paymentAutoType, err = strconv.Atoi(c.Query("payment_auto_type")); err != nil {
			utils.Log.Errorf("payment_auto_type [%v] is invalid, expect a integer", c.Param("payment_auto_type"))
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
		if paymentAutoType == 0 || paymentAutoType == 1 {
			// pass
		} else {
			utils.Log.Errorf("payment_auto_type [%v] is invalid, expect 0 or 1", c.Param("payment_auto_type"))
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
	}

	var enable int = 1 // 是否启用账号，默认为启用
	if strings.TrimSpace(c.Query("enable")) != "" {
		if enable, err = strconv.Atoi(c.Query("enable")); err != nil {
			utils.Log.Errorf("enable [%v] is invalid, expect a integer", c.Param("enable"))
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
		if enable == 0 || enable == 1 {
			// pass
		} else {
			utils.Log.Errorf("enable [%v] is invalid, expect 0 or 1", c.Param("enable"))
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
	}

	name := c.Query("name")
	account := c.Query("account")
	userPayId := c.Query("user_pay_id")

	var amountFloat float64
	bank := c.Query("bank")
	bankBranch := c.Query("bank_branch")
	accountDefault := c.Query("account_default")
	var accountDefaultInt int

	var qrCodeTxt = ""
	var qrCode = ""
	var qrCodeOrigin = ""
	if payType == models.PaymentTypeWeixin || payType == models.PaymentTypeAlipay {
		// 检测方式为Weixin或者Alipay时的参数

		if paymentAutoType == 0 { // 手动账号
			file, err := c.FormFile("file")
			if err != nil {
				utils.Log.Errorf("get form err: [%v]", err)
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
				return ret
			}

			src, err := file.Open()
			if err != nil {
				utils.Log.Errorf("open form file err: [%v]", err)
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
				return ret
			}
			defer src.Close()

			var imgBytes []byte
			if imgBytes, err = ioutil.ReadAll(src); err != nil {
				utils.Log.Errorf("read form file err: [%v]", err)
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
				return ret
			}

			randStr := utils.GetRandomString(6)
			// 保存到本地文件系统
			saveImgLocally := utils.Config.GetBool("qrcode.savelocally")
			if saveImgLocally {
				var imgPath = utils.Config.GetString("qrcode.imgpath")
				if imgPath == "" {
					utils.Log.Errorf("missing configuration qrcode.imgpath")
					ret.Status = response.StatusFail
					ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
					return ret
				}
				imgLocalFilename := fmt.Sprintf("%s-%s-%s%s", strconv.Itoa(int(uid)), strconv.Itoa(payType), randStr, filepath.Ext(file.Filename))
				// 下面把上传的图片（收款二维码）保存到本地文件中
				if err := ioutil.WriteFile(filepath.Join(imgPath, imgLocalFilename), imgBytes, 0664); err != nil {
					utils.Log.Errorf("save qrcode image to local fail, err = [%+v]", err)
				}
			}

			// 把用户上传的原始图片上传到阿里云OSS
			originalObjectKey := fmt.Sprintf("original/merchant/%s/%s/%s%s", strconv.Itoa(int(uid)), strconv.Itoa(payType), randStr, filepath.Ext(file.Filename))
			var originalImgUrl = ""
			if originalImgUrl, err = utils.UploadQrcode2AliyunOss(originalObjectKey, ioutil.NopCloser(bytes.NewBuffer(imgBytes))); err != nil {
				utils.Log.Errorf("upload original qrcode img to cloud storage err: [%v]", err)
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrCloudStorageFail.Data()
				return ret
			}

			// 分析二维码，调用二维码及图片金额识别服务
			var qrCodeInfo utils.QrcodeRespMsg
			if qrCodeInfo, err = utils.GetQrCodeInfo(ioutil.NopCloser(bytes.NewBuffer(imgBytes)), randStr+filepath.Ext(file.Filename), c.Query("qr_code_txt")); err != nil {
				utils.Log.Errorf("func GetQrCodeInfo fail. err: [%v]", err)
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
				return ret
			}
			utils.Log.Debugf("qrcode info, amount = [%v], qrCodeTxt = [%v]", qrCodeInfo.Amount, qrCodeInfo.QrCodeTxt)

			if qrCodeInfo.QrCodeTxt == "" {
				utils.Log.Errorf("can not decode qrcode")
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrQrCodeDecodeFail.Data()
				return ret
			}

			// 新生成的二维码图片是base64编码，先解码
			var generatedImgUrl = ""
			decoded, err := base64.StdEncoding.DecodeString(qrCodeInfo.NewQrCodeBase64)
			if err != nil {
				utils.Log.Errorf("generated image is not encoded with base64. error:", err)
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
				return ret
			}

			// 把新生成的二维码（都是png格式）上传到阿里云OSS
			generatedObjectKey := fmt.Sprintf("generated/merchant/%s/%s/%s.png", strconv.Itoa(int(uid)), strconv.Itoa(payType), randStr)
			if generatedImgUrl, err = utils.UploadQrcode2AliyunOss(generatedObjectKey, ioutil.NopCloser(bytes.NewBuffer(decoded))); err != nil {
				utils.Log.Errorf("upload generated qrcode img to cloud storage err: [%v]", err)
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrCloudStorageFail.Data()
				return ret
			}

			if qrCodeInfo.Amount == "" {
				// 图片中没有找到金额
				utils.Log.Debugf("Don't find amount in qrcode image")
			} else {
				if amountFloat, err = strconv.ParseFloat(qrCodeInfo.Amount, 64); err != nil {
					utils.Log.Errorf("qrCodeInfo.Amount [%v] is invalid", qrCodeInfo.Amount)
					ret.Status = response.StatusFail
					ret.ErrCode, ret.ErrMsg = err_code.AppErrSvrInternalFail.Data()
					return ret
				}
			}

			qrCodeTxt = qrCodeInfo.QrCodeTxt
			qrCodeOrigin = originalImgUrl
			qrCode = generatedImgUrl
		} else { // 自动账号
			// 如果是修改自动账号，则需要提供name（只有实名信息可以修改，user_pay_id是不能修改的）
			// 如果是增加自动账号，则必需提供user_pay_id和name
			if isUpdate {
				if strings.TrimSpace(name) == "" {
					utils.Log.Errorf("Update auto payment info, but missing argument name")
					ret.Status = response.StatusFail
					ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
					return ret
				}
			} else {
				if strings.TrimSpace(name) == "" {
					utils.Log.Errorf("Add auto payment info, but missing argument name")
					ret.Status = response.StatusFail
					ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
					return ret
				}
				if strings.TrimSpace(userPayId) == "" {
					utils.Log.Errorf("Add auto payment info, but missing argument user_pay_id")
					ret.Status = response.StatusFail
					ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
					return ret
				}
			}
		}
	} else {
		// 检测方式为银行时的参数
		if strings.TrimSpace(name) == "" {
			utils.Log.Errorf("Add bank payment info, but missing name")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
		if strings.TrimSpace(account) == "" {
			utils.Log.Errorf("Add bank payment info, but missing account")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
		if strings.TrimSpace(bank) == "" {
			utils.Log.Errorf("Add bank payment info, but missing bank")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
		if strings.TrimSpace(bankBranch) == "" {
			utils.Log.Errorf("Add bank payment info, but missing bank_branch")
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}

		if accountDefaultInt, err = strconv.Atoi(accountDefault); err != nil {
			utils.Log.Errorf("account_default [%v] is invalid", accountDefault)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
			return ret
		}
	}

	if paymentAutoType == 0 { // 手动收款账号
		if isUpdate {
			return updatePaymentInfoToDB(uid, paymentId, payType, name, amountFloat, qrCodeTxt, qrCodeOrigin, qrCode, account, bank, bankBranch, accountDefaultInt)
		} else {
			return addPaymentInfoToDB(uid, payType, name, amountFloat, qrCodeTxt, qrCodeOrigin, qrCode, account, bank, bankBranch, accountDefaultInt)
		}
	} else { // 自动收款账号
		if isUpdate {
			return updateAutoPaymentInfoToDB(uid, paymentId, name, account, enable)
		} else {
			return addAutoPaymentInfoToDB(uid, payType, name, account, userPayId)
		}
	}
}

func addPaymentInfoToDB(uid int64, payType int, name string, amount float64, qrCodeTxt, qrCodeOrigin, qrCode, account, bank, bankBranch string, accountDefault int) response.AddPaymentRet {
	var ret response.AddPaymentRet

	var paymentInfo models.PaymentInfo
	paymentInfo.Uid = uid
	paymentInfo.PayType = payType
	paymentInfo.Name = name
	paymentInfo.EAmount = amount
	if payType == models.PaymentTypeWeixin {
		paymentInfo.EAccount = account
		paymentInfo.BankAccount = ""
		if utils.IsWeixinQrCode(qrCodeTxt) {
			paymentInfo.AuditStatus = models.PaymentAuditPass
		} else {
			paymentInfo.AuditStatus = models.PaymentAuditNopass
		}
	} else if payType == models.PaymentTypeAlipay {
		paymentInfo.EAccount = account
		paymentInfo.BankAccount = ""

		if utils.IsAlipayQrCode(qrCodeTxt) {
			paymentInfo.AuditStatus = models.PaymentAuditPass
		} else {
			paymentInfo.AuditStatus = models.PaymentAuditNopass
		}
	} else {
		paymentInfo.EAccount = ""
		paymentInfo.BankAccount = account
		paymentInfo.AuditStatus = models.PaymentAuditPass
	}
	paymentInfo.Bank = bank
	paymentInfo.BankBranch = bankBranch
	paymentInfo.QrCodeTxt = qrCodeTxt
	paymentInfo.QrCodeOrigin = qrCodeOrigin
	paymentInfo.QrCode = qrCode
	paymentInfo.AccountDefault = accountDefault

	if err := utils.DB.Create(&paymentInfo).Error; err != nil {
		utils.Log.Errorf("AddPaymentInfo, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, paymentInfo)
	return ret
}

func updatePaymentInfoToDB(uid int64, paymentId int64, payType int, name string, amount float64, qrCodeTxt, qrCodeOrigin, qrCode, account, bank, bankBranch string, accountDefault int) response.AddPaymentRet {
	var ret response.AddPaymentRet

	var paymentInfo models.PaymentInfo
	if err := utils.DB.First(&paymentInfo, "uid = ? and id = ?", uid, paymentId).Error; err != nil {
		utils.Log.Warnf("not found payment[id=%v] for user[%v], err", paymentId, uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}

	paymentInfo.PayType = payType
	paymentInfo.Name = name
	paymentInfo.EAmount = amount
	if payType == models.PaymentTypeWeixin {
		paymentInfo.EAccount = account
		paymentInfo.BankAccount = ""
		if utils.IsWeixinQrCode(qrCodeTxt) {
			paymentInfo.AuditStatus = models.PaymentAuditPass
		} else {
			paymentInfo.AuditStatus = models.PaymentAuditNopass
		}
	} else if payType == models.PaymentTypeAlipay {
		paymentInfo.EAccount = account
		paymentInfo.BankAccount = ""

		if utils.IsAlipayQrCode(qrCodeTxt) {
			paymentInfo.AuditStatus = models.PaymentAuditPass
		} else {
			paymentInfo.AuditStatus = models.PaymentAuditNopass
		}
	} else {
		paymentInfo.EAccount = ""
		paymentInfo.BankAccount = account
		paymentInfo.AuditStatus = models.PaymentAuditPass
	}
	paymentInfo.Bank = bank
	paymentInfo.BankBranch = bankBranch
	paymentInfo.QrCodeTxt = qrCodeTxt
	paymentInfo.QrCodeOrigin = qrCodeOrigin
	paymentInfo.QrCode = qrCode
	paymentInfo.AccountDefault = accountDefault

	// TODO Save方法会改变所有字段，下面最好改为update，
	if err := utils.DB.Save(&paymentInfo).Error; err != nil {
		utils.Log.Errorf("updatePaymentInfoToDB fail, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	utils.Log.Infof("update payment[id=%v] for user[%v] successful", paymentId, uid)
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, paymentInfo)
	return ret
}

// 增加自动收款信息（仅适用于支付宝或微信）
func addAutoPaymentInfoToDB(uid int64, payType int, name string, account string, userPayId string) response.AddPaymentRet {
	var ret response.AddPaymentRet

	var paymentInfo models.PaymentInfo
	paymentInfo.Uid = uid
	paymentInfo.PayType = payType
	paymentInfo.PaymentAutoType = 1
	paymentInfo.Name = name
	paymentInfo.UserPayId = userPayId
	paymentInfo.AuditStatus = models.PaymentAuditPass
	paymentInfo.LastUseTime = time.Now()

	paymentInfo.EAccount = account

	tx := utils.DB.Begin()
	if tx.Error != nil {
		utils.Log.Errorf("tx in func addAutoPaymentInfoToDB begin fail, tx=[%v]", tx)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	if err := tx.Create(&paymentInfo).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("tx in func addAutoPaymentInfoToDB rollback, tx=[%v]", tx)
		utils.Log.Errorf("addAutoPaymentInfoToDB fail, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	var merchant models.Merchant
	if err := dbcache.GetMerchantById(int64(uid), &merchant); err != nil {
		tx.Rollback()
		utils.Log.Errorf("tx in func addAutoPaymentInfoToDB rollback, tx=[%v]", tx)
		utils.Log.Errorf("addAutoPaymentInfoToDB, call GetMerchantById fail. [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	// 更新 curr_auto_weixin_payment_id 或者 curr_auto_alipay_payment_id 字段为新增加的payment info的Id
	if paymentInfo.PayType == models.PaymentTypeWeixin {
		// 修改preferences表中curr_auto_weixin_payment_id字段
		if err := tx.Table("preferences").Where("id = ?", merchant.PreferencesId).Update(
			"curr_auto_weixin_payment_id", paymentInfo.Id,
		).Error; err != nil {
			tx.Rollback()
			utils.Log.Errorf("tx in func addAutoPaymentInfoToDB rollback, tx=[%v]", tx)
			utils.Log.Errorf("addAutoPaymentInfoToDB, update curr_auto_weixin_payment_id for merchant(uid=[%d]) fail. [%v]", uid, err)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
			return ret
		}
	} else if paymentInfo.PayType == models.PaymentTypeAlipay {
		// 修改preferences表中curr_auto_alipay_payment_id字段
		if err := tx.Table("preferences").Where("id = ?", merchant.PreferencesId).Update(
			"curr_auto_alipay_payment_id", paymentInfo.Id,
		).Error; err != nil {
			utils.Log.Errorf("addAutoPaymentInfoToDB, update curr_auto_alipay_payment_id for merchant(uid=[%d]) fail. [%v]", uid, err)
			tx.Rollback()
			utils.Log.Errorf("tx in func addAutoPaymentInfoToDB rollback, tx=[%v]", tx)
			ret.Status = response.StatusFail
			ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
			return ret
		}
	} else {
		utils.Log.Errorf("addAutoPaymentInfoToDB, invalid payType %d", payType)
		utils.Log.Errorf("tx in func addAutoPaymentInfoToDB rollback, tx=[%v]", tx)
		tx.Rollback()
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("error tx in func addAutoPaymentInfoToDB commit, err=[%v]", err)
	}

	// 修改了Preference表，应使其缓存失效
	if err := dbcache.InvalidatePreference(int64(merchant.PreferencesId)); err != nil {
		utils.Log.Warnf("InvalidatePreference fail, err [%v]", err)
	}

	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, paymentInfo)
	return ret
}

// 修改当前使用的自动收款信息（仅适用于支付宝或微信）
func updateAutoPaymentInfoToDB(uid int64, paymentId int64, name string, account string, enable int) response.AddPaymentRet {
	var ret response.AddPaymentRet

	var paymentInfo models.PaymentInfo
	if err := utils.DB.First(&paymentInfo, "uid = ? and id = ?", uid, paymentId).Error; err != nil {
		utils.Log.Warnf("not found payment[id=%v] for user[%v], err", paymentId, uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}

	var merchant models.Merchant
	if err := dbcache.GetMerchantById(uid, &merchant); err != nil {
		utils.Log.Errorf("updateAutoPaymentInfoToDB, call GetMerchantById fail. [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	// 检测这个自动收款信息是不是当前正在使用的，如果不是，则报错
	var pref models.Preferences
	if err := dbcache.GetPreferenceById(int64(merchant.PreferencesId), &pref); err != nil {
		utils.Log.Errorf("updateAutoPaymentInfoToDB, can't find preference record in db for merchant(uid=[%d]),  err [%v]", uid, err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}
	//if paymentInfo.PayType == models.PaymentTypeWeixin {
	//	if pref.CurrAutoWeixinPaymentId != paymentId {
	//		utils.Log.Errorf("only payment info with CurrAutoWeixinPaymentId can be updated")
	//		ret.Status = response.StatusFail
	//		ret.ErrCode, ret.ErrMsg = err_code.AppErrUpdateRealNameFail.Data()
	//		return ret
	//	}
	//} else if paymentInfo.PayType == models.PaymentTypeAlipay {
	//	if pref.CurrAutoAlipayPaymentId != paymentId {
	//		utils.Log.Errorf("only payment info with CurrAutoAlipayPaymentId can be updated")
	//		ret.Status = response.StatusFail
	//		ret.ErrCode, ret.ErrMsg = err_code.AppErrUpdateRealNameFail.Data()
	//		return ret
	//	}
	//} else {
	//	utils.Log.Errorf("addAutoPaymentInfoToDB, invalid payType %d", paymentInfo.PayType)
	//	ret.Status = response.StatusFail
	//	ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
	//	return ret
	//}

	paymentInfo.Name = name
	paymentInfo.AuditStatus = models.PaymentAuditPass
	paymentInfo.EAccount = account // 不会是银行卡
	paymentInfo.Enable = enable

	utils.Log.Debugf("func updateAutoPaymentInfoToDB, paymentInfo = %+v", paymentInfo)
	if err := utils.DB.Save(&paymentInfo).Error; err != nil {
		utils.Log.Errorf("updateAutoPaymentInfoToDB fail, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	utils.Log.Infof("update payment[id=%v] for user[%v] successful", paymentId, uid)
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, paymentInfo)
	return ret
}

func GetPaymentInfo(uid int64, c *gin.Context) response.GetPaymentsPageRet {
	var ret response.GetPaymentsPageRet

	var err error

	db := utils.DB.Model(&models.PaymentInfo{})
	db = db.Where("uid = ?", uid)

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
	if pageSize > 50 {
		utils.Log.Errorf("page_size [%v] is too large, must <= 50", pageSize)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrPageSizeTooLarge.Data()
		return ret
	}

	queryPaymentAutoType := c.Query("payment_auto_type")
	if strings.EqualFold(queryPaymentAutoType, "0") {
		// 仅查询手动收款账号
		db = db.Where("payment_auto_type = 0")
	} else if strings.EqualFold(queryPaymentAutoType, "1") {
		// 仅查询自动收款账号
		db = db.Where("payment_auto_type = 1")
	}

	// 前端的query参数type可以是wechat/alipay/bank/all
	queryType := c.Query("type")
	if strings.EqualFold(queryType, "wechat") {
		db = db.Where("pay_type = ? ", models.PaymentTypeWeixin)
	} else if strings.EqualFold(queryType, "alipay") {
		db = db.Where("pay_type = ? ", models.PaymentTypeAlipay)
	} else if strings.EqualFold(queryType, "bank") {
		// 银行卡，pay_type >= 4
		db = db.Where("pay_type >= 4")
	} else if strings.EqualFold(queryType, "all") {
		// 不加过滤条件
	} else if queryType == "" {
		// 当query参数type为空时，进入兼容模式。
		// query参数pay_type是以前的设计，旧的Android app可能会使用。为兼容，代码暂时保留。适当时候，可以删除下面兼容代码。
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
				// 增加一个查询条件
				db = db.Where("pay_type = ? ", payType)
			}
		}
	} else {
		utils.Log.Warnln("query param type [%v] is invalid, only wechat/alipay/bank/all is accepted", queryType)
	}

	// 设置page参数，返回给前端
	db.Count(&ret.TotalCount)

	db = db.Order("payment_infos.enable desc, payment_infos.updated_at desc") // 把enable的排前面， enable中把最后更新的排前面
	db = db.Offset(pageSize * (pageNum - 1)).Limit(pageSize)

	ret.PageCount = int(math.Ceil(float64(ret.TotalCount) / float64(pageSize)))
	ret.PageNum = pageNum   // 前端提供的查询参数或默认参数，返回给前端
	ret.PageSize = pageSize // 前端提供的查询参数或默认参数，返回给前端

	var payments []models.PaymentInfo
	if err := db.Find(&payments).Error; err != nil {
		utils.Log.Errorf("GetPaymentInfo, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	} else {
		if len(payments) == 0 {
			utils.Log.Infof("GetPaymentInfo, can't find payment info for merchant(uid=[%d]).", uid)
			// 查不到没必要报错给前端，返回空即可
			ret.Status = response.StatusSucc
			ret.Data = make([]models.PaymentInfo, 0, 1)
			return ret
		} else {

			var merchant models.Merchant
			if err := dbcache.GetMerchantById(uid, &merchant); err != nil {
				utils.Log.Errorf("call GetMerchantById fail. [%v]", uid, err)
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
				return ret
			}

			var pref models.Preferences
			if err := dbcache.GetPreferenceById(int64(merchant.PreferencesId), &pref); err != nil {
				utils.Log.Errorf("can't find preference record in db for merchant(uid=[%d]),  err [%v]", uid, err)
				ret.Status = response.StatusFail
				ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
				return ret
			}

			//currAutoWechatPaymentId := pref.CurrAutoWeixinPaymentId
			//currAutoAlipayPaymentId := pref.CurrAutoAlipayPaymentId

			ret.Status = response.StatusSucc
			for _, payment := range payments {
				//if payment.PayType == models.PaymentTypeWeixin {
				//	if payment.Id == currAutoWechatPaymentId {
				//		payment.CurrAutoPayment = 1 // 设置该记录为当前使用的自动收款方式
				//	}
				//} else if payment.PayType == models.PaymentTypeAlipay {
				//	if payment.Id == currAutoAlipayPaymentId {
				//		payment.CurrAutoPayment = 1 // 设置该记录为当前使用的自动收款方式
				//	}
				//}
				ret.Data = append(ret.Data, payment)
			}
		}
		return ret
	}
}

func DeletePaymentInfo(uid int, paymentId int) response.DeletePaymentRet {
	var ret response.DeletePaymentRet

	var payment models.PaymentInfo
	rowAffected := utils.DB.Table("payment_infos").Where("uid = ? and id = ? and in_use = 0", uid, paymentId).Delete(&payment).RowsAffected
	if rowAffected == 0 {
		utils.Log.Errorf("DeletePaymentInfo, db err, record not found, may be qrcode is in_use")
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrQrCodeInUseError.Data()
		return ret
	}

	// TODO 删除阿里云OSS中数据
	ret.Status = response.StatusSucc
	return ret
}

// GetBankList - get bank list.
func GetBankList() (result response.GetBankListRet) {
	banks := utils.Config.GetStringMapString("banks")
	result.Status = response.StatusSucc
	bankArr := []models.BankInfo{}
	for bank := range banks {
		bankArr = append(bankArr, models.BankInfo{
			Name: bank,
			ID:   banks[bank],
		})
	}
	result.Data = bankArr
	return result
}
