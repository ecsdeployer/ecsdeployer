package config_test

import (
	"fmt"
	"testing"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestCronJob(t *testing.T) {
	t.Run("IsTaskStruct", func(t *testing.T) {
		require.Implements(t, (*config.IsTaskStruct)(nil), &config.CronJob{})
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
			// {&config.CronJob{Schedule: "rate(1 minute)", CommonTaskAttrs: config.CommonTaskAttrs{Architecture: util.Ptr(config.Architecture("wrong"))}}, "not a valid arch"},
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

	t.Run("UnmarshalYAML", func(t *testing.T) {

		tzMinus8UTC := time.FixedZone("UTC-8", -8*3600)

		tables := []struct {
			str      string
			expStart *time.Time
			expEnd   *time.Time
			expSched string
			errMatch string
		}{
			{
				str: `
				name: test
				schedule: rate(1 hour)
				start_date: 2023-99-02`,
				errMatch: "Invalid format for start",
			},
			{
				str: `
				name: test
				schedule: rate(1 hour)
				end_date: 2023-01-02T09:12:55-08:00
				start_date: 2024-01-02T09:12:55-08:00`,
				errMatch: "end date cannot be before the start date",
			},
			{
				str: `
				name: test
				schedule: rate(1 hour)`,
				expSched: "rate(1 hour)",
			},
			{
				str: `
				name: test
				schedule: rate(1 hour)
				start_date: 2023-01-02T09:12:55Z`,
				expSched: "rate(1 hour)",
				expStart: util.Ptr(time.Date(2023, 1, 2, 9, 12, 55, 0, time.UTC)),
			},
			{
				str: `
				name: test
				schedule: rate(1 hour)
				start_date: 2023-01-02T09:12:55-08:00`,
				expSched: "rate(1 hour)",
				expStart: util.Ptr(time.Date(2023, 1, 2, 9, 12, 55, 0, tzMinus8UTC)),
			},
			{
				str: `
				name: test
				schedule: rate(1 hour)
				start_date: 2023-01-02T09:12:55-08:00
				end_date: 2024-01-02T09:12:55-08:00`,
				expSched: "rate(1 hour)",
				expStart: util.Ptr(time.Date(2023, 1, 2, 9, 12, 55, 0, tzMinus8UTC)),
				expEnd:   util.Ptr(time.Date(2024, 1, 2, 9, 12, 55, 0, tzMinus8UTC)),
			},
		}
		for i, table := range tables {
			t.Run(fmt.Sprintf("test_%02d", i+1), func(t *testing.T) {
				cleanStr := testutil.CleanTestYaml(table.str)

				obj, err := yaml.ParseYAMLString[config.CronJob](cleanStr)

				if table.errMatch != "" {
					require.Error(t, err)
					require.ErrorContains(t, err, table.errMatch)
					return
				}

				require.NoError(t, err)

				if table.expStart != nil {
					require.NotNil(t, obj.StartDate)
					require.Equal(t, "UTC", obj.StartDate.Location().String())
					require.Equal(t, table.expStart.UTC(), obj.StartDate.UTC())
				}

				if table.expEnd != nil {
					require.NotNil(t, obj.EndDate)
					require.Equal(t, "UTC", obj.EndDate.Location().String())
					require.Equal(t, table.expEnd.UTC(), obj.EndDate.UTC())
				}

				if table.expSched != "" {
					require.Equal(t, table.expSched, obj.Schedule)
				}

			})
		}
	})
}
