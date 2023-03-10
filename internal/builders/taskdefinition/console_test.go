package taskdefinition_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/rshell"
	"ecsdeployer.com/ecsdeployer/internal/testutil/buildtestutils"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/stretchr/testify/require"
)

func TestConsole(t *testing.T) {

	buildtestutils.StartMocker(t)

	t.Run("default", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/everything.yml", buildtestutils.OptSetNumSSMVars(4))

		consoleTask := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, consoleTask)

		require.NotNil(t, taskDefinition)

		primary, err := buildtestutils.GetContainer(taskDefinition, "console")
		require.NoError(t, err)

		require.NotNil(t, primary.LinuxParameters)
		require.NotNil(t, primary.LinuxParameters.InitProcessEnabled)
		require.True(t, *primary.LinuxParameters.InitProcessEnabled)

		expectedRshell, _ := util.Jsonify(rshell.DockerLabel{
			Cluster:          "fakecluster",
			SubnetIds:        []string{"subnet-2222222222"},
			SecurityGroupIds: []string{"sg-11111111111"},
			AssignPublicIp:   false,
			Port:             8722,
		})

		require.Contains(t, primary.DockerLabels, rshell.LabelName)
		require.JSONEq(t, expectedRshell, primary.DockerLabels[rshell.LabelName])

	})

	t.Run("console with awslogs", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/awslogs.yml", buildtestutils.OptSetNumSSMVars(2))

		task := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, task)
		require.NotNil(t, taskDefinition)
	})

	t.Run("console with splunk", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/customlog.yml", buildtestutils.OptSetNumSSMVars(2))

		task := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, task)
		require.NotNil(t, taskDefinition)
	})
}
