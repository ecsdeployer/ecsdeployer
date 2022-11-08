package steps

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestPreflightStep(t *testing.T) {
	project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
		testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("faketg"),
	})

	err := PreflightStep(project).Apply(ctx)
	require.NoError(t, err)

}
