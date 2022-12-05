package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestCommonTaskAttrs_Smoke(t *testing.T) {
	common := &config.CommonTaskAttrs{
		Architecture: util.Ptr(config.ArchitectureARM64),
		CommonContainerAttrs: config.CommonContainerAttrs{
			Name: "test",
		},
		Network: &config.NetworkConfiguration{
			Subnets: []config.NetworkFilter{
				{
					ID: util.Ptr("subnet-1111111"),
				},
			},
		},
	}

	require.NoError(t, common.Validate())

	require.Truef(t, common.IsTaskStruct(), "IsTaskStruct")

	fields := common.TemplateFields()
	require.Equalf(t, "test", fields["Name"], "Name")
	require.Equalf(t, "arm64", fields["Arch"], "Arch")
}

func TestCommonTaskAttrs_Validate(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		goodArch := config.ArchitectureAMD64
		badArch := config.Architecture("badbad")
		tables := []struct {
			obj         *config.CommonTaskAttrs
			expectedErr string
		}{
			{&config.CommonTaskAttrs{Architecture: &goodArch}, ""},
			{&config.CommonTaskAttrs{Architecture: &badArch}, "not a valid arch"},
		}

		for _, table := range tables {

			err := table.obj.Validate()

			if table.expectedErr == "" {
				require.NoError(t, err)
				continue
			}
			require.Error(t, err)
			require.ErrorContains(t, err, table.expectedErr)
		}
	})
}
