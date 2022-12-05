package config_test

import (
	"fmt"
	"testing"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNewDuration(t *testing.T) {
	t.Run("NewDurationFromString", func(t *testing.T) {
		tables := []struct {
			str     string
			seconds int32
		}{
			{"1h", 3600},
			{"2v", -1},
		}

		for _, table := range tables {
			t.Run(table.str, func(t *testing.T) {
				dur, err := config.NewDurationFromString(table.str)
				if table.seconds == -1 {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.Equal(t, table.seconds, dur.ToAwsInt32())
				require.Equal(t, float64(table.seconds), dur.ToDuration().Seconds())
			})
		}
	})

	t.Run("NewDurationFromUint", func(t *testing.T) {
		tables := []struct {
			seconds int32
		}{
			{3600},
		}

		for _, table := range tables {
			t.Run(fmt.Sprintf("%d", table.seconds), func(t *testing.T) {
				dur, _ := config.NewDurationFromUint(uint32(table.seconds))

				require.Equal(t, table.seconds, dur.ToAwsInt32())
				require.Equal(t, float64(table.seconds), dur.ToDuration().Seconds())
			})
		}
	})

	t.Run("NewDurationFromTDuration", func(t *testing.T) {
		tables := []struct {
			dur     time.Duration
			seconds int32
		}{
			{1 * time.Hour, 3600},
			{(1 * time.Hour) + (400 * time.Millisecond), 3600},
		}

		for _, table := range tables {
			t.Run(table.dur.String(), func(t *testing.T) {
				dur := config.NewDurationFromTDuration(table.dur)
				require.Equal(t, table.seconds, dur.ToAwsInt32())
				require.Equal(t, float64(table.seconds), dur.ToDuration().Seconds())
			})
		}
	})

}

func TestDuration_Unmarshal(t *testing.T) {
	tables := []struct {
		str     string
		seconds int32
	}{

		// valid
		{"0", 0},
		{"3600", 3600},
		{"1h", 3600},
		{"1h10s", 3610},
		{"2h400ms", 7200},

		// known edge case... too lazy to fix
		{"null", 0},

		// failures
		{"1d", -1},
		{"-100", -1},
		{"x", -1},
		{"true", -1},
		{`"blah"`, -1},
		{`"3600"`, -1},
		{"thing: blah\nyar: test", -1},
	}

	for _, table := range tables {
		t.Run(table.str, func(t *testing.T) {

			dur, err := yaml.ParseYAMLString[config.Duration](table.str)
			if table.seconds == -1 {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, table.seconds, dur.ToAwsInt32())
			require.Equal(t, float64(table.seconds), dur.ToDuration().Seconds())

		})
	}
}
