package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCronjobStep(t *testing.T) {

	t.Run("scheduler", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			// testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("faketg"),

			testutil.Mock_Scheduler_GetSchedule_Missing("dummy", "ecsd-cron-dummy-cron1"),
			testutil.Mock_Scheduler_CreateSchedule("dummy", "ecsd-cron-dummy-cron1"),

			testutil.Mock_Logs_CreateLogGroup_AllowAny(),
			testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
			// testutil.Mock_Events_PutRule_Generic(),
			// testutil.Mock_Events_PutTargets_Generic(),
		})

		err := CronjobStep(project.CronJobs[0], false).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("eventbridge", func(t *testing.T) {

		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			// testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("faketg"),
			testutil.Mock_Logs_CreateLogGroup_AllowAny(),
			testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
			testutil.Mock_Events_PutRule_Generic(),
			testutil.Mock_Events_PutTargets_Generic(),
		})
		ctx.Project.Settings.CronUsesEventing = true

		err := CronjobStep(project.CronJobs[0], true).Apply(ctx)
		require.NoError(t, err)
	})

}
