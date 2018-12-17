package utils

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"strconv"
)

func TransformTypeToString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case int:
		return strconv.Itoa(v.(int))
	case string:
		return v.(string)
	case float64:
		return strconv.FormatFloat(v.(float64),
			'f', -1, 64)
	case bool:
		if v.(bool) {
			return "true"
		}
		return "false"
	case int64:
		return strconv.Itoa(v.(int))
	default:
		println(reflect.TypeOf(t).String())
		return ""
	}
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
