package taskdefinition

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type Builder struct {
	ctx    *config.Context
	entity config.IsTaskStruct

	commonTask   *config.CommonTaskAttrs
	project      *config.Project
	taskDefaults *config.FargateDefaults
	templates    *config.NameTemplates

	taskDef *ecs.RegisterTaskDefinitionInput
}

func NewBuilder(ctx *config.Context, entity config.IsTaskStruct) *Builder {
	builder := &Builder{
		ctx:    ctx,
		entity: entity,
	}

	builder.init()

	return builder
}

// Deprecated
func (td *Builder) commonTaskAttrs() *config.CommonTaskAttrs {
	common, err := config.ExtractCommonTaskAttrs(td.entity)
	if err != nil {
		panic(fmt.Errorf("no CommonTaskAttrs attribute?"))
	}
	return common
}

func (builder *Builder) init() {
	builder.project = builder.ctx.Project
	builder.taskDefaults = builder.project.TaskDefaults
	builder.templates = builder.project.Templates
	builder.commonTask = builder.commonTaskAttrs()
}

type taskLevelFunc func() error

func (builder *Builder) Build() (*ecs.RegisterTaskDefinitionInput, error) {
	// builder.init()

	builder.taskDef = &ecs.RegisterTaskDefinitionInput{
		Family:                  new(string),
		NetworkMode:             ecsTypes.NetworkModeAwsvpc,
		ContainerDefinitions:    []ecsTypes.ContainerDefinition{},
		RequiresCompatibilities: []ecsTypes.Compatibility{ecsTypes.CompatibilityFargate},
		RuntimePlatform:         &ecsTypes.RuntimePlatform{},
		// ProxyConfiguration:      &ecsTypes.ProxyConfiguration{},
		// Tags:                    []ecsTypes.Tag{},
		// TaskRoleArn:             new(string),
		// ExecutionRoleArn:        new(string),
		// Volumes:                 []ecsTypes.Volume{},
		// Cpu:                     new(string),
		// Memory:                  new(string),
		// EphemeralStorage:        &ecsTypes.EphemeralStorage{},
	}

	for _, funcName := range []taskLevelFunc{
		builder.applyRoles,
		builder.applyTaskResources,
		builder.applyTags,
	} {
		if err := funcName(); err != nil {
			return nil, err
		}
	}

	return builder.taskDef, nil
}
