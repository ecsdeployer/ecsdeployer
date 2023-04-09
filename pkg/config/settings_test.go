package config_test

import (
	"testing"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestSettings(t *testing.T) {
	t.Run("ApplyDefaults", func(t *testing.T) {
		obj := &config.Settings{}
		obj.ApplyDefaults()

		require.Equal(t, 90*time.Minute, obj.PreDeployTimeout.ToDuration(), "PreDeployTimeout")
		require.NotNil(t, obj.KeepInSync, "KeepInSync")
		require.NotNil(t, obj.WaitForStable, "WaitForStable")
		require.NotNil(t, obj.SSMImport, "SSMImport")
	})

	t.Run("Validate", func(t *testing.T) {
		obj := &config.Settings{
			DisableMarkerTag: true,
			KeepInSync:       util.Ptr(config.NewKeepInSyncFromBool(true)),
		}
		obj.ApplyDefaults()
		err := obj.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, config.ErrValidation)
		require.ErrorContains(t, err, "If you disable the marker")
	})
}
