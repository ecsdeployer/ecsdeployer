package preflight

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestPreflightStep(t *testing.T) {
	_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
		testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("faketg"),
	})

	err := Step{}.Run(ctx)
	require.NoError(t, err)
}

// TODO: Need to add failure cases for each preflight step
