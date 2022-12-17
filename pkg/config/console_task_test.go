package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

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
	}

	for _, table := range tables {
		con, err := yaml.ParseYAMLString[config.ConsoleTask](table.str)

		sc.CheckYAML(t, table.str)

		require.NoErrorf(t, err, "unexpected error for <%s> %s", table.str, err)

		if table.enabled != con.IsEnabled() {
			t.Errorf("expected <%s> to %v console", table.str, table.enabled)
		}
	}
}
