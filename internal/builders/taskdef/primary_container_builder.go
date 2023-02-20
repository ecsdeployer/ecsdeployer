package taskdef

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/builders/pipeline"
	"ecsdeployer.com/ecsdeployer/internal/builders/taskdef/containers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type PrimaryContainerBuilder struct {
	Tpl *tmpl.Template

	Resource config.IsTaskStruct
}

func (pc *PrimaryContainerBuilder) Apply(obj *pipeline.PipeItem[ecs.RegisterTaskDefinitionInput]) error {

	primaryContainerPipe := pipeline.NewPipeItem(obj.Context, &ecsTypes.ContainerDefinition{})

	common, err := config.ExtractCommonTaskAttrs(pc.Resource)
	if err != nil {
		return err
	}

	ctx := obj.Context
	project := ctx.Project
	taskDefaults := project.TaskDefaults

	containerName, err := pc.Tpl.Apply(*project.Templates.ContainerName)
	if err != nil {
		return err
	}
	pc.Tpl = pc.Tpl.WithExtraField("ContainerName", containerName)

	err = primaryContainerPipe.Apply(
		&containers.ContainerDefaultsBuilder{Resource: common.CommonContainerAttrs},
		&containers.ConsoleBuilder{Resource: pc.Resource},
		&containers.PortMappingsBuilder{Resource: pc.Resource},
	)
	if err != nil {
		return err
	}

	image := util.Coalesce(common.Image, taskDefaults.Image, project.Image)
	if image == nil {
		return errors.New("You have not specified an image to deploy")
	}

	primaryContainerPipe.Data.Image = aws.String(image.Value())

	primaryContainerPipe.Data.Name = aws.String(containerName)

	if obj.Data.ContainerDefinitions == nil {
		obj.Data.ContainerDefinitions = make([]ecsTypes.ContainerDefinition, 0)
	}
	obj.Data.ContainerDefinitions = append(obj.Data.ContainerDefinitions, *primaryContainerPipe.GetData())

	return nil
}
