package err_code


var  (
	AppErrArgInvalid   = Err{10400,"request param is invalid"}
	AppErrEmailInvalid = Err{10401,"email address is invalid"}
	AppErrPhoneInvalid = Err{10402,"phone number is invalid"}
	AppErrNationCodeInvalid = Err{10403,"nation code is invalid"}
	AppErrSendSMSFail = Err{10404,"send sms fail"}
	AppErrSendEmailFail = Err{10404,"send email fail"}
	AppErrSvrInternalFail = Err{10405,"server internal error"}
	AppErrDBAccessFail = Err{10405,"database error"}
	AppErrRandomCodeVerifyFail = Err{10406,"random code verify fail"}
	AppErrUserPasswordError = Err{10402, "user name or password is invalid"}
	AppErrOldPasswordError = Err{10402, "your old password is invalid"}
	AppErrPhoneAlreadyRegister = Err{10402, "phone already registered"}
	AppErrEmailAlreadyRegister = Err{10402, "email already registered"}
	AppErrNicknameTooLong = Err{10402, "nickname is too long"}
)