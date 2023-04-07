package crondeployment

import (
	"bytes"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
	"github.com/webdestroya/go-log"
)

func TestCronDeploymentStep(t *testing.T) {

	testutil.DisableLoggingForTest(t)

	t.Run("String", func(t *testing.T) {
		require.Equal(t, "deploying cronjobs", Step{}.String())
	})

	t.Run("Skip", func(t *testing.T) {
		t.Run("no cronjobs", func(t *testing.T) {
			ctx := config.New(&config.Project{})
			require.True(t, Step{}.Skip(ctx))
		})

		t.Run("has cronjobs", func(t *testing.T) {
			ctx := config.New(&config.Project{
				CronJobs: []*config.CronJob{
					{},
				},
			})
			require.False(t, Step{}.Skip(ctx))
		})
	})

	t.Run("Run", func(t *testing.T) {
		t.Run("legacy", func(t *testing.T) {
			oldLog := log.Log
			t.Cleanup(func() {
				log.Log = oldLog
			})

			var w bytes.Buffer
			log.Log = log.New(&w)

			_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
				// testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("faketg"),
				testutil.Mock_Logs_CreateLogGroup_AllowAny(),
				testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
				testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
				testutil.Mock_Events_PutRule_Generic(),
				testutil.Mock_Events_PutTargets_Generic(),
			})
			ctx.Project.Settings.CronUsesEventing = true

			require.NoError(t, Step{}.Run(ctx))

			require.Contains(t, "DEPRECATED", w.String())

		})
	})
}
