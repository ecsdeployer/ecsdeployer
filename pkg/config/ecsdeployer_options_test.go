package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestEcsDeployerOptions(t *testing.T) {
	t.Run("UnmarshalYAML", func(t *testing.T) {
		sc := testutil.NewSchemaChecker(&config.EcsDeployerOptions{})

		tables := []struct {
			valid bool
			str   string
			exp   *config.VersionConstraint
		}{
			{true, `1.2.3`, util.Must(config.NewVersionConstraint("1.2.3"))},
			{true, `v1.2.3`, util.Must(config.NewVersionConstraint("v1.2.3"))},
			{true, `<=v1.2.3`, util.Must(config.NewVersionConstraint("<= v1.2.3"))},

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

			require.Equal(t, table.exp.String(), obj.String())

			require.NoError(t, sc.CheckYAML(t, table.str))
		}
	})

	t.Run("IsAllowedAccountId", func(t *testing.T) {

		ecsdBlank := &config.EcsDeployerOptions{}

		ecsd := &config.EcsDeployerOptions{
			AllowedAccountId: util.Ptr("111111111111"),
		}
		tables := []struct {
			allowed bool
			acctid  string
		}{
			{true, "111111111111"},

			{false, "222222222222"},
			{false, "11111111111"},
			{false, "333333333333"},
		}

		for _, table := range tables {
			require.True(t, ecsdBlank.IsAllowedAccountId(table.acctid))

			require.Equal(t, table.allowed, ecsd.IsAllowedAccountId(table.acctid))
		}

	})

	t.Run("IsVersionAllowed", func(t *testing.T) {
		ecsdBlank := &config.EcsDeployerOptions{}

		ecsdRange := &config.EcsDeployerOptions{
			RequiredVersion: util.Must(config.NewVersionConstraint("~> 1.2.3")),
		}

		ecsdGT := &config.EcsDeployerOptions{
			RequiredVersion: util.Must(config.NewVersionConstraint(">= 1.2.3")),
		}

		tables := []struct {
			ver        string
			allowRange bool
			allowGT    bool
		}{
			{"1.2.3", true, true},
			{"development", false, true},
		}

		for _, table := range tables {
			blankAllowed, _ := ecsdBlank.IsVersionAllowed(table.ver)
			require.Truef(t, blankAllowed, "Blank")

			rngAllowed, _ := ecsdRange.IsVersionAllowed(table.ver)
			require.Equalf(t, table.allowRange, rngAllowed, "Ranged")

			gtAllowed, _ := ecsdGT.IsVersionAllowed(table.ver)
			require.Equalf(t, table.allowGT, gtAllowed, "GreaterThan")
		}
	})
}
