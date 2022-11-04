package builders

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestBuildUpdateService_Basic(t *testing.T) {

	// just a basic test to make sure we can pass the common stuff thru it
	closeMock := testutil.MockSimpleStsProxy(t)
	defer closeMock()

	ctx, err := testutil.LoadProjectConfig("testdata/dummy.yml")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	tables := []struct {
		thing    *config.Service
		expGrace int32
		lbCount  int
	}{
		{ctx.Project.Services[0], -1, 1},
		{ctx.Project.Services[1], -1, 0},
		{ctx.Project.Services[2], 55, 1},
		{ctx.Project.Services[3], 122, 3},
	}

	for _, table := range tables {
		svcInput, err := BuildUpdateService(ctx, table.thing)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}

		if !*svcInput.EnableECSManagedTags {
			t.Errorf("Got incorrect ECSManagedTags")
		}

		if len(svcInput.LoadBalancers) != table.lbCount {
			t.Errorf("Expected %d LoadBalancers, but got %d", table.lbCount, len(svcInput.LoadBalancers))
		}

		if table.expGrace >= 0 {
			if svcInput.HealthCheckGracePeriodSeconds == nil {
				t.Errorf("Expected HealthCheckGrace to exist, but got nil")
			}

			if *svcInput.HealthCheckGracePeriodSeconds != table.expGrace {
				t.Errorf("Expected HealthCheckGrace to be %d, but got %d", table.expGrace, *svcInput.HealthCheckGracePeriodSeconds)
			}

		} else if svcInput.HealthCheckGracePeriodSeconds != nil {
			t.Errorf("Expected HealthCheckGrace to be nil, but got value")
		}

	}

}
