package service

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/testutil/buildtestutils"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/stretchr/testify/require"
)

func TestBuildUpdate_Basic(t *testing.T) {

	// just a basic test to make sure we can pass the common stuff thru it
	buildtestutils.StartMocker(t)

	ctx := buildtestutils.LoadProjectConfig(t, "../testdata/dummy.yml")

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
			svcInput := genUpdateService(t, ctx, table.thing)

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

func genUpdateService(t *testing.T, ctx *config.Context, entity *config.Service) *ecs.UpdateServiceInput {
	t.Helper()
	obj, err := BuildUpdate(ctx, entity)
	require.NoError(t, err)

	_, err = awsclients.ECSClient().UpdateService(ctx.Context, obj)
	require.NoError(t, err)

	return obj
}
