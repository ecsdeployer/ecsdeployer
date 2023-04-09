package cronjob

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCronjobStepScheduler(t *testing.T) {
	t.Run("schedule missing", func(t *testing.T) {
		project, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Scheduler_GetSchedule_Missing("dummy", "ecsd-cron-dummy-cron1"),
			testutil.Mock_Scheduler_CreateSchedule("dummy", "ecsd-cron-dummy-cron1"),
			testutil.Mock_Logs_CreateLogGroup_AllowAny(),
			testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		})
		ctx.Project.Settings.CronUsesEventing = false
		err := New(project.CronJobs[0]).Run(ctx)
		require.NoError(t, err)
	})

	t.Run("schedule exists", func(t *testing.T) {
		project, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Scheduler_GetSchedule("dummy", "ecsd-cron-dummy-cron1"),
			testutil.Mock_Scheduler_UpdateSchedule("dummy", "ecsd-cron-dummy-cron1"),
			testutil.Mock_Logs_CreateLogGroup_AllowAny(),
			testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		})

		ctx.Project.Settings.CronUsesEventing = false
		err := New(project.CronJobs[0]).Run(ctx)
		require.NoError(t, err)
	})
}
