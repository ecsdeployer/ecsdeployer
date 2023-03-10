package service_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/builders/service"
	"ecsdeployer.com/ecsdeployer/internal/testutil/buildtestutils"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
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
}
