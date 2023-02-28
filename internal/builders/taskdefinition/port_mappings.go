package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyServicePortMappings() error {

	service, isService := (b.entity).(*config.Service)
	if !isService {
		return nil
	}

	if !service.IsLoadBalanced() {
		return nil
	}

	if b.primaryContainer.PortMappings == nil {
		b.primaryContainer.PortMappings = make([]ecsTypes.PortMapping, 0)
	}

	for _, lb := range service.LoadBalancers {
		b.primaryContainer.PortMappings = append(b.primaryContainer.PortMappings, lb.PortMapping.ToAwsPortMapping())
	}

	return nil
}

func (b *Builder) applySidecarPortMappings(cdef *ecsTypes.ContainerDefinition, sidecar *config.Sidecar) error {

	if sidecar.PortMappings == nil || len(sidecar.PortMappings) == 0 {
		return nil
	}

	if cdef.PortMappings == nil {
		cdef.PortMappings = make([]ecsTypes.PortMapping, 0)
	}

	for _, pm := range sidecar.PortMappings {
		cdef.PortMappings = append(cdef.PortMappings, pm.ToAwsPortMapping())
	}

	return nil
}
