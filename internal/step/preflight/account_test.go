package preflight

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCheckAccount(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		require.Equal(t, "aws account", checkAccount{}.String())
	})

	t.Run("CheckAndSkip", func(t *testing.T) {

		testutil.MockSimpleStsProxy(t)

		t.Run("not set", func(t *testing.T) {
			ctx := config.New(&config.Project{EcsDeployerOptions: &config.EcsDeployerOptions{}})
			require.True(t, checkAccount{}.Skip(ctx))
			require.NoError(t, checkAccount{}.Check(ctx))
		})

		t.Run("set same account", func(t *testing.T) {
			ctx := config.New(&config.Project{EcsDeployerOptions: &config.EcsDeployerOptions{
				AllowedAccountId: util.Ptr(awsmocker.DefaultAccountId),
			}})
			require.False(t, checkAccount{}.Skip(ctx))
			require.NoError(t, checkAccount{}.Check(ctx))
		})

		t.Run("set diff account", func(t *testing.T) {
			ctx := config.New(&config.Project{EcsDeployerOptions: &config.EcsDeployerOptions{
				AllowedAccountId: util.Ptr("111111111111"),
			}})
			require.False(t, checkAccount{}.Skip(ctx))
			err := checkAccount{}.Check(ctx)
			require.Error(t, err)
			require.ErrorContains(t, err, "is not an allowed account")
		})
	})
}
