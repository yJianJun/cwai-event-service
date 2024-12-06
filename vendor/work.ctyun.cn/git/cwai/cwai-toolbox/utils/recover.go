package utils

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/golang/glog"
)

// WithRecover 带recover panic调用方法,避免goroutine 里panic
// fn 必须是func, args是fn的参数
// 使用方法： go util.WithRecover(foo)
func WithRecover(fn interface{}, args ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			if n < len(buf) {
				buf = buf[:n]
			}

			glog.Warningf("panic %+v, stack: %s", err, strings.ReplaceAll(string(buf), "\n", "\t"))
		}
	}()

	argsv := make([]reflect.Value, len(args))
	for i, v := range args {
		argsv[i] = reflect.ValueOf(v)
	}
	fnv := reflect.ValueOf(fn)
	fnv.Call(argsv)
}
