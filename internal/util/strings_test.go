package util

import (
	"testing"

	"github.com/stretchr/testify/require"
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
		require.Equalf(t, table.expected, actual, table.str)
	}

	require.Truef(t, IsBlank(nil), "IsBlank(nil)")
}
