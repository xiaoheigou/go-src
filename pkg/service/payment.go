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
	"yuudidi.com/pkg/models"
	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/protocol/response/err_code"
	"yuudidi.com/pkg/utils"
)

func AddPaymentInfo(c *gin.Context) response.AddPaymentRet {
	return addOrUpdatePaymentInfo(c, false)
}

func UpdatePaymentInfo(c *gin.Context) response.AddPaymentRet {
	return addOrUpdatePaymentInfo(c, true)
}

func addOrUpdatePaymentInfo(c *gin.Context, isUpdate bool) response.AddPaymentRet {
	var ret response.AddPaymentRet

	var uid int
	var err error
	if uid, err = strconv.Atoi(c.Param("uid")); err != nil {
		utils.Log.Errorf("uid [%v] is invalid, expect a integer", c.Param("uid"))
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}

	var paymentId int
	if isUpdate {
		// id仅在更新信息时需要
		if paymentId, err = strconv.Atoi(c.Param("id")); err != nil {
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
	if !(payType == models.PaymentTypeWeixin || payType == models.PaymentTypeAlipay || payType == models.PaymentTypeBanck) {
		utils.Log.Errorf("pay_type [%v] is invalid", payType)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrArgInvalid.Data()
		return ret
	}
	name := c.Query("name")
	account := c.Query("account")
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
			imgLocalFilename := fmt.Sprintf("%s-%s-%s%s", strconv.Itoa(uid), strconv.Itoa(payType), randStr, filepath.Ext(file.Filename))
			// 下面把上传的图片（收款二维码）保存到本地文件中
			if err := ioutil.WriteFile(filepath.Join(imgPath, imgLocalFilename), imgBytes, 0664); err != nil {
				utils.Log.Errorf("save qrcode image to local fail, err = [%+v]", err)
			}
		}

		// 把用户上传的原始图片上传到阿里云OSS
		originalObjectKey := fmt.Sprintf("original/merchant/%s/%s/%s%s", strconv.Itoa(uid), strconv.Itoa(payType), randStr, filepath.Ext(file.Filename))
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
		generatedObjectKey := fmt.Sprintf("generated/merchant/%s/%s/%s.png", strconv.Itoa(uid), strconv.Itoa(payType), randStr)
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

	if isUpdate {
		return updatePaymentInfoToDB(uid, paymentId, payType, name, amountFloat, qrCodeTxt, qrCodeOrigin, qrCode, account, bank, bankBranch, accountDefaultInt)
	} else {
		return addPaymentInfoToDB(uid, payType, name, amountFloat, qrCodeTxt, qrCodeOrigin, qrCode, account, bank, bankBranch, accountDefaultInt)
	}
}

func addPaymentInfoToDB(uid int, payType int, name string, amount float64, qrCodeTxt, qrCodeOrigin, qrCode, account, bank, bankBranch string, accountDefault int) response.AddPaymentRet {
	var ret response.AddPaymentRet

	var paymentInfo models.PaymentInfo
	paymentInfo.Uid = int64(uid)
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

func updatePaymentInfoToDB(uid int, id int, payType int, name string, amount float64, qrCodeTxt, qrCodeOrigin, qrCode, account, bank, bankBranch string, accountDefault int) response.AddPaymentRet {
	var ret response.AddPaymentRet

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

	if err := utils.DB.Save(&paymentInfo).Error; err != nil {
		utils.Log.Errorf("updatePaymentInfoToDB fail, db err [%v]", err)
		ret.Status = response.StatusFail
		ret.ErrCode, ret.ErrMsg = err_code.AppErrDBAccessFail.Data()
		return ret
	}

	utils.Log.Infof("update payment[id=%v] for user[%v] successful", id, uid)
	ret.Status = response.StatusSucc
	ret.Data = append(ret.Data, paymentInfo)
	return ret
}

func GetPaymentInfo(uid int, c *gin.Context) response.GetPaymentsPageRet {
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
			if !(payType == models.PaymentTypeWeixin || payType == models.PaymentTypeAlipay || payType == models.PaymentTypeBanck) {
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
			ret.Status = response.StatusSucc
			for _, payment := range payments {
				ret.Data = append(ret.Data, payment)
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

	// TODO 删除阿里云OSS中数据
	ret.Status = response.StatusSucc
	return ret
}
