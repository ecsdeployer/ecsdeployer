package service

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func (b *Builder) applyCapacityStrategy() error {

	spotOverride := util.Coalesce(b.entity.SpotOverride, b.taskDefaults.SpotOverride, &config.SpotOverrides{})

	b.serviceDef.CapacityProviderStrategy = spotOverride.ExportCapacityStrategy()

	return nil
}
