package taskdefinition

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type templater interface {
	Apply(string) (string, error)
}

type Task struct {
	Context *config.Context
	Entity  config.IsTaskStruct

	Name          string
	Role          *config.RoleArn
	ExecutionRole *config.RoleArn

	Arch config.Architecture

	Cpu     *config.CpuSpec
	Memory  *config.MemorySpec
	Storage *config.StorageSpec

	Tags []config.NameValuePair

	PrimaryContainer *Container

	Sidecars []*Container

	Tpl templater
}

func (td *Task) CommonAttrs() *config.CommonTaskAttrs {
	common, err := config.ExtractCommonTaskAttrs(td.Entity)
	if err != nil {
		panic(fmt.Errorf("no CommonTaskAttrs attribute?"))
	}
	return common
}

// func NewTask(ctx *config.Context, entity config.IsTaskStruct) *Task {
// 	tpl, err := tmpl.New(ctx).
// 	return &Task{
// 		Context: ctx,
// 		Entity:  entity,
// 	}
// }

func (td *Task) NewContainer() *Container {
	return td.LinkContainer(&Container{})
}

func (td *Task) LinkContainer(cd *Container) *Container {
	cd.ParentTask = td
	return cd
}

func (td *Task) EvalTpl(tplStr string) (string, error) {
	return td.Tpl.Apply(tplStr)
}

func (td *Task) Export() (*ecs.RegisterTaskDefinitionInput, error) {
	taskdef := &ecs.RegisterTaskDefinitionInput{
		NetworkMode:             ecsTypes.NetworkModeAwsvpc,
		ContainerDefinitions:    make([]ecsTypes.ContainerDefinition, 0, len(td.Sidecars)+1),
		Family:                  new(string),
		Cpu:                     new(string),
		Memory:                  new(string),
		EphemeralStorage:        &ecsTypes.EphemeralStorage{},
		ExecutionRoleArn:        new(string),
		RequiresCompatibilities: []ecsTypes.Compatibility{ecsTypes.CompatibilityFargate},
		RuntimePlatform:         &ecsTypes.RuntimePlatform{},
		ProxyConfiguration:      &ecsTypes.ProxyConfiguration{},
		Tags:                    []ecsTypes.Tag{},
		TaskRoleArn:             new(string),
		Volumes:                 []ecsTypes.Volume{},
	}

	return taskdef, nil
}
