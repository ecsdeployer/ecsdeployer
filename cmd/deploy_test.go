package cmd

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/steps"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestDeployCmd(t *testing.T) {
	silenceLogging(t)

	// t.Run("checks function reference", func(t *testing.T) {
	// 	require.(t, steps.CleanupOnlyStep, stepDeploymentStepFunc)
	// })

	t.Run("calls correct function", func(t *testing.T) {
		oldRef := stepDeploymentStepFunc
		t.Cleanup(func() {
			stepDeploymentStepFunc = oldRef
		})

		awsmocker.Start(t, nil)

		wasCalled := false
		stepDeploymentStepFunc = func(_ *config.Project) *steps.Step {
			wasCalled = true
			return steps.NoopStep()
		}

		cmd := newDeployCmd(defaultCmdMetadata()).cmd
		cmd.Root().SetArgs([]string{"-q"})
		cmd.SetArgs([]string{"-c", "testdata/info_simple.yml"})

		_, _, err := executeCmdAndReturnOutput(cmd)

		require.NoError(t, err)

		require.True(t, wasCalled)
	})

	// tagMatcher := map[string]string{
	// 	"ecsdeployer/project": "dummy/fancy",
	// }

	// awsmocker.Start(t, &awsmocker.MockerOptions{
	// 	Mocks: []*awsmocker.MockedEndpoint{
	// 		testutil.Mock_Logs_DescribeLogGroups(nil),
	// 		testutil.Mock_Tagging_GetResources("events:rule", tagMatcher, []string{}),
	// 		testutil.Mock_Tagging_GetResources("ecs:task-definition", tagMatcher, []string{}),
	// 		testutil.Mock_Tagging_GetResources("ecs:service", tagMatcher, []string{}),
	// 	},
	// })

}
