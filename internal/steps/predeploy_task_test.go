package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestPreDeployTaskStep(t *testing.T) {

	commonMocks := []*awsmocker.MockedEndpoint{
		testutil.Mock_Logs_CreateLogGroup_AllowAny(),
		testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
		testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		testutil.Mock_ECS_RunTask(),
		testutil.Mock_ECS_DescribeTasks_Pending("PENDING", 1),
	}

	t.Run("failed to launch", func(t *testing.T) {
		// Prepend the failure, (since the common has a success)
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", append([]*awsmocker.MockedEndpoint{
			testutil.Mock_ECS_RunTask_FailToLaunch(1),
		},
			commonMocks...,
		))
		err := PreDeployTaskStep(project.PreDeployTasks[0]).Apply(ctx)
		require.ErrorContains(t, err, "task failed to launch")
	})

	t.Run("successful run", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", append(commonMocks,
			testutil.Mock_ECS_DescribeTasks_Pending("RUNNING", 1),
			testutil.Mock_ECS_DescribeTasks_Stopped(ecsTypes.TaskStopCodeEssentialContainerExited, 0, 2), // 1 for waiter, 1 for recheck
		))
		err := PreDeployTaskStep(project.PreDeployTasks[0]).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("failed task", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", append(commonMocks,
			testutil.Mock_ECS_DescribeTasks_Stopped(ecsTypes.TaskStopCodeEssentialContainerExited, 1, 0),
		))
		t.Run("throws error", func(t *testing.T) {
			pdTask := project.PreDeployTasks[0]
			pdTask.IgnoreFailure = false
			err := PreDeployTaskStep(pdTask).Apply(ctx)
			require.Error(t, err)
			require.ErrorContains(t, err, "Container exited with code: 1")
		})

		t.Run("ignores error when requested", func(t *testing.T) {
			pdTask := project.PreDeployTasks[0]
			pdTask.IgnoreFailure = true
			err := PreDeployTaskStep(pdTask).Apply(ctx)
			require.NoError(t, err)
		})
	})

	t.Run("describeTask failure", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", append(commonMocks,
			testutil.Mock_ECS_DescribeTasks_Stopped(ecsTypes.TaskStopCodeEssentialContainerExited, 1, 1),
			testutil.Mock_ECS_DescribeTasks_Failure(1),
		))

		err := PreDeployTaskStep(project.PreDeployTasks[0]).Apply(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "Task failed to describe")
	})

	t.Run("failure stop codes", func(t *testing.T) {
		tables := []struct {
			stopCode  ecsTypes.TaskStopCode
			errString string
		}{
			{ecsTypes.TaskStopCodeUserInitiated, "User killed the task"},
			{ecsTypes.TaskStopCodeTaskFailedToStart, "Failed to Start"},

			// rare
			{ecsTypes.TaskStopCodeServiceSchedulerInitiated, "Some very weird stop code was given: ServiceSchedulerInitiated"},
			{ecsTypes.TaskStopCodeTerminationNotice, "Some very weird stop code was given: TerminationNotice"},

			// this one really isnt possible since we dont use spot for predeploys
			{ecsTypes.TaskStopCodeSpotInterruption, "Some very weird stop code was given: SpotInterruption"},
		}

		for _, table := range tables {
			t.Run(string(table.stopCode), func(t *testing.T) {
				project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", append(commonMocks,
					testutil.Mock_ECS_DescribeTasks_Pending("RUNNING", 1),
					testutil.Mock_ECS_DescribeTasks_Stopped(table.stopCode, 1, 2),
				))

				err := PreDeployTaskStep(project.PreDeployTasks[0]).Apply(ctx)
				require.Error(t, err)
				require.ErrorContains(t, err, table.errString)
			})
		}
	})

}
