package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestService_Unmarshal_Basic(t *testing.T) {
	tables := []struct {
		file  string
		valid bool
		isLB  bool
	}{
		{"testdata/service_single_loadbalancer.yml", true, true},
		{"testdata/service_multi_loadbalancer.yml", true, true},
		{"testdata/service_worker.yml", true, false},
	}

	for _, table := range tables {
		t.Run(table.file, func(t *testing.T) {
			svc, err := yaml.ParseYAMLFile[config.Service](table.file)

			if table.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			require.Truef(t, svc.IsTaskStruct(), "IsTaskStruct")
			require.Equalf(t, table.isLB, svc.IsLoadBalanced(), "IsLoadBalanced")
			require.Equalf(t, !table.isLB, svc.IsWorker(), "IsWorker")
		})

	}
}

func TestService_Validate(t *testing.T) {
	require.ErrorContains(t, (&config.Service{DesiredCount: -1}).Validate(), "desired count cannot")
	require.ErrorContains(t, (&config.Service{DesiredCount: 1, RolloutConfig: &config.RolloutConfig{Minimum: util.Ptr(int32(100)), Maximum: util.Ptr(int32(110))}}).Validate(), "impossible")
}
