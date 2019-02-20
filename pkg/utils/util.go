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

// 取出仅在第一个数组中出现过的元素
func DiffSet(list ...[]int64) []int64 {
	Log.Debugf("func DiffSet arguments: %v", list)
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

// 取出同时出现在两个数组中的元素，且保证结果集中元素顺序和第一个数组保持一致
// 如果list1中元素不唯一，这个函数不保证结果集中元素的唯一性
func InterSetInt64(list1, list2 []int64) []int64 {
	var result []int64

	// 先把list2保存到hash表myMap中
	var myMap = make(map[int64]int)
	for _, v := range list2 {
		myMap[v] = 1
	}

	// 遍历list1，如果元素存在于myMap中，则增加到result中
	for _, v := range list1 {
		if _, ok := myMap[v]; ok {
			result = append(result, v)
		}
	}

	return result
}

//取交集
func IntersectList(list ...[]int64) []int64 {
	var result []int64
	tempMap := make(map[int64]int)

	sign := len(list) - 1

	for i, v := range list {
		if len(v) == 0 {
			return result
		}
		for _, v1 := range v {
			if i < sign && tempMap[v1] == i {
				tempMap[v1] = tempMap[v1] + 1
			} else if i == sign && sign == tempMap[v1] {
				tempMap[v1] = i + 1
				result = append(result, v1)
			}
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

//移除数组元素
func RemoveElement(nums []int64, val int64) []int64 {
	//如果是空切片，那就返回0
	if len(nums) == 0 {
		return nil
	}
	//用一个索引
	//循环去比较
	//当一样的时候就删除对应下标的值
	//当不一样的时候，索引加1
	index := 0
	for index < len(nums) {
		if nums[index] == val {
			nums = append(nums[:index], nums[index+1:]...)
			continue
		}
		index++
	}
	return nums
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
