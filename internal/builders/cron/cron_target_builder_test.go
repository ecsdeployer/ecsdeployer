package cron

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/aws/aws-sdk-go-v2/aws"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func TestBuildCronTarget(t *testing.T) {
	closeMock := testutil.MockSimpleStsProxy(t)
	defer closeMock()

	ctx, err := testutil.LoadProjectConfig("testdata/dummy.yml")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// MUST MATCH THE ORDER OF THE dummy.yml FILE
	tables := []struct {
		name       string
		state      eventTypes.RuleState
		schedule   string
		desciption *string
	}{
		{"dummy-rule-job1", eventTypes.RuleStateEnabled, "cron(0 9 * * ? *)", nil},
		{"dummy-rule-job2", eventTypes.RuleStateEnabled, "rate(1 hour)", aws.String("somedesc")},
		{"dummy-rule-job3", eventTypes.RuleStateDisabled, "rate(1 hour)", nil},
		{"dummy-rule-job4", eventTypes.RuleStateEnabled, "rate(1 hour)", nil},
	}

	for i, table := range tables {
		putTargets, err := BuildCronTarget(ctx, ctx.Project.CronJobs[i], "faketask:1")
		if err != nil {
			t.Errorf("Unexpected error <index#%d>: %s", i, err)
		}

		if !testutil.AssertStringEquals(table.name, putTargets.Rule) {
			t.Errorf("Incorrect Rule. <index#%d> Expected=%s, got=%s", i, table.name, *putTargets.Rule)
		}

		// if !testutil.AssertEquals[eventtypes.RuleState](table.state, putTargets.State) {
		// 	t.Errorf("Incorrect State. <index#%d> Expected=%s, got=%s", i, table.state, putTargets.State)
		// }

		// if !testutil.AssertStringEquals(table.schedule, putTargets.ScheduleExpression) {
		// 	t.Errorf("Incorrect ScheduleExpression. <index#%d> Expected=%s, got=%s", i, table.schedule, *putTargets.ScheduleExpression)
		// }

		// if !testutil.AssertStringEquals(table.desciption, putTargets.Description) {
		// 	t.Errorf("Incorrect Desc. <index#%d> Expected=%v, got=%v", i, table.desciption, putTargets.Description)
		// }
	}

}
