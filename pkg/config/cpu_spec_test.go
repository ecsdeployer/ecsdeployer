package config_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestCpuSpec_Unmarshal(t *testing.T) {

	sc := testutil.NewSchemaChecker(config.CpuSpec(0))

	tables := []struct {
		str      string
		valid    bool
		expected int32
	}{
		{"1234", true, 1234},
		{"10", true, 10},
		{`16 vcpu`, true, 16384},
		{`0.5 vcpu`, true, 512},
		{`.5 vcpu`, true, 512},
		{`.25 vCPU`, true, 256},

		{`16 vcpus`, true, 16384},
		{`0.5 vcpus`, true, 512},
		{`.5 vcpus`, true, 512},
		{`.25 vCPUs`, true, 256},

		{`16 cores`, true, 16384},
		{`0.5 cores`, true, 512},
		{`.5 cores`, true, 512},
		{`.25 cores`, true, 256},

		{`-50`, false, 0},
		{"false", false, 0},
		{`. cores`, false, 0},
	}

	for testNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d_%s", testNum+1, table.str), func(t *testing.T) {
			cpu, err := yaml.ParseYAMLString[config.CpuSpec](table.str)

			if !table.valid {
				require.Errorf(t, err, "Parse failure")
				require.ErrorIs(t, err, config.ErrValidation)

				// we allow "invalid" things in the schema, but then error when parsing.
				// require.Errorf(t, sc.CheckYAML(t, table.str), "Schema Validation")
				return
			}
			require.NoError(t, err)
			require.NoError(t, sc.CheckYAML(t, table.str))

			require.NoError(t, cpu.Validate())
			require.EqualValues(t, table.expected, *cpu)
			require.Equal(t, table.expected, cpu.Shares())
		})
	}
}

func TestCpuSpec_NewCpuSpec(t *testing.T) {
	tables := []struct {
		value    int32
		valid    bool
		expected int32
	}{
		{1024, true, 1024},
		{-10, false, 0},
	}

	for _, table := range tables {
		cpu, err := config.NewCpuSpec(table.value)

		if !table.valid {
			require.Error(t, err)
			break
		}

		require.NoError(t, err)
		require.Equal(t, table.expected, cpu.Shares())

	}
}
