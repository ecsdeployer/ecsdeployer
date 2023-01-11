package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestCronJob(t *testing.T) {
	t.Run("IsTaskStruct", func(t *testing.T) {
		require.True(t, (&config.CronJob{}).IsTaskStruct())
	})

	t.Run("ApplyDefaults", func(t *testing.T) {
		require.NotPanics(t, func() {
			(&config.CronJob{}).ApplyDefaults()
		})
	})

	t.Run("IsDisabled", func(t *testing.T) {
		tables := []struct {
			expected bool
			obj      *config.CronJob
		}{
			{true, &config.CronJob{Disabled: true}},

			{false, &config.CronJob{Disabled: false}},
			{false, &config.CronJob{}},
		}
		for _, table := range tables {
			require.Equal(t, table.expected, table.obj.IsDisabled())
		}
	})

	t.Run("Validate", func(t *testing.T) {
		tables := []struct {
			obj    *config.CronJob
			errStr string
		}{
			// requires schedule
			{&config.CronJob{}, "must provide a cron schedule"},
			{&config.CronJob{Schedule: "rate(1 minute)"}, ""},

			// also validates the common stuff
			{&config.CronJob{Schedule: "rate(1 minute)", CommonTaskAttrs: config.CommonTaskAttrs{Architecture: util.Ptr(config.Architecture("wrong"))}}, "not a valid arch"},
		}
		for _, table := range tables {
			err := table.obj.Validate()
			if table.errStr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, table.errStr)
				require.ErrorIs(t, err, config.ErrValidation)
			} else {
				require.NoError(t, err)
			}
		}
	})
}
