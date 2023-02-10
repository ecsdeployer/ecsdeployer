package steps

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCleanupCronjobsStep(t *testing.T) {

	t.Run("scheduler", func(t *testing.T) {

		t.Run("when no results", func(t *testing.T) {
			project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
				testutil.Mock_Scheduler_ListSchedules("dummy", []testutil.MockListScheduleEntry{}),
			})
			ctx.Project.Settings.CronUsesEventing = false
			err := CleanupCronjobsStep(project.Settings.KeepInSync).Apply(ctx)
			require.NoError(t, err)
		})

		t.Run("when only relevant cronjobs", func(t *testing.T) {
			project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
				testutil.Mock_Scheduler_ListSchedules("dummy", []testutil.MockListScheduleEntry{
					{Name: "ecsd-cron-dummy-cron1"},
				}),
			})
			ctx.Project.Settings.CronUsesEventing = false
			err := CleanupCronjobsStep(project.Settings.KeepInSync).Apply(ctx)
			require.NoError(t, err)
		})

		t.Run("when old cronjobs", func(t *testing.T) {
			project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
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
			err := CleanupCronjobsStep(project.Settings.KeepInSync).Apply(ctx)
			require.NoError(t, err)
		})
	})

	/////// DEPRECATED

	t.Run("eventbridge", func(t *testing.T) {

		tagMatcher := map[string]string{
			"ecsdeployer/project": "dummy",
		}
		ruleArnPrefix := fmt.Sprintf("arn:aws:events:%s:%s:rule/", awsmocker.DefaultRegion, awsmocker.DefaultAccountId)

		t.Run("when no results", func(t *testing.T) {
			project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
				testutil.Mock_Tagging_GetResources("events:rule", tagMatcher, []string{}),
			})
			ctx.Project.Settings.CronUsesEventing = true
			err := CleanupCronjobsStep(project.Settings.KeepInSync).Apply(ctx)
			require.NoError(t, err)
		})

		t.Run("when only relevant cronjobs", func(t *testing.T) {
			project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
				testutil.Mock_Tagging_GetResources("events:rule", tagMatcher, []string{
					ruleArnPrefix + "dummy-rule-cron1",
				}),
			})
			ctx.Project.Settings.CronUsesEventing = true
			err := CleanupCronjobsStep(project.Settings.KeepInSync).Apply(ctx)
			require.NoError(t, err)
		})

		t.Run("when old cronjobs", func(t *testing.T) {
			project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
				testutil.Mock_Tagging_GetResources("events:rule", tagMatcher, []string{
					ruleArnPrefix + "dummy-rule-cron1",
					ruleArnPrefix + "dummy-rule-oldcron",
					ruleArnPrefix + "custombus/dummy-rule-oldcustcron",
				}),

				testutil.Mock_Events_ListTargetsByRule("dummy-rule-oldcron", "", []string{"dummy-target-oldcron"}),
				testutil.Mock_Events_RemoveTargets("dummy-rule-oldcron", "", "dummy-target-oldcron"),
				testutil.Mock_Events_DeleteRule("dummy-rule-oldcron", ""),

				testutil.Mock_Events_ListTargetsByRule("dummy-rule-oldcustcron", "custombus", []string{"dummy-target-oldcustcron"}),
				testutil.Mock_Events_RemoveTargets("dummy-rule-oldcustcron", "custombus", "dummy-target-oldcustcron"),
				testutil.Mock_Events_DeleteRule("dummy-rule-oldcustcron", "custombus"),
			})
			ctx.Project.Settings.CronUsesEventing = true
			err := CleanupCronjobsStep(project.Settings.KeepInSync).Apply(ctx)
			require.NoError(t, err)
		})
	})
}
