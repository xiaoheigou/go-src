package utils

import "strings"

func IsValidPhone(nationCode string, phone string) bool {
	// 每个国家手机号码不一样，目前仅做简单的校验

	// 检查是否都是数字
	for _, c := range phone {
		if c < '0' || c > '9' {
			Log.Warnf("Invalid phone [%v], phone can only contains 0-9", phone)
			return false
		}
	}

	if nationCode == "86" {
		// 中国手机号，仅检测长度是否为11
		if len(phone) != 11 {
			Log.Warnf("Invalid phone [%v], length of phone must be 11", phone)
			return false
		} else {
			return true
		}
	} else {
		// 其它国家手机号，简单地限制一下长度
		minLen := 3
		maxLen := 20
		if len(phone) < minLen || len(phone) > maxLen {
			Log.Warnf("Invalid phone [%v], length of phone is too small or long", phone)
			return false
		}
		return true
	}
}

func IsValidNationCode(nationCode string) bool {
	minLen := 1
	maxLen := 10
	if len(nationCode) < minLen || len(nationCode) > maxLen {
		Log.Warnf("Invalid nation code [%v], length of nation code is too small or long", nationCode)
		return false
	}

	// 检查是否都是数字
	for _, c := range nationCode {
		if c < '0' || c > '9' {
			Log.Warnf("Invalid nation code [%v], nation code can only contains 0-9", nationCode)
			return false
		}
	}

	return true
}

func IsValidEmail(email string) bool {
	// 目前对邮箱仅做简单的校验

	if len(email) == 0 {
		Log.Warnf("Invalid email [%v]", email)
		return false
	}

	index := strings.Index(email, "@")
	if index <= 0 {
		// 不存在@或者第一个字符为@，则非法
		Log.Warnf("Invalid email [%v]", email)
		return false
	}

	// 最后一个字符为@，则非法
	if index == len(email) - 1 {
		Log.Warnf("Invalid email [%v]", email)
		return false
	}

	return true
}