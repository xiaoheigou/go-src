package utils

import (
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
