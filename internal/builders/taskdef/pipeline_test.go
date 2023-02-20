package taskdef

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func TestPipelineBuild(t *testing.T) {
	testutil.MockSimpleStsProxy(t)

	ctx, err := config.NewFromYAML("testdata/dummy.yml")
	require.NoError(t, err)

	t.Run("load balanced service", func(t *testing.T) {
		lbService := ctx.Project.Services[0]

		taskDefinition, err := PipelineBuild(ctx, lbService)
		require.NoError(t, err)

		taskDefinitionOld, err := Build(ctx, lbService)
		require.NoError(t, err)

		require.NotNil(t, taskDefinition.Family, "Family")
		require.EqualValues(t, "dummy-svc1", *taskDefinition.Family)

		require.NotNil(t, taskDefinition.TaskRoleArn, "TaskRoleArn")
		require.Equal(t, "arn:aws:iam::555555555555:role/faketask", *taskDefinition.TaskRoleArn)

		require.NotNil(t, taskDefinition.ExecutionRoleArn, "ExecutionRoleArn")
		require.Equal(t, "arn:aws:iam::555555555555:role/fakeexec", *taskDefinition.ExecutionRoleArn)

		require.Equal(t, ecsTypes.NetworkModeAwsvpc, taskDefinition.NetworkMode)
		require.Contains(t, taskDefinition.RequiresCompatibilities, ecsTypes.CompatibilityFargate)

		require.NotNil(t, taskDefinition.Cpu, "Cpu")
		require.Equal(t, "1024", *taskDefinition.Cpu)

		require.NotNil(t, taskDefinition.Memory, "Memory")
		require.Equal(t, "2048", *taskDefinition.Memory)

		require.GreaterOrEqual(t, len(taskDefinition.ContainerDefinitions), 1, "number of containers")
		primaryCont := taskDefinition.ContainerDefinitions[0]
		require.NotNil(t, primaryCont, "primaryCont")
		require.Equal(t, []string{"bundle", "exec", "puma", "-C", "config/puma.rb"}, primaryCont.Command)

		require.Len(t, primaryCont.PortMappings, 1)
		require.NotNil(t, primaryCont.PortMappings[0].HostPort, "HostPort")
		require.EqualValues(t, 1234, *primaryCont.PortMappings[0].HostPort)
		require.EqualValues(t, "tcp", primaryCont.PortMappings[0].Protocol)

		require.NotNil(t, primaryCont.Image, "ImageURI")
		require.Equal(t, "fake:latest", *primaryCont.Image)

		require.EqualValues(t, 0, primaryCont.Cpu)
		require.Nil(t, primaryCont.Memory)
		require.Nil(t, primaryCont.MemoryReservation)

		jsonNew, err := util.Jsonify(taskDefinition)
		require.NoError(t, err)
		jsonOld, err := util.Jsonify(taskDefinitionOld)
		require.NoError(t, err)

		require.JSONEq(t, jsonOld, jsonNew)

	})

}
