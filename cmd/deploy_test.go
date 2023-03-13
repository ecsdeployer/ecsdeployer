package cmd

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/steps"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/caarlos0/log"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestDeployCmd(t *testing.T) {
	silenceLogging(t)

	t.Run("calls correct function", func(t *testing.T) {
		oldRef := stepDeploymentStepFunc
		t.Cleanup(func() {
			stepDeploymentStepFunc = oldRef
		})

		testutil.StartMocker(t, nil)

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
}

func TestDeploySmoke(t *testing.T) {
	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: []*awsmocker.MockedEndpoint{
			testutil.Mock_EC2_DescribeSecurityGroups_Simple(),
			testutil.Mock_EC2_DescribeSubnets_Simple(),
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
			testutil.Mock_ELBv2_DescribeTargetGroups_Generic_Success(),
			testutil.Mock_Logs_DescribeLogGroups(map[string]int32{}),
			testutil.Mock_SSM_GetParametersByPath("/ecsdeployer/dummy/", []string{"SSM_VAR1", "SSM_VAR2"}),
			// testutil.Mock_ECS_DeregisterTaskDefinition()
		},
	})

	cmd := newRootCmd("fake", func(i int) {}).cmd
	log.Strings[log.DebugLevel] = "%"

	// cmd := newDeployCmd(defaultCmdMetadata()).cmd
	// cmd.Root().SetArgs([]string{"-q"})
	// log.SetLevel(log.DebugLevel)
	// cmd.Root().SetArgs([]string{"--debug"})
	// cmd.SetArgs([]string{"deploy", "-c", "../internal/builders/testdata/smoke.yml", "--debug"})
	cmd.SetArgs([]string{"deploy", "-c", "../internal/builders/testdata/everything.yml", "--debug"})

	_, _, err := executeCmdAndReturnOutput(cmd)
	require.NoError(t, err)

}
