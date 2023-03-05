package service

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/util"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyNetworking() error {
	network := util.Coalesce(b.entity.Network, b.taskDefaults.Network, b.project.Network)
	if network == nil {
		return errors.New("Unable to resolve network configuration!")
	}

	b.serviceDef.NetworkConfiguration = &ecsTypes.NetworkConfiguration{}
	if err := network.Resolve(b.ctx, b.serviceDef.NetworkConfiguration); err != nil {
		return err
	}

	return nil
}
