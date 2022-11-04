package yaml

import (
	"io"
	"os"
)

// This is mainly used in tests, normal app operations should not need this
func ParseYAMLFile[T interface{}](file string) (*T, error) {

	f, err := os.Open(file) // #nosec
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return ParseYAMLString[T](string(data))
}

func ParseYAMLString[T interface{}](data string) (*T, error) {

	obj := *new(T)

	if err := UnmarshalStrict([]byte(data), &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}
