package config_test

import (
	"testing"

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
		}
	})

}
