package err_code

var (
	AppErrArgInvalid                     = Err{10400, "请求参数错误"}
	AppErrEmailInvalid                   = Err{10401, "Email格式错误"}
	AppErrPhoneInvalid                   = Err{10402, "手机号码格式错误"}
	AppErrNationCodeInvalid              = Err{10403, "国家码格式错误"}
	AppErrSendSMSFail                    = Err{10404, "发送短信失败"}
	AppErrSendEmailFail                  = Err{10405, "发送邮件失败"}
	AppErrPageSizeTooLarge               = Err{10406, "参数page_size过大"}
	AppErrSvrInternalFail                = Err{10407, "服务器内部错误"}
	AppErrCloudStorageFail               = Err{10408, "云存储错误"}
	AppErrQrCodeDecodeFail               = Err{10409, "二维码解码失败"}
	AppErrDBAccessFail                   = Err{10410, "数据库错误"}
	AppErrRandomCodeVerifyFail           = Err{10411, "验证码错误"}
	AppErrRegisterRandomCodeReVerifyFail = Err{10411, "整个注册过程要在15分钟内完成，请重新注册"}
	AppErrCaptchaVerifyFail              = Err{10412, "验证码未通过，不能发送"}
	AppErrUserPasswordError              = Err{10413, "用户名或密码错误"}
	AppErrOldPasswordError               = Err{10414, "您的旧密码不正确"}
	AppErrPhoneAlreadyRegister           = Err{10415, "手机号码已经被注册"}
	AppErrNicknameTooLong                = Err{10416, "昵称太长"}
	AppErrGeetestVerifyFail              = Err{10417, "极验认证失败"}
	AppErrQrCodeInUseError               = Err{10440, "二维码删除失败，可能正在被使用"}
	AppErrLoginTryTooManyTimes           = Err{10441, "连续登录3次失败，请24小时后再试"}
	AppErrUpdateRealNameFail             = Err{10450, "您只能更新当前使用的自动收款信息"}
)
