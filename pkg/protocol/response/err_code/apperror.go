package err_code

var (
	AppErrArgInvalid           = Err{10400, "request param is invalid or missing"}
	AppErrEmailInvalid         = Err{10401, "email address is invalid"}
	AppErrPhoneInvalid         = Err{10402, "phone number is invalid"}
	AppErrNationCodeInvalid    = Err{10403, "nation code is invalid"}
	AppErrSendSMSFail          = Err{10404, "send sms fail"}
	AppErrSendEmailFail        = Err{10405, "send email fail"}
	AppErrPageSizeTooLarge     = Err{10406, "page size is too large"}
	AppErrSvrInternalFail      = Err{10407, "server internal error"}
	AppErrCloudStorageFail     = Err{10408, "cloud storage error"}
	AppErrQrCodeDecodeFail     = Err{10409, "decode qrcode fail"}
	AppErrDBAccessFail         = Err{10410, "database error"}
	AppErrRandomCodeVerifyFail = Err{10411, "random code verify fail"}
	AppErrCaptchaVerifyFail    = Err{10412, "captcha verify fail, can not send sms/email"}
	AppErrUserPasswordError    = Err{10413, "user name or password is invalid"}
	AppErrOldPasswordError     = Err{10414, "your old password is invalid"}
	AppErrPhoneAlreadyRegister = Err{10415, "phone already registered"}
	AppErrNicknameTooLong      = Err{10416, "nickname is too long"}
	AppErrGeetestVerifyFail    = Err{10417, "Geetest verify fail"}
	AppErrQrCodeInUseError     = Err{10440, "may be your qrcode is in used"}
)
