package taskdefinition

import ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"

func (b *Builder) applyContainers() error {

	if b.taskDef.ContainerDefinitions == nil {
		b.taskDef.ContainerDefinitions = make([]ecsTypes.ContainerDefinition, 0)
	}

	// first container is always primary
	b.taskDef.ContainerDefinitions = append(b.taskDef.ContainerDefinitions, *b.primaryContainer)

	// process sidecars
	if len(b.sidecars) > 0 {
		for _, sidecar := range b.sidecars {
			sidecar := sidecar
			b.taskDef.ContainerDefinitions = append(b.taskDef.ContainerDefinitions, *sidecar)
		}
	}

	// add logging
	if b.loggingContainer != nil {
		b.taskDef.ContainerDefinitions = append(b.taskDef.ContainerDefinitions, *b.loggingContainer)
	}

	return nil
}
