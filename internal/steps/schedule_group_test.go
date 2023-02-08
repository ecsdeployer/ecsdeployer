package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestScheduleGroupStep(t *testing.T) {
	t.Run("when no cronjobs", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{})
		ctx.Project.CronJobs = nil

		err := ScheduleGroupStep(project).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("schedule group not created", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Scheduler_GetScheduleGroup_Missing("dummy"),
			testutil.Mock_Scheduler_CreateScheduleGroup("dummy"),
		})

		err := ScheduleGroupStep(project).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("schedule group exists", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Scheduler_GetScheduleGroup("dummy"),
		})
		err := ScheduleGroupStep(project).Apply(ctx)
		require.NoError(t, err)
	})
}
