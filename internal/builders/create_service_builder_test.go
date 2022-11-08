package builders

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestBuildCreateService_Basic(t *testing.T) {
	testutil.MockSimpleStsProxy(t)
	// just a basic test to make sure we can pass the common stuff thru it

	ctx, err := config.NewFromYAML("testdata/dummy.yml")
	require.NoError(t, err)

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
		createSvcInput, err := BuildCreateService(ctx, table.thing)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}

		if !createSvcInput.EnableECSManagedTags {
			t.Errorf("Got incorrect ECSManagedTags")
		}
		if len(createSvcInput.Tags) != 1 {
			t.Errorf("Expected 1 tag, got %d", len(createSvcInput.Tags))
		}

		if len(createSvcInput.LoadBalancers) != table.lbCount {
			t.Errorf("Expected %d LoadBalancers, but got %d", table.lbCount, len(createSvcInput.LoadBalancers))
		}

		if table.expGrace >= 0 {
			if createSvcInput.HealthCheckGracePeriodSeconds == nil {
				t.Errorf("Expected HealthCheckGrace to exist, but got nil")
			}

			if *createSvcInput.HealthCheckGracePeriodSeconds != table.expGrace {
				t.Errorf("Expected HealthCheckGrace to be %d, but got %d", table.expGrace, *createSvcInput.HealthCheckGracePeriodSeconds)
			}

		} else if createSvcInput.HealthCheckGracePeriodSeconds != nil {
			t.Errorf("Expected HealthCheckGrace to be nil, but got value")
		}

	}

}
