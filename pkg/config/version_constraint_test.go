package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	hcVersion "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestVersionConstraint(t *testing.T) {

	t.Run("ValidateAndCheck", func(t *testing.T) {

		curECSDVersion := hcVersion.Must(hcVersion.NewVersion("1.2.3"))

		tables := []struct {
			passes bool
			str    string
		}{
			{true, "= 1.2.3"},
			{true, "~> 1.2"},
			{true, "< 2"},
			{true, ">= 1"},

			{false, "> 1.2.3"},
		}

		for _, table := range tables {
			t.Run(table.str, func(t *testing.T) {

				vc := util.Must(config.NewVersionConstraint(table.str))

				require.Equalf(t, table.passes, vc.Check(curECSDVersion), "Check")

				require.Equal(t, table.passes, vc.Check(curECSDVersion))
			})
		}

	})

	t.Run("UnmarshalYAML", func(t *testing.T) {

		sc := testutil.NewSchemaChecker(&config.VersionConstraint{})

		tables := []struct {
			valid bool
			str   string
			exp   *config.VersionConstraint
		}{
			{true, `1.2.3`, util.Must(config.NewVersionConstraint("1.2.3"))},
			{true, `1.2.3`, util.Must(config.NewVersionConstraint("=1.2.3"))},
			{true, `<=1.2.3`, util.Must(config.NewVersionConstraint("<= 1.2.3"))},

			{false, `bad`, nil},
			{false, `1`, nil},
			{false, `false`, nil},
		}

		for _, table := range tables {

			obj, err := yaml.ParseYAMLString[config.VersionConstraint](table.str)

			if !table.valid {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, obj)

			require.True(t, table.exp.Constraints().Equals(obj.Constraints()))

			require.NoError(t, sc.CheckYAML(t, table.str))
		}
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		tables := []struct {
			ver *config.VersionConstraint
			exp string
		}{
			{util.Must(config.NewVersionConstraint("1.2.3")), `1.2.3`},
			{util.Must(config.NewVersionConstraint("<= 1.2.3")), `<= 1.2.3`},
		}

		for _, table := range tables {
			// marshal the expected because json adds \uXXXXX characters
			expected, err := util.Jsonify(table.exp)
			require.NoError(t, err)

			res, err := util.Jsonify(table.ver)
			require.NoError(t, err)
			require.JSONEq(t, expected, res)
		}
	})
}
