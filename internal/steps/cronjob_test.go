package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCronjobStep(t *testing.T) {
	project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
		// testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("faketg"),
		testutil.Mock_Logs_CreateLogGroup(),
		testutil.Mock_Logs_PutRetentionPolicy(),
		testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		testutil.Mock_Events_PutRule_Generic(),
		testutil.Mock_Events_PutTargets_Generic(),
	})

	err := CronjobStep(project.CronJobs[0]).Apply(ctx)
	require.NoError(t, err)

}
