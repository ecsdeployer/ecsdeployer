package cronjob

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCronjobStepLegacy(t *testing.T) {
	project, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
		// testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("faketg"),
		testutil.Mock_Logs_CreateLogGroup_AllowAny(),
		testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
		testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		testutil.Mock_Events_PutRule_Generic(),
		testutil.Mock_Events_PutTargets_Generic(),
	})
	ctx.Project.Settings.CronUsesEventing = true

	err := New(project.CronJobs[0]).Run(ctx)
	require.NoError(t, err)
}
