package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestConsoleTask_Unmarshal(t *testing.T) {

	type conDummy struct {
		Console *config.ConsoleTask `yaml:"console,omitempty" json:"console,omitempty"`
	}

	tables := []struct {
		str     string
		enabled bool
	}{
		{"console: true", true},
		{"console: false", false},
		{"console:\n  enabled: true", true},
		{"console:\n  enabled: false", false},
	}

	for _, table := range tables {
		con, err := yaml.ParseYAMLString[conDummy](table.str)

		require.NoErrorf(t, err, "unexpected error for <%s> %s", table.str, err)

		if table.enabled != con.Console.IsEnabled() {
			t.Errorf("expected <%s> to %v console", table.str, table.enabled)
		}
	}
}
