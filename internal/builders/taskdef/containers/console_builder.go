package containers

import (
	"ecsdeployer.com/ecsdeployer/internal/builders/pipeline"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ConsoleBuilder struct {
}

func (cb *ConsoleBuilder) Apply(obj *pipeline.PipeItem[ecsTypes.ContainerDefinition]) error {
	return ErrNotImplemented
}
