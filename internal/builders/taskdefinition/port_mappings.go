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

	pmList := make([]config.PortMapping, 0)

	for _, lb := range service.LoadBalancers {
		pmList = append(pmList, *lb.PortMapping)
	}

	return assignPortMappings(b.primaryContainer, pmList)
}

func (b *Builder) applySidecarPortMappings(cdef *ecsTypes.ContainerDefinition, sidecar *config.Sidecar) error {
	return assignPortMappings(cdef, sidecar.PortMappings)
}

func assignPortMappings(cdef *ecsTypes.ContainerDefinition, mappings []config.PortMapping) error {
	if len(mappings) == 0 {
		return nil
	}

	if cdef.PortMappings == nil {
		cdef.PortMappings = make([]ecsTypes.PortMapping, 0, len(mappings))
	}

	for _, pm := range mappings {
		if !containerHasPortMapping(cdef, pm) {
			cdef.PortMappings = append(cdef.PortMappings, pm.ToAwsPortMapping())
		}
	}

	return nil
}

func containerHasPortMapping(cdef *ecsTypes.ContainerDefinition, mapping config.PortMapping) bool {
	if cdef.PortMappings == nil || len(cdef.PortMappings) == 0 {
		return false
	}

	for _, pm := range cdef.PortMappings {
		if mapping.Protocol == pm.Protocol && *pm.ContainerPort == *mapping.Port {
			return true
		}
	}

	return false
}
