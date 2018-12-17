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
	UpdateMerchantStatusErr       = Err{20202, "update merchant status is failed"}
	CreateMerchantRechargeErr     = Err{20203, "create recharge apply is failed"}
	NotFoundAssetApplyErr         = Err{20204, "not found recharge apply"}

	//订单相关错误码
	NoAccountIdOrDistributorIdErr = Err{20501, "accountId or distributorId is null"}
	NoOrderNumberErr              = Err{20502, "orderNumber is null"}
	NoOrderFindErr                = Err{20503, "no order found by orderNumber"}
	CreateOrderErr                = Err{20504, "create order failed"}
	UpdateOrderErr                = Err{20505, "update order failed"}
	DeleteOrderErr                = Err{20506, "delete order failed"}
)
