// this will delete services that are no longer used
package cleanupservices

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct{}

func (Step) String() string {
	return "cleaning orphaned services"
}

func (Step) Skip(ctx *config.Context) bool {
	return !ctx.Project.Settings.KeepInSync.GetServices()
}

func (Step) Clean(ctx *config.Context) error {

	return nil
}
