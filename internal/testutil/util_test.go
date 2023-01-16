package testutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStripIndentation(t *testing.T) {
	tables := []struct {
		str      string
		expected string
	}{
		{"", ""},
		{"test", "test"},
		{"    test", "test"},
		{" \t test", "test"},
		{"\t\ttest", "test"},
		{"\n\ttest", "\ntest"},
		{"\ttest\n\ttest", "test\ntest"},
		{"\ttest\n\ttest        ", "test\ntest        "},
		{"\n\ttest\n\ttest        ", "\ntest\ntest        "},
	}

	for tNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", tNum+1), func(t *testing.T) {
			require.Equal(t, table.expected, StripIndentation(table.str))
		})
	}
}
