package preflight

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	hcVersion "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestVersionCheck(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		require.Equal(t, "version requirements", checkVersion{}.String())
	})
}

func TestIsVersionAllowed(t *testing.T) {

	fakedVersion := hcVersion.Must(hcVersion.NewSemver("1.2.3"))

	tables := []struct {
		str    string
		passes bool
	}{
		{`1.2.3`, true},
		{`~> 1.2`, true},

		{`1.2.4`, false},
	}

	for _, table := range tables {
		t.Run(table.str, func(t *testing.T) {

			obj, err := yaml.ParseYAMLString[config.VersionConstraint](table.str)
			require.NoError(t, err, "If you want to test parse errors, do it in config package")

			require.Equal(t, table.passes, isVersionAllowed(obj, fakedVersion))
		})
	}
}
