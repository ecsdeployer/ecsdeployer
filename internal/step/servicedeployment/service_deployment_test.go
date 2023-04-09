package servicedeployment

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestServiceDeploymentStep(t *testing.T) {
	testutil.DisableLoggingForTest(t)

	t.Run("String", func(t *testing.T) {
		require.Equal(t, "deploying services", Step{}.String())
	})

	t.Run("Skip", func(t *testing.T) {
		t.Run("no services", func(t *testing.T) {
			ctx := config.New(&config.Project{})
			require.True(t, Step{}.Skip(ctx))
		})

		t.Run("has services", func(t *testing.T) {
			ctx := config.New(&config.Project{
				Services: []*config.Service{
					{},
				},
			})
			require.False(t, Step{}.Skip(ctx))
		})
	})

	// Run is tested via all the other tests
}
