package taskdef

import (
	"ecsdeployer.com/ecsdeployer/internal/builders/pipeline"
	"ecsdeployer.com/ecsdeployer/internal/builders/taskdef/containers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type PrimaryContainerBuilder struct {
	Tpl *tmpl.Template

	Resource config.IsTaskStruct
}

func (pc *PrimaryContainerBuilder) Apply(obj *pipeline.PipeItem[ecsTypes.ContainerDefinition]) error {

	cont := obj.Data
	if cont == nil {
		cont = &ecsTypes.ContainerDefinition{}
	}

	common, err := config.ExtractCommonTaskAttrs(pc.Resource)
	if err != nil {
		return err
	}
	_ = common

	ctx := obj.Context
	project := ctx.Project

	containerName, err := pc.Tpl.Apply(*project.Templates.ContainerName)
	if err != nil {
		return err
	}
	pc.Tpl = pc.Tpl.WithExtraField("ContainerName", containerName)

	contPi := pipeline.NewPipeItem(cont)

	err = contPi.Apply(&containers.ConsoleBuilder{})
	if err != nil {
		return err
	}

	cont = contPi.GetData()

	cont.Name = aws.String(containerName)

	obj.Data = cont
	return ErrNotImplemented
}
