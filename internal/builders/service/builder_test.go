package service

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestBuildCreate_Basic(t *testing.T) {
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
		t.Run(table.thing.Name, func(t *testing.T) {
			createSvcInput, err := BuildCreate(ctx, table.thing)
			require.NoError(t, err)

			require.Truef(t, createSvcInput.EnableECSManagedTags, "ECSManagedTags")

			require.Lenf(t, createSvcInput.Tags, 1, "Tags")

			require.Lenf(t, createSvcInput.LoadBalancers, table.lbCount, "LoadBalancers")

			if table.expGrace >= 0 {
				require.NotNil(t, createSvcInput.HealthCheckGracePeriodSeconds)
				require.Equal(t, table.expGrace, *createSvcInput.HealthCheckGracePeriodSeconds)
			} else {
				require.Nil(t, createSvcInput.HealthCheckGracePeriodSeconds)
			}
		})

	}

}
