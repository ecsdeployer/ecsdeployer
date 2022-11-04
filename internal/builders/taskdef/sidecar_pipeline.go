package taskdef

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func SidecarPipeline(input *pipelineInput) error {

	// TODO: add sidecar containers
	common := input.Common

	if len(common.Sidecars) == 0 {
		return nil
	}

	if true {
		// TODO: FINISH THIS
		return errors.New("SIDECAR PIPELINE NOT YET FINISHED")
	}

	contDefs := []ecsTypes.ContainerDefinition{}

	for _, sc := range common.Sidecars {
		conDef, err := buildSidecarContainer(input, sc)
		if err != nil {
			return err
		}

		if conDef == nil {
			continue
		}

		contDefs = append(contDefs, *conDef)
	}

	// Do this as the absolute last step so that we aren't modifying it if we encounter errors
	input.TaskDef.ContainerDefinitions = append(input.TaskDef.ContainerDefinitions, contDefs...)

	return nil
}

func buildSidecarContainer(input *pipelineInput, sidecar *config.Sidecar) (*ecsTypes.ContainerDefinition, error) {

	scDef, err := ContainerDefinitionBuilder(input.Context, sidecar.CommonContainerAttrs)
	if err != nil {
		return nil, err
	}

	// PORT MAPPINGS
	if len(sidecar.PortMappings) > 0 {
		for _, pm := range sidecar.PortMappings {
			scDef.PortMappings = append(scDef.PortMappings, pm.ToAwsPortMapping())
		}
	}

	return scDef, nil
}
