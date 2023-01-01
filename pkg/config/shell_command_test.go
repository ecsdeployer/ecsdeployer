package config_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestShellCommand_Parsing(t *testing.T) {
	sc := testutil.NewSchemaChecker(&config.ShellCommand{})
	tables := []struct {
		str      string
		expected config.ShellCommand
	}{
		{`"test blah"`, config.ShellCommand{"test", "blah"}},
		{`""`, config.ShellCommand{""}},
		{`["test", "blah"]`, config.ShellCommand{"test", "blah"}},
		{`["test", true]`, config.ShellCommand{"test", "true"}},
		{`["test", 123]`, config.ShellCommand{"test", "123"}},
		{`["test", ""]`, config.ShellCommand{"test", ""}},
		{`"test -c 'something something'"`, config.ShellCommand{"test", "-c", "something something"}},
		{"- test\n- blah", config.ShellCommand{"test", "blah"}},
		{"- test\n- 1234", config.ShellCommand{"test", "1234"}},

		{"foo: bar", nil},
	}

	for testNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", testNum+1), func(t *testing.T) {

			obj, err := yaml.ParseYAMLString[config.ShellCommand](table.str)
			if table.expected == nil {
				require.Error(t, err)
				require.Error(t, sc.CheckYAML(t, table.str))
				return
			}

			require.NoError(t, err)
			require.NotNil(t, obj)

			require.Equal(t, table.expected.String(), obj.String())

			require.NoError(t, sc.CheckYAML(t, table.str))
		})
	}
}
