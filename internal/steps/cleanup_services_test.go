package steps

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCleanupServicesStep(t *testing.T) {

	tagMatcher := map[string]string{
		"ecsdeployer/project": "dummy",
	}
	serviceArnPrefix := fmt.Sprintf("arn:aws:ecs:%s:%s:service/testcluster/", awsmocker.DefaultRegion, awsmocker.DefaultAccountId)

	t.Run("when no results", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:service", tagMatcher, []string{}),
		})
		err := CleanupServicesStep(project.Settings.KeepInSync).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when only relevant services", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:service", tagMatcher, []string{
				serviceArnPrefix + "dummy-web",
				serviceArnPrefix + "dummy-worker",
			}),
		})
		err := CleanupServicesStep(project.Settings.KeepInSync).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when old services", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:service", tagMatcher, []string{
				serviceArnPrefix + "dummy-web",
				serviceArnPrefix + "dummy-worker",
				serviceArnPrefix + "dummy-somethingelse",
			}),
			testutil.Mock_ECS_DescribeServices_jmespath(map[string]any{
				"services[0]": serviceArnPrefix + "dummy-somethingelse",
			}, ecsTypes.Service{DesiredCount: 3}, 0),

			testutil.Mock_ECS_UpdateService_jmespath(map[string]any{
				"service":      serviceArnPrefix + "dummy-somethingelse",
				"desiredCount": 0,
			}, ecsTypes.Service{DesiredCount: 0}),

			testutil.Mock_ECS_DeleteService_jmespath(map[string]any{
				"service": serviceArnPrefix + "dummy-somethingelse",
				"force":   true,
			}, ecsTypes.Service{
				DesiredCount: 0,
				Status:       aws.String("DRAINING"),
			}),
		})
		err := CleanupServicesStep(project.Settings.KeepInSync).Apply(ctx)
		require.NoError(t, err)
	})

}
