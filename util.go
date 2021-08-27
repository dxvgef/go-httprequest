package httprequest

import (
	"unsafe"
)

// []byte转string
func BytesToStr(value []byte) string {
	return *(*string)(unsafe.Pointer(&value))
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
