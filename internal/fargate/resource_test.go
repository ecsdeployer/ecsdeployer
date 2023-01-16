package fargate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFargateResourceStrings(t *testing.T) {
	tables := []struct {
		res FargateResource
		cpu string
		mem string
	}{
		{FargateResource{1024, 2048}, "1024", "2048"},
	}

	for _, table := range tables {
		require.Equal(t, table.cpu, table.res.CpuString())
		require.Equal(t, table.mem, table.res.MemoryString())
	}
}
