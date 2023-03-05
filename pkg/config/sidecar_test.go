package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestSidecar(t *testing.T) {
	sc := testutil.NewSchemaChecker(&config.Sidecar{})

	tables := []struct {
		label         string
		str           string
		invalid       bool
		errorContains string
	}{
		{
			label: "normal props",
			str: `
			name: sidecar
			inherit_env: true
			memory_reservation: 256
			port_mappings:
				- 8080/tcp
			essential: true`,
		},

		{
			label:         "missing name",
			str:           `essential: true`,
			invalid:       true,
			errorContains: "you must set a name",
		},
	}

	for _, table := range tables {
		t.Run(table.label, func(t *testing.T) {
			yamlStr := testutil.CleanTestYaml(table.str)
			obj, err := yaml.ParseYAMLString[config.Sidecar](yamlStr)

			if table.invalid {
				require.Error(t, err)
				require.ErrorIs(t, err, config.ErrValidation)
				if table.errorContains != "" {
					require.ErrorContains(t, err, table.errorContains)
				}
				return
			}

			require.NoError(t, err)
			require.NoError(t, sc.CheckYAML(t, yamlStr))

			require.NotNil(t, obj)
			require.NotNil(t, obj.GetCommonContainerAttrs())

		})
	}
}
