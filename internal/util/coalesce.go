package util

import "errors"

func Coalesce[T any](values ...*T) *T {
	for _, value := range values {
		if value != nil {
			return value
		}
	}

	return nil
}

func ShouldCoalesce[T any](values ...*T) (*T, error) {
	ret := Coalesce(values...)
	if ret == nil {
		return nil, errors.New("expected non-null, but got null")
	}
	return ret, nil
}

func MustCoalesce[T any](values ...*T) *T {
	ret, err := ShouldCoalesce(values...)
	if err != nil {
		panic(err)
	}
	return ret
}

type CoalesceFunc[T any] func(T) bool

// Returns the first non-nil element that the provided function returns true for
func CoalesceWithFunc[T any, PtrT *T](checkFunc CoalesceFunc[PtrT], values ...PtrT) (PtrT, bool) {
	for _, value := range values {
		if value == nil {
			continue
		}

		if checkFunc(value) {
			return value, true
		}
	}

	return nil, false
}
