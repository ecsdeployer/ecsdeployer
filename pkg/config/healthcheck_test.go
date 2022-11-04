package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
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
		if table.valid != (err == nil) {
			t.Errorf("Entry<%d> was expected to have Validate()==%t but it wasnt: %s", i, table.valid, err)

		}
	}
}
