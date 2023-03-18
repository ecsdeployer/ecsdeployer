package cronjob

import (
	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct {
	cronjob *config.CronJob
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

	return step.Skipf("NOT FINISHED - %s", s.cronjob.Name)
}

func (s *Step) getSchedule(ctx *config.Context) (bool, error) {
	return false, nil
}

func (s *Step) createSchedule(ctx *config.Context) error {
	return nil
}

func (s *Step) updateSchedule(ctx *config.Context) error {
	return nil
}
