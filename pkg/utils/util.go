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
		return strconv.FormatInt(v.(int64), 10)
	default:
		println(reflect.TypeOf(t).String())
		return ""
	}
}

//取出和第一个数组出现过的元素
func DiffSet(list ...[]int64) []int64 {
	Log.Debugf("%v", list)
	if len(list) == 1 {
		return list[0]
	}
	if len(list) <= 0 {
		return nil
	}
	first := list[0]
	temp := make(map[int64]int)
	for _, v := range first {
		temp[v] = 1
	}

	for i, v := range list {
		if i > 0 {
			for _, v1 := range v {
				if temp[v1] == 1 {
					temp[v1] = 2
				}
			}
		}
	}
	var result []int64
	for k, v := range temp {
		if v == 1 {
			result = append(result, k)
		}
	}
	return result
}

func MergeList(l1, l2, l3 []int64) []int64 {
	var result []int64
	tempMap := make(map[int64]int)

	for _, v := range l1 {
		tempMap[v] = 1
	}

	for _, v := range l2 {
		if tempMap[v] == 1 {
			tempMap[v] = 2
		}
	}

	for _, v := range l3 {
		if tempMap[v] == 2 {
			tempMap[v] = 3
			result = append(result, v)
		}
	}

	return result
}

func ConvertStringToInt(ids []string, results *[]int64) error {
	var result []int64
	for _, id := range ids {
		if temp, err := strconv.ParseInt(id, 10, 64); err != nil {
			return err
		} else {
			result = append(result, temp)
		}
	}
	*results = result
	return nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
