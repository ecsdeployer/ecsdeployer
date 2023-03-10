package config_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestPreDeployTask(t *testing.T) {
	require.Implements(t, (*config.IsTaskStruct)(nil), &config.PreDeployTask{})
	require.NotPanics(t, func() {
		(&config.PreDeployTask{}).ApplyDefaults()
	})
}

func TestPreDeployTask_Validate(t *testing.T) {
	require.ErrorContains(t, (&config.PreDeployTask{}).Validate(), "need to name")
}

func TestPreDeployTask_Parsing(t *testing.T) {
	sc := testutil.NewSchemaChecker(&config.PreDeployTask{})
	tables := []struct {
		str     string
		valid   bool
		checker func(*testing.T, *config.PreDeployTask)
	}{
		{
			str:   "name: somepd\ncommand: bundle exec ruby stuff.rb",
			valid: true,
			checker: func(t *testing.T, pdt *config.PreDeployTask) {
				require.Equal(t, "somepd", pdt.Name)
				require.ElementsMatch(t, config.ShellCommand{"bundle", "exec", "ruby", "stuff.rb"}, *pdt.Command)
			},
		},

		{
			str:   "name: somepd",
			valid: true,
			checker: func(t *testing.T, pdt *config.PreDeployTask) {
				require.Equal(t, "somepd", pdt.Name)
				require.Empty(t, pdt.Command)
			},
		},

		{
			str:     "command: sleep 123",
			valid:   false,
			checker: nil,
		},
	}

	for testNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", testNum+1), func(t *testing.T) {

			obj, err := yaml.ParseYAMLString[config.PreDeployTask](table.str)
			if !table.valid {
				require.True(t, err != nil || sc.CheckYAML(t, table.str) != nil)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, obj)

			require.NoError(t, obj.Validate())

			if table.checker != nil {
				table.checker(t, obj)
			}

			require.NoError(t, sc.CheckYAML(t, table.str))
		})
	}
}
