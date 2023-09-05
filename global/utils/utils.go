package utils

import (
	"unsafe"
)

// StrToBytes 字符串转换字节
func StrToBytes(s string) []byte {
	t := (*[2]uintptr)(unsafe.Pointer(&s))
	b := [3]uintptr{t[0], t[1], t[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}

// BytesToStr 字节转字符串
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// JudgeDuplication 判断是否重复
func JudgeDuplication(validValue string, validList []string) bool {
	found := false
	for _, item := range validList {
		if item == validValue {
			found = true
			break
		}
	}

	return found
}

// RemoveDuplicates 去重
func RemoveDuplicates(slice []string) []string {
	var result []string

	for _, item := range slice {
		found := false
		for _, val := range result {
			if val == item {
				found = true
				break
			}
		}
		if !found {
			result = append(result, item)
		}
	}

	return result
}
