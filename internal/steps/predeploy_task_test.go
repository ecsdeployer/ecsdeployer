package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestPreDeployTaskStep(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
		// testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("faketg"),
		testutil.Mock_Logs_CreateLogGroup(),
		testutil.Mock_Logs_PutRetentionPolicy(),
		testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		testutil.Mock_ECS_RunTask(),
		testutil.Mock_ECS_DescribeTasks_Pending("PENDING", 1),
		testutil.Mock_ECS_DescribeTasks_Pending("RUNNING", 1),
		// testutil.Mock_ECS_DescribeTasks_Pending("STOPPED", 1),
	})

	t.Run("successful run", func(t *testing.T) {
		err := PreDeployTaskStep(project.PreDeployTasks[0]).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("failed task", func(t *testing.T) {
		err := PreDeployTaskStep(project.PreDeployTasks[0]).Apply(ctx)
		require.NoError(t, err)
	})

}
