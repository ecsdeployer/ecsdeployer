package taskdef

import (
	"ecsdeployer.com/ecsdeployer/internal/builders/pipeline"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

// This takes a task definition and just resolves all the container images within it
type ContainerImageResolver struct {
	Tpl *tmpl.Template
	// Resource config.IsTaskStruct
}

func (builder *ContainerImageResolver) Apply(pi *pipeline.PipeItem[ecs.RegisterTaskDefinitionInput]) error {

	for i, container := range pi.Data.ContainerDefinitions {
		if container.Image == nil {
			// not sure how this will ever work
			continue
		}

		templatedImageUri, err := builder.Tpl.Apply(*container.Image)
		if err != nil {
			return err
		}

		pi.Data.ContainerDefinitions[i].Image = &templatedImageUri
	}

	return nil
}
