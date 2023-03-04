package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestProxyConfig_Unmarshal(t *testing.T) {

	//

	t.Run("shorthand false", func(t *testing.T) {
		hc, err := yaml.ParseYAMLString[config.ProxyConfig](`false`)
		require.NoError(t, err)
		require.True(t, hc.Disabled)
		require.NoError(t, hc.Validate())
	})

	t.Run("shorthand true", func(t *testing.T) {
		_, err := yaml.ParseYAMLString[config.ProxyConfig](`true`)
		require.Error(t, err)
		require.ErrorIs(t, err, config.ErrValidation)
	})

	t.Run("normal", func(t *testing.T) {
		sc := testutil.NewSchemaChecker(&config.ProxyConfig{})

		tables := []struct {
			label         string
			str           string
			invalid       bool
			errorContains string
			disabled      bool
			ptype         string
			cname         string
			props         map[string]string
		}{
			{
				label:    "explicit disable",
				str:      `disabled: true`,
				disabled: true,
			},
			{
				label: "everything",
				str: `
				type: APPMESH
				container_name: proxywoxy
				properties:
					BlahThing: 1234`,
				cname: "proxywoxy",
				props: map[string]string{
					"BlahThing": "1234",
				},
			},

			{
				label: "defaults",
				str: `
				properties:
					BlahThing: 1234
					TplThing: {template: "{{.Project}}"}`,
				props: map[string]string{
					"BlahThing": "1234",
					"TplThing":  "{{.Project}}",
				},
			},

			{
				label: "failure ssm in prop",
				str: `
				properties:
					BlahThing: 1234
					SsmThing: {ssm: /path}`,
				invalid:       true,
				errorContains: "cannot reference SSM",
			},

			{
				label: "blank container",
				str: `
				container_name: ""
				properties:
					BlahThing: 1234`,
				invalid:       true,
				errorContains: "container_name is required",
			},

			{
				label: "blank type",
				str: `
				type: ""
				properties:
					BlahThing: 1234`,
				invalid:       true,
				errorContains: "type is required",
			},
		}

		for _, table := range tables {
			t.Run(table.label, func(t *testing.T) {
				hcYaml := testutil.CleanTestYaml(table.str)
				hc, err := yaml.ParseYAMLString[config.ProxyConfig](hcYaml)

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

				if table.disabled {
					require.True(t, hc.Disabled)
					return
				}

				if table.ptype != "" {
					require.EqualValues(t, table.ptype, *hc.Type, "Type")
				}

				if table.cname == "" {
					require.Equal(t, "envoy", *hc.ContainerName, "ContainerName")
				} else {
					require.EqualValues(t, table.cname, *hc.ContainerName, "ContainerName")
				}

				if table.props != nil {
					require.NotNil(t, hc.Properties)
					for k, v := range table.props {
						propVal, _ := hc.Properties[k].GetValue(testutil.TplDummy)

						require.Equal(t, v, propVal)
					}
				}

			})
		}

	})
}
