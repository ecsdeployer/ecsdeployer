package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestVolume(t *testing.T) {
	sc := testutil.NewSchemaChecker(&config.Volume{})

	tables := []struct {
		label         string
		str           string
		invalid       bool
		errorContains string
		name          string
		hasEFS        bool
	}{
		{
			label: "shorthand",
			str:   `bindvol`,
			name:  "bindvol",
		},
		{
			label: "efs block",
			str: `
			name: efstest
			efs:
				file_system_id: fs-123`,
			name:   "efstest",
			hasEFS: true,
		},

		{
			label: "bad efs block",
			str: `
			name: efstest
			efs:
				access_point_id: fs-123`,
			invalid:       true,
			errorContains: "provide a FileSystemID",
		},
		{
			label:         "missing name",
			str:           `name: ""`,
			invalid:       true,
			errorContains: "volume name cannot be empty",
		},
	}

	for _, table := range tables {
		t.Run(table.label, func(t *testing.T) {
			yamlStr := testutil.CleanTestYaml(table.str)
			obj, err := yaml.ParseYAMLString[config.Volume](yamlStr)

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

			require.EqualValues(t, table.name, *awsObj.Name, "Name")
			if table.hasEFS {
				require.NotNil(t, awsObj.EfsVolumeConfiguration)

			} else {
				require.Nil(t, awsObj.EfsVolumeConfiguration)
			}
		})
	}
}
