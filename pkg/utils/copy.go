package utils

import (
	"fmt"
	"reflect"
)

// CopyProperties 将 source 结构体的字段值拷贝到 dst 结构体的相同名称字段
func CopyProperties(dst, src interface{}) {
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
		//如果source字段值没有值，则忽略不拷贝
		if IsInvalid(value) {
			continue
		}
		if value.Type() == dvalue.Type() && value.IsValid() {
			dvalue.Set(value)
		}
	}
}

func IsInvalid(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.Array:
		n := v.Len()
		for i := 0; i < n; i++ {
			if !IsInvalid(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.String:
		return v.Len() == 0
	case reflect.Struct:
		n := v.NumField()
		for i := 0; i < n; i++ {
			if !IsInvalid(v.Field(i)) && v.Type().Field(i).Name != "_" {
				return false
			}
		}
		return true
	case reflect.Invalid:
		return true
	default:
		panic(&reflect.ValueError{"reflect.Value.IsInvalid", v.Kind()})
	}
}
