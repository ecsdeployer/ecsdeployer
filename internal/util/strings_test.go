package util

import (
	"fmt"
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

func TestLongestCommonPrefix(t *testing.T) {
	tables := []struct {
		prefix string
		strs   []string
	}{
		{
			prefix: "dummy-",
			strs: []string{
				"dummy-thing",
				"dummy-thing2",
				"dummy-svc2",
				"dummy-svc2",
			},
		},
		{
			prefix: "dum",
			strs: []string{
				"dummy-thing",
				"dummy-thing2",
				"dum",
			},
		},

		{"", nil},
		{"", []string{}},
	}

	for tNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", tNum+1), func(t *testing.T) {
			require.Equal(t, table.prefix, LongestCommonPrefix(table.strs))
		})
	}
}
