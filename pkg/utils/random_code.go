package utils

import (
	"crypto/rand"
	mathrand "math/rand"
	"math/big"
	"strconv"
	"time"
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


// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var src = mathrand.NewSource(time.Now().UnixNano())
// 下面函数用于生成随机字符串。注：它不是密码学安全的，不能用于随机性要求高的场景！
func GetRandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}