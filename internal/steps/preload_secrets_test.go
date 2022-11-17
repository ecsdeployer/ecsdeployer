package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
	"golang.org/x/exp/maps"
)

func TestPreloadSecretsStep(t *testing.T) {

	t.Run("when disabled", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_simple.yml", []*awsmocker.MockedEndpoint{})
		err := PreloadSecretsStep(project).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when access denied", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			awsmocker.Mock_Failure("ssm", "GetParametersByPath"),
		})
		err := PreloadSecretsStep(project).Apply(ctx)
		require.ErrorContains(t, err, "AccessDenied")
	})

	t.Run("when no secrets", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_SSM_GetParametersByPath("/ecsdeployer/secrets/dummy/", []string{}),
		})
		err := PreloadSecretsStep(project).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when secrets returned", func(t *testing.T) {
		secretNames := []string{
			"SOMEVAR_NAME",
			"DATABASE_URL",
			"VAR1",
			"VAR2",
			"SOME_REALLY_LONG_VARIABLE_NAME",
			"WHATEVER",
		}
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_SSM_GetParametersByPath("/ecsdeployer/secrets/dummy/", secretNames),
		})
		err := PreloadSecretsStep(project).Apply(ctx)
		require.NoError(t, err)

		require.Subset(t, maps.Keys(ctx.Cache.SSMSecrets), secretNames)

		for _, name := range secretNames {
			require.NotNil(t, ctx.Cache.SSMSecrets[name])
			require.IsType(t, config.EnvVar{}, ctx.Cache.SSMSecrets[name])

			require.NotNil(t, ctx.Cache.SSMSecrets[name].ValueSSM)

			secretArn := *ctx.Cache.SSMSecrets[name].ValueSSM
			arn, err := arn.Parse(secretArn)
			require.NoError(t, err)
			require.Equal(t, "parameter/ecsdeployer/secrets/dummy/"+name, arn.Resource)
		}

	})
}
