package yaml

import (
	"io"
	"os"
	"strings"
)

// This is mainly used in tests, normal app operations should not need this
func ParseYAMLFile[T any](file string) (*T, error) {

	f, err := os.Open(file) // #nosec
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParseYAML[T](f)
}

// This is mainly used in tests, normal app operations should not need this
func ParseYAML[T any](f io.Reader) (*T, error) {

	obj := *new(T)

	if err := unmarshalReader(f, true, &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

func ParseYAMLString[T any](data string) (*T, error) {
	return ParseYAML[T](strings.NewReader(data))
}
