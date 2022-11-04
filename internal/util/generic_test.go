package util

import (
	"errors"
	"testing"
)

func TestFirstParam(t *testing.T) {

	if FirstParam(1, 2, 3, 4, 5) != 1 {
		t.Fatal("it's broken")
	}
}

func mustGood() (string, error) {
	return "yay", nil
}

func mustBad() (string, error) {
	return "", errors.New("oopsies")
}

func TestMust(t *testing.T) {

	if Must(mustGood()) != "yay" {
		t.Fatal("Must() is broken")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Must() should have panicked")
		}
	}()

	_ = Must(mustBad())
}

func TestPtr(t *testing.T) {
	value := "test"

	resp := Ptr(value)

	if value != *resp {
		t.Fatal("Ptr() is broken")
	}

}
