package err_code

type Err struct {
	ErrCode int
	ErrMsg  string
}

func (e Err) Data() (int, string) {
	return e.ErrCode, e.ErrMsg
}

var (
	DistributorErr                = Err{20100, "create distributor is failed."}
	RequestParamErr               = Err{20001, "request param is error."}
	CreateUserErr                 = Err{20400, "create user is failed."}
	NotFoundMerchant              = Err{20201, "not found merchant"}
	NotFoundUser                  = Err{20401, "not found user"}
	UserPasswordError             = Err{20402, "user password is invalid"}
	UpdateMerchantStatusErr       = Err{20205, "update merchant status is failed"}
	NoAccountIdOrDistributorIdErr = Err{20403, "accountId or distributorId is null"}
)
