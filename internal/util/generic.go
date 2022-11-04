package util

// https://github.com/icza/gog/blob/main/gog.go

func Ptr[T any](v T) *T {
	return &v
}

// returns the first value of a multi-return statement
func FirstParam[T any, U any](fp T, _ ...U) T {
	return fp
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
