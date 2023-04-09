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

func TestCommonTaskAttrs(t *testing.T) {
	t.Run("InterfaceUsage", func(t *testing.T) {

		type hasCommonTaskAttrs interface {
			GetCommonTaskAttrs() config.CommonTaskAttrs
		}

		require.Implements(t, (*hasCommonTaskAttrs)(nil), &config.Service{}, "Service")
		require.Implements(t, (*hasCommonTaskAttrs)(nil), &config.ConsoleTask{}, "ConsoleTask")
		require.Implements(t, (*hasCommonTaskAttrs)(nil), &config.PreDeployTask{}, "PreDeployTask")
		require.Implements(t, (*hasCommonTaskAttrs)(nil), &config.CronJob{}, "CronJob")
		require.Implements(t, (*hasCommonTaskAttrs)(nil), &config.CommonTaskAttrs{}, "CommonTaskAttrs")
		require.Implements(t, (*hasCommonTaskAttrs)(nil), &config.FargateDefaults{}, "FargateDefaults")
	})

	t.Run("CanOverride", func(t *testing.T) {
		require.True(t, (&config.CommonTaskAttrs{}).CanOverride())
		strVal := "x"
		require.False(t, (&config.CommonTaskAttrs{Architecture: util.Ptr(config.ArchitectureARM64)}).CanOverride())
		require.False(t, (&config.CommonTaskAttrs{PlatformVersion: &strVal}).CanOverride())
		require.False(t, (&config.CommonTaskAttrs{ProxyConfig: &config.ProxyConfig{}}).CanOverride())
		require.False(t, (&config.CommonTaskAttrs{Network: &config.NetworkConfiguration{}}).CanOverride())
		require.False(t, (&config.CommonTaskAttrs{Sidecars: []*config.Sidecar{{}}}).CanOverride())
		require.False(t, (&config.CommonTaskAttrs{Tags: []config.NameValuePair{{}}}).CanOverride())
		require.False(t, (&config.CommonTaskAttrs{Volumes: config.VolumeList{"x": config.Volume{}}}).CanOverride())

		require.False(t, (&config.CommonTaskAttrs{CommonContainerAttrs: config.CommonContainerAttrs{DependsOn: []config.DependsOn{{}}}}).CanOverride())
	})
}

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
	require.Implements(t, (*config.IsTaskStruct)(nil), common)

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
