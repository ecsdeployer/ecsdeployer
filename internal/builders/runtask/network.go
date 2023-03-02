package runtask

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/util"
)

func (b *Builder) applyNetworking() error {
	network := util.Coalesce(b.entity.Network, b.taskDefaults.Network, b.project.Network)
	if network == nil {
		return errors.New("Unable to resolve network configuration!")
	}
	ecsNetworkConfig, err := network.ResolveECS(b.ctx)
	if err != nil {
		return err
	}
	b.runTaskDef.NetworkConfiguration = ecsNetworkConfig

	return nil
}
