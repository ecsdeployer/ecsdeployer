package servicedeployment

import (
	"ecsdeployer.com/ecsdeployer/internal/semerrgroup"
	"ecsdeployer.com/ecsdeployer/internal/step/service"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct{}

func (Step) String() string {
	return "deploying services"
}

func (Step) Skip(ctx *config.Context) bool {
	return len(ctx.Project.Services) == 0
}

func (Step) Run(ctx *config.Context) error {

	g := semerrgroup.New(ctx.Concurrency(5))

	for _, svc := range ctx.Project.Services {
		svc := svc
		g.Go(func() error {
			return service.New(svc).Run(ctx)
		})
	}

	return g.Wait()
}
