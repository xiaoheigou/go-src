package err_code


var  (
	AppErrArgInvalid   = Err{10400,"request param is invalid"}
	AppErrEmailInvalid = Err{10401,"email address is invalid"}
	AppErrPhoneInvalid = Err{10402,"phone number is invalid"}
	AppErrNationCodeInvalid = Err{10403,"nation code is invalid"}
	AppErrSendSMSFail = Err{10404,"send sms fail"}
	AppErrSvrInternalFail = Err{10405,"server internal error"}

)