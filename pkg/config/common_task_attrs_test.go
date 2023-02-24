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
		tables := []struct {
			str         string
			expectedErr string
		}{
			{
				str: `name: testing`,
			},
			{
				str: `
				name: thing2
				arch: fake`,
				expectedErr: config.ErrInvalidArchitecture.Reason,
			},
		}

		for testNum, table := range tables {
			t.Run(fmt.Sprintf("test_%02d", testNum+1), func(t *testing.T) {
				cleanStr := testutil.CleanTestYaml(table.str)
				_, err := yaml.ParseYAMLString[config.CommonTaskAttrs](cleanStr)

				if table.expectedErr == "" {
					require.NoError(t, err)
					return
				}
				require.Error(t, err)
				require.ErrorContains(t, err, table.expectedErr)
				require.ErrorIs(t, err, config.ErrValidation)
			})
		}
	})
}
