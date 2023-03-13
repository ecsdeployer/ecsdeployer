package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestKeepInSyncTaskDefinitions_Unmarshal(t *testing.T) {
	tables := []struct {
		str       string
		expected  config.KeepInSyncTaskDefinitions
		errString string
	}{

		{`true`, config.KeepInSyncTaskDefinitionsEnabled, ""},
		{`false`, config.KeepInSyncTaskDefinitionsDisabled, ""},
		{`only_managed`, config.KeepInSyncTaskDefinitionsOnlyManaged, ""},

		{`"true"`, config.KeepInSyncTaskDefinitionsEnabled, ""},
		{`"false"`, config.KeepInSyncTaskDefinitionsDisabled, ""},
		{`ONLY_MANAGED`, config.KeepInSyncTaskDefinitionsOnlyManaged, ""},
	}

	for _, table := range tables {
		t.Run(table.str, func(t *testing.T) {
			obj, err := yaml.ParseYAMLString[config.KeepInSyncTaskDefinitions](table.str)
			if table.errString != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, table.errString)
				return
			}

			require.Equal(t, table.expected, *obj)

		})
	}
}
