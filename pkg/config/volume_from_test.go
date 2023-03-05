package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestVolumeFrom(t *testing.T) {
	sc := testutil.NewSchemaChecker(&config.VolumeFrom{})

	tables := []struct {
		label         string
		str           string
		invalid       bool
		errorContains string
		source        string
		readonly      bool
	}{
		{
			label:  "shorthand",
			str:    `somecontainer`,
			source: "somecontainer",
		},
		{
			label: "blockmode",
			str: `
			source: cont1
			readonly: false`,
			source:   "cont1",
			readonly: false,
		},

		{
			label: "blockmode-ro",
			str: `
			source: cont1
			readonly: true`,
			source:   "cont1",
			readonly: true,
		},

		{
			label: "blank cont",
			str: `
			source: ""
			readonly: false`,
			invalid:       true,
			errorContains: "source container cannot be empty",
		},
	}

	for _, table := range tables {
		t.Run(table.label, func(t *testing.T) {
			yamlStr := testutil.CleanTestYaml(table.str)
			obj, err := yaml.ParseYAMLString[config.VolumeFrom](yamlStr)

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

			require.Equal(t, table.source, obj.SourceContainer, "Obj: SourceContainer")
			require.Equal(t, table.readonly, obj.ReadOnly, "Obj: ReadOnly")

			awsObj := obj.ToAws()
			require.Equal(t, table.source, *awsObj.SourceContainer, "Aws: SourceContainer")
			require.Equal(t, table.readonly, *awsObj.ReadOnly, "Aws: ReadOnly")
		})
	}
}
