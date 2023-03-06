package buildtestutils

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/builders/taskdefinition"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func GenTaskDef(t *testing.T, ctx *config.Context, entity config.IsTaskStruct) *ecs.RegisterTaskDefinitionInput {
	t.Helper()
	taskDefinition, err := taskdefinition.Build(ctx, entity)
	require.NoError(t, err)

	_, err = awsclients.ECSClient().RegisterTaskDefinition(ctx.Context, taskDefinition)
	require.NoError(t, err)

	return taskDefinition
}

func GetContainer(taskDef *ecs.RegisterTaskDefinitionInput, containerName string) (ecsTypes.ContainerDefinition, error) {

	for _, container := range taskDef.ContainerDefinitions {
		if container.Name == nil {
			continue
		}

		if *container.Name == containerName {
			return container, nil
		}
	}

	return ecsTypes.ContainerDefinition{}, fmt.Errorf("could not find container '%s'", containerName)
}

func GetPredeployTask(project *config.Project, name string) *config.PreDeployTask {
	for _, task := range project.PreDeployTasks {
		if task.Name == name {
			return task
		}
	}

	panic("FAILED TO FIND PREDEPLOY TASK")
}

func GetServiceTask(project *config.Project, name string) *config.Service {
	for _, task := range project.Services {
		if task.Name == name {
			return task
		}
	}

	panic("FAILED TO FIND SERVICE TASK")
}
