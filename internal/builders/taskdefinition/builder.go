package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

const (
	containerNameTplField = "ContainerName"
)

type hasContainerAttrs interface {
	GetCommonContainerAttrs() config.CommonContainerAttrs
}

type hasTaskAttrs interface {
	hasContainerAttrs
	GetCommonTaskAttrs() config.CommonTaskAttrs
}

type Builder struct {
	ctx    *config.Context
	entity config.IsTaskStruct

	commonTask   *config.CommonTaskAttrs
	project      *config.Project
	taskDefaults *config.FargateDefaults
	templates    *config.NameTemplates

	deploymentEnvVars config.EnvVarMap
	baseEnvVars       config.EnvVarMap

	primaryContainer *ecsTypes.ContainerDefinition
	loggingContainer *ecsTypes.ContainerDefinition

	sidecars []*ecsTypes.ContainerDefinition

	commonTplFields tmpl.Fields

	taskDef *ecs.RegisterTaskDefinitionInput
}

func NewBuilder(ctx *config.Context, entity config.IsTaskStruct) (*Builder, error) {
	builder := &Builder{
		ctx:    ctx,
		entity: entity,
	}

	if err := builder.init(); err != nil {
		return nil, err
	}

	return builder, nil
}

func (builder *Builder) init() error {
	builder.project = builder.ctx.Project
	builder.taskDefaults = builder.project.TaskDefaults
	builder.templates = builder.project.Templates
	builder.commonTask = util.Ptr(builder.entity.GetCommonTaskAttrs())
	builder.sidecars = make([]*ecsTypes.ContainerDefinition, 0)

	if commonTplFields, err := helpers.GetDefaultTaskTemplateFields(builder.ctx, builder.commonTask); err != nil {
		return err
	} else {
		builder.commonTplFields = commonTplFields
	}

	if err := builder.createDeploymentEnvVars(); err != nil {
		return err
	}
	if err := builder.createTaskEnvVars(); err != nil {
		return err
	}

	return nil
}

func (builder *Builder) tpl() *tmpl.Template {
	// inefficient, but safer and we dont need efficiency
	return tmpl.New(builder.ctx).WithExtraFields(builder.commonTplFields)
}

func (builder *Builder) tplEval(tplStr string) (string, error) {
	// inefficient, but safer and we dont need efficiency
	retval, err := builder.tpl().Apply(tplStr)
	if err != nil {
		return "", err
	}
	return retval, nil
}

func (builder *Builder) containerTplEval(containerName, tplStr string) (string, error) {
	// inefficient, but safer and we dont need efficiency
	retval, err := builder.tpl().WithExtraField(containerNameTplField, containerName).Apply(tplStr)
	if err != nil {
		return "", err
	}
	return retval, nil
}

type taskLevelFunc func() error

func (builder *Builder) Build() (*ecs.RegisterTaskDefinitionInput, error) {
	// builder.init()

	builder.taskDef = &ecs.RegisterTaskDefinitionInput{}

	for _, funcName := range []taskLevelFunc{
		builder.applyTaskDefaults,
		builder.applyRoles,
		builder.applyTaskResources,
		builder.applyTags,

		builder.applyLoggingFirelensContainer, // must be before any other containers

		builder.applyPrimaryContainer,
		builder.applyServicePortMappings,
		builder.applyRemoteShell,

		builder.applySidecarContainers,
		builder.applyContainers,

		builder.applyVolumes,

		builder.applyCleanup,
	} {
		if err := funcName(); err != nil {
			return nil, err
		}
	}

	return builder.taskDef, nil
}
