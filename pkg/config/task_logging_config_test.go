package config_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
)

func TestLoggingDisableFlag(t *testing.T) {
	require.Equal(t, "none", config.LoggingDisableFlag)
}

func TestTaskLoggingConfig(t *testing.T) {
	t.Run("ApplyDefaults", func(t *testing.T) {
		obj := &config.TaskLoggingConfig{}
		obj.ApplyDefaults()

		require.Nil(t, obj.Driver)
		require.Empty(t, obj.Options)

	})

	t.Run("IsDisabled", func(t *testing.T) {
		tables := []struct {
			driver   string
			disabled bool
		}{
			{config.LoggingDisableFlag, true},
			{"awslogs", false},
			{"splunk", false},
			{"firelens", false},
			{"", false},
		}

		for _, table := range tables {
			obj := &config.TaskLoggingConfig{}
			if table.driver != "" {
				obj.Driver = aws.String(table.driver)
			}

			require.Equal(t, table.disabled, obj.IsDisabled())
			require.NoError(t, obj.Validate())
		}
	})

}

func TestTaskLoggingConfig_Marshalling(t *testing.T) {

	sc := testutil.NewSchemaChecker(&config.TaskLoggingConfig{})

	tables := []struct {
		str      string
		failure  bool
		disabled bool
		driver   string
		numOpts  int
	}{

		{
			str:      `false`,
			disabled: true,
		},
		{
			str:     `true`,
			failure: true,
		},

		{
			str: ``,
			// no driver, not disabled. should inherit everything
			numOpts: 0,
		},
		{
			str: `null`,
			// no driver, not disabled. should inherit everything
			numOpts: 0,
		},
		{
			str:      "driver: none",
			disabled: true,
			driver:   "none",
		},
		{
			str:      "driver: awslogs",
			disabled: false,
			driver:   "awslogs",
		},
		{
			str:      "driver: awslogs\noptions:\n  thing: blah",
			disabled: false,
			driver:   "awslogs",
			numOpts:  1,
		},
	}

	for i, table := range tables {
		t.Run(fmt.Sprintf("entry_%02d", i), func(t *testing.T) {

			obj, err := yaml.ParseYAMLString[config.TaskLoggingConfig](table.str)
			if table.failure {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, obj)

			require.Equalf(t, table.disabled, obj.IsDisabled(), "IsDisabled")

			require.NoError(t, sc.CheckYAML(t, table.str))

			if table.driver != "" {
				require.Equalf(t, table.driver, *obj.Driver, "Driver")
			}
			require.Lenf(t, obj.Options, table.numOpts, "NumOptions")

		})
	}
}
