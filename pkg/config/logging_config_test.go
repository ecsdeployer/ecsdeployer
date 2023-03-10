package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestLoggingConfig(t *testing.T) {
	t.Run("ApplyDefaults", func(t *testing.T) {
		obj := &config.LoggingConfig{}
		obj.ApplyDefaults()

		require.False(t, obj.IsDisabled())
		require.Equal(t, config.LoggingTypeAwslogs, obj.Type())
	})
}

func TestLoggingConfig_Unmarshal(t *testing.T) {

	tables := []struct {
		label string

		str string

		// will fail Validate()
		invalid                 bool
		validationErrorContains string

		// doesnt match the schema definition
		badSchema bool

		// run these checks
		checker func(*testing.T, *config.LoggingConfig)
	}{
		// {"true", false, defaultRetention, nil},
		// {"false", true, defaultRetention, nil},
		// {"retention: forever", defaultDisabled, -1, nil},
		// {"retention: 1", defaultDisabled, 1, nil},
		// {"disabled: true\nretention: 14", true, 14, nil},
		// {"disabled: false\nretention: 14", false, 14, nil},
		// {"disabled: false", false, defaultRetention, nil},
		// {"disabled: true", true, defaultRetention, nil},

		{
			label: "Shorthand enable logging",
			str:   "true",
			checker: func(t *testing.T, obj *config.LoggingConfig) {
				require.False(t, obj.IsDisabled())
				require.False(t, obj.AwsLogConfig.IsDisabled())
				require.Equal(t, config.LoggingTypeAwslogs, obj.Type())
			},
		},

		{
			label: "Shorthand disable logging",
			str:   "false",
			checker: func(t *testing.T, obj *config.LoggingConfig) {
				require.True(t, obj.IsDisabled())
				require.True(t, obj.FirelensConfig.IsDisabled())
				require.Equal(t, config.LoggingTypeDisabled, obj.Type())
			},
		},

		{
			label: "Enable and set retention",
			str:   "awslogs:\n  retention: 14",
			checker: func(t *testing.T, obj *config.LoggingConfig) {
				require.False(t, obj.IsDisabled())
				require.Equal(t, config.LoggingTypeAwslogs, obj.Type())
				require.EqualValues(t, 14, obj.AwsLogConfig.Retention.Days())
			},
		},

		{
			label: "Using firelens",
			str:   "firelens: true",
			checker: func(t *testing.T, obj *config.LoggingConfig) {
				require.False(t, obj.IsDisabled())
				require.False(t, obj.FirelensConfig.IsDisabled())
				require.Equal(t, config.LoggingTypeFirelens, obj.Type())
			},
		},

		{
			label: "Using custom",
			str: `
			custom:
				driver: splunk
			`,
			checker: func(t *testing.T, obj *config.LoggingConfig) {
				require.False(t, obj.IsDisabled())
				require.False(t, obj.Custom.IsDisabled())
				require.Equal(t, config.LoggingTypeCustom, obj.Type())
			},
		},

		{
			label: "Using firelens with ssm router opts",
			str: `
			firelens:
				router_options:
					Thing:
						ssm: someval
			`,
			invalid:                 true,
			validationErrorContains: "you cannot have SSM options",
		},

		{
			label: "Using firelens with ssm opts",
			str: `
			firelens:
			  options:
			    Thing:
			      ssm: someval`,
		},

		// {
		// 	label: "Disable all parts but not parent",
		// 	str: `
		// 	awslogs:
		// 	  disabled: true
		// 	firelens: false`,
		// 	invalid:                 true,
		// 	validationErrorContains: "if you want to disable logging",
		// },
	}

	sc := testutil.NewSchemaChecker(&config.LoggingConfig{})

	for _, table := range tables {
		t.Run(table.label, func(t *testing.T) {

			yamlStr := testutil.CleanTestYaml(table.str)

			schemaErr := sc.CheckYAML(t, yamlStr)
			if table.badSchema {
				require.Error(t, schemaErr)
				return
			}

			require.NoError(t, schemaErr)

			obj, err := yaml.ParseYAMLString[config.LoggingConfig](yamlStr)

			if table.invalid {
				require.Error(t, err)
				require.ErrorIs(t, err, config.ErrValidation)

				if table.validationErrorContains != "" {
					require.ErrorContains(t, err, table.validationErrorContains)
				}
				return
			}

			require.NoError(t, err)

			if table.checker != nil {
				table.checker(t, obj)
			}
		})
	}
}
