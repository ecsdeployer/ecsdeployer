package util

import (
	"reflect"
)

// Returns whether or not a given value could be set to nil
func IsNilable(v reflect.Value) bool {
	switch v.Type().Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return true
	default:
		return false
	}
}

type DeepFindInStructOptions struct {
	// If true, this will return struct VALUE fields that have a zero value for the requested type
	// By default if you have a value field (vs a pointer field), with a zero struct value,
	// it will not be returned as it is assumed to have not been provided
	IncludeZeroStruct bool
}

// Must provide a pointer as the haystack
// See DeepFindInStructAdvanced
// This version will ignore zero values for structs
func DeepFindInStruct[T interface{}](haystack interface{}) []*T {
	return DeepFindInStructAdvanced[T](haystack, &DeepFindInStructOptions{})
}

// Must provide a pointer as the haystack
//
// this will recurse into a struct and find all the values of type T within that struct,
// as well as any child structs
//
// Note: this will stop searching once it finds a match. If you have self-referencing structs,
// it will not find child structs within a matching struct
func DeepFindInStructAdvanced[T interface{}](haystack interface{}, options *DeepFindInStructOptions) []*T {

	needleV := reflect.TypeOf(new(T)).Elem()

	haystackVal := reflect.Indirect(reflect.ValueOf(haystack))
	haystackType := haystackVal.Type()

	results := make([]*T, 0)

	for _, field := range reflect.VisibleFields(haystackType) {
		f := haystackVal.FieldByIndex(field.Index)
		// fmt.Printf("FIELD Name=%s Type=%v Kind=%v Exp=%t\n", field.Name, f.Type(), f.Kind(), field.IsExported())

		if !field.IsExported() {
			// ignore unexported
			continue
		}

		value := f

		// cant have anything if it is null
		if IsNilable(f) {
			if f.IsNil() {
				continue
			}

			if f.Kind() == reflect.Ptr || f.Kind() == reflect.Interface {
				value = f.Elem()
			}
		}

		// value := reflect.Indirect(f)

		if value.Type() == needleV {
			switch {
			case value.CanAddr():
				results = append(results, value.Addr().Interface().(*T))

			case options.IncludeZeroStruct && value.Kind() == reflect.Struct && value.IsZero():
				zeroStruct := value.Interface().(T)
				results = append(results, &zeroStruct)
			}
			continue
		}

		if value.Kind() == reflect.Struct {
			res := DeepFindInStructAdvanced[T](value.Interface(), options)
			results = append(results, res...)
			continue
		}

		if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
			for i := 0; i < value.Len(); i++ {
				item := value.Index(i)

				if IsNilable(item) {
					if item.IsNil() {
						continue
					}

					if item.Kind() == reflect.Ptr || item.Kind() == reflect.Interface {
						item = item.Elem()
					}
				}

				if item.Type() == needleV {
					results = append(results, item.Addr().Interface().(*T))
					continue
				}

				if item.Kind() == reflect.Struct {
					res := DeepFindInStructAdvanced[T](item.Interface(), options)
					results = append(results, res...)
					continue
				}

				if item.Kind() == reflect.Array || item.Kind() == reflect.Slice {
					// TODO: Add ability to search thru nested lists
					panic("array/slice within slice not supported")
				}
			}
			continue
		}
	}

	return results
}

type thinger struct{}

var _ = DeepFindInStruct[thinger](&thinger{})
var _ = DeepFindInStructAdvanced[thinger](&thinger{}, &DeepFindInStructOptions{})
