package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck_Validate(t *testing.T) {
	tables := []struct {
		obj   config.HealthCheck
		valid bool
	}{
		{config.HealthCheck{Command: []string{"CMD", "test"}}, true},
		{config.HealthCheck{Command: []string{"CMD", "test"}, Retries: util.Ptr[int32](5)}, true},

		{config.HealthCheck{}, false},
		{config.HealthCheck{Command: []string{"test"}, Retries: util.Ptr[int32](5)}, false},
		{config.HealthCheck{Command: []string{"CMD", "test"}, Retries: util.Ptr[int32](-1)}, false},
	}

	for i, table := range tables {
		table.obj.ApplyDefaults()

		err := table.obj.Validate()

		if table.valid {
			require.NoErrorf(t, err, "entry#%d", i)
		} else {
			require.Error(t, err, "entry#%d", i)
			require.ErrorIs(t, err, config.ErrValidation)
		}

	}
}

func TestHealthCheck_Unmarshal(t *testing.T) {

	//

	t.Run("shorthand false", func(t *testing.T) {
		hc, err := yaml.ParseYAMLString[config.HealthCheck](`false`)
		require.NoError(t, err)
		require.True(t, hc.Disabled)
		require.NoError(t, hc.Validate())
	})

	t.Run("shorthand true", func(t *testing.T) {
		_, err := yaml.ParseYAMLString[config.HealthCheck](`true`)
		require.Error(t, err)
		require.ErrorIs(t, err, config.ErrValidation)
	})

	t.Run("normal", func(t *testing.T) {
		sc := testutil.NewSchemaChecker(&config.HealthCheck{})

		tables := []struct {
			label         string
			str           string
			invalid       bool
			errorContains string
			disabled      bool
			cmdParse      []string
			retries       int
			interval      int
			startp        int
			timeout       int
		}{
			{
				label:    "explicit disable",
				str:      `disabled: true`,
				disabled: true,
			},
			{
				label: "everything",
				str: `
				command: CMD test
				retries: 5
				interval: 10
				timeout: 8
				start_period: 30`,
				cmdParse: []string{"CMD", "test"},
				retries:  5,
				interval: 10,
				startp:   30,
				timeout:  8,
			},

			{
				label:         "bad command",
				str:           `command: test`,
				invalid:       true,
				errorContains: "command MUST start",
			},

			{
				label:         "missing command",
				str:           `retries: 5`,
				invalid:       true,
				errorContains: "command cannot be empty",
			},

			{
				label:    "only command",
				str:      `command: CMD test`,
				cmdParse: []string{"CMD", "test"},
			},
		}

		for _, table := range tables {
			t.Run(table.label, func(t *testing.T) {
				hcYaml := testutil.CleanTestYaml(table.str)
				hc, err := yaml.ParseYAMLString[config.HealthCheck](hcYaml)

				if table.invalid {
					require.Error(t, err)
					require.ErrorIs(t, err, config.ErrValidation)
					if table.errorContains != "" {
						require.ErrorContains(t, err, table.errorContains)
					}
					return
				}

				require.NoError(t, err)

				require.NoError(t, sc.CheckYAML(t, hcYaml))

				if table.disabled {
					require.True(t, hc.Disabled)
					return
				}

				if table.cmdParse != nil {
					require.EqualValues(t, table.cmdParse, hc.Command, "Command")
				}

				if table.retries == 0 {
					require.Nil(t, hc.Retries, "RetriesNil")
				} else {
					require.EqualValues(t, table.retries, *hc.Retries, "Retries")
				}

				if table.interval == 0 {
					require.Nil(t, hc.Interval, "IntervalNil")
				} else {
					require.EqualValues(t, table.interval, hc.Interval.ToAwsInt32(), "Interval")
				}

				if table.startp == 0 {
					require.Nil(t, hc.StartPeriod, "StartPeriodNil")
				} else {
					require.EqualValues(t, table.startp, hc.StartPeriod.ToAwsInt32(), "StartPeriod")
				}

				if table.timeout == 0 {
					require.Nil(t, hc.Timeout, "TimeoutNil")
				} else {
					require.EqualValues(t, table.timeout, hc.Timeout.ToAwsInt32(), "Timeout")
				}

			})
		}

	})
}
