package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestUlimit(t *testing.T) {
	sc := testutil.NewSchemaChecker(&config.Ulimit{})

	tables := []struct {
		label         string
		str           string
		invalid       bool
		errorContains string
		name          string
		hard          int
		soft          int
	}{
		{
			label: "shorthand",
			str: `
			name: nofile
			limit: 1234`,
			name: "nofile",
			soft: 1234,
			hard: 1234,
		},
		{
			label: "common",
			str: `
			name: nofile
			soft: 8192
			hard: 8192`,
			name: "nofile",
			soft: 8192,
			hard: 8192,
		},

		{
			label: "missing name",
			str: `
			name: ""
			soft: 8192
			hard: 8192`,
			invalid:       true,
			errorContains: "must provide a name",
		},

		{
			label: "only soft",
			str: `
			name: nofile
			soft: 8192`,
			name: "nofile",
			soft: 8192,
			hard: 8192,
		},

		{
			label: "only hard",
			str: `
			name: nofile
			hard: 8192`,
			name: "nofile",
			soft: 0,
			hard: 8192,
		},

		{
			label: "soft bigger than hard",
			str: `
			name: nofile
			soft: 16384
			hard: 8192`,
			invalid:       true,
			errorContains: "soft limit cannot be higher than hard",
		},

		{
			label:         "missing values",
			str:           `name: nofile`,
			invalid:       true,
			errorContains: "must provide a value for the hard",
		},
	}

	for _, table := range tables {
		t.Run(table.label, func(t *testing.T) {
			yamlStr := testutil.CleanTestYaml(table.str)
			obj, err := yaml.ParseYAMLString[config.Ulimit](yamlStr)

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

			awsObj := obj.ToAws()

			require.EqualValues(t, table.name, awsObj.Name, "Name")
			require.EqualValues(t, table.soft, awsObj.SoftLimit, "Soft")
			require.EqualValues(t, table.hard, awsObj.HardLimit, "Hard")
		})
	}
}
