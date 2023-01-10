package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestAwsLogConfig(t *testing.T) {
	t.Run("IsDisabled", func(t *testing.T) {
		tables := []struct {
			expected bool
			obj      *config.AwsLogConfig
		}{
			{false, &config.AwsLogConfig{}},
			{false, &config.AwsLogConfig{Disabled: false}},
			{true, &config.AwsLogConfig{Disabled: true}},
		}

		for _, table := range tables {
			require.Equal(t, table.expected, table.obj.IsDisabled())
		}
	})

	t.Run("ApplyDefaults", func(t *testing.T) {
		obj := &config.AwsLogConfig{}
		obj.ApplyDefaults()

		require.False(t, obj.IsDisabled())
		require.EqualValues(t, 180, obj.Retention.Days())
	})

}

func TestAwsLogConfig_Unmarshal(t *testing.T) {

	def := &config.AwsLogConfig{}
	def.ApplyDefaults()

	defaultDisabled := def.IsDisabled()
	defaultRetention := def.Retention.Days()

	tables := []struct {
		str string

		disabled bool
		retDays  int32
		checker  func(*testing.T, *config.AwsLogConfig)
	}{
		{"true", false, defaultRetention, nil},
		{"false", true, defaultRetention, nil},
		{"retention: forever", defaultDisabled, -1, nil},
		{"retention: 1", defaultDisabled, 1, nil},
		{"disabled: true\nretention: 14", true, 14, nil},
		{"disabled: false\nretention: 14", false, 14, nil},
		{"disabled: false", false, defaultRetention, nil},
		{"disabled: true", true, defaultRetention, nil},

		{
			str:      "retention: 14\noptions:\n  SomeVal: testing",
			disabled: false,
			retDays:  14,
			checker: func(t *testing.T, alc *config.AwsLogConfig) {
				require.NotNil(t, alc.Options)
				require.Len(t, alc.Options, 1)
				require.Contains(t, alc.Options, "SomeVal")
			},
		},

		{
			str:      "retention: 14\noptions:\n  SomeVal: testing\n  Something: {ssm: yar}",
			disabled: false,
			retDays:  14,
			checker: func(t *testing.T, alc *config.AwsLogConfig) {
				require.NotNil(t, alc.Options)
				require.Len(t, alc.Options, 2)
				require.Contains(t, alc.Options, "SomeVal")
				require.Contains(t, alc.Options, "Something")

				val := alc.Options["Something"]
				require.True(t, val.IsSSM())
			},
		},
	}

	sc := testutil.NewSchemaChecker(&config.AwsLogConfig{})

	for _, table := range tables {
		obj, err := yaml.ParseYAMLString[config.AwsLogConfig](table.str)
		require.NoError(t, err)

		require.NoError(t, sc.CheckYAML(t, table.str))

		require.Equal(t, table.disabled, obj.IsDisabled())
		require.Equal(t, table.retDays, obj.Retention.Days())

		if table.checker != nil {
			table.checker(t, obj)
		}
	}
}
