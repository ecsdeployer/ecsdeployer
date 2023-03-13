package preflight

import (
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type checkContainerImages struct{}

func (checkContainerImages) String() string {
	return "container images"
}

func (checkContainerImages) Check(ctx *config.Context) error {

	for _, image := range util.DeepFindInStruct[config.ImageUri](ctx.Project) {
		if _, err := helpers.ResolveImageUri(ctx, image); err != nil {
			return err
		}
	}

	return nil
}
