package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestFirelensAwsLogGroup(t *testing.T) {
	t.Run("Enabled", func(t *testing.T) {
		require.True(t, (&config.FirelensAwsLogGroup{Path: "test"}).Enabled())
		require.False(t, (&config.FirelensAwsLogGroup{Path: ""}).Enabled())
		require.False(t, (&config.FirelensAwsLogGroup{}).Enabled())
	})

	t.Run("UnmarshalYAML", func(t *testing.T) {
		tables := []struct {
			str     string
			path    string
			invalid bool
			enabled bool
		}{
			{`false`, "", false, false},
			{`/test/log`, "/test/log", false, true},

			{`true`, "", true, false},
			{`something: wrong`, "", true, false},
		}

		sc := testutil.NewSchemaChecker(&config.FirelensAwsLogGroup{})

		for _, table := range tables {
			t.Run(table.str, func(t *testing.T) {
				obj, err := yaml.ParseYAMLString[config.FirelensAwsLogGroup](table.str)
				if table.invalid {
					require.Error(t, err)
					// require.ErrorIs(t, err, config.ErrValidation) // some of them are yaml.TypeError
					return
				}

				require.Equal(t, table.path, obj.Path)
				require.Equal(t, table.enabled, obj.Enabled())

				require.NoError(t, sc.CheckYAML(t, table.str))
			})
		}
	})
}
