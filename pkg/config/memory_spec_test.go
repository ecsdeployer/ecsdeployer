package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestMemorySpec_Parse_Valid(t *testing.T) {

	cpuSpec := util.Must(config.NewCpuSpec(1024))

	tables := []struct {
		str      string
		expected int32
		isMultip bool
	}{
		{"1x", 1024, true},
		{"2x", 2048, true},
		{"x2", 2048, true},
		{"0.5x", 512, true},
		{"0.25x", 256, true},
		{"0.125x", 128, true},
		{"0.03125x", 32, true},

		{"512", 512, false},
		{"0.5gb", 512, false},
		{"0.5 gb", 512, false},
		{"0.5 GB", 512, false},
		{"0.25 GB", 256, false},
		{"2 GB", 2048, false},
		{"2g", 2048, false},
		{"2gb", 2048, false},
		{"2 gb", 2048, false},
	}

	sc := testutil.NewSchemaChecker(&config.MemorySpec{})

	for _, table := range tables {
		t.Run(table.str, func(t *testing.T) {
			memSpec, err := config.ParseMemorySpec(table.str)
			require.NoError(t, err)

			memValue, err := memSpec.MegabytesFromCpu(cpuSpec)
			require.NoError(t, err)
			require.Equal(t, table.expected, memValue)

			memValuePtr, err := memSpec.MegabytesPtrFromCpu(cpuSpec)
			require.NoError(t, err)
			require.Equal(t, table.expected, *memValuePtr)

			if !table.isMultip {
				require.Equal(t, table.expected, memSpec.GetValueOnly())
			}

			// attempt to parse it as well
			obj, err := yaml.ParseYAMLString[config.MemorySpec](table.str)
			require.NoError(t, err)
			require.NoError(t, sc.CheckYAML(t, table.str))
			require.True(t, memSpec.Equals(obj))
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

func TestMemorySpec_MegabytesFromCpu(t *testing.T) {

	cpuSpec := util.Must(config.NewCpuSpec(1024))

	_, err := (&config.MemorySpec{}).MegabytesFromCpu(cpuSpec)
	require.ErrorContains(t, err, "No value or mult")

	_, err = util.Must(config.ParseMemorySpec("2x")).MegabytesFromCpu(nil)
	require.ErrorContains(t, err, "CPU value needed")
}

func TestMemorySpec_Validate(t *testing.T) {

	err := (&config.MemorySpec{}).Validate()
	require.ErrorContains(t, err, "must specify memory")

}
