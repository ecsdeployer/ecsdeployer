package config_test

import (
	"testing"
	"time"

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

	t.Run("CanOverride", func(t *testing.T) {
		require.True(t, (&config.CommonContainerAttrs{}).CanOverride())
		strVal := "x"
		durVal := config.NewDurationFromTDuration(2 * time.Millisecond)
		require.False(t, (&config.CommonContainerAttrs{User: &strVal}).CanOverride())
		require.False(t, (&config.CommonContainerAttrs{Workdir: &strVal}).CanOverride())
		require.False(t, (&config.CommonContainerAttrs{EntryPoint: &config.ShellCommand{"X"}}).CanOverride())
		require.False(t, (&config.CommonContainerAttrs{Credentials: &strVal}).CanOverride())
		require.False(t, (&config.CommonContainerAttrs{StartTimeout: &durVal}).CanOverride())
		require.False(t, (&config.CommonContainerAttrs{StopTimeout: &durVal}).CanOverride())
		require.False(t, (&config.CommonContainerAttrs{MountPoints: []config.Mount{{}}}).CanOverride())
		require.False(t, (&config.CommonContainerAttrs{Ulimits: []config.Ulimit{{}}}).CanOverride())
		require.False(t, (&config.CommonContainerAttrs{VolumesFrom: []config.VolumeFrom{{}}}).CanOverride())
		require.False(t, (&config.CommonContainerAttrs{DependsOn: []config.DependsOn{{}}}).CanOverride())
		require.False(t, (&config.CommonContainerAttrs{DockerLabels: []config.NameValuePair{{}}}).CanOverride())
		require.False(t, (&config.CommonContainerAttrs{HealthCheck: &config.HealthCheck{}}).CanOverride())
		require.False(t, (&config.CommonContainerAttrs{LoggingConfig: &config.TaskLoggingConfig{}}).CanOverride())
	})
}
