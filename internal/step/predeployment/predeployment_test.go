package predeployment

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestPredeploymentStep(t *testing.T) {
	testutil.DisableLoggingForTest(t)
	t.Run("String", func(t *testing.T) {
		require.Equal(t, "predeploy tasks", Step{}.String())
	})

	t.Run("Skip", func(t *testing.T) {
		t.Run("no predeploy tasks", func(t *testing.T) {
			ctx := config.New(&config.Project{})
			require.True(t, Step{}.Skip(ctx))
		})

		t.Run("has predeploy tasks", func(t *testing.T) {
			ctx := config.New(&config.Project{
				PreDeployTasks: []*config.PreDeployTask{
					{},
				},
			})
			require.False(t, Step{}.Skip(ctx))
		})
	})
}
