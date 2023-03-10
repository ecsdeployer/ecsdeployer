package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestCommonContainerAttrs(t *testing.T) {
	t.Run("InterfaceUsage", func(t *testing.T) {

		type hasCommonContainerAttrs interface {
			GetCommonContainerAttrs() config.CommonContainerAttrs
		}

		require.Implements(t, (*hasCommonContainerAttrs)(nil), &config.Service{}, "Service")
		require.Implements(t, (*hasCommonContainerAttrs)(nil), &config.ConsoleTask{}, "ConsoleTask")
		require.Implements(t, (*hasCommonContainerAttrs)(nil), &config.PreDeployTask{}, "PreDeployTask")
		require.Implements(t, (*hasCommonContainerAttrs)(nil), &config.CronJob{}, "CronJob")
		require.Implements(t, (*hasCommonContainerAttrs)(nil), &config.CommonTaskAttrs{}, "CommonTaskAttrs")
		require.Implements(t, (*hasCommonContainerAttrs)(nil), &config.CommonContainerAttrs{}, "CommonContainerAttrs")
		require.Implements(t, (*hasCommonContainerAttrs)(nil), &config.FargateDefaults{}, "FargateDefaults")
		require.Implements(t, (*hasCommonContainerAttrs)(nil), &config.Sidecar{}, "Sidecar")
	})
}
