package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type TaskDefinition struct {
	Context *config.Context
	Entity  config.IsTaskStruct

	Name          string
	Role          *config.RoleArn
	ExecutionRole *config.RoleArn

	Cpu     config.CpuSpec
	Memory  config.MemorySpec
	Storage config.StorageSpec

	Tags []config.NameValuePair

	PrimaryContainer *ContainerDefinition

	Sidecars []*ContainerDefinition
}

func NewTaskDefinition()

func (td *TaskDefinition) LinkContainer(cd *ContainerDefinition) *ContainerDefinition {
	cd.ParentTaskDef = td
	return cd
}

func (td *TaskDefinition) Export() (*ecs.RegisterTaskDefinitionInput, error) {
	return nil, nil
}
