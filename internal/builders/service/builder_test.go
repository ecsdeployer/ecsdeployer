package service

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/builders/buildtestutils"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
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

func TestBuildCreate(t *testing.T) {
	buildtestutils.StartMocker(t)
	t.Run("smoketest", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/smoke.yml", buildtestutils.OptSetNumSSMVars(4))

		lbService := buildtestutils.GetServiceTask(ctx.Project, "web")
		svcInput := genCreateService(t, ctx, lbService)

		require.NotNil(t, svcInput)

		require.NotNil(t, svcInput.LoadBalancers)
		require.Len(t, svcInput.LoadBalancers, 1)
		lbRecord := svcInput.LoadBalancers[0]
		require.Equal(t, "web", *lbRecord.ContainerName, "LB Container")
		require.Equal(t, "arn:aws:elasticloadbalancing:us-east-1:555555555555:targetgroup/c87-deployer-test-web/73e2d6bc24d8a067", *lbRecord.TargetGroupArn, "LB TargetGroup")

	})
}

func genCreateService(t *testing.T, ctx *config.Context, entity *config.Service) *ecs.CreateServiceInput {
	t.Helper()
	obj, err := BuildCreate(ctx, entity)
	require.NoError(t, err)

	_, err = awsclients.ECSClient().CreateService(ctx.Context, obj)
	require.NoError(t, err)

	return obj
}

func genUpdateService(t *testing.T, ctx *config.Context, entity *config.Service) *ecs.UpdateServiceInput {
	t.Helper()
	obj, err := BuildUpdate(ctx, entity)
	require.NoError(t, err)

	_, err = awsclients.ECSClient().UpdateService(ctx.Context, obj)
	require.NoError(t, err)

	return obj
}
