package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestMemorySpec_Parse_Valid(t *testing.T) {

	cpuSpec := util.Must(config.NewCpuSpec(1024))

	tables := []struct {
		str      string
		expected int32
	}{
		{"1x", 1024},
		{"2x", 2048},
		{"x2", 2048},
		{"0.5x", 512},
		{"0.25x", 256},
		{"0.125x", 128},
		{"0.03125x", 32},
		{"512", 512},
		{"0.5gb", 512},
		{"0.5 gb", 512},
		{"0.5 GB", 512},
		{"0.25 GB", 256},
		{"2 GB", 2048},
		{"2g", 2048},
		{"2gb", 2048},
		{"2 gb", 2048},
	}

	for _, table := range tables {
		t.Run(table.str, func(t *testing.T) {
			memSpec, err := config.ParseMemorySpec(table.str)
			require.NoError(t, err)

			memValue, err := memSpec.MegabytesFromCpu(cpuSpec)
			require.NoError(t, err)

			require.Equal(t, table.expected, memValue)
		})
	}
}

func TestMemorySpec_Parse_Invalid(t *testing.T) {
	tables := []struct {
		str string
	}{
		{"1xx"},
		{"xx2"},
	}

	for _, table := range tables {
		_, err := config.ParseMemorySpec(table.str)
		require.ErrorContains(t, err, "Invalid memory specification format")

	}
}

func TestMemorySpec_MarshalJSON(t *testing.T) {
	tables := []struct {
		str      string
		expected string
	}{
		{"1x", `"1x"`},
		{"1024", `1024`},
		{"0.5gb", `512`},
	}

	for _, table := range tables {
		memSpec, err := config.ParseMemorySpec(table.str)
		require.NoError(t, err)

		jsonStr, err := util.Jsonify(memSpec)
		require.NoError(t, err)

		require.Equal(t, table.expected, jsonStr)

	}
}
