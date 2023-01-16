package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestPreDeployTaskStep(t *testing.T) {

	// if testing.Short() {
	// 	t.SkipNow()
	// }

	project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
		// testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("faketg"),
		testutil.Mock_Logs_CreateLogGroup_AllowAny(),
		testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
		testutil.Mock_ECS_RegisterTaskDefinition_Generic(),

		testutil.Mock_ECS_RunTask_FailToLaunch(1),

		testutil.Mock_ECS_RunTask(),
		testutil.Mock_ECS_DescribeTasks_Pending("PENDING", 1),
		testutil.Mock_ECS_DescribeTasks_Pending("RUNNING", 1),
		testutil.Mock_ECS_DescribeTasks_Stopped(ecsTypes.TaskStopCodeEssentialContainerExited, 0, 2), // 1 for waiter, 1 for recheck

		testutil.Mock_ECS_DescribeTasks_Pending("PENDING", 1),
		testutil.Mock_ECS_DescribeTasks_Stopped(ecsTypes.TaskStopCodeEssentialContainerExited, 1, 2),
	})

	t.Run("failed to launch", func(t *testing.T) {
		err := PreDeployTaskStep(project.PreDeployTasks[0]).Apply(ctx)
		require.ErrorContains(t, err, "task failed to launch")
	})

	t.Run("successful run", func(t *testing.T) {
		err := PreDeployTaskStep(project.PreDeployTasks[0]).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("failed task", func(t *testing.T) {
		err := PreDeployTaskStep(project.PreDeployTasks[0]).Apply(ctx)
		require.Error(t, err)
	})

}
