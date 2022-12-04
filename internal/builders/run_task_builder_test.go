package builders

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestBuildRunTask_Basic(t *testing.T) {

	// just a basic test to make sure we can pass the common stuff thru it

	testutil.MockSimpleStsProxy(t)

	ctx, err := config.NewFromYAML("testdata/dummy.yml")
	require.NoError(t, err)

	tables := []struct {
		thing *config.PreDeployTask
	}{
		{ctx.Project.PreDeployTasks[0]},
		{ctx.Project.PreDeployTasks[1]},
	}

	for _, table := range tables {
		runTask, err := BuildRunTask(ctx, table.thing)
		require.NoError(t, err)
		require.True(t, runTask.EnableECSManagedTags)

	}

}
