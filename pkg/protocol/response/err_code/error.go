package err_code

type Err struct {
	ErrCode int
	ErrMsg  string
}

func (e Err) Data() (int, string) {
	return e.ErrCode, e.ErrMsg
}

var (
	LoginErr                  = Err{20001, "invalid username or password"}
	RequestParamErr           = Err{20002, "request param is error."}
	CreateDistributorErr      = Err{20100, "create distributor is failed."}
	CreateUserErr             = Err{20200, "create user is failed."}
	UpdateUserErr             = Err{20201, "update user is failed"}
	ResetUserPasswordErr      = Err{20202, "reset user password is failed"}
	OriginUserPasswordErr     = Err{20203, "user origin password is invalid"}
	UpdateUserPasswordErr     = Err{20204, "update user password is failed"}
	NotFoundUser              = Err{20205, "not found user"}
	UpdateMerchantStatusErr   = Err{20300, "update merchant status is failed"}
	NotFoundMerchant          = Err{20301, "not found merchant"}
	CreateMerchantRechargeErr = Err{20302, "create recharge apply is failed"}
	NotFoundAssetApplyErr     = Err{20400, "not found recharge apply"}
	NotFoundTicketErr         = Err{20601, "not found ticket"}

	//订单相关错误码
	NoAccountIdOrDistributorIdErr = Err{20501, "accountId or distributorId is null"}
	NoOrderNumberErr              = Err{20502, "orderNumber is null"}
	NoOrderFindErr                = Err{20503, "no order found by orderNumber"}
	CreateOrderErr                = Err{20504, "create order failed"}
	UpdateOrderErr                = Err{20505, "update order failed"}
	DeleteOrderErr                = Err{20506, "delete order failed"}
	IllegalSignErr                = Err{20507, "can not pass signing "}
	QuantityNotEnoughErr          = Err{20508, "distributor do not have enough number of coin"}
	NoSecretKeyFindErr            = Err{20509, "can not get secretkey according to apiKey"}
)
