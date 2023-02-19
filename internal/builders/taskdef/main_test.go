package taskdef

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
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
		require.Len(t, taskDefinition.Tags, 2)
	}

}

func TestBuild(t *testing.T) {
	testutil.MockSimpleStsProxy(t)

	ctx, err := config.NewFromYAML("testdata/dummy.yml")
	require.NoError(t, err)

	t.Run("load balanced service", func(t *testing.T) {
		lbService := ctx.Project.Services[0]

		taskDefinition, err := Build(ctx, lbService)
		require.NoError(t, err)

		require.EqualValues(t, "dummy-svc1", *taskDefinition.Family)

		require.Equal(t, "arn:aws:iam::555555555555:role/faketask", *taskDefinition.TaskRoleArn)
		require.Equal(t, "arn:aws:iam::555555555555:role/fakeexec", *taskDefinition.ExecutionRoleArn)
		require.Equal(t, ecsTypes.NetworkModeAwsvpc, taskDefinition.NetworkMode)
		require.Contains(t, taskDefinition.RequiresCompatibilities, ecsTypes.CompatibilityFargate)

		require.Equal(t, "1024", *taskDefinition.Cpu)
		require.Equal(t, "2048", *taskDefinition.Memory)

		primaryCont := taskDefinition.ContainerDefinitions[0]
		require.NotNil(t, primaryCont)
		require.Equal(t, []string{"bundle", "exec", "puma", "-C", "config/puma.rb"}, primaryCont.Command)

		require.Len(t, primaryCont.PortMappings, 1)
		require.EqualValues(t, 1234, *primaryCont.PortMappings[0].HostPort)
		require.EqualValues(t, "tcp", primaryCont.PortMappings[0].Protocol)

		require.Equal(t, "fake:latest", *primaryCont.Image)

		require.EqualValues(t, 0, primaryCont.Cpu)
		require.Nil(t, primaryCont.Memory)
		require.Nil(t, primaryCont.MemoryReservation)
	})

}
