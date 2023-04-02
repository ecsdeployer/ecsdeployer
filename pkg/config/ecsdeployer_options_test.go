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
			str     string
			expVC   *config.VersionConstraint
			expAcct string
		}{
			{"required_version: 1.2.3\nallowed_account_id: 123456789000", util.Must(config.NewVersionConstraint("1.2.3")), "123456789000"},
			{"required_version: 1.2.3", util.Must(config.NewVersionConstraint("1.2.3")), ""},
			{"allowed_account_id: 123456789000", nil, "123456789000"},
		}

		for _, table := range tables {

			obj, err := yaml.ParseYAMLString[config.EcsDeployerOptions](table.str)

			require.NoError(t, err)
			require.NotNil(t, obj)

			if table.expVC == nil {
				require.Nil(t, obj.RequiredVersion)
			} else {
				require.NotNil(t, obj.RequiredVersion)
				require.Equal(t, table.expVC.String(), obj.RequiredVersion.String())
			}

			if table.expAcct == "" {
				require.Nil(t, obj.AllowedAccountId)
			} else {
				require.NotNil(t, obj.AllowedAccountId)
				require.Equal(t, table.expAcct, *obj.AllowedAccountId)
			}

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

	// t.Run("IsVersionAllowed", func(t *testing.T) {
	// 	ecsdBlank := &config.EcsDeployerOptions{}

	// 	ecsdRange := &config.EcsDeployerOptions{
	// 		RequiredVersion: util.Must(config.NewVersionConstraint("~> 1.2.3")),
	// 	}

	// 	ecsdGT := &config.EcsDeployerOptions{
	// 		RequiredVersion: util.Must(config.NewVersionConstraint(">= 1.2.3")),
	// 	}

	// 	tables := []struct {
	// 		ver        string
	// 		allowRange bool
	// 		allowGT    bool
	// 	}{
	// 		{"1.2.3", true, true},
	// 		{version.DevVersionID, false, true},
	// 	}

	// 	for _, table := range tables {
	// 		blankAllowed, _ := ecsdBlank.IsVersionAllowed(table.ver)
	// 		require.Truef(t, blankAllowed, "Blank")

	// 		rngAllowed, _ := ecsdRange.IsVersionAllowed(table.ver)
	// 		require.Equalf(t, table.allowRange, rngAllowed, "Ranged")

	// 		gtAllowed, _ := ecsdGT.IsVersionAllowed(table.ver)
	// 		require.Equalf(t, table.allowGT, gtAllowed, "GreaterThan")
	// 	}
	// })
}
