package util

import (
	"testing"
)

func TestIsBlank(t *testing.T) {
	tables := []struct {
		str      string
		expected bool
	}{
		{"", true},
		{" ", true},
		{"   ", true},

		{" x  ", false},
		{"test", false},
		{"x", false},
	}

	for _, table := range tables {
		actual := IsBlank(&table.str)
		if actual != table.expected {
			t.Errorf("Expected IsBlank(%s) to be %t, but it was %t", table.str, table.expected, actual)
		}
	}

	if !IsBlank(nil) {
		t.Error("Expected IsBlank(nil) to be true, but it was false")
	}
}
