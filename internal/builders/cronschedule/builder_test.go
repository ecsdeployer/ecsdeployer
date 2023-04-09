package cronschedule

import (
	"encoding/json"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	schedulerTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/stretchr/testify/require"
)

func TestBuildSchedule(t *testing.T) {
	testutil.MockSimpleStsProxy(t)

	ctx, err := config.NewFromYAML("../cron/testdata/dummy.yml")
	require.NoError(t, err)

	ctxNoCronEnv, err := config.NewFromYAML("../cron/testdata/dummy.yml")
	require.NoError(t, err)
	ctxNoCronEnv.Project.Settings.SkipCronEnvVars = true

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

	t.Run("with_cron_envvars", func(t *testing.T) {
		for i, table := range tables {
			schedule, err := BuildCreate(ctx, ctx.Project.CronJobs[i], "faketask:1")
			require.NoErrorf(t, err, "Index#%d", i)
			require.EqualValuesf(t, table.name, *schedule.Name, "Index#%d, Name", i)

			input := make(map[string]interface{})

			err = json.Unmarshal([]byte(*schedule.Target.Input), &input)
			require.NoError(t, err)

			require.Len(t, input["containerOverrides"], 1)
			require.Len(t, (input["containerOverrides"].([]interface{}))[0].(map[string]interface{})["environment"], len(config.DefaultCronEnvVars))

		}
	})

	t.Run("without_cron_envvars", func(t *testing.T) {
		for i, table := range tables {
			schedule, err := BuildCreate(ctxNoCronEnv, ctxNoCronEnv.Project.CronJobs[i], "faketask:1")
			require.NoErrorf(t, err, "Index#%d", i)
			require.EqualValuesf(t, table.name, *schedule.Name, "Index#%d, Name", i)

			require.Equal(t, "{}", *schedule.Target.Input)

		}
	})

}
