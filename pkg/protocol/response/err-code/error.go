package err_code

type Err struct {
	ErrCode int
	ErrMsg  string
}

func (e Err) Data() (int,string) {
	return e.ErrCode,e.ErrMsg
}

var  (
	DistributorErr = Err{20100,"123"}

)