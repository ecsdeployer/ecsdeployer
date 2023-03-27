package cleanupcronjobs

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct{}

func (Step) String() string {
	return "cleaning orphaned cronjobs"
}

func (Step) Skip(ctx *config.Context) bool {
	return !ctx.Project.Settings.KeepInSync.GetCronjobs()
}

func (Step) Clean(ctx *config.Context) error {

	if ctx.Project.Settings.CronUsesEventing {
		return runLegacyCleanup(ctx)
	}

	return runSchedulerCleanup(ctx)
}
