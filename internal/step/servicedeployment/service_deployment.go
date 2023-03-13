package servicedeployment

import (
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

	return nil
}
