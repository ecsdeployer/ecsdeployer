package console

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestConsoleStep(t *testing.T) {
	require.Equal(t, "registering console task", Step{}.String())

	t.Run("happy path", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Logs_CreateLogGroup_AllowAny(),
			testutil.Mock_Logs_PutRetentionPolicy_AllowAny(),
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		})
		err := Step{}.Run(ctx)
		require.NoError(t, err)
	})
}
