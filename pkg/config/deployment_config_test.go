package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestDeploymentConfig_ValidateDesiredCountSuccess(t *testing.T) {

	tables := []struct {
		min     int32
		max     int32
		desired int32
	}{
		{0, 100, 1},
		{0, 150, 1},
		{0, 150, 2},
	}

	for _, table := range tables {
		dc := &config.RolloutConfig{
			Minimum: &table.min,
			Maximum: &table.max,
		}

		err := dc.ValidateWithDesiredCount(table.desired)

		require.NoError(t, err)
	}
}

func TestDeploymentConfig_ValidateDesiredCountFailures(t *testing.T) {

	tables := []struct {
		min     int32
		max     int32
		desired int32
	}{
		{100, 150, 1},
		{100, 125, 1},
	}

	for _, table := range tables {
		dc := &config.RolloutConfig{
			Minimum: &table.min,
			Maximum: &table.max,
		}

		err := dc.ValidateWithDesiredCount(table.desired)
		require.Error(t, err)
	}
}
