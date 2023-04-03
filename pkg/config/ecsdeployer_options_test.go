package config_test

import (
	"fmt"
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
			errMsg  string
		}{
			{
				str:     "required_version: 1.2.3\nallowed_account_id: 123456789000",
				expVC:   util.Must(config.NewVersionConstraint("1.2.3")),
				expAcct: "123456789000",
			},
			{
				str:   "required_version: 1.2.3",
				expVC: util.Must(config.NewVersionConstraint("1.2.3")),
			},
			{
				str:     "allowed_account_id: 123456789000",
				expVC:   nil,
				expAcct: "123456789000",
			},

			{
				str:    "allowed_account_id: xxxx",
				errMsg: "AWS AccountIDs must be",
			},
		}

		for tnum, table := range tables {
			t.Run(fmt.Sprintf("test_%02d", tnum+1), func(t *testing.T) {
				obj, err := yaml.ParseYAMLString[config.EcsDeployerOptions](table.str)

				if table.errMsg != "" {
					require.Error(t, err)
					require.ErrorContains(t, err, table.errMsg)
					return
				}

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
			})
		}
	})

}
