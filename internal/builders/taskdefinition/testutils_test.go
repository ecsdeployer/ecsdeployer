package taskdefinition_test

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

func genTaskDef(t *testing.T, ctx *config.Context, entity config.IsTaskStruct) *ecs.RegisterTaskDefinitionInput {
	t.Helper()
	taskDefinition, err := taskdefinition.Build(ctx, entity)
	require.NoError(t, err)

	_, err = awsclients.ECSClient().RegisterTaskDefinition(ctx.Context, taskDefinition)
	require.NoError(t, err)

	return taskDefinition
}

func getContainer(taskDef *ecs.RegisterTaskDefinitionInput, containerName string) (ecsTypes.ContainerDefinition, error) {

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

type kvToMapFunc[T any] func(T) (string, string)

func kvListToMap[T any](kvList []T, mapFunc kvToMapFunc[T]) map[string]string {
	newMap := make(map[string]string, len(kvList))

	for _, entry := range kvList {
		k, v := mapFunc(entry)
		newMap[k] = v
	}

	return newMap
}

func kvListToMap_KVP(val ecsTypes.KeyValuePair) (string, string) {
	return *val.Name, *val.Value
}

type kvToSliceFunc[T any, K any] func(T) K

func kvListToSlice[T any, K any](kvList []T, sliceFunc kvToSliceFunc[T, K]) []K {
	newMap := make([]K, 0, len(kvList))

	for _, entry := range kvList {
		newMap = append(newMap, sliceFunc(entry))
	}

	return newMap
}

func kvListToSlice_PortMaps(val ecsTypes.PortMapping) string {
	return fmt.Sprintf("%d/%s", *val.ContainerPort, val.Protocol)
}

func getPredeployTask(project *config.Project, name string) *config.PreDeployTask {
	for _, task := range project.PreDeployTasks {
		if task.Name == name {
			return task
		}
	}

	panic("FAILED TO FIND PREDEPLOY TASK")
}

func getServiceTask(project *config.Project, name string) *config.Service {
	for _, task := range project.Services {
		if task.Name == name {
			return task
		}
	}

	panic("FAILED TO FIND SERVICE TASK")
}
