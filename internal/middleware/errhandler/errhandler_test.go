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
		require.NoError(t, Handle(func(ctx *config.Context) error {
			return nil
		})(nil))
	})

	t.Run("step skipped", func(t *testing.T) {
		require.NoError(t, Handle(func(ctx *config.Context) error {
			return step.Skip("some skipper error")
		})(nil))
	})

	t.Run("some err", func(t *testing.T) {
		require.Error(t, Handle(func(ctx *config.Context) error {
			return fmt.Errorf("step errored")
		})(nil))
	})
}
