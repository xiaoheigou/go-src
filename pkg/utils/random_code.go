package utils

import (
	"crypto/rand"
	"math/big"
	"strconv"
)
import "sync/atomic"

type count32 int32
// 用户来标记产生的第多少个随机码
var RandomSeq count32
func (c *count32) GetCount() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

// 产生安全随机码
func GetSecuRandomCode() (string, error) {
	var max int64 = 99999
	var min int64 = 10000

	nBig, err := rand.Int(rand.Reader, big.NewInt(max - min))
	if err != nil {
		Log.Errorf("Can't generate secure random code")
		return "", err
	}

	return strconv.Itoa(int(nBig.Int64() + min)), nil
}