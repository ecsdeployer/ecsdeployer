package buildtestutils

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

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
