package config_test

import (
	"testing"
	"time"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func intPtr(v int) *int { return &v }

func TestSettings(t *testing.T) {
	t.Run("ApplyDefaults", func(t *testing.T) {
		obj := &config.Settings{}
		obj.ApplyDefaults()

		require.Equal(t, 90*time.Minute, obj.PreDeployTimeout.ToDuration(), "PreDeployTimeout")
		require.NotNil(t, obj.KeepInSync, "KeepInSync")
		require.NotNil(t, obj.WaitForStable, "WaitForStable")
		require.NotNil(t, obj.SSMImport, "SSMImport")
		require.NotNil(t, obj.Concurrency, "Concurrency")
		require.Equal(t, 2, *obj.Concurrency, "Concurrency default")
	})

	t.Run("Validate", func(t *testing.T) {
		obj := &config.Settings{
			DisableMarkerTag: true,
			KeepInSync:       new(config.NewKeepInSyncFromBool(true)),
		}
		obj.ApplyDefaults()
		err := obj.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, config.ErrValidation)
		require.ErrorContains(t, err, "If you disable the marker")
	})

	t.Run("Concurrency", func(t *testing.T) {
		t.Run("custom value", func(t *testing.T) {
			obj := &config.Settings{Concurrency: intPtr(5)}
			obj.ApplyDefaults()
			require.NoError(t, obj.Validate())
			require.Equal(t, 5, *obj.Concurrency)
		})

		t.Run("too low", func(t *testing.T) {
			obj := &config.Settings{Concurrency: intPtr(0)}
			obj.ApplyDefaults()
			err := obj.Validate()
			require.Error(t, err)
			require.ErrorIs(t, err, config.ErrValidation)
			require.ErrorContains(t, err, "concurrency must be between")
		})

		t.Run("too high", func(t *testing.T) {
			obj := &config.Settings{Concurrency: intPtr(11)}
			obj.ApplyDefaults()
			err := obj.Validate()
			require.Error(t, err)
			require.ErrorIs(t, err, config.ErrValidation)
			require.ErrorContains(t, err, "concurrency must be between")
		})

		t.Run("minimum valid", func(t *testing.T) {
			obj := &config.Settings{Concurrency: intPtr(1)}
			obj.ApplyDefaults()
			require.NoError(t, obj.Validate())
		})

		t.Run("maximum valid", func(t *testing.T) {
			obj := &config.Settings{Concurrency: intPtr(10)}
			obj.ApplyDefaults()
			require.NoError(t, obj.Validate())
		})
	})
}
