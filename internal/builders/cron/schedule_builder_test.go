package cron

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	schedulerTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/stretchr/testify/require"
)

func TestBuildSchedule(t *testing.T) {
	testutil.MockSimpleStsProxy(t)

	ctx, err := config.NewFromYAML("testdata/dummy.yml")
	require.NoError(t, err)

	// MUST MATCH THE ORDER OF THE dummy.yml FILE
	tables := []struct {
		name       string
		state      schedulerTypes.ScheduleState
		schedule   string
		desciption *string
	}{
		{"ecsd-cron-dummy-job1", schedulerTypes.ScheduleStateEnabled, "cron(0 9 * * ? *)", nil},
		{"ecsd-cron-dummy-job2", schedulerTypes.ScheduleStateEnabled, "rate(1 hour)", aws.String("somedesc")},
		{"ecsd-cron-dummy-job3", schedulerTypes.ScheduleStateDisabled, "rate(1 hour)", nil},
		{"ecsd-cron-dummy-job4", schedulerTypes.ScheduleStateEnabled, "rate(1 hour)", nil},
	}

	for i, table := range tables {
		schedule, err := BuildSchedule(ctx, ctx.Project.CronJobs[i], "faketask:1")
		require.NoErrorf(t, err, "Index#%d", i)
		require.EqualValuesf(t, table.name, *schedule.Name, "Index#%d, Name", i)
	}

}
