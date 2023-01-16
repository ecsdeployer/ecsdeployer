package taskdef

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestBuild_Basic(t *testing.T) {

	// just a basic test to make sure we can pass the common stuff thru it
	testutil.MockSimpleStsProxy(t)

	ctx, err := config.NewFromYAML("testdata/dummy.yml")
	require.NoError(t, err)

	tables := []struct {
		thing config.IsTaskStruct
	}{
		{ctx.Project.ConsoleTask},

		{ctx.Project.PreDeployTasks[0]},
		{ctx.Project.PreDeployTasks[1]},

		{ctx.Project.Services[0]},
		{ctx.Project.Services[1]},

		{ctx.Project.CronJobs[0]},
	}

	for _, table := range tables {
		taskDefinition, err := Build(ctx, table.thing)
		require.NoError(t, err)
		require.Equal(t, "fake:latest", *taskDefinition.ContainerDefinitions[0].Image)
		require.Len(t, taskDefinition.Tags, 1)
	}

}
