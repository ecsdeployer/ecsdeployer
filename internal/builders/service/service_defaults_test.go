package service_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/builders/service"
	"ecsdeployer.com/ecsdeployer/internal/testutil/buildtestutils"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func TestServiceDefaults(t *testing.T) {
	buildtestutils.StartMocker(t)
	t.Run("service_role", func(t *testing.T) {
		t.Run("no role", func(t *testing.T) {
			ctx := buildtestutils.LoadProjectConfig(t, "../testdata/dummy.yml")

			createSvcInput, err := service.BuildCreate(ctx, ctx.Project.Services[0])
			require.NoError(t, err)

			require.Nil(t, createSvcInput.Role)
		})

		t.Run("with role", func(t *testing.T) {
			ctx := buildtestutils.LoadProjectConfig(t, "../testdata/dummy.yml")

			roleArn, err := yaml.ParseYAMLString[config.RoleArn]("fakeservicerole")
			require.NoError(t, err)

			ctx.Project.ServiceRole = roleArn

			createSvcInput, err := service.BuildCreate(ctx, ctx.Project.Services[0])
			require.NoError(t, err)

			require.NotNil(t, createSvcInput.Role)
			require.Equal(t, "arn:aws:iam::555555555555:role/fakeservicerole", *createSvcInput.Role)

		})
	})

	t.Run("az_rebalancing", func(t *testing.T) {
		t.Run("disabled when max is 100", func(t *testing.T) {
			ctx := buildtestutils.LoadProjectConfig(t, "../testdata/dummy.yml")

			svc := ctx.Project.Services[1] // svc2: no LB, desired=1
			svc.RolloutConfig = &config.RolloutConfig{
				Minimum: aws.Int32(50),
				Maximum: aws.Int32(100),
			}

			createSvcInput, err := service.BuildCreate(ctx, svc)
			require.NoError(t, err)

			require.Equal(t, ecsTypes.AvailabilityZoneRebalancingDisabled, createSvcInput.AvailabilityZoneRebalancing)
		})

		t.Run("disabled when max is below 100", func(t *testing.T) {
			ctx := buildtestutils.LoadProjectConfig(t, "../testdata/dummy.yml")

			svc := ctx.Project.Services[1]
			svc.RolloutConfig = &config.RolloutConfig{
				Minimum: aws.Int32(0),
				Maximum: aws.Int32(50),
			}

			createSvcInput, err := service.BuildCreate(ctx, svc)
			require.NoError(t, err)

			require.Equal(t, ecsTypes.AvailabilityZoneRebalancingDisabled, createSvcInput.AvailabilityZoneRebalancing)
		})

		t.Run("not set when max is above 100", func(t *testing.T) {
			ctx := buildtestutils.LoadProjectConfig(t, "../testdata/dummy.yml")

			svc := ctx.Project.Services[1]
			svc.RolloutConfig = &config.RolloutConfig{
				Minimum: aws.Int32(50),
				Maximum: aws.Int32(200),
			}

			createSvcInput, err := service.BuildCreate(ctx, svc)
			require.NoError(t, err)

			// When max > 100, AZ rebalancing should not be explicitly set (zero value),
			// allowing AWS to use its default (ENABLED for create, existing value for update).
			require.Empty(t, createSvcInput.AvailabilityZoneRebalancing)
		})

		t.Run("passed through to update", func(t *testing.T) {
			ctx := buildtestutils.LoadProjectConfig(t, "../testdata/dummy.yml")

			svc := ctx.Project.Services[1]
			svc.RolloutConfig = &config.RolloutConfig{
				Minimum: aws.Int32(50),
				Maximum: aws.Int32(100),
			}

			updateSvcInput, err := service.BuildUpdate(ctx, svc)
			require.NoError(t, err)

			require.Equal(t, ecsTypes.AvailabilityZoneRebalancingDisabled, updateSvcInput.AvailabilityZoneRebalancing)
		})
	})
}
