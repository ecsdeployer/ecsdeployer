package cleanupcronjobs

import (
	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct{}

func (Step) String() string {
	return "cleaning orphaned cronjobs"
}

func (Step) Skip(ctx *config.Context) bool {
	return !ctx.Project.Settings.KeepInSync.GetCronjobs() || ctx.Project.Settings.CronUsesEventing
}

func (Step) Clean(ctx *config.Context) error {
	return step.Skip("NOT FINISHED")
	// return nil
}
