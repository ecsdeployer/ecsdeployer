package cronjob

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct {
	groupName    string
	scheduleName string
	taskDefArn   string
	cronjob      *config.CronJob
}

func New(cronjob *config.CronJob) *Step {
	return &Step{
		cronjob: cronjob,
	}
}

func (Step) String() string {
	return "cronjob"
}

func (s *Step) Run(ctx *config.Context) error {

	if ctx.Project.Settings.CronUsesEventing {
		return s.runLegacyVariant(ctx)
	}

	return s.runSchedulerVariant(ctx)
}
