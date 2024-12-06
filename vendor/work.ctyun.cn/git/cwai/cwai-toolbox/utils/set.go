package utils

import (
	"reflect"
)

// IntSet 整数集合
type IntSet map[int]bool

// NewIntSet 创建新的集合
func NewIntSet(ints []int) (s IntSet) {
	s = make(IntSet)
	for _, i := range ints {
		s[i] = true
	}
	return s
}

// NewIntSetWith 创建新集合
func NewIntSetWith(f func(interface{}) int, values interface{}) (s IntSet) {
	s = make(IntSet)
	switch reflect.TypeOf(values).Kind() {
	case reflect.Slice:
		v := reflect.ValueOf(values)

		for i := 0; i < v.Len(); i++ {
			s[f(v.Index(i))] = true
		}
		break
	default:
		s[f(values)] = true
	}
	return s
}

// Add 添加元素
func (s IntSet) Add(ele int) {
	s[ele] = true
}

// Difference 差集
func (s IntSet) Difference(s2 IntSet) (diff []int) {
	for k := range s {
		if !s2[k] {
			diff = append(diff, k)
		}
	}
	return diff
}

// Intersection 交集
func (s IntSet) Intersection(s2 IntSet) (inter []int) {
	for k := range s {
		if s2[k] {
			inter = append(inter, k)
		}
	}
	return inter
}

// Has 是否包含
func (s IntSet) Has(i int) bool {
	return s[i]
}

// Slice 返回slice
func (s IntSet) Slice() []int {
	var ints []int
	for k := range s {
		ints = append(ints, k)
	}
	return ints
}
