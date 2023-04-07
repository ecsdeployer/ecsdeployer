package preflight

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestCheckCluster(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		require.Equal(t, "cluster", checkCluster{}.String())
	})

	t.Run("Check", func(t *testing.T) {
		testutil.MockSimpleStsProxy(t)

		t.Run("not set", func(t *testing.T) {
			err := checkCluster{}.Check(config.New(&config.Project{}))
			require.Error(t, err)
			require.ErrorContains(t, err, "No cluster information")
		})
	})
}
