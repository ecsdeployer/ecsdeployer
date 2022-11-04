package util

import (
	"reflect"
)

func IsNilable(v reflect.Value) bool {
	switch v.Type().Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return true
	default:
		return false
	}
}

// Must provide a pointer as the haystack
//
// this will recurse into a struct and find all the values of type T within that struct,
// as well as any child structs
func DeepFindInStruct[T interface{}](haystack interface{}) []*T {

	needleV := reflect.TypeOf(new(T)).Elem()
	// needleP := reflect.PointerTo(needleV)

	// needleVKind := needleV.Kind()
	// needlePKind := needleP.Kind()

	// fmt.Printf("needleVType=%v, needlePtype=%v nvKind=%v npKind=%v\n", needleV, needleP, needleVKind, needlePKind)

	haystackVal := reflect.Indirect(reflect.ValueOf(haystack))
	// haystackVal := reflect.ValueOf(haystack).Elem()
	haystackType := haystackVal.Type()

	// fmt.Printf("Haystacktype=%v\n", haystackType)

	results := make([]*T, 0)
	// resultValue := reflect.ValueOf(results)

	for _, field := range reflect.VisibleFields(haystackType) {
		f := haystackVal.FieldByIndex(field.Index)
		// fmt.Printf("FIELD Name=%s Type=%v Kind=%v Exp=%t\n", field.Name, f.Type(), f.Kind(), field.IsExported())

		if !field.IsExported() {
			// ignore unexported
			continue
		}

		// ignore unexported fields
		// if !f.CanSet() {
		// 	fmt.Println("CAN NOT SET")
		// 	continue
		// }

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
			results = append(results, value.Addr().Interface().(*T))
			continue
		}

		if value.Kind() == reflect.Struct {
			// fmt.Printf("STRUCT!\n")
			res := DeepFindInStruct[T](value.Interface())
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
					res := DeepFindInStruct[T](item.Interface())
					results = append(results, res...)
					continue
				}

				if item.Kind() == reflect.Array || item.Kind() == reflect.Slice {
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
