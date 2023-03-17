package servicedeployment

import (
	"ecsdeployer.com/ecsdeployer/internal/semerrgroup"
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

	g := semerrgroup.New(5)

	for _, service := range ctx.Project.Services {
		service := service
		g.Go(func() error {
			return deployService(ctx, service)
		})
	}

	return g.Wait()
}

func deployService(ctx *config.Context, service *config.Service) error {
	// log.WithField("name", service.Name).Debug("deploying")

	existingSvc, err := describeService(ctx, service)
	if err != nil {
		return err
	}

	if existingSvc == nil {
		return createService(ctx, service)
	}

	return updateService(ctx, service)
}
