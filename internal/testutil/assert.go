package testutil

import (
	"reflect"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"golang.org/x/exp/slices"
)

func AssertStringEquals(val1 any, val2 any) bool {
	return AssertEquals[string](val1, val2)
}

func AssertSliceEquals[T comparable](val1 []T, val2 []T) bool {
	return slices.CompareFunc(val1, val2, func(v1 T, v2 T) int {
		res := AssertEquals[T](v1, v2)
		if res {
			return 0
		}
		return -1
	}) == 0
}

func AssertEquals[T comparable](val1 any, val2 any) bool {

	// if val1 == nil || val2 == nil {
	// 	// fmt.Printf("CHECK: v1=%t v2=%t\n", val1 == nil, val2 == nil)
	// 	return val1 == nil && val2 == nil
	// }

	v1 := reflect.ValueOf(val1)
	v2 := reflect.ValueOf(val2)

	if util.IsNilable(v1) && util.IsNilable(v2) {
		if v1.IsNil() != v2.IsNil() {
			return false
		}
	}

	iv1 := reflect.Indirect(v1)
	iv2 := reflect.Indirect(v2)

	if reflect.TypeOf(iv1) != reflect.TypeOf(iv2) {
		panic("type mismatch for AssertEquals")
	}

	if iv1.IsValid() != iv2.IsValid() {
		return false
	}

	if !iv1.IsValid() {
		return true
	}

	if iv1.IsZero() != iv2.IsZero() {
		return false
	}

	if iv1.CanInterface() != iv2.CanInterface() {
		return false
	}

	value1 := iv1.Interface().(T)
	value2 := iv2.Interface().(T)

	return value1 == value2
}
