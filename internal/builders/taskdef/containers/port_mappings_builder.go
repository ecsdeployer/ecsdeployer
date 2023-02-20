package containers

import (
	"ecsdeployer.com/ecsdeployer/internal/builders/pipeline"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type PortMappingsBuilder struct {
	Resource config.IsTaskStruct
}

func (builder *PortMappingsBuilder) Apply(obj *pipeline.PipeItem[ecsTypes.ContainerDefinition]) error {

	service, isService := (builder.Resource).(*config.Service)
	if !isService {
		return nil
	}

	if !service.IsLoadBalanced() {
		return nil
	}

	if obj.Data.PortMappings == nil {
		obj.Data.PortMappings = make([]ecsTypes.PortMapping, 0)
	}

	for _, lb := range service.LoadBalancers {
		obj.Data.PortMappings = append(obj.Data.PortMappings, lb.PortMapping.ToAwsPortMapping())
	}

	return nil
}
