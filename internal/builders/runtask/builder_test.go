package runtask

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/buildtestutils"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func TestBuildRunTask_Basic(t *testing.T) {

	// just a basic test to make sure we can pass the common stuff thru it

	buildtestutils.StartMocker(t)

	ctx := buildtestutils.LoadProjectConfig(t, "../testdata/dummy.yml")

	tables := []struct {
		thing *config.PreDeployTask
	}{
		{ctx.Project.PreDeployTasks[0]},
		{ctx.Project.PreDeployTasks[1]},
	}

	for _, table := range tables {
		runTask := genRunTaskDef(t, ctx, table.thing)
		require.True(t, runTask.EnableECSManagedTags)

	}

}

func TestBuild_Detailed(t *testing.T) {
	buildtestutils.StartMocker(t)

	ctx := buildtestutils.LoadProjectConfig(t, "../testdata/dummy.yml")

	t.Run("normal", func(t *testing.T) {

		pdTest1Yaml := `
		name: testpd1
		command: "something something"
		`

		pdTask, err := yaml.ParseYAMLString[config.PreDeployTask](testutil.CleanTestYaml(pdTest1Yaml))
		require.NoError(t, err)

		runTask := genRunTaskDef(t, ctx, pdTask)

		require.EqualValues(t, 1, *runTask.Count)

		require.True(t, runTask.EnableECSManagedTags, "ECSManagedTags")
		require.False(t, runTask.EnableExecuteCommand, "EnableExecuteCommand")
		require.EqualValues(t, ecsTypes.PropagateTagsTaskDefinition, runTask.PropagateTags, "PropagateTags")
		require.EqualValues(t, ecsTypes.LaunchTypeFargate, runTask.LaunchType, "LaunchType")
		require.Equal(t, "LATEST", *runTask.PlatformVersion, "PlatformVersion")

		require.Equal(t, "ecsd:dummy:deployer", *runTask.StartedBy, "StartedBy")
		require.Equal(t, "ecsd:dummy:pd:testpd1", *runTask.Group, "Group")

		require.Equal(t, "arn:aws:ecs:us-east-1:555555555555:cluster/fakecluster", *runTask.Cluster, "Cluster")

		require.Equal(t, ecsTypes.AssignPublicIpDisabled, runTask.NetworkConfiguration.AwsvpcConfiguration.AssignPublicIp, "AssignPublicIp")
		require.Contains(t, runTask.NetworkConfiguration.AwsvpcConfiguration.Subnets, "subnet-2222222222", "subnets")
		require.Contains(t, runTask.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups, "sg-11111111111", "security groups")
	})
}

func genRunTaskDef(t *testing.T, ctx *config.Context, entity *config.PreDeployTask) *ecs.RunTaskInput {
	t.Helper()
	runTaskInput, err := Build(ctx, entity)
	require.NoError(t, err)

	_, err = awsclients.ECSClient().RunTask(ctx.Context, runTaskInput)
	require.NoError(t, err)

	return runTaskInput
}
