package crondeployment

import (
	"ecsdeployer.com/ecsdeployer/internal/deprecate"
	"ecsdeployer.com/ecsdeployer/internal/semerrgroup"
	"ecsdeployer.com/ecsdeployer/internal/step/cronjob"
	"ecsdeployer.com/ecsdeployer/internal/step/schedulegroup"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct{}

func (Step) String() string {
	return "deploying cronjobs"
}

func (Step) Skip(ctx *config.Context) bool {
	return len(ctx.Project.CronJobs) == 0
}

func (Step) Run(ctx *config.Context) error {

	if ctx.Project.Settings.CronUsesEventing {
		deprecate.Deprecate_LegacyCron(ctx)
	}

	if err := (schedulegroup.Step{}).Run(ctx); err != nil {
		return err
	}

	g := semerrgroup.NewSkipAware(semerrgroup.New(ctx.Concurrency(5)))

	for _, job := range ctx.Project.CronJobs {
		job := job
		g.Go(func() error {
			return cronjob.New(job).Run(ctx)
		})
	}

	return g.Wait()
}
