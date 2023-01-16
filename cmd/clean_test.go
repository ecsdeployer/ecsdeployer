package cmd

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/steps"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestCleanCmd(t *testing.T) {
	silenceLogging(t)

	t.Run("calls correct function", func(t *testing.T) {
		oldRef := stepCleanupOnlyStepFunc
		t.Cleanup(func() {
			stepCleanupOnlyStepFunc = oldRef
		})

		testutil.StartMocker(t, nil)

		wasCalled := false
		stepCleanupOnlyStepFunc = func(_ *config.Project) *steps.Step {
			wasCalled = true
			return steps.NoopStep()
		}

		cmd := newCleanCmd(defaultCmdMetadata()).cmd
		cmd.Root().SetArgs([]string{"-q"})
		cmd.SetArgs([]string{"-c", "testdata/info_simple.yml"})

		_, _, err := executeCmdAndReturnOutput(cmd)

		require.NoError(t, err)

		require.True(t, wasCalled)
	})
}
