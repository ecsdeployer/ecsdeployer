package service

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestBuildUpdate_Basic(t *testing.T) {

	// just a basic test to make sure we can pass the common stuff thru it
	testutil.MockSimpleStsProxy(t)

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
		t.Run(table.thing.Name, func(t *testing.T) {
			svcInput, err := BuildUpdate(ctx, table.thing)
			require.NoError(t, err)

			require.Truef(t, *svcInput.EnableECSManagedTags, "ECSManagedTags")

			require.Lenf(t, svcInput.LoadBalancers, table.lbCount, "LoadBalancer count mismatch")

			if table.expGrace >= 0 {
				require.NotNilf(t, svcInput.HealthCheckGracePeriodSeconds, "Expected HealthCheckGrace to exist, but got nil")
				require.Equalf(t, table.expGrace, *svcInput.HealthCheckGracePeriodSeconds, "Expected HealthCheckGrace to be %d, but got %d", table.expGrace, *svcInput.HealthCheckGracePeriodSeconds)

			} else if svcInput.HealthCheckGracePeriodSeconds != nil {
				require.Nil(t, svcInput.HealthCheckGracePeriodSeconds)
			}
		})
	}

}
