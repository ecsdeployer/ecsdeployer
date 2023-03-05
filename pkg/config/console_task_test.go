package config_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func TestConsoleTask(t *testing.T) {
	t.Run("IsEnabled", func(t *testing.T) {
		require.False(t, (&config.ConsoleTask{}).IsEnabled())
	})

	t.Run("IsTaskStruct", func(t *testing.T) {
		require.Implements(t, (*config.IsTaskStruct)(nil), &config.ConsoleTask{})
	})

	t.Run("ApplyDefaults", func(t *testing.T) {
		obj := &config.ConsoleTask{}
		obj.ApplyDefaults()

		require.False(t, *obj.Enabled)
		require.Equal(t, config.ConsoleTaskContainerName, obj.Name)
		require.EqualValues(t, 8722, *obj.PortMapping.Port)
		require.EqualValues(t, ecsTypes.TransportProtocolTcp, obj.PortMapping.Protocol)
	})

	t.Run("Validate", func(t *testing.T) {
		tables := []struct {
			obj      *config.ConsoleTask
			valid    bool
			errMatch string
		}{
			{&config.ConsoleTask{}, false, "must provide name"},
			{
				obj: &config.ConsoleTask{
					CommonTaskAttrs: config.CommonTaskAttrs{
						CommonContainerAttrs: config.CommonContainerAttrs{
							Name: "test",
						},
					},
				},
				valid:    false,
				errMatch: "must provide port",
			},
			{
				obj: &config.ConsoleTask{
					CommonTaskAttrs: config.CommonTaskAttrs{
						CommonContainerAttrs: config.CommonContainerAttrs{
							Name: "test",
						},
					},
					PortMapping: util.FirstParam(config.NewPortMappingFromString("8080/tcp")),
				},
				valid: true,
			},
		}

		for tNum, table := range tables {
			t.Run(fmt.Sprintf("test_%02d", tNum+1), func(t *testing.T) {
				err := table.obj.Validate()

				if table.valid {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
					if table.errMatch != "" {
						require.ErrorContains(t, err, table.errMatch)
					}
				}
			})
		}
	})
}

func TestConsoleTask_Unmarshal(t *testing.T) {

	sc := testutil.NewSchemaChecker(&config.ConsoleTask{})

	tables := []struct {
		str     string
		enabled bool
	}{
		{"true", true},
		{"false", false},
		{"enabled: true", true},
		{"enabled: false", false},
		{"port: 1234", false},
	}

	for _, table := range tables {
		t.Run(table.str, func(t *testing.T) {

			con, err := yaml.ParseYAMLString[config.ConsoleTask](table.str)

			sc.CheckYAML(t, table.str)

			require.NoError(t, err)
			require.Equal(t, table.enabled, con.IsEnabled())

			res, err := util.Jsonify(con)
			require.NoError(t, err)

			if table.enabled {
				require.Contains(t, res, `"enabled":true`)
				require.Contains(t, res, `"name":"console"`)
			} else {
				// require.Equal(t, "false", res)
				require.JSONEq(t, "false", res)
			}
		})
	}
}
