package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFirstParam(t *testing.T) {
	require.Equal(t, 1, FirstParam(1, 2, 3, 4, 5))
}

func mustGood() (string, error) {
	return "yay", nil
}

func mustBad() (string, error) {
	return "", errors.New("oopsies")
}

func TestMust(t *testing.T) {

	require.Equal(t, "yay", Must(mustGood()))

	require.Panics(t, func() {
		_ = Must(mustBad())
	})
}

func TestPtr(t *testing.T) {
	value := "test"
	resp := Ptr(value)
	require.Equal(t, value, *resp)
}
