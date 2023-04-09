package errhandler

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {

	testutil.DisableLoggingForTest(t)

	t.Run("no errors", func(t *testing.T) {
		nilThrower := func(ctx *config.Context) error {
			return nil
		}

		require.NoError(t, Handle(nilThrower)(nil))
		require.NoError(t, Ignore(nilThrower)(nil))
	})

	t.Run("step skipped", func(t *testing.T) {
		skipThrower := func(ctx *config.Context) error {
			return step.Skip("some skipper error")
		}

		require.NoError(t, Handle(skipThrower)(nil))
		require.NoError(t, Ignore(skipThrower)(nil))
	})

	t.Run("some err", func(t *testing.T) {
		errThrower := func(ctx *config.Context) error {
			return fmt.Errorf("step errored")
		}

		require.Error(t, Handle(errThrower)(nil))
		require.NoError(t, Ignore(errThrower)(nil))
	})
}
