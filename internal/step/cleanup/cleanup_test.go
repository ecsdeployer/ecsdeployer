package cleanup

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestCleanupStep(t *testing.T) {
	testutil.DisableLoggingForTest(t)

	t.Run("String", func(t *testing.T) {
		require.Equal(t, "cleanup", Step{}.String())
	})

	t.Run("Skip", func(t *testing.T) {
		t.Run("sync disabled", func(t *testing.T) {
			kis := config.NewKeepInSyncFromBool(false)
			ctx := config.New(&config.Project{
				Settings: &config.Settings{
					KeepInSync: &kis,
				},
			})
			require.True(t, Step{}.Skip(ctx))
		})

		t.Run("not disabled", func(t *testing.T) {
			kis := config.NewKeepInSyncFromBool(true)
			ctx := config.New(&config.Project{
				Settings: &config.Settings{
					KeepInSync: &kis,
				},
			})
			require.False(t, Step{}.Skip(ctx))
		})
	})
}
