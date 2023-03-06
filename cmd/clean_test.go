package cmd

import (
	"errors"
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

	t.Run("handles failures", func(t *testing.T) {
		oldRef := stepCleanupOnlyStepFunc
		t.Cleanup(func() {
			stepCleanupOnlyStepFunc = oldRef
		})

		testutil.StartMocker(t, nil)

		expectedErr := errors.New("explode")

		wasCalled := false
		stepCleanupOnlyStepFunc = func(_ *config.Project) *steps.Step {
			wasCalled = true
			return steps.NewStep(&steps.Step{
				Label: "Failure",
				Create: func(ctx *config.Context, s *steps.Step, sm *steps.StepMetadata) (steps.OutputFields, error) {
					return nil, expectedErr
				},
			})
		}

		cmd := newCleanCmd(defaultCmdMetadata()).cmd
		cmd.Root().SetArgs([]string{"-q"})
		cmd.SetArgs([]string{"-c", "testdata/info_simple.yml"})

		_, _, err := executeCmdAndReturnOutput(cmd)

		require.Error(t, err)

		var checkErr *exitError

		require.ErrorAs(t, err, &checkErr)
		require.ErrorIs(t, checkErr.err, expectedErr)

		require.True(t, wasCalled)
	})
}
