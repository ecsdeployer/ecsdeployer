package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestVolumeList(t *testing.T) {
	sc := testutil.NewSchemaChecker(&config.VolumeList{})

	tables := []struct {
		label         string
		str           string
		invalid       bool
		errorContains string
		length        int
	}{
		{
			label: "normal array",
			str: `
			- bindvol
			- bind2`,
			length: 2,
		},

		{
			label: "duped array",
			str: `
			- bindvol
			- bindvol
			- bind2`,
			invalid:       true,
			errorContains: "Duplicate volume name: bindvol",
		},

		{
			label: "array with child error",
			str: `
			- bindvol
			- bind2
			- name: badvol
				efs:
					access_point_id: ap123`,
			invalid:       true,
			errorContains: "FileSystemID",
		},
	}

	for _, table := range tables {
		t.Run(table.label, func(t *testing.T) {
			yamlStr := testutil.CleanTestYaml(table.str)
			objPtr, err := yaml.ParseYAMLString[config.VolumeList](yamlStr)

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

			// the parser returns a ptr
			obj := *objPtr

			require.Len(t, obj, table.length)

			awsObj := obj.ToAws()
			require.Len(t, awsObj, table.length)
		})
	}
}

func TestMergeVolumeLists(t *testing.T) {
	volMap := config.VolumeList{
		"vol1": config.Volume{Name: "vol1"},
		"vol2": config.Volume{Name: "vol2"},
	}

	volMap2 := config.VolumeList{
		"vol3": config.Volume{Name: "vol3"},
		"vol4": config.Volume{Name: "vol4"},
	}

	volMap3 := config.VolumeList{
		"vol1": config.Volume{Name: "volx"},
		"vol5": config.Volume{Name: "vol5"},
	}

	t.Run("all nils", func(t *testing.T) {
		result := config.MergeVolumeLists(nil, nil)
		require.NotNil(t, result)
		require.Len(t, result, 0)
	})

	t.Run("value with nil", func(t *testing.T) {
		result := config.MergeVolumeLists(volMap, nil)
		require.NotNil(t, result)
		require.Len(t, result, 2)
		require.Contains(t, result, "vol1")
		require.Contains(t, result, "vol2")
	})

	t.Run("nil with value", func(t *testing.T) {
		result := config.MergeVolumeLists(nil, volMap)
		require.NotNil(t, result)
		require.Len(t, result, 2)
		require.Contains(t, result, "vol1")
		require.Contains(t, result, "vol2")
	})

	t.Run("overrides", func(t *testing.T) {
		result := config.MergeVolumeLists(volMap, volMap3)
		require.NotNil(t, result)
		require.Len(t, result, 3)
		require.Contains(t, result, "vol1")
		require.Contains(t, result, "vol2")
		require.Contains(t, result, "vol5")
		require.Equal(t, volMap3["vol1"].Name, result["vol1"].Name)
	})

	t.Run("no conflicts", func(t *testing.T) {
		result := config.MergeVolumeLists(volMap, volMap2)
		require.NotNil(t, result)
		require.Len(t, result, 4)
		require.Contains(t, result, "vol1")
		require.Contains(t, result, "vol2")
		require.Contains(t, result, "vol3")
		require.Contains(t, result, "vol4")
	})
}
