package service

import (
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

const (
	containerNameTplField = "Container"
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
	entity *config.Service

	commonTask   *config.CommonTaskAttrs
	project      *config.Project
	taskDefaults *config.FargateDefaults
	templates    *config.NameTemplates

	commonTplFields tmpl.Fields

	serviceDef *ecs.CreateServiceInput
}

func BuildCreate(ctx *config.Context, entity *config.Service) (*ecs.CreateServiceInput, error) {
	builder, err := newBuilder(ctx, entity)
	if err != nil {
		return nil, err
	}

	return builder.Build()
}

func newBuilder(ctx *config.Context, entity *config.Service) (*Builder, error) {
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

	if commonTplFields, err := helpers.GetDefaultTaskTemplateFields(builder.ctx, builder.commonTask); err != nil {
		return err
	} else {
		builder.commonTplFields = commonTplFields
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

type applierFunc func() error

func (builder *Builder) Build() (*ecs.CreateServiceInput, error) {
	builder.serviceDef = &ecs.CreateServiceInput{}

	for _, funcName := range []applierFunc{
		builder.applyServiceDefaults,
		builder.applyCapacityStrategy,
		builder.applyLoadBalancers,
		builder.applyNetworking,
		builder.applyTags,
	} {
		if err := funcName(); err != nil {
			return nil, err
		}
	}

	return builder.serviceDef, nil
}
