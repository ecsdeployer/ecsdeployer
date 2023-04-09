package taskdefinition

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestTaskDefinitionSubstep(t *testing.T) {
	expectedTaskDefArn := fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/dummy-web:999", awsmocker.DefaultRegion, awsmocker.DefaultAccountId)

	t.Run("happy path", func(t *testing.T) {
		project, ctx := steptestutil.StepTestAwsMocker(t, "../../step/testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{

			testutil.Mock_Logs_CreateLogGroup_AllowAny(),
			testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		})

		taskDefArn, err := Register(ctx, project.Services[0])
		require.NoError(t, err)
		require.Equal(t, expectedTaskDefArn, taskDefArn)

	})

	t.Run("unhappy path", func(t *testing.T) {
		t.Run("log denied", func(t *testing.T) {
			project, ctx := steptestutil.StepTestAwsMocker(t, "../../step/testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
				testutil.Mock_Logs_CreateLogGroup_Deny("/ecsdeployer/app/dummy/web"),
			})

			_, err := Register(ctx, project.Services[0])
			require.Error(t, err)
			require.False(t, step.IsSkip(err))
			require.ErrorContains(t, err, "Failed to provision log group")
		})

		t.Run("log already exists", func(t *testing.T) {
			project, ctx := steptestutil.StepTestAwsMocker(t, "../../step/testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
				testutil.Mock_Logs_CreateLogGroup_AlreadyExists("/ecsdeployer/app/dummy/web"),
				testutil.Mock_Logs_DescribeLogGroups_Single("/ecsdeployer/app/dummy/web", 30),
				testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
				testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
			})

			taskDefArn, err := Register(ctx, project.Services[0])
			require.NoError(t, err)
			require.Equal(t, expectedTaskDefArn, taskDefArn)
		})
	})
}
