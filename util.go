package httprequest

import (
	"unsafe"
)

// []byte转string
func BytesToStr(value []byte) string {
	return *(*string)(unsafe.Pointer(&value))
}

// 字符串转[]byte
func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// 在[]int中判断是否存在指定的值
func inIntSlice(status int, retryStatus []int) bool {
	for k := range retryStatus {
		if status == retryStatus[k] {
			return true
		}
	}
	return false
}