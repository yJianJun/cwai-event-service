package util

import (
	"fmt"
	"reflect"
)

// CopyProperties 将 src 结构体的字段值拷贝到 dst 结构体的相同名称字段
func CopyProperties(dst, src interface{}) error {
	sval := reflect.ValueOf(src).Elem()
	dval := reflect.ValueOf(dst).Elem()
	for i := 0; i < sval.NumField(); i++ {
		value := sval.Field(i)
		name := sval.Type().Field(i).Name
		if dval.Kind() == reflect.Ptr {
			dval = dval.Elem()
		}
		dvalue := dval.FieldByName(name)
		if !dvalue.IsValid() {
			continue
		}
		stype := uint(sval.Type().Field(i).Type.Kind())
		dtype := uint(dval.Type().Field(i).Type.Kind())
		fmt.Println(stype, "->", dtype)
		if value.IsZero() {
			continue
		}
		if value.Type() == dvalue.Type() && value.IsValid() {
			dvalue.Set(value)
		}
	}
	return nil
}
