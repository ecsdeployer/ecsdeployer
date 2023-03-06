package taskdefinition_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/rshell"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestConsole(t *testing.T) {

	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: []*awsmocker.MockedEndpoint{
			// Mock_ECS_RegisterTaskDefinition_Dump(t),
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		},
	})

	t.Run("default", func(t *testing.T) {
		ctx := loadProjectConfig(t, "everything.yml", optSetNumSSMVars(4))

		consoleTask := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, consoleTask)

		require.NotNil(t, taskDefinition)

		primary, err := getContainer(taskDefinition, "console")
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
		ctx := loadProjectConfig(t, "awslogs.yml", optSetNumSSMVars(2))

		task := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, task)
		require.NotNil(t, taskDefinition)
	})

	t.Run("console with splunk", func(t *testing.T) {
		ctx := loadProjectConfig(t, "customlog.yml", optSetNumSSMVars(2))

		task := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, task)
		require.NotNil(t, taskDefinition)
	})
}
