package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
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
		svc, err := yaml.ParseYAMLFile[config.Service](table.file)
		if table.valid && err != nil {
			t.Errorf("<%s> unexpected error: %s", table.file, err)
		}

		if !table.valid && err == nil {
			t.Errorf("<%s> expected error, but didnt get one", table.file)
		}

		if svc.IsLoadBalanced() != table.isLB {
			t.Errorf("<%s> expected IsLoadBalanced=%t, but was not", table.file, table.isLB)
		}

	}
}
