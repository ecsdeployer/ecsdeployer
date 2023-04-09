package cleanupservices

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
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
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:service", tagMatcher, []string{}),
		})
		err := Step{}.Clean(ctx)
		require.NoError(t, err)
	})

	t.Run("when only relevant services", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:service", tagMatcher, []string{
				serviceArnPrefix + "dummy-web",
				serviceArnPrefix + "dummy-worker",
			}),
		})
		err := Step{}.Clean(ctx)
		require.NoError(t, err)
	})

	t.Run("when old services", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
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
				"service": "dummy-somethingelse",
				"force":   true,
			}, ecsTypes.Service{
				DesiredCount: 0,
				Status:       aws.String("DRAINING"),
			}),
		})
		err := Step{}.Clean(ctx)
		require.NoError(t, err)
	})

}
