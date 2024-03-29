package cleanupcronjobs

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCleanupCronjobsLegacy(t *testing.T) {
	tagMatcher := map[string]string{
		"ecsdeployer/project": "dummy",
	}
	ruleArnPrefix := fmt.Sprintf("arn:aws:events:%s:%s:rule/", awsmocker.DefaultRegion, awsmocker.DefaultAccountId)

	t.Run("when no results", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("events:rule", tagMatcher, []string{}),
		})
		ctx.Project.Settings.CronUsesEventing = true
		err := Step{}.Clean(ctx)
		require.NoError(t, err)
	})

	t.Run("when only relevant cronjobs", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("events:rule", tagMatcher, []string{
				ruleArnPrefix + "dummy-rule-cron1",
			}),
		})
		ctx.Project.Settings.CronUsesEventing = true
		err := Step{}.Clean(ctx)
		require.NoError(t, err)
	})

	t.Run("when old cronjobs", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
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
		err := Step{}.Clean(ctx)
		require.NoError(t, err)
	})
}
