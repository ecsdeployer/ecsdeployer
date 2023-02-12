package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCronSchedulesStep(t *testing.T) {
	t.Run("when no cronjobs", func(t *testing.T) {
		project, _ := stepTestAwsMocker(t, "testdata/project_advanced.yml", nil)
		project.CronJobs = nil
		step := CronSchedulesStep(project)
		require.Equal(t, "Noop", step.Label)
	})

	t.Run("with multiple jobs", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_multicron.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Scheduler_GetScheduleGroup("dummy"),
			testutil.Mock_Scheduler_GetSchedule("dummy", "ecsd-cron-dummy-cron1"),
			testutil.Mock_Scheduler_UpdateSchedule("dummy", "ecsd-cron-dummy-cron1"),
			testutil.Mock_Scheduler_GetSchedule("dummy", "ecsd-cron-dummy-cron2"),
			testutil.Mock_Scheduler_UpdateSchedule("dummy", "ecsd-cron-dummy-cron2"),
			testutil.Mock_Scheduler_GetSchedule("dummy", "ecsd-cron-dummy-cron3"),
			testutil.Mock_Scheduler_UpdateSchedule("dummy", "ecsd-cron-dummy-cron3"),
			testutil.Mock_Scheduler_GetSchedule("dummy", "ecsd-cron-dummy-cron4"),
			testutil.Mock_Scheduler_UpdateSchedule("dummy", "ecsd-cron-dummy-cron4"),
			testutil.Mock_Logs_CreateLogGroup_AllowAny(),
			testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		})

		err := CronSchedulesStep(project).Apply(ctx)
		require.NoError(t, err)
	})
}
