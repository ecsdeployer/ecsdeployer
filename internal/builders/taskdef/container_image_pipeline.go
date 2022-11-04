package taskdef

import (
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func ContainerImagePipeline(input *pipelineInput) error {

	common := input.Common

	tpl := tmpl.New(input.Context).WithExtraFields(tmpl.Fields{
		"Name":     common.Name,
		"TaskName": common.Name,
	})

	for i, cont := range input.TaskDef.ContainerDefinitions {

		if cont.Image == nil {
			// not sure how this will ever work
			continue
		}

		imgUri := *cont.Image

		imgTpld, err := tpl.Apply(imgUri)
		if err != nil {
			return err
		}

		input.TaskDef.ContainerDefinitions[i].Image = aws.String(imgTpld)
	}

	return nil
}
