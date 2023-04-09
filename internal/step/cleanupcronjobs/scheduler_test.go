package cleanupcronjobs

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCleanupCronjobsScheduler(t *testing.T) {
	t.Run("when no results", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Scheduler_ListSchedules("dummy", []testutil.MockListScheduleEntry{}),
		})
		ctx.Project.Settings.CronUsesEventing = false
		err := Step{}.Clean(ctx)
		// err := CleanupCronjobsStep(project.Settings.KeepInSync).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when only relevant cronjobs", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Scheduler_ListSchedules("dummy", []testutil.MockListScheduleEntry{
				{Name: "ecsd-cron-dummy-cron1"},
			}),
		})
		ctx.Project.Settings.CronUsesEventing = false
		err := Step{}.Clean(ctx)
		// err := CleanupCronjobsStep(project.Settings.KeepInSync).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when old cronjobs", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Scheduler_ListSchedules("dummy", []testutil.MockListScheduleEntry{
				{Name: "ecsd-cron-dummy-cron1"},
				{Name: "ecsd-cron-dummy-oldcron"},
				{Name: "ecsd-cron-dummy-othercron"},
				{Name: "ignored-cron"},
			}),

			testutil.Mock_Scheduler_DeleteSchedule("dummy", "ecsd-cron-dummy-oldcron"),
			testutil.Mock_Scheduler_DeleteSchedule("dummy", "ecsd-cron-dummy-othercron"),
		})
		ctx.Project.Settings.CronUsesEventing = false
		err := Step{}.Clean(ctx)
		// err := CleanupCronjobsStep(project.Settings.KeepInSync).Apply(ctx)
		require.NoError(t, err)
	})
}
