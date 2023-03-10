package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestMounts_Unmarshal(t *testing.T) {
	sc := testutil.NewSchemaChecker(&config.Mount{})

	tables := []struct {
		label         string
		str           string
		invalid       bool
		errorContains string
		path          string
		source        string
		readonly      bool
	}{
		{
			label: "everything",
			str: `
			path: /mnt/test
			source: testvol
			readonly: true`,
			path:     "/mnt/test",
			source:   "testvol",
			readonly: true,
		},

		{
			label: "normal readwrite",
			str: `
			path: /mnt/rwtest
			source: somevol`,
			path:   "/mnt/rwtest",
			source: "somevol",
		},

		{
			label:         "missing path",
			str:           `source: testvol2`,
			invalid:       true,
			errorContains: "path cannot be empty",
		},

		{
			label:         "missing source",
			str:           `path: /blah`,
			invalid:       true,
			errorContains: "source cannot be empty",
		},
	}

	for _, table := range tables {
		t.Run(table.label, func(t *testing.T) {
			hcYaml := testutil.CleanTestYaml(table.str)
			hc, err := yaml.ParseYAMLString[config.Mount](hcYaml)

			if table.invalid {
				require.Error(t, err)
				require.ErrorIs(t, err, config.ErrValidation)
				if table.errorContains != "" {
					require.ErrorContains(t, err, table.errorContains)
				}
				return
			}

			require.NoError(t, err)

			require.NoError(t, sc.CheckYAML(t, hcYaml))

			require.Equal(t, table.readonly, hc.ReadOnly, "ReadOnly")
			require.Equal(t, table.path, hc.ContainerPath, "ContainerPath")
			require.Equal(t, table.source, hc.SourceVolume, "SourceVolume")

			ecsMP := hc.ToAws()
			require.Equal(t, table.readonly, *ecsMP.ReadOnly, "AWS ReadOnly")
			require.Equal(t, table.path, *ecsMP.ContainerPath, "AWS ContainerPath")
			require.Equal(t, table.source, *ecsMP.SourceVolume, "AWS SourceVolume")

		})
	}
}
