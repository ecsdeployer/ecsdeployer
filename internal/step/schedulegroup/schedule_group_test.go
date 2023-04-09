package schedulegroup

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestScheduleGroupStep(t *testing.T) {

	t.Run("using legacy cron flow", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Scheduler_GetScheduleGroup("dummy"),
		})
		ctx.Project.Settings.CronUsesEventing = true
		err := Step{}.Run(ctx)
		require.True(t, step.IsSkip(err))
	})

	t.Run("when no cronjobs", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{})
		ctx.Project.CronJobs = nil

		err := Step{}.Run(ctx)
		require.True(t, step.IsSkip(err))
	})

	t.Run("schedule group not created", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Scheduler_GetScheduleGroup_Missing("dummy"),
			testutil.Mock_Scheduler_CreateScheduleGroup("dummy"),
		})

		err := Step{}.Run(ctx)
		require.NoError(t, err)
	})

	t.Run("schedule group exists", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Scheduler_GetScheduleGroup("dummy"),
		})
		err := Step{}.Run(ctx)
		require.NoError(t, err)
	})
}
