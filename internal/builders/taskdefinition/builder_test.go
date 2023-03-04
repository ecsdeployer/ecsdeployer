package taskdefinition_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

/*
TESTS TO DO:

- service
- console
- cronjob
- predeploy task

- task with firelens
- task with awslogs
- task with SSM vars
- task without ssm vars
- task with/without deployment vars
*/

func TestTaskDefinitionBuilder(t *testing.T) {

	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: []*awsmocker.MockedEndpoint{
			// mock_ECS_RegisterTaskDefinition_Dump(t),
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		},
	})

	t.Run("load balanced service", func(t *testing.T) {
		ctx := loadProjectConfig(t, "dummy.yml", optSetNumSSMVars(4))

		lbService := ctx.Project.Services[0]

		taskDefinition := genTaskDef(t, ctx, lbService)

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

	t.Run("everything", func(t *testing.T) {
		ctx := loadProjectConfig(t, "everything.yml", optSetNumSSMVars(2))

		tables := []struct {
			entity config.IsTaskStruct
		}{
			{ctx.Project.ConsoleTask},
			{ctx.Project.PreDeployTasks[0]},
			{ctx.Project.PreDeployTasks[1]},
			{ctx.Project.Services[0]},
			{ctx.Project.Services[1]},
			{ctx.Project.Services[2]},
			{ctx.Project.Services[3]},
			{ctx.Project.CronJobs[0]},
		}

		for _, table := range tables {
			t.Run(fmt.Sprintf("sub_%T_%s", table.entity, table.entity.GetCommonContainerAttrs().Name), func(t *testing.T) {
				taskDefinition := genTaskDef(t, ctx, table.entity)
				require.NotNil(t, taskDefinition)
			})
		}

	})

	t.Run("with proxy", func(t *testing.T) {
		ctx := loadProjectConfig(t, "everything.yml", optSetNumSSMVars(2))

		pdTest1Yaml := `
		name: testpd1
		command: "something something"
		proxy:
			properties:
				Blah: yar
		`

		pdTask, err := yaml.ParseYAMLString[config.PreDeployTask](testutil.CleanTestYaml(pdTest1Yaml))
		require.NoError(t, err)

		taskDefinition := genTaskDef(t, ctx, pdTask)

		require.NotNil(t, taskDefinition.ProxyConfiguration)
		require.EqualValues(t, "APPMESH", taskDefinition.ProxyConfiguration.Type, "ProxyType")
		require.EqualValues(t, "envoy", *taskDefinition.ProxyConfiguration.ContainerName, "ProxyContainer")

		propMap := kvListToMap(taskDefinition.ProxyConfiguration.Properties, kvListToMap_KVP)

		require.Contains(t, propMap, "Blah")
		require.Equal(t, "yar", propMap["Blah"])

	})

	t.Run("storage", func(t *testing.T) {
		ctx := loadProjectConfig(t, "everything.yml")

		taskDefinition := genTaskDef(t, ctx, getPredeployTask(ctx.Project, "pd-storage"))

		require.NotNil(t, taskDefinition.EphemeralStorage)
		require.EqualValues(t, 50, taskDefinition.EphemeralStorage.SizeInGiB)
	})
}
