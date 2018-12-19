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
	CreateUserErr             = Err{20400, "create user is failed."}
	UpdateUserErr             = Err{20202, "update user is failed"}
	ResetUserPasswordErr      = Err{20202, "reset user password is failed"}
	OriginUserPasswordErr     = Err{20202, "user origin password is invalid"}
	UpdateUserPasswordErr     = Err{20202, "update user password is failed"}
	NotFoundUser              = Err{20401, "not found user"}
	UpdateMerchantStatusErr   = Err{20202, "update merchant status is failed"}
	NotFoundMerchant          = Err{20201, "not found merchant"}
	CreateMerchantRechargeErr = Err{20203, "create recharge apply is failed"}
	NotFoundAssetApplyErr     = Err{20204, "not found recharge apply"}

	//订单相关错误码
	NoAccountIdOrDistributorIdErr = Err{20501, "accountId or distributorId is null"}
	NoOrderNumberErr              = Err{20502, "orderNumber is null"}
	NoOrderFindErr                = Err{20503, "no order found by orderNumber"}
	CreateOrderErr                = Err{20504, "create order failed"}
	UpdateOrderErr                = Err{20505, "update order failed"}
	DeleteOrderErr                = Err{20506, "delete order failed"}
)
