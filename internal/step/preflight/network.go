package preflight

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type checkNetwork struct{}

func (checkNetwork) String() string {
	return "network"
}

func (checkNetwork) Check(ctx *config.Context) error {

	for _, network := range util.DeepFindInStruct[config.NetworkConfiguration](ctx.Project) {
		if err := network.Resolve(ctx, nil); err != nil {
			return err
		}
	}

	return nil
}
