package utils

import (
	"runtime"
)

// GetCallFuncName 获取调用函数名称
func GetCallFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
