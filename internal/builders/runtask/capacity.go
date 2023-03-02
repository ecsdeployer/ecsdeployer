package runtask

import (
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyCapacityStrategy() error {

	b.runTaskDef.LaunchType = ecsTypes.LaunchTypeFargate

	// OR
	// b.runTaskDef.CapacityProviderStrategy = config.NewSpotOnDemand().ExportCapacityStrategy()

	return nil
}
