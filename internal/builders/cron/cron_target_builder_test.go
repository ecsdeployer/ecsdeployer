package cron

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/stretchr/testify/require"
)

func TestBuildCronTarget(t *testing.T) {
	testutil.MockSimpleStsProxy(t)

	ctx, err := config.NewFromYAML("testdata/dummy.yml")
	require.NoError(t, err)

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
		require.NoErrorf(t, err, "Index#%d", i)
		require.EqualValuesf(t, table.name, *putTargets.Rule, "Index#%d, Name", i)
	}

}
