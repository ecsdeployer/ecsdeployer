package config_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestRolloutConfig(t *testing.T) {

	type minmax struct {
		count int32
		min   int32
		max   int32
	}

	tables := []struct {
		str     string
		min     int32
		max     int32
		invalid bool
		checks  []minmax
	}{
		{
			str: "min: 0\nmax: 100",
			min: 0,
			max: 100,
			checks: []minmax{
				{10, 0, 10},
				{125, 0, 125},
				{0, 0, 0},
			},
		},

		{
			str: "min: 100\nmax: 150",
			min: 100,
			max: 150,
			checks: []minmax{
				{10, 10, 15},
				{100, 100, 150},
				{0, 0, 0},
			},
		},

		{
			str: "min: 100\nmax: 115",
			min: 100,
			max: 115,
			checks: []minmax{
				{10, 10, 11},
				{0, 0, 0},
				{1, -1, -1},
			},
		},

		{
			str:     "min: 100\nmax: 100",
			invalid: true,
		},
		{
			str:     "min: 10",
			invalid: true,
		},
		{
			str:     "max: 100",
			invalid: true,
		},
		{
			str:     "max: -1",
			invalid: true,
		},
		{
			str:     "min: -1",
			invalid: true,
		},
		{
			str:     "min: 10\nmax: 9",
			invalid: true,
		},
		{
			str:     "blah: 10\nmax: 9",
			invalid: true,
		},
		{
			str:     "min: -1\nmax: 9",
			invalid: true,
		},
	}

	sc := testutil.NewSchemaChecker(&config.RolloutConfig{})

	for tNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d_%s", tNum, table.str), func(t *testing.T) {
			obj, err := yaml.ParseYAMLString[config.RolloutConfig](table.str)
			if table.invalid {
				require.True(t, err != nil || sc.CheckYAML(t, table.str) != nil)
				return
			}
			require.NoError(t, err)

			require.NoError(t, obj.Validate())

			require.InDelta(t, table.min, obj.MinimumPercent()*100.0, 0.1, "MinimumPercent")
			require.InDelta(t, table.max, obj.MaximumPercent()*100.0, 0.1, "MaximumPercent")

			require.Equal(t, table.min, *(obj.GetAwsConfig().MinimumHealthyPercent))
			require.Equal(t, table.max, *(obj.GetAwsConfig().MaximumPercent))

			for _, check := range table.checks {

				t.Run(fmt.Sprintf("check_%d", check.count), func(t *testing.T) {
					if check.max < 0 || check.min < 0 {
						require.Error(t, obj.ValidateWithDesiredCount(check.count))
						return
					} else {
						require.NoError(t, obj.ValidateWithDesiredCount(check.count))
					}

					actualMin, actualMax := obj.GetMinMaxCount(check.count)
					require.Equal(t, check.min, actualMin)
					require.Equal(t, check.max, actualMax)
				})

			}

		})
	}
}

func TestRolloutConfig_NewDeploymentConfigFromService(t *testing.T) {

	// this just tricks the service to think it is a load balanced service
	fakeLbs := []config.LoadBalancer{{}, {}}

	tables := []struct {
		min int32
		max int32
		svc *config.Service
	}{
		{0, 100, &config.Service{DesiredCount: 0}},
		{100, 200, &config.Service{DesiredCount: 1}},

		// Load balanced
		{100, 200, &config.Service{DesiredCount: 1, LoadBalancers: fakeLbs}},
		{100, 200, &config.Service{DesiredCount: 2, LoadBalancers: fakeLbs}},
		{100, 200, &config.Service{DesiredCount: 3, LoadBalancers: fakeLbs}},
		{100, 200, &config.Service{DesiredCount: 4, LoadBalancers: fakeLbs}},

		{100, 150, &config.Service{DesiredCount: 5, LoadBalancers: fakeLbs}},
		{100, 150, &config.Service{DesiredCount: 100, LoadBalancers: fakeLbs}},
	}

	for tNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", tNum), func(t *testing.T) {
			ro := config.NewDeploymentConfigFromService(table.svc)

			require.Equal(t, table.min, *ro.Minimum)
			require.Equal(t, table.max, *ro.Maximum)
		})
	}
}
